package main

import (
	"fmt"
	"strings"

	"github.com/ps78674/zabbix-raidstat/plugins/internal/functions"
)

// GetControllersIDs - get number of controllers in the system
func GetControllersIDs(execPath string) []string {
	inputData := functions.GetCommandOutput(execPath, "info", "-o", "hba")
	return functions.GetRegexpAllSubmatch(inputData, "Adapter ID:[\\s]+(.*)")
}

// GetLogicalDrivesIDs - get number of logical drives for controller with ID 'controllerID'
func GetLogicalDrivesIDs(execPath string, controllerID string) []string {
	functions.GetCommandOutput(execPath, "adapter", "-i", controllerID)
	inputData := functions.GetCommandOutput(execPath, "info", "-o", "ld")
	return functions.GetRegexpAllSubmatch(inputData, "id:[\\s]+(.*)")
}

// GetPhysicalDrivesIDs - get number of physical drives for controller with ID 'controllerID'
func GetPhysicalDrivesIDs(execPath string, controllerID string) []string {
	functions.GetCommandOutput(execPath, "adapter", "-i", controllerID)
	inputData := functions.GetCommandOutput(execPath, "info", "-o", "pd")
	return functions.GetRegexpAllSubmatch(inputData, "PD ID:[\\s]+(.*)")
}

// GetControllerStatus - get controller status
func GetControllerStatus(execPath string, controllerID string, indent int) []byte {
	type ReturnData struct {
		Status      string `json:"status"`
		ModelNumber string `json:"modelnumber"`
		PartNumber  string `json:"partnumber"`
	}

	inputData := functions.GetCommandOutput(execPath, "info", "-o", "hba", "-i", controllerID)

	healthStatuses := []string{}
	for _, v := range []string{
		"Image health",
		"Autoload image health",
		"Boot loader image health",
		"Firmware image health",
		"Boot ROM image health",
		"HBA info image health",
	} {
		if s := functions.GetRegexpSubmatch(inputData, fmt.Sprintf("%s:[\\s]+(.*)", v)); s != "Healthy" {
			healthStatuses = append(healthStatuses, fmt.Sprintf("%s is %s", v, s))
		}
	}

	var status string
	if len(healthStatuses) == 0 {
		status = "OK"
	} else {
		status = strings.Join(healthStatuses, ", ")
	}

	modelnumber := functions.GetRegexpSubmatch(inputData, "ModelNumber:[\\s]+(.*)")
	partnumber := functions.GetRegexpSubmatch(inputData, "PartNumber:[\\s]+(.*)")

	data := ReturnData{
		Status:      functions.TrimSpacesLeftAndRight(status),
		ModelNumber: functions.TrimSpacesLeftAndRight(modelnumber),
		PartNumber:  functions.TrimSpacesLeftAndRight(partnumber),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetLDStatus - get logical drive status
func GetLDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status   string `json:"status"`
		Name     string `json:"name"`
		Size     string `json:"size"`
		RaidMode string `json:"raidmode"`
	}

	functions.GetCommandOutput(execPath, "adapter", "-i", controllerID) // set adapter for next commands (mvcli-specific)
	inputData := functions.GetCommandOutput(execPath, "info", "-o", "ld", "-i", deviceID)
	status := functions.GetRegexpSubmatch(inputData, "VD status:[\\s]+(.*)")
	name := functions.GetRegexpSubmatch(inputData, "name:[\\s]+(.*)")
	size := functions.GetRegexpSubmatch(inputData, "size:[\\s]+(.*)")
	raidmode := functions.GetRegexpSubmatch(inputData, "RAID mode:[\\s]+(.*)")

	if status == "optimal" {
		status = "OK"
	}

	data := ReturnData{
		Status:   functions.TrimSpacesLeftAndRight(status),
		Name:     functions.TrimSpacesLeftAndRight(name),
		Size:     functions.TrimSpacesLeftAndRight(size),
		RaidMode: functions.TrimSpacesLeftAndRight(raidmode),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

// GetPDStatus - get physical drive status
func GetPDStatus(execPath string, controllerID string, deviceID string, indent int) []byte {
	type ReturnData struct {
		Status          string `json:"status"`
		Model           string `json:"model"`
		FirmwareVersion string `json:"firmwareversion"`
		Size            string `json:"size"`
		CurrentSpeed    string `json:"currentspeed"`
	}

	inputData := functions.GetCommandOutput(execPath, "info", "-o", "pd", "-i", deviceID)
	status := functions.GetRegexpSubmatch(inputData, "PD status:[\\s]+(.*)")
	model := functions.GetRegexpSubmatch(inputData, "model:[\\s]+(.*)")
	firmwareversion := functions.GetRegexpSubmatch(inputData, "Firmware version:[\\s]+(.*)")
	size := functions.GetRegexpSubmatch(inputData, "Size:[\\s]+(.*)")
	currentspeed := functions.GetRegexpSubmatch(inputData, "Current speed:[\\s]+(.*)")

	if status == "online" {
		status = "OK"
	}

	data := ReturnData{
		Status:          functions.TrimSpacesLeftAndRight(status),
		Model:           functions.TrimSpacesLeftAndRight(model),
		FirmwareVersion: functions.TrimSpacesLeftAndRight(firmwareversion),
		Size:            functions.TrimSpacesLeftAndRight(size),
		CurrentSpeed:    functions.TrimSpacesLeftAndRight(currentspeed),
	}

	return append(functions.MarshallJSON(data, indent), "\n"...)
}

func main() {}
