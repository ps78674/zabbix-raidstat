package main

import (
	"fmt"
	"strings"
	"regexp"
	"os"

	"github.com/ps78674/zabbix-raidstat/plugins/internal/functions"
)

// GetControllersIDs - get number of controllers in the system
func GetControllersIDs(execPath string) []string {
	inputData := functions.GetCommandOutput(execPath, "-AdpGetPciInfo", "-aALL")
	return functions.GetRegexpAllSubmatch(inputData, "for Controller (\\d*)")
}

// GetLogicalDrivesIDs - get number of logical drives for controller with ID 'controllerID'
func GetLogicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetCommandOutput(execPath, "-LdInfo", fmt.Sprintf("-a%s", controllerID), "-Lall", "-NoLog")
	return functions.GetRegexpAllSubmatch(inputData, "Virtual Drive: (.*?)[\\s]")
}

// GetPhysicalDrivesIDs - get number of physical drives for controller with ID 'controllerID'
func GetPhysicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetCommandOutput(execPath, "-PDList", fmt.Sprintf("-a%s", controllerID), "-NoLog")

	result := regexp.MustCompile("Enclosure Device ID: (\\d+)\\nSlot Number: (\\d+)").FindAllStringSubmatch(string(inputData), -1)

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Regexp is '%s'\n", "Enclosure Device ID: (\\d+)\\nSlot Number: (\\d+)")
		fmt.Printf("Result is '%s'\n", result)
	}

	data := []string{}

	if len(result) > 0 {
		for _, v := range result {
			data = append(data, fmt.Sprintf("%s:%s", v[1], v[2]))
		}
	}

	return data

	return functions.GetRegexpAllSubmatch(inputData, "Enclosure Device ID: (\\d+)\\nSlot Number: (\\d+)")
}

// GetControllerStatus - get controller status
func GetControllerStatus(execPath string, controllerID string, indent int) []byte {
	type ReturnData struct {
		Status        string `json:"status"`
		Model         string `json:"model"`
		BatteryStatus string `json:"batterystatus"`
	}

	inputData := functions.GetCommandOutput(execPath, "-AdpAllInfo", fmt.Sprintf("-a%s", controllerID), "-NoLog")
	model := functions.GetRegexpSubmatch(inputData, "roduct Name[\\s]+: (.*)")

	healthStatuses := []string{}
	for _, v := range []string{
		"Degraded",
		"Offline",
		"Critical Disks",
		"Failed Disks",
	} {
		s := functions.GetRegexpSubmatch(inputData, fmt.Sprintf("%s[\\s]+: (.*)", v));

		if s != "0" {
			healthStatuses = append(healthStatuses, fmt.Sprintf("%s is %s", v, s))
		}
	}

	var status string
	if len(healthStatuses) == 0 {
		status = "OK"
	} else {
		status = strings.Join(healthStatuses, ", ")
	}

	inputData = functions.GetCommandOutput(execPath, "-AdpBbuCmd", fmt.Sprintf("-a%s", controllerID), "-GetBbuStatus", "-NoLog")
	batteryStatus := functions.GetRegexpSubmatch(inputData, "Battery State: (.*)")

	data := ReturnData{
		Status:        functions.TrimSpacesLeftAndRight(status),
		Model:         functions.TrimSpacesLeftAndRight(model),
		BatteryStatus: functions.TrimSpacesLeftAndRight(batteryStatus),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetLDStatus - get logical drive status
func GetLDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status string `json:"status"`
		Size   string `json:"size"`
	}

	inputData := functions.GetCommandOutput(execPath, "-LdInfo", fmt.Sprintf("-L%s", deviceID), fmt.Sprintf("-a%s", controllerID), "-NoLog")
	status := functions.GetRegexpSubmatch(inputData, "State *: (.*)")
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
		Status             string `json:"status"`
		Model              string `json:"model"`
		Size               string `json:"size"`
		CurrentTemperature string `json:"currenttemperature"`
	}

	inputData := functions.GetCommandOutput(execPath, "-pdInfo", fmt.Sprintf("-PhysDrv[%s]", deviceID), fmt.Sprintf("-a%s", controllerID), "-NoLog")
	status := functions.GetRegexpSubmatch(inputData, "Firmware state: (.*)")
	model := functions.GetRegexpSubmatch(inputData, "Inquiry Data: (.*)")
	size := functions.GetRegexpSubmatch(inputData, "Raw Size: (.*) \\[")
	currentTemperature := functions.GetRegexpSubmatch(inputData, "Drive Temperature :(\\d+)C")

	data := ReturnData{
		Status:             functions.TrimSpacesLeftAndRight(status),
		Model:              functions.TrimSpacesLeftAndRight(model),
		Size:               functions.TrimSpacesLeftAndRight(size),
		CurrentTemperature: functions.TrimSpacesLeftAndRight(currentTemperature),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

func main() {}
