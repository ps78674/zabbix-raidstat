package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/ps78674/docopt.go"
)

const configFile = "config.json"

var (
	indent       int
	toolVendor   string
	toolBinary   string
	operation    string
	argOption    string
	controllerID string
	deviceID     string
)

func init() {
	type Config struct {
		Vendors interface{} `json:"vendors"`
	}

	var (
		configJSON      Config
		vendors         []string
		discoveryOption string
		statusOption    string
		options         []string
	)

	ex, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configFile, err := os.Open(fmt.Sprintf("%s/%s", filepath.Dir(ex), configFile))
	if err != nil {
		fmt.Printf("Error opening config file: %s\n", err)
		os.Exit(1)
	}

	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(configData, &configJSON); err != nil {
		fmt.Printf("Error unmarshalling JSON data: %s\n", err)
		os.Exit(1)
	}

	if configJSON.Vendors == nil {
		fmt.Println("Failed to get vendors from config file.")
		os.Exit(1)
	}

	for v := range configJSON.Vendors.(map[string]interface{}) {
		vendors = append(vendors, v)
	}

	discoveryOptions := []string{"ct", "ld", "pd"}
	statusOptions := []string{"ct,<CONTROLLER_ID>", "ld,<CONTROLLER_ID>,<LD_ID>", "pd,<CONTROLLER_ID>,<PD_ID>"}

	var programName = filepath.Base(os.Args[0])
	var usage = fmt.Sprintf(`%[1]s: parse raid vendor tool output and format it as json

Usage:
  %[1]s (-v <VENDOR>) (-d <OPTION> | -s <OPTION>) [-i <INT>]

Options:
  -v, --vendor <VENDOR>    raid tool vendor, one of: %[2]s
  -d, --discover <OPTION>  discovery option, one of: %[3]s
  -s, --status <OPTION>    status option, one of: %[4]s
  -i, --indent <INT>       indent json output level [default: 0]

  -h, --help               show this screen
	`, programName, strings.Join(vendors, " | "), strings.Join(discoveryOptions, " | "), strings.Join(statusOptions, " | "))

	cmdOpts, err := docopt.ParseDoc(usage)
	if err != nil {
		fmt.Printf("error parsing options: %s\n", err)
		os.Exit(1)
	}

	toolVendor, _ = cmdOpts.String("--vendor")
	discoveryOption, _ = cmdOpts.String("--discover")
	statusOption, _ = cmdOpts.String("--status")
	indent, _ = cmdOpts.Int("--indent")

	for i, v := range vendors {
		if v != toolVendor {
			if i == len(vendors)-1 {
				fmt.Printf("Vendors must be one of '%s' (ex.: -v adaptec), got '%s'.\n", strings.Join(vendors, " | "), toolVendor)
				docopt.PrintHelpOnly(nil, usage)
				os.Exit(1)
			}
			continue
		}
		break
	}

	toolBinary = configJSON.Vendors.(map[string]interface{})[toolVendor].(string)

	if len(discoveryOption) != 0 {
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
				docopt.PrintHelpOnly(nil, usage)
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
	type (
		Element struct {
			CT string `json:"{#CT_ID}"`
		}
		Reply struct {
			Data []Element `json:"data"`
		}
	)

	var (
		d    []Element
		JSON []byte
		jErr error
	)

	pGetControllersIDs, err := p.Lookup("GetControllersIDs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	controllersIDs := pGetControllersIDs.(func(string) []string)(toolBinary)

	for _, v := range controllersIDs {
		d = append(d, Element{CT: v})
	}

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
	type (
		Element struct {
			CT string `json:"{#CT_ID}"`
			LD string `json:"{#LD_ID}"`
		}
		Reply struct {
			Data []Element `json:"data"`
		}
	)

	var (
		d    []Element
		JSON []byte
		jErr error
	)

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
	type (
		Element struct {
			CT string `json:"{#CT_ID}"`
			PD string `json:"{#PD_ID}"`
		}
		Reply struct {
			Data []Element `json:"data"`
		}
	)

	var (
		d    []Element
		JSON []byte
		jErr error
	)

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
		fmt.Printf("Error opening plugin '%s.so': %s\n", toolVendor, err)
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
