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

	result := regexp.MustCompile("Enclosure #\\s+: (\\d+)\\n \\s+Slot #\\s+: (\\d+)").FindAllStringSubmatch(string(inputData), -1)

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Regexp is '%s'\n", "Enclosure #\\s+: (\\d+)\\n \\s+Slot #\\s+: (\\d+)")
		fmt.Printf("Result is '%s'\n", result)
	}

	data := []string{}

	if len(result) > 0 {
		for _, v := range result {
			data = append(data, fmt.Sprintf("%s:%s", v[1], v[2]))
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

	lines := strings.Split(string(inputData), "\n")
	capture := false
	status := ""
	model := ""
	totalSize := ""
	enclosure := ""
	slot := ""
	var sliceData []byte

	if len(lines) > 0 {
		for _, v := range lines {
			if strings.Contains(v, "Device is a Hard disk") {
				capture = true
			}

			if capture {
				sliceData = append(sliceData, v+"\n"...)

				if strings.Contains(v, "Enclosure") {
					enclosure = functions.GetRegexpSubmatch([]byte(v), "Enclosure # *: (.*)")
				}

				if strings.Contains(v, "Slot") {
					slot = functions.GetRegexpSubmatch([]byte(v), "Slot # *: (.*)")
				}

				if strings.Contains(v, "Drive Type") {
					if enclosure == deviceData[0] && slot == deviceData[1] {
						status = functions.GetRegexpSubmatch(sliceData, "[\\s]{2}State *: (.*)")
						model = functions.GetRegexpSubmatch(sliceData, "Model Number *: (.*)")
						totalSize = functions.GetRegexpSubmatch(sliceData, "Size \\(in MB\\)/\\(in sectors\\) *: (\\d+)/\\d+")

						break
					}

					sliceData = nil
				}
			}
		}
	}

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

func main() {}
