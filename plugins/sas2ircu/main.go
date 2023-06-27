package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ps78674/zabbix-raidstat/plugins/internal/functions"
)

// GetControllersIDs - get number of controllers in the system
func GetControllersIDs(execPath string) []string {
	inputData := functions.GetCommandOutput(execPath, "list")
	return functions.GetRegexpAllSubmatch(inputData, "\\s+(\\d+)\\s+.*")
}

// GetLogicalDrivesIDs - get number of logical drives for controller with ID 'controllerID'
func GetLogicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetCommandOutput(execPath, controllerID, "display")
	return functions.GetRegexpAllSubmatch(inputData, "IR volume (\\d+)")
}

// GetPhysicalDrivesIDs - get number of physical drives for controller with ID 'controllerID'
func GetPhysicalDrivesIDs(execPath string, controllerID string) []string {
	inputData := functions.GetCommandOutput(execPath, controllerID, "display")
	sliceArr := functions.GetArraySliceByte(inputData, "Device is a Hard disk", "Drive Type")
	data := []string{}

	if len(sliceArr) > 0 {
		for _, v := range sliceArr {
			enclosure := functions.GetRegexpSubmatch([]byte(v), "Enclosure # *: (.*)")
			slot := functions.GetRegexpSubmatch([]byte(v), "Slot # *: (.*)")

			if len(enclosure) > 0 && len(slot) > 0 {
				data = append(data, fmt.Sprintf("%s:%s", enclosure, slot))
			}
		}
	}

	return data
}

// GetControllerStatus - get controller status
func GetControllerStatus(execPath string, controllerID string, indent int) []byte {
	type ReturnData struct {
		Status string `json:"status"`
		Model  string `json:"model"`
		// Temperature string `json:"temperature"`
	}

	inputData := functions.GetCommandOutput(execPath, controllerID, "display")
	model := functions.GetRegexpSubmatch(inputData, "Controller type *: (.*)")

	result := regexp.MustCompile("Status of volume\\s+: .*\\((.*)\\)").FindAllStringSubmatch(string(inputData), -1)

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Regexp is '%s'\n", "Enclosure #\\s+: (\\d+)\\n \\s+Slot #\\s+: (\\d+)")
		fmt.Printf("Result is '%s'\n", result)
	}

	healthStatuses := []string{}

	if len(result) > 0 {
		for _, v := range result {
			if v[1] != "OKY" {
				healthStatuses = append(healthStatuses, fmt.Sprintf("%s", v))
			}
		}
	}

	var status string
	if len(healthStatuses) == 0 {
		status = "OK"
	} else {
		status = strings.Join(healthStatuses, ", ")
	}

	data := ReturnData{
		Status: functions.TrimSpacesLeftAndRight(status),
		Model:  functions.TrimSpacesLeftAndRight(model),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetLDStatus - get logical drive status
func GetLDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status string `json:"status"`
		Size   string `json:"size"`
	}

	inputData := functions.GetCommandOutput(execPath, controllerID, "display")
	sliceData := functions.GetSliceByte(inputData, "IR volume "+deviceID, "Physical")

	status := functions.GetRegexpSubmatch(sliceData, "Status of volume *: (.*)")
	size := functions.GetRegexpSubmatch(sliceData, "Size \\(in MB\\) *: (.*)")

	if status == "Okay (OKY)" {
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
		Status    string `json:"status"`
		Model     string `json:"model"`
		TotalSize string `json:"totalsize"`
	}

	deviceData := strings.Split(deviceID, ":")
	if len(deviceData) < 2 {
		fmt.Printf("Error - wrong device id '%s'.\n", deviceID)
		os.Exit(1)
	}

	inputData := functions.GetCommandOutput(execPath, controllerID, "display")
	sliceArr := functions.GetArraySliceByte(inputData, "Device is a Hard disk", "Drive Type")

	if len(sliceArr) > 0 {
		for _, v := range sliceArr {
			enclosure := functions.GetRegexpSubmatch([]byte(v), "Enclosure # *: (.*)")
			slot := functions.GetRegexpSubmatch([]byte(v), "Slot # *: (.*)")

			if enclosure == deviceData[0] && slot == deviceData[1] {
				status := functions.GetRegexpSubmatch([]byte(v), "[\\s]{2}State *: (.*)")
				model := functions.GetRegexpSubmatch([]byte(v), "Model Number *: (.*)")
				totalSize := functions.GetRegexpSubmatch([]byte(v), "Size \\(in MB\\)/\\(in sectors\\) *: (\\d+)/\\d+")

				if status == "Optimal (OPT)" {
					status = "OK"
				}

				data := ReturnData{
					Status:    functions.TrimSpacesLeftAndRight(status),
					Model:     functions.TrimSpacesLeftAndRight(model),
					TotalSize: functions.TrimSpacesLeftAndRight(totalSize),
				}

				return append(functions.MarshallJSON(data, indent), "\n"...)
			}
		}
	}

	return []byte("")
}

func main() {}
