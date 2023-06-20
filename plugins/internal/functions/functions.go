package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// TrimSpacesLeftAndRight - trim leading and trailing spaces
func TrimSpacesLeftAndRight(input string) string {
	return strings.TrimLeft(strings.TrimRight(input, " "), " ")
}

// GetCommandOutput - get input data from RAID tool
func GetCommandOutput(execPath string, args ...string) []byte {
	timeout := 10
	execContext, contextCancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer contextCancel()

	cmd := exec.CommandContext(execContext, execPath, args...)
	data, err := cmd.Output()

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Command '%s %s' output is:\n'''\n%s\n'''\n", execPath, strings.Join(args, " "), string(data))
	}

	if err != nil {
		if execContext.Err() == context.DeadlineExceeded {
			fmt.Printf("Command '%s' timed out.\n", cmd)
		} else {
			fmt.Printf("Error executing command '%s %s': %s\n", execPath, strings.Join(args, " "), err)
		}

		os.Exit(1)
	}

	return data
}

// GetRegexpSubmatch - returns string from 1st capture group
func GetRegexpSubmatch(buf []byte, re string) (data string) {
	result := regexp.MustCompile(re).FindStringSubmatch(string(buf))

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Regexp is '%s'\n", re)
		if len(result) > 0 {
			fmt.Printf("Result is '%s'\n", result[0])
		}
		fmt.Printf("Input data is:\n'''\n%s\n'''\n", string(buf))
	}

	if len(result) > 0 {
		data = result[1]
	}

	return
}

// GetRegexpAllSubmatch - returns strings from all capture groups
func GetRegexpAllSubmatch(buf []byte, re string) (data []string) {
	result := regexp.MustCompile(re).FindAllStringSubmatch(string(buf), -1)

	if os.Getenv("RAIDSTAT_DEBUG") == "y" {
		fmt.Printf("Regexp is '%s'\n", re)
		fmt.Printf("Result is '%s'\n", result)
		fmt.Printf("Input data is:\n'''\n%s\n'''\n", string(buf))
	}

	if len(result) > 0 {
		for _, v := range result {
			data = append(data, v[1])
		}
	}

	return
}

// MarshallJSON - returns json object
func MarshallJSON(data interface{}, indent int) []byte {
	var (
		JSON []byte
		jErr error
	)

	if indent > 0 {
		JSON, jErr = json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	} else {
		JSON, jErr = json.Marshal(data)
	}

	if jErr != nil {
		fmt.Printf("Error marshalling JSON: %s\n", jErr.Error())
		os.Exit(1)
	}

	return JSON
}

func GetSliceByte(buf []byte, start string, end string) []byte {
	lines := strings.Split(string(buf), "\n")
	capture := false
	var sliceData []byte

	if len(lines) > 0 {
		for _, v := range lines {
			if strings.Contains(v, start) {
				capture = true
			}

			if capture {
				sliceData = append(sliceData, v+"\n"...)

				if strings.Contains(v, end) {
					break
				}
			}
		}
	}

	return sliceData
}

func GetArraySliceByte(buf []byte, start string, end string) (data []string) {
	lines := strings.Split(string(buf), "\n")
	capture := false
	var sliceData []byte

	if len(lines) > 0 {
		for _, v := range lines {
			if strings.Contains(v, start) {
				if capture {
					data = append(data, string(sliceData))
					sliceData = nil
				} else {
					capture = true
				}
			}

			if strings.Contains(v, end) {
				if capture {
					data = append(data, string(sliceData))
					sliceData = nil
				}

				capture = false
			}

			if capture {
				sliceData = append(sliceData, v+"\n"...)
			}
		}
	}

	return
}
