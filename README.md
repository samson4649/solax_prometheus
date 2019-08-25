# Prometheus Metric Exporter for Solax Solar Systems

## Introduction
I wanted to attempt to build a go application and had a solar system that was exposing metrics via an API on a closed network. This application is designed to run on the closed wifi network (from a raspberry pi in this case) and expose prometheus metrics on the internal LAN on port 4444

## Building
	# Native
	go build solax_exporter.go
	
	# Raspberry Pi
	GOOS=linux GOARCH=arm GOARM=5 go build solax_exporter.go

## How it works

The application assumes a network connection to the Solax wifi network and connects to it at 11.11.11.1 on URI /api/realTimeData.htm returning a json object. 

This json object is then parsed according to the exports below (thanks GitHobi) and update the metrics accordingly.

This loop is in a go routine and has a sleep timer of 10s between runs.  

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
