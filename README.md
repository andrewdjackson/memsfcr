# MemsFCR 
The Mems FCR supports Mac and Windows. The application connects using the serial interface and supports the serial FTDI cables you can purchase for the Android MemsDiag application. Instructions on how to create this cable or alternatively build a bluetooth wireless interface are also included.

All MEMS data points are read and can be viewed directly. If any faults have been detected these will be displayed and can be cleared.

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

### Dataframes
![Dataframes](./resources/screenshots/dataframes.png)
A live view of the calculated parameters received from the ECU. Since not all the values in the MEMS dataframe have been mapped this may be updated as more information comes to light.

### Information
Useful information on the operating characteristics of the MEMS ECU