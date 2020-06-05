package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func (d *discovery) discoveryRaid(line string) error {
	split := strings.SplitN(line, ":", 2)
	match := strings.ToLower(strings.TrimSpace(split[0]))
	if strings.HasPrefix(match, "mcontroller") {
		d.DeviceID = strings.TrimSpace(split[1])
	}
	if strings.HasPrefix(match, "product name") {
		d.DeviceAlias = strings.TrimSpace(split[1])
	}
	d.DeviceType = "Raid"
	return nil
}

func getRaidStats(device string) {
	bin, err := getBin("megacli")
	if err != nil {
		if debug {
			log.Printf("getStatsPhysicalDrive: %v", err)
		}
		os.Exit(1)
	}
	args := []string{"-AdpAllInfo", "-aAll"}
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		if debug {
			log.Printf("getStatsPhysicalDrive: %v", err)
		}
		os.Exit(1)
	}

	replaced := strings.ReplaceAll(string(out), "Adapter #", "Adapter #MController: ")
	parts := strings.Split(replaced, "Adapter #")

	var prevLine string
	var upLevel string
	var raidID string
	result := make(map[string]map[string]map[string]string)
	for _, raid := range parts {
		for _, line := range strings.Split(raid, "\n") {
			if strings.HasPrefix(line, "MController:") {
				data := processLine(line, valString)
				raidID = data.(string)
				result[raidID] = make(map[string]map[string]string)
			}
			if strings.Contains(line, "====") && len(prevLine) > 0 {
				upLevel = strings.TrimSpace(prevLine)
				result[raidID][upLevel] = map[string]string{}
			}
			split := strings.SplitN(line, ":", 2)
			if len(split) == 2 && len(upLevel) > 0 {
				value := strings.TrimSpace(split[0])
				result[raidID][upLevel][value] = strings.TrimSpace(split[1])
			}
			prevLine = line
		}
	}

	if _, ok := result[device]; ok {
		//r, _ := json.MarshalIndent(result[device], "", " ")
		r, _ := json.Marshal(result[raidID])
		fmt.Print(string(r))
	} else {
		fmt.Printf("Raid not exist %v", device)
		os.Exit(1)
	}
}

func discoveryRaidControllers() {
	//megacli -AdpAllInfo -aAll
	bin, err := getBin("megacli")
	if err != nil {
		if debug {
			log.Printf("discoveryRaidControllers: %v", err)
		}
		os.Exit(1)
	}
	args := []string{"-AdpAllInfo", "-aAll"}
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		if debug {
			log.Printf("discoveryRaidControllers: %v", err)
		}
		os.Exit(1)
	}

	result := []discovery{}

	replaced := strings.ReplaceAll(string(out), "Adapter #", "Adapter # MController: ")
	parts := strings.Split(replaced, "Adapter #")
	for _, dev := range parts {
		device := discovery{}
		for _, line := range strings.Split(dev, "\n") {
			device.discoveryRaid(line)
		}
		if len(device.DeviceID) > 0 {
			result = append(result, device)
		}
	}
	r, _ := json.Marshal(result)
	fmt.Println(string(r))

}
