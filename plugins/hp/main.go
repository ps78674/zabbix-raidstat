package main

import (
	"fmt"

	functions ".."
)

// GetControllersIDs - get number of controllers in the system
func GetControllersIDs(execPath string) []string {
	inputData := functions.GetInputData(execPath, "ctrl", "all", "show")
	return functions.GetRegexpAllSubmatch(inputData, "in Slot (.*?)[\\s]")
}

// GetLogicalDrivesIDs - get number of logical drives for controller with ID 'controllerID'
func GetLogicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetInputData(execPath, "ctrl", fmt.Sprintf("slot=%s", controllerID), "ld", "all", "show")
	return functions.GetRegexpAllSubmatch(inputData, "logicaldrive (.*?)[\\s]")
}

// GetPhysicalDrivesIDs - get number of physical drives for controller with ID 'controllerID'
func GetPhysicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetInputData(execPath, "ctrl", fmt.Sprintf("slot=%s", controllerID), "pd", "all", "show")
	return functions.GetRegexpAllSubmatch(inputData, "physicaldrive (.*?)[\\s]")
}

// GetControllerStatus - get controller status
func GetControllerStatus(execPath string, controllerID string, indent int) []byte {
	type ReturnData struct {
		Status        string `json:"status"`
		Model         string `json:"model"`
		BatteryStatus string `json:"batterystatus"`
	}

	inputData := functions.GetInputData(execPath, "ctrl", fmt.Sprintf("slot=%s", controllerID), "show", "status")
	status := functions.GetRegexpSubmatch(inputData, "Controller Status *: (.*)")
	model := functions.GetRegexpSubmatch(inputData, "(.*) in Slot")
	batteryStatus := functions.GetRegexpSubmatch(inputData, "Battery/Capacitor Status *: (.*)")
	data := ReturnData{Status: status, Model: model, BatteryStatus: batteryStatus}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetLDStatus - get logical drive status
func GetLDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status string `json:"status"`
		Size   string `json:"size"`
	}

	inputData := functions.GetInputData(execPath, "ctrl", fmt.Sprintf("slot=%s", controllerID), "ld", deviceID, "show", "detail")
	status := functions.GetRegexpSubmatch(inputData, "Status *: (.*)")
	size := functions.GetRegexpSubmatch(inputData, "Size *: (.*)")
	data := ReturnData{Status: status, Size: size}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetPDStatus - get physical drive status
func GetPDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status             string `json:"status"`
		Model              string `json:"model"`
		Size               string `json:"size"`
		CurrentTemperature string `json:"currenttemperature"`
		MaximumTemperature string `json:"maximumtemperature"`
	}

	inputData := functions.GetInputData(execPath, "ctrl", fmt.Sprintf("slot=%s", controllerID), "pd", deviceID, "show", "detail")
	status := functions.GetRegexpSubmatch(inputData, "[\\s]{2}Status: (.*)")
	model := functions.GetRegexpSubmatch(inputData, "Model: (.*)")
	size := functions.GetRegexpSubmatch(inputData, "[\\s]{2}Size: (.*)")
	currentTemperature := functions.GetRegexpSubmatch(inputData, "Current Temperature \\(C\\): (.*)")
	maximumTemperature := functions.GetRegexpSubmatch(inputData, "Maximum Temperature \\(C\\): (.*)")
	data := ReturnData{Status: status, Model: model, Size: size, CurrentTemperature: currentTemperature, MaximumTemperature: maximumTemperature}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

func main() {}
