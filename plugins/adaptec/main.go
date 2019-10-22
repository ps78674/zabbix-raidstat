package main

import (
	"fmt"
	"os"
	"strings"

	functions ".."
)

// GetControllersIDs - get number of controllers in the system
func GetControllersIDs(execPath string) []string {
	inputData := functions.GetCommandOutput(execPath, "list")
	return functions.GetRegexpAllSubmatch(inputData, "Controller ([^a-zA-Z].*?):")
}

// GetLogicalDrivesIDs - get number of logical drives for controller with ID 'controllerID'
func GetLogicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetCommandOutput(execPath, "getconfig", controllerID, "ld")
	return functions.GetRegexpAllSubmatch(inputData, "Logical Device number (.*)[\\s]")
}

// GetPhysicalDrivesIDs - get number of physical drives for controller with ID 'controllerID'
func GetPhysicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetCommandOutput(execPath, "getconfig", controllerID, "pd")
	return functions.GetRegexpAllSubmatch(inputData, "Device is a Hard drive[\\s\\S]*?Reported Channel,Device\\(T:L\\)[\\s]*[:][\\s](.*?)\\(.*\\)[\\s]")
}

// GetControllerStatus - get controller status
func GetControllerStatus(execPath string, controllerID string, indent int) []byte {
	type ReturnData struct {
		Status      string `json:"status"`
		Model       string `json:"model"`
		Temperature string `json:"temperature"`
	}

	inputData := functions.GetCommandOutput(execPath, "getconfig", controllerID, "ad")
	status := functions.GetRegexpSubmatch(inputData, "Controller Status *: (.*)")
	model := functions.GetRegexpSubmatch(inputData, "Controller Model *: (.*)")
	temperature := functions.GetRegexpSubmatch(inputData, "Temperature *: (.*) C")

	if status == "Optimal" {
		status = "OK"
	}

	data := ReturnData{
		Status:      functions.TrimSpacesLeftAndRight(status),
		Model:       functions.TrimSpacesLeftAndRight(model),
		Temperature: functions.TrimSpacesLeftAndRight(temperature),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetLDStatus - get logical drive status
func GetLDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status string `json:"status"`
		Size   string `json:"size"`
	}

	inputData := functions.GetCommandOutput(execPath, "getconfig", controllerID, "ld", deviceID)
	status := functions.GetRegexpSubmatch(inputData, "Status of Logical Device *: (.*)")
	size := functions.GetRegexpSubmatch(inputData, "Size *: (.*)")

	if status == "Optimal" {
		status = "OK"
	}

	data := ReturnData{
		Status: functions.TrimSpacesLeftAndRight(status),
		Size:   functions.TrimSpacesLeftAndRight(size),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetPDStatus - get physical drive status
func GetPDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status      string `json:"status"`
		Model       string `json:"model"`
		Smart       string `json:"smart"`
		SmartWarn   string `json:"smartwarnings"`
		TotalSize   string `json:"totalsize"`
		Temperature string `json:"temperature"`
	}

	deviceData := strings.Split(deviceID, ",")
	if len(deviceData) < 2 {
		fmt.Printf("Error - wrong device id '%s'.\n", deviceID)
		os.Exit(1)

	}

	inputData := functions.GetCommandOutput(execPath, "getconfig", controllerID, "pd", deviceData[0], deviceData[1])
	status := functions.GetRegexpSubmatch(inputData, "[\\s]{2}State *: (.*)")
	model := functions.GetRegexpSubmatch(inputData, "Model *: (.*)")
	smart := functions.GetRegexpSubmatch(inputData, "S.M.A.R.T. *: (.*)")
	smartWarn := functions.GetRegexpSubmatch(inputData, "S.M.A.R.T. warnings *: (.*)")
	totalSize := functions.GetRegexpSubmatch(inputData, "Total Size *: (.*)")
	temperature := functions.GetRegexpSubmatch(inputData, "Temperature *: (.*) C")

	if status == "Online" {
		status = "OK"
	}

	if smart == "No" {
		smart = "OK"
	}

	data := ReturnData{
		Status:      functions.TrimSpacesLeftAndRight(status),
		Model:       functions.TrimSpacesLeftAndRight(model),
		Smart:       functions.TrimSpacesLeftAndRight(smart),
		SmartWarn:   functions.TrimSpacesLeftAndRight(smartWarn),
		TotalSize:   functions.TrimSpacesLeftAndRight(totalSize),
		Temperature: functions.TrimSpacesLeftAndRight(temperature),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

func main() {}
