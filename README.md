# Prometheus Metric Exporter for Solax Solar Systems

# Description
Read Solax Inverter real-time API and exposes as Prometheus metric endpoint.

Real time: 
- Power
- Current
- Voltage
- Grid power data
- Battery data
- Health data
- Daily/Total energy summaries

## Building
	# Native
	go build solax_exporter.go
	
	# Raspberry Pi
	GOOS=linux GOARCH=arm GOARM=5 go build solax_exporter.go

## How it works

The application assumes a network connection to the Solax wifi network and connects to it at 11.11.11.1 on URI /api/realTimeData.htm returning a json object. 

This json object is then parsed according to the exports below (thanks GitHobi) and update the metrics accordingly.

## What this exports

- PV - PV1 Current
- PV - PV2 Current
- PV - PV1 Voltage
- PV - PV2 Voltage
- PV - PV1 Input Power
- PV - PV2 Input Power
- Grid - Output Current
- Grid - Network Voltage
- Grid - Power
- Grid - Feed in Power
- Grid - Frequency
- Grid - Exported
- Grid - Imported
- Battery Voltage
- Dis/Charge Current
- Battery Power
- Battery Temperature
- Remaining Capacity
- Inverter Yield - Today
- Inverter Yield - This Month
- Battery Yield - Total

## Thanks 
Thanks to GitHobi (https://github.com/GitHobi/solax) for providing information on the data array - saved me a lot of time!
