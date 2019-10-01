package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

var toolVendor string
var toolBinary string
var discoveryOption string
var statusOption string
var indent int

var vendors []string
var discoveryOptions []string
var statusOptions []string
var vendorTools map[string]string

var operation string
var options []string
var argOption string
var controllerID string
var deviceID string

func init() {
	vendors = []string{"adaptec", "hp"}
	discoveryOptions = []string{"ct", "ld", "pd"}
	statusOptions = []string{"ct,<CONTROLLER_ID>", "ld,<CONTROLLER_ID>,<LD_ID>", "pd,<CONTROLLER_ID>,<PD_ID>"}
	vendorTools = map[string]string{
		"adaptec": "arcconf",
		"hp":      "ssacli",
	}

	flag.StringVar(&toolVendor, "vendor", "", fmt.Sprintf("RAID tool vendor, one of '%s'", strings.Join(vendors, " | ")))
	flag.StringVar(&toolBinary, "path", "", "RAID tool full path, like '/opt/<BINARY>'")
	flag.StringVar(&discoveryOption, "d", "", fmt.Sprintf("Discovery option, one of '%s'", strings.Join(discoveryOptions, " | ")))
	flag.StringVar(&statusOption, "s", "", fmt.Sprintf("Status option, one of '%s'", strings.Join(statusOptions, " | ")))
	flag.IntVar(&indent, "indent", 0, "Indent JSON output for <INT>")

	flag.Parse()

	if len(toolVendor) == 0 {
		fmt.Printf("RAID vendor must be set.\n")
		flag.Usage()
		os.Exit(1)
	}

	if len(toolBinary) == 0 {
		toolBinary = vendorTools[toolVendor]
	}

	for i, v := range vendors {
		if v != toolVendor {
			if i == len(vendors)-1 {
				fmt.Printf("Vendors must be one of '%s' (ex.: -vendor adaptec), got '%s'.\n", strings.Join(vendors, " | "), toolVendor)
				flag.Usage()
				os.Exit(1)
			}

			continue
		}

		break
	}

	if len(discoveryOption) == 0 && len(statusOption) == 0 {
		fmt.Println("Operation ('-d' or '-s') must be provided.")
		flag.Usage()
		os.Exit(1)
	}

	if len(discoveryOption) != 0 && len(statusOption) != 0 {
		fmt.Println("Only '-d' or '-s' must be provided.")
		flag.Usage()
		os.Exit(1)
	} else if len(discoveryOption) != 0 {
		operation = "Discovery"
		options = discoveryOptions
		argOption = discoveryOption
	} else if len(statusOption) != 0 {
		operation = "Status"
		options = statusOptions
		argOption = statusOption
	}

	for i, v := range options {
		rangeValues := strings.Split(v, ",")
		argOptionValues := strings.SplitN(argOption, ",", 3)

		if argOptionValues[0] != rangeValues[0] || len(argOptionValues) != len(rangeValues) {
			if i == len(options)-1 {
				fmt.Printf("%s option must be one of '%s', got '%s'.\n", operation, strings.Join(options, " | "), argOption)
				flag.Usage()
				os.Exit(1)
			}

			continue
		}

		if len(argOptionValues) > 1 {
			argOption = argOptionValues[0]
		}

		if len(argOptionValues) == 2 || len(argOptionValues) == 3 {
			controllerID = argOptionValues[1]
		}

		if len(argOptionValues) == 3 {
			controllerID = argOptionValues[1]
			deviceID = argOptionValues[2]
		}

		break
	}
}

func discoverControllers(p *plugin.Plugin) {
	pGetControllersIDs, err := p.Lookup("GetControllersIDs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	controllersIDs := pGetControllersIDs.(func(string) []string)(toolBinary)

	type Element struct {
		CT string `json:"{#CT_ID}"`
	}

	type Reply struct {
		Data []Element `json:"data"`
	}

	var d []Element

	for _, v := range controllersIDs {
		d = append(d, Element{CT: v})
	}

	var JSON []byte
	var jErr error

	if indent > 0 {
		JSON, jErr = json.MarshalIndent(Reply{d}, "", strings.Repeat(" ", indent))
	} else {
		JSON, jErr = json.Marshal(Reply{d})
	}

	if jErr != nil {
		fmt.Println(jErr)
		os.Exit(1)
	}

	os.Stdout.Write(append(JSON, "\n"...))
}

func discoverLogicalDrives(p *plugin.Plugin) {
	type Element struct {
		CT string `json:"{#CT_ID}"`
		LD string `json:"{#LD_ID}"`
	}

	type Reply struct {
		Data []Element `json:"data"`
	}

	pGetControllersIDs, err := p.Lookup("GetControllersIDs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pGetLogicalDrivesIDs, err := p.Lookup("GetLogicalDrivesIDs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var d []Element
	var JSON []byte
	var jErr error

	controllersIDs := pGetControllersIDs.(func(string) []string)(toolBinary)
	for _, ctID := range controllersIDs {
		logicalDrivesIDs := pGetLogicalDrivesIDs.(func(string, string) []string)(toolBinary, ctID)
		for _, ldID := range logicalDrivesIDs {
			d = append(d, Element{CT: ctID, LD: ldID})
		}

	}

	if indent > 0 {
		JSON, jErr = json.MarshalIndent(Reply{d}, "", strings.Repeat(" ", indent))
	} else {
		JSON, jErr = json.Marshal(Reply{d})
	}

	if err != nil {
		fmt.Println(jErr)
		os.Exit(1)
	}

	os.Stdout.Write(append(JSON, "\n"...))
}

func discoverPhysicalDrives(p *plugin.Plugin) {
	type Element struct {
		CT string `json:"{#CT_ID}"`
		PD string `json:"{#PD_ID}"`
	}

	type Reply struct {
		Data []Element `json:"data"`
	}

	var d []Element
	var JSON []byte
	var jErr error

	pGetControllersIDs, err := p.Lookup("GetControllersIDs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pGetPhysicalDrivesIDs, err := p.Lookup("GetPhysicalDrivesIDs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	controllersIDs := pGetControllersIDs.(func(string) []string)(toolBinary)
	for _, ctID := range controllersIDs {
		logicalDrivesIDs := pGetPhysicalDrivesIDs.(func(string, string) []string)(toolBinary, ctID)
		for _, pdID := range logicalDrivesIDs {
			d = append(d, Element{CT: ctID, PD: pdID})
		}

	}

	if indent > 0 {
		JSON, jErr = json.MarshalIndent(Reply{d}, "", strings.Repeat(" ", indent))
	} else {
		JSON, jErr = json.Marshal(Reply{d})
	}

	if err != nil {
		fmt.Println(jErr)
		os.Exit(1)
	}

	os.Stdout.Write(append(JSON, "\n"...))
}

func getControllerStatus(p *plugin.Plugin, controllerID string) {
	pGetControllerStatus, err := p.Lookup("GetControllerStatus")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Stdout.Write(pGetControllerStatus.(func(string, string, int) []byte)(toolBinary, controllerID, indent))
}

func getLDStatus(p *plugin.Plugin, controllerID string, deviceID string) {
	pGetLDStatus, err := p.Lookup("GetLDStatus")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Stdout.Write(pGetLDStatus.(func(string, string, string, int) []byte)(toolBinary, controllerID, deviceID, indent))
}

func getPDStatus(p *plugin.Plugin, controllerID string, deviceID string) {
	pGetPDStatus, err := p.Lookup("GetPDStatus")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Stdout.Write(pGetPDStatus.(func(string, string, string, int) []byte)(toolBinary, controllerID, deviceID, indent))
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p, err := plugin.Open(filepath.Dir(ex) + "/" + toolVendor + ".so")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch argOption {
	case "ct":
		switch operation {
		case "Discovery":
			discoverControllers(p)
		case "Status":
			getControllerStatus(p, controllerID)
		}
	case "ld":
		switch operation {
		case "Discovery":
			discoverLogicalDrives(p)
		case "Status":
			getLDStatus(p, controllerID, deviceID)
		}
	case "pd":
		switch operation {
		case "Discovery":
			discoverPhysicalDrives(p)
		case "Status":
			getPDStatus(p, controllerID, deviceID)
		}
	}
}
