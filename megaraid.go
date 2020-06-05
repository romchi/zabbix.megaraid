package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	valString int = iota
	valInt
)

type discovery struct {
	DeviceID    string `json:"{#DEVICE_ID}"`
	DeviceType  string `json:"{#DEVICE_TYPE}"`
	DeviceAlias string `json:"{#DEVICE_ALIAS}"`
}

var debug = false

func main() {
	discoveryCommand := flag.NewFlagSet("discover", flag.ExitOnError)
	statsCommand := flag.NewFlagSet("stats", flag.ExitOnError)

	discoveryDeviceType := discoveryCommand.String("type", "", "device type {raid, ld, pd} (Required)")

	statsDeviceType := statsCommand.String("type", "", "device type {raid, ld, pd} (Required)")
	statsDeviceName := statsCommand.String("name", "", `Device "name" to get stats (Required)`)

	if len(os.Args) < 2 {
		fmt.Println("[discovery, stats, check] - required one command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "discovery":
		discoveryCommand.Parse(os.Args[2:])
	case "stats":
		statsCommand.Parse(os.Args[2:])
	case "check":
		checkControllers()
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	if discoveryCommand.Parsed() {
		metricChoices := map[string]bool{"raid": true, "ld": true, "pd": true}
		if _, validChoice := metricChoices[*discoveryDeviceType]; !validChoice {
			discoveryCommand.PrintDefaults()
			os.Exit(1)
		}
		switch *discoveryDeviceType {
		case "raid":
			discoveryRaidControllers()
		case "ld":
			discoveryLogicalDrive()
		case "pd":
			discoveryPhysicalDrive()
		default:
			discoveryCommand.PrintDefaults()
			os.Exit(1)
		}
	}

	if statsCommand.Parsed() {
		metricChoices := map[string]bool{"raid": true, "ld": true, "pd": true}
		if _, validChoice := metricChoices[*statsDeviceType]; !validChoice {
			statsCommand.PrintDefaults()
			os.Exit(0)
		}
		if len(*statsDeviceName) < 1 {
			statsCommand.PrintDefaults()
			os.Exit(0)
		}
		switch *statsDeviceType {
		case "ad":
			getRaidStats(*statsDeviceName)
		case "ld":
			getStatsLogicalDrive(*statsDeviceName)
		case "pd":
			getStatsPhysicalDrive(*statsDeviceName)
		default:
			statsCommand.PrintDefaults()
			os.Exit(0)
		}
	}
}

func checkControllers() {
	raid, err := isMegaRAID()
	if err != nil {
		if debug {
			log.Printf("Unable check raid controllers: %v", err)
		}
		fmt.Print(0)
		os.Exit(0)
	}
	megacli, err := isMegacli()
	if err != nil {
		if debug {
			log.Printf("Unable check megacli controllers: %v", err)
		}
		fmt.Print(0)
		os.Exit(0)
	}
	if raid && megacli {
		if debug {
			fmt.Printf("raid installed - %t\nmegacli installed - %t\n", raid, megacli)
		}
		fmt.Print(1)
		os.Exit(0)
	}
	if debug {
		fmt.Printf("raid installed - %t\nmegacli installed - %t\n", raid, megacli)
	}
	fmt.Print(0)
}

func isMegaRAID() (bool, error) {
	bin, err := getBin("lspci")
	if err != nil {
		return false, err
	}

	cmd := exec.Command(bin)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	err = cmd.Run()
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(out)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "LSI") &&
			strings.Contains(scanner.Text(), "MegaRAID") {
			return true, nil
		}
	}
	return false, nil
}

func isMegacli() (bool, error) {
	bin, err := getBin("megacli")
	if err != nil {
		return false, err
	}
	if len(bin) > 7 {
		return true, nil
	}
	return false, nil
}

func isLSPCI() (bool, error) {
	bin, err := getBin("lspci")
	if err != nil {
		return false, err
	}
	if len(bin) > 5 {
		return true, nil
	}
	return false, nil
}

func getBin(binFile string) (string, error) {
	locations := []string{"/bin", "/sbin", "/usr/bin", "/usr/sbin", "/usr/local/bin", "/usr/local/sbin"}

	for _, path := range locations {
		lookup := path + "/" + binFile
		fileInfo, err := os.Stat(path + "/" + binFile)
		if err != nil {
			continue
		}
		if !fileInfo.IsDir() {
			return lookup, nil
		}
	}
	return "", fmt.Errorf("Not found: '%v'", binFile)
}

func processLine(line string, valueType int) interface{} {
	split := strings.SplitN(line, ":", 2)
	if len(split) != 2 {
		if debug {
			log.Printf("processLine line not have 2 values: %v", line)
		}
		if valueType == valString {
			return ""
		}
		if valueType == valInt {
			return 0
		}
	}

	data := strings.TrimSpace(split[1])

	if valueType == valString {
		return string(data)
	}
	if valueType == valInt {
		value, err := strconv.ParseInt(data, 10, 0)
		if err != nil {
			if debug {
				log.Printf("processLine valInt: %v", err)
			}
			return 0
		}
		return int(value)
	}
	if debug {
		log.Printf("processLine unable process line: %v", line)
	}
	return errors.New("processLine: valueType not supported")
}
