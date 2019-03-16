#include "rosco.h"
#include <libgen.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/time.h>
#include <time.h>
#include <unistd.h>

enum command_idx {
    MC_Read = 0,
    MC_Read_Raw = 1,
    MC_Read_IAC = 2,
    MC_PTC = 3,
    MC_FuelPump = 4,
    MC_IAC_Close = 5,
    MC_IAC_Open = 6,
    MC_AC = 7,
    MC_Coil = 8,
    MC_Injectors = 9,
    MC_Interactive = 10,
    MC_MemsGauge_Read = 11,
    MC_Num_Commands = 12
};

static const char* commands[] = {
    "read",
    "read-raw",
    "read-iac",
    "ptc",
    "fuelpump",
    "iac-close",
    "iac-open",
    "ac",
    "coil",
    "injectors",
    "interactive",
    "memsgauge"
};

void printbuf(uint8_t* buf, unsigned int count)
{
    unsigned int idx = 0;

    while (idx < count) {
        idx += 1;
        printf("%02X ", buf[idx - 1]);
        if (idx % 16 == 0) {
            printf("\n");
        }
    }
    printf("\n");
}

int16_t readserial(mems_info* info, uint8_t* buffer, uint16_t quantity)
{
    int16_t bytesread = -1;

#if defined(WIN32)
    DWORD w32BytesRead = 0;

    if ((ReadFile(info->sd, (UCHAR*)buffer, quantity, &w32BytesRead, NULL) == TRUE) && (w32BytesRead > 0)) {
        bytesread = w32BytesRead;
    }
#else
    bytesread = read(info->sd, buffer, quantity);
#endif

    return bytesread;
}

int16_t writeserial(mems_info* info, uint8_t* buffer, uint16_t quantity)
{
    int16_t byteswritten = -1;

#if defined(WIN32)
    DWORD w32BytesWritten = 0;

    if ((WriteFile(info->sd, (UCHAR*)buffer, quantity, &w32BytesWritten, NULL) == TRUE) && (w32BytesWritten == quantity)) {
        byteswritten = w32BytesWritten;
    }
#else
    byteswritten = write(info->sd, buffer, quantity);
#endif

    return byteswritten;
}

char* current_time(void)
{
    time_t t = time(NULL);
    struct tm tm = *localtime(&t);
    struct timeval tv;
    static char buffer[50];

    gettimeofday(&tv, NULL);

    sprintf(buffer, "%4d-%02d-%02dT%02d:%02d:%02d.%03d", tm.tm_year + 1900, tm.tm_mon + 1, tm.tm_mday, tm.tm_hour, tm.tm_min, tm.tm_sec, tv.tv_usec / 1000);
    return buffer;
}

char* current_date(void)
{
    time_t t = time(NULL);
    struct tm tm = *localtime(&t);
    static char buffer[50];

    sprintf(buffer, "%4d-%02d-%02d-%02d%02d%02d", tm.tm_year + 1900, tm.tm_mon + 1, tm.tm_mday, tm.tm_hour, tm.tm_min, tm.tm_sec);
    return buffer;
}

bool prefix(const char pre, const char* str)
{
    return (str[0] == pre);
}

int read_config(readmems_config* config)
{
    FILE* file = fopen("readmems.cfg", "r"); /* should check the result */
    char line[256];
    char* key;
    char* value;
    char* search = "=";
    char comment = '#';
    int cmd_idx = -1;
    int len;

    if (file) {
        while (fgets(line, sizeof(line), file)) {
            if (strlen(line) > 1) {
                if (prefix(comment, line)) {
                    // skip comments
                } else {
                    key = strtok(line, search);
                    value = strtok(NULL, search);

                    // remove newline
                    len = strlen(value);
                    if (value[len - 1] == '\n')
                        value[len - 1] = 0;

                    if (strcasecmp(key, "port") == 0) {
                        config->port = strdup(value);
                    }

                    if (strcasecmp(key, "command") == 0) {
                        config->command = strdup(value);
                    }

                    if (strcasecmp(key, "output") == 0) {
                        config->output = strdup(value);
                    }

                    if (strcasecmp(key, "loop") == 0) {
                        config->loop = strdup(value);
                    }
                }
            }
        }

        cmd_idx = find_command(config->command);
    }

    fclose(file);

    return cmd_idx;
}

int find_command(char* command)
{
    int cmd_idx = 0;

    while ((cmd_idx < MC_Num_Commands) && (strcasecmp(command, commands[cmd_idx]) != 0)) {
        cmd_idx += 1;
    }

    return cmd_idx;
}

int read_command_line_config(readmems_config* config, int argc, char** argv)
{
    int cmd_idx = 0;
    librosco_version ver;

    ver = mems_get_lib_version();

    if (argc < 3) {
        printf("readmems using librosco v%d.%d.%d\n", ver.major, ver.minor, ver.patch);
        printf("Diagnostic utility using ROSCO protocol for MEMS 1.6 systems\n");
        printf("Usage: %s <serial device> <command> [read-loop-count]\n", basename(argv[0]));
        printf(" where <command> is one of the following:\n");

        for (cmd_idx = 0; cmd_idx < MC_Num_Commands; ++cmd_idx) {
            printf("\t%s\n", commands[cmd_idx]);
        }
        printf(" and [read-loop-count] is either a number or 'inf' to read forever.\n");

        return -1;
    }

    // locate the index of the current command
    config->command = strdup(argv[2]);
    cmd_idx = find_command(config->command);

    if (cmd_idx >= MC_Num_Commands) {
        printf("Invalid command: %s\n", argv[2]);
        return -1;
    }

    config->port = strdup(argv[1]);
    config->output = strdup("stdout");

    if (argc >= 4) {
        config->loop = strdup(argv[3]);
    }

    if (cmd_idx != MC_Interactive) {
        printf("Running command: %s\n", commands[cmd_idx]);
    }

    return cmd_idx;
}

char* open_file(FILE** fp)
{
    static char filename[256];

    sprintf(filename, "readmems-%s.log", current_date());

    // open the file for writing
    *fp = fopen(filename, "w");

    return filename;
}

int write_log(FILE** fp, char* line)
{
    if (*fp)
        return fprintf(*fp, "%s", line);
    else
        return -1;
}

void delete_file(char* filename)
{
    remove(filename);
}

bool interactive_mode(mems_info* info, uint8_t* response_buffer)
{
    size_t icmd_size = 8;
    char* icmd_buf_ptr;
    uint8_t icmd;
    ssize_t bytes_read = 0;
    ssize_t total_bytes_read = 0;
    bool quit = false;

    if ((icmd_buf_ptr = (char*)malloc(icmd_size)) != NULL) {
        printf("Enter a command (in hex) or 'quit'.\n> ");
        while (!quit && (fgets(icmd_buf_ptr, icmd_size, stdin) != NULL)) {
            if ((strncmp(icmd_buf_ptr, "q", 1) == 0) || (strncmp(icmd_buf_ptr, "quit", 4) == 0)) {
                quit = true;
            } else if (icmd_buf_ptr[0] != '\n' && icmd_buf_ptr[1] != '\r') {
                icmd = strtoul(icmd_buf_ptr, NULL, 16);
                if ((icmd >= 0) && (icmd <= 0xff)) {
                    if (writeserial(info, &icmd, 1) == 1) {
                        bytes_read = 0;
                        total_bytes_read = 0;
                        do {
                            bytes_read = readserial(info, response_buffer + total_bytes_read, 1);
                            total_bytes_read += bytes_read;
                        } while (bytes_read > 0);

                        if (total_bytes_read > 0) {
                            printbuf(response_buffer, total_bytes_read);
                        } else {
                            printf("No response from ECU.\n");
                        }
                    } else {
                        printf("Error: failed to write command byte to serial port.\n");
                    }
                } else {
                    printf("Error: command must be between 0x00 and 0xFF.\n");
                }
                printf("> ");
            } else {
                printf("> ");
            }
        }

        free(icmd_buf_ptr);
    } else {
        printf("Error allocating command buffer memory.\n");
    }

    return 0;
}

int main(int argc, char** argv)
{
    bool success = false;
    int cmd_idx = 0;
    mems_data data;
    mems_data_frame_80 frame80;
    mems_data_frame_7d frame7d;
    mems_info info;

    uint8_t* frameptr;
    uint8_t bufidx;
    uint8_t readval = 0;
    uint8_t iac_limit_count = 80; // number of times to re-send an IAC move command when
    FILE* fp = NULL;
    char log_line[256];
    bool connected = false;

    // the ECU is already reporting that the valve has
    // reached its requested position
    int read_loop_count = 1;
    bool read_inf = false;

    // this is twice as large as the micro's on-chip ROM, so it's probably sufficient
    uint8_t response_buffer[16384];

    char win32devicename[16];
    char* port;

    // read the config file for defaults
    readmems_config config;
    cmd_idx = read_config(&config);

    if (argc > 1) {
        // process the command line arguments
        cmd_idx = read_command_line_config(&config, argc, argv);

        if (cmd_idx < 0)
            return -1;
    }

    if (strcmp(config.loop, "inf") == 0) {
        read_inf = true;
    } else {
        read_loop_count = strtoul(config.loop, NULL, 0);
    }

    if (strcmp(config.output, "stdout") != 0) {
        config.output = open_file(&fp);
    }

    printf("Using config:\nport: %s\ncommand: %s (%d)\noutput: %s\nloop: %s\n", config.port, config.command, cmd_idx, config.output, config.loop);

    mems_init(&info);

#if defined(WIN32)
    // correct for microsoft's legacy nonsense by prefixing with "\\.\"
    strcpy(win32devicename, "\\\\.\\");
    strncat(win32devicename, config.port, 16);
    port = win32devicename;
#else
    port = config.port;
#endif

    printf("attempting to connect to %s\n", port);
    connected = mems_connect(&info, port);

    if (connected) {
        if (mems_init_link(&info, response_buffer)) {
            printf("ECU responded to D0 command with: %02X %02X %02X %02X\n", response_buffer[0], response_buffer[1], response_buffer[2], response_buffer[3]);

            switch (cmd_idx) {
            case MC_MemsGauge_Read:
                sprintf(log_line, "#time,engineSpeed,waterTemp,intakeAirTemp,throttleVoltage,manifoldPressure,idleBypassPos,mainVoltage,idleswitch,closedloop,lambdaVoltage_mV,intakeAirTempSensorFault,coolantTempSensorFault,fuelpumpCircuitFault,throttlepotCircuitFault\n");
                printf("%s", log_line);

                if (fp)
                    write_log(&fp, log_line);

                while (read_inf || (read_loop_count-- > 0)) {
                    if (mems_read(&info, &data)) {
                        sprintf(log_line, "%s,%u,%u,%u,%f,%f,%u,%f,%u,%u,%u,%u,%u,%u,%u\n",
                            current_time(),
                            data.engine_rpm,
                            data.coolant_temp_c,
                            data.intake_air_temp_c,
                            data.throttle_pot_voltage,
                            data.map_kpa,
                            data.iac_position,
                            data.battery_voltage,
                            data.idle_switch,
                            data.closed_loop,
                            data.lambda_voltage_mv,
                            data.intake_air_temp_sensor_fault,
                            data.coolant_temp_sensor_fault,
                            data.fuel_pump_circuit_fault,
                            data.throttle_pot_circuit_fault);

                        printf("%s", log_line);

                        if (fp)
                            write_log(&fp, log_line);

                        success = true;
                    }
                }
                break;

            case MC_Read:
                printf("executing read\n");

                while (read_inf || (read_loop_count-- > 0)) {
                    if (mems_read(&info, &data)) {
                        printf("RPM: %u\nCoolant (deg C): %u\nAmbient (deg C): %u\nIntake air (deg C): %u\n"
                               "Fuel temp (deg C): %u\nMAP (kPa): %f\nMain voltage: %f\nThrottle pot voltage: %f\n"
                               "Idle switch: %u\nPark/neutral switch: %u\nFault codes: %u\nIAC position: %u\n"
                               "-------------\n",
                            data.engine_rpm, data.coolant_temp_c, data.ambient_temp_c,
                            data.intake_air_temp_c, data.fuel_temp_c, data.map_kpa, data.battery_voltage,
                            data.throttle_pot_voltage, data.idle_switch, data.park_neutral_switch,
                            data.fault_codes, data.iac_position);
                        success = true;
                    }
                }
                break;

            case MC_Read_Raw:
                while (read_inf || (read_loop_count-- > 0)) {
                    if (mems_read_raw(&info, &frame80, &frame7d)) {
                        frameptr = (uint8_t*)&frame80;
                        printf("80: ");
                        for (bufidx = 0; bufidx < sizeof(mems_data_frame_80); ++bufidx) {
                            printf("%02X ", frameptr[bufidx]);
                        }
                        printf("\n");

                        frameptr = (uint8_t*)&frame7d;
                        printf("7D: ");
                        for (bufidx = 0; bufidx < sizeof(mems_data_frame_7d); ++bufidx) {
                            printf("%02X ", frameptr[bufidx]);
                        }
                        printf("\n");

                        success = true;
                    }
                }
                break;

            case MC_Read_IAC:
                if (mems_read_iac_position(&info, &readval)) {
                    printf("0x%02X\n", readval);
                    success = true;
                }
                break;

            case MC_PTC:
                if (mems_test_actuator(&info, MEMS_PTCRelayOn, NULL)) {
                    sleep(2);
                    success = mems_test_actuator(&info, MEMS_PTCRelayOff, NULL);
                }
                break;

            case MC_FuelPump:
                if (mems_test_actuator(&info, MEMS_FuelPumpOn, NULL)) {
                    sleep(2);
                    success = mems_test_actuator(&info, MEMS_FuelPumpOff, NULL);
                }
                break;

            case MC_IAC_Close:
                do {
                    success = mems_test_actuator(&info, MEMS_CloseIAC, &readval);

                    // For some reason, diagnostic tools will continue to send send the
                    // 'close' command many times after the IAC has already reached the
                    // fully-closed position. Emulate that behavior here.
                    if (success && (readval == 0x00)) {
                        iac_limit_count -= 1;
                    }
                } while (success && iac_limit_count);
                break;

            case MC_IAC_Open:
                // The SP Rover 1 pod considers a value of 0xB4 to represent an opened
                // IAC valve, so repeat the open command until the valve is opened to
                // that point.
                do {
                    success = mems_test_actuator(&info, MEMS_OpenIAC, &readval);
                } while (success && (readval < 0xB4));
                break;

            case MC_AC:
                if (mems_test_actuator(&info, MEMS_ACRelayOn, NULL)) {
                    sleep(2);
                    success = mems_test_actuator(&info, MEMS_ACRelayOff, NULL);
                }
                break;

            case MC_Coil:
                success = mems_test_actuator(&info, MEMS_FireCoil, NULL);
                break;

            case MC_Injectors:
                success = mems_test_actuator(&info, MEMS_TestInjectors, NULL);
                break;

            case MC_Interactive:
                success = interactive_mode(&info, response_buffer);
                break;

            default:
                printf("Error: invalid command\n");
                break;
            }
        } else {
            printf("Error in initialization sequence.\n");
        }
        mems_disconnect(&info);
    } else {
        printf("Error: could not open serial device (%s).\n", port);
        if (fp) {
            fclose(fp);
            delete_file(config.output);
        }
    }

    mems_cleanup(&info);

    if (fp)
        fclose(fp);

    return success ? 0 : -2;
}
