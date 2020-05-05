# MemsFCR 
MemsFCR is a MEMS 1.6 - 1.9 Fault Code Reader and Analytics application.

MemsFCR allows you to connect to your MEMS ECU and monitor the diagnostic parameters.
All MEMS data points are read and can be viewed directly. If any faults have been detected these will be displayed and can be cleared.
The data can be logged to a CSV file for easy analysis using Excel or Google Sheets.

The Make file supports building on Mac and Windows, although I have included pre-build binaries for Windows (64bit) and MacOS.
The application connects using the serial interface and supports the serial FTDI cables you can purchase for the Android MemsDiag application. Instructions on how to create this cable or alternatively build a bluetooth wireless interface are also included.

[Download MacOS MemsFCR](https://github.com/andrewdjackson/memsfcr/raw/master/dist/MacOS-MemsFCR.zip)

[Download Windows MemFCR](https://github.com/andrewdjackson/memsfcr/raw/master/dist/Windows-MemsFCR.zip)

### The Dashboard
![Dashboard](./resources/screenshots/dashboard.png)
This tab shows the live running parameters from the ECU. If a fault is detected the ECU Fault indicator will light. The list of detected fault codes is displayed below the gauges. Any faults detected can be cleared by clicking Clear Faults.

### Profiling
![Profiling](./resources/screenshots/profiling.png)
The profiling shows a running graph of the following metrics:
* Engine RPM
* Lambda Voltage
* Loop Indicator
* Coolant Temperature

These 4 parameters are very important in determining the running condition of the engine and are used a primary parameters by the ECU for adjusting fueling and timing. Status indicators at the top show faults or warnings calculated from expected events of anomalies in the data. 

### Adjustments
![Adjustments](./resources/screenshots/adjustments.png)
* Reset ECU
This resets all adjustments and learnt parameters such as Short and Long Term Fuel Trim. The ECU will re-learn these parameters. 

* Reset Adjustments
Any changes made to the adjustments can be reset by clicking this button.

* Idle Speed
Increases / decreases the idle running speed. Hot idle should settle arounf 850 RPM.
This is a facility which is built into the MEMS ECU to overcome some situations during the service life of a vehicle where it might be wished to eradicate a problem such as a rattle or engine wear by slightly increasing the idle speed of the engine. This offset adds approximately 50 RPM to the Idle target speed. The function can be removed by resetting the adaptive values.

* Idle Hot (Decay) Position
 This is the number of IACV steps from fully closed (0) which the ECU has learned as the correct position to maintain the target idle speed with a fully warmed up engine. If this value is outside the range 10 - 50 steps, then this is an indication of a possible fault condition or poor adjustment. This value can be forced for a short time using this function.

* Fuel Trim
This allows the adjustment of the Long Term Fuel Trim (LTFT) which is used to offset fuel injection of rich or lean running

* Ignition Advance
This is a facility to overcome some situations during the service life of a vehicle where it might be wished to eradicate a problem such as a low octane fuel being constantly used or engine wear by slightly advancing the ignition timing the idle speed of the engine. The function can be removed by resetting the adaptive values.

### Dataframes
![Dataframes](./resources/screenshots/dataframes.png)
A live view of the calculated parameters received from the ECU. Since not all the values in the MEMS dataframe have been mapped this may be updated as more information comes to light.

### Information
Useful information on the operating characteristics of the MEMS ECU

### To Do..
Still to do - the command / response loop can get blocked if multiple overlapping commands are sent. I need to disable the buttons whilst a command is in progress.

### How to build a cable

To build a cable that will connect to the ECU's diagnostic port, you will need three things:

A 5V (TTL level) USB to serial converter. I strongly recommend the FTDI TTL-232R-5V. Note that it's important to have a 5V converter -- a normal RS-232 port will not work, and neither will a converter that uses regular RS-232 voltage levels. The FTDI cable (or an equivalent, such as the GearMo 5V cable) is available from different retailers:

* FTDI TTL-232R-5V-WE from Mouser (US)
* GearMo TTL-232R-5V equivalent from Amazon (US)
If you're in the UK/Europe, you may want to check your local Amazon site for one of these. Remember that it is important to get one that uses 0V-5V voltage levels.

A TE Connectivity type 172201 connector, which will mate to the connector on the car's diagnostic port. The shell for this connector is available from different retailers, although stock quantities may vary:

* Mouser (US)
* Tencell (UK)

Three pins for the connector shell. The pins are also manufactured by TE Connectivity, part number 170280 for strips of 50. I've been told that they're also available singly, part number 170293.

* Allied Electronics (US)
* RS (UK)

Once you have the above components, solder leads onto the pins, insert them into the connector shell, and solder the FTDI cable wires to the pin leads according to the following table. For the pin numbering on the 172201 connector, look at the face of the connector with the key (notch) on the bottom. The first pin (C549-1) will be on the top, and the numbering continues clockwise. (If you're looking at the mating connector in the car, the numbering goes counter-clockwise.)

When looking for the diagnostic connector on the vehicle, note that cars with a security/immobilizer module (such as the Lucas 5AS) will often have a second connector of the same type. On the Mini SPi, the engine ECU connector is beige, while the security module connector is green. Make sure you're connecting to the right one.

Pin assignment for USB PC interface cable

Pin number|FTDI wire color	|Pin assignment	|Wire color on mating connector in car
----------|-----------------|---------------|-------------------------------------
C549-1| Black	|Signal ground	|Pink w/ black
C549-2|	Yellow	|Rx (car ECU to PC)	|White w/ yellow
C549-3|	Orange	|Tx (PC to car ECU)	|Black w/ green


You can buy a pre-built cable, when I find the link I'll add it here

I use an HC-05 TTL Serial / Bluetooth board which you can connect to wirelessly, when I have a bit of time I'll add the instructions here as well.
It's very easy but needs the round connection parts in Colin's instructions.
