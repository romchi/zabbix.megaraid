package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type ldStats struct {
	DeviceID              string
	Name                  string
	RaidLevel             string
	Size                  string
	SectorSize            int
	IsVDemulated          string
	ParitySize            string
	State                 string
	StripSize             string
	NumberOfDrivesPerSpan int
	SpanDepth             int
	DiskCachePolicy       string
	CurrentAccessPolicy   string
	CurrentCachePolicy    string
	EncryptionType        string
	BadBlocksExist        string
	IsVDCached            string
}

func (d *discovery) discoveryLD(line string) error {
	if strings.HasPrefix(line, "LogicalDrive") {
		data := processLine(line, valString)
		d.DeviceID = data.(string)
	}
	if strings.HasPrefix(line, "Name") {
		data := processLine(line, valString)
		d.DeviceAlias = data.(string)
	}
	d.DeviceType = "LD"
	return nil
}

func (stats *ldStats) collect(line string) error {
	if strings.HasPrefix(line, "LogicalDrive") {
		data := processLine(line, valString)
		stats.DeviceID = data.(string)
	}
	if strings.HasPrefix(line, "Name") {
		data := processLine(line, valString)
		stats.Name = data.(string)
	}
	if strings.HasPrefix(line, "RAID Level") {
		data := processLine(line, valString)
		stats.RaidLevel = data.(string)
	}
	if strings.HasPrefix(line, "Size") {
		data := processLine(line, valString)
		stats.Size = data.(string)
	}
	if strings.HasPrefix(line, "Sector Size") {
		data := processLine(line, valInt)
		stats.SectorSize = data.(int)
	}
	if strings.HasPrefix(line, "Is VD emulated") {
		data := processLine(line, valString)
		stats.IsVDemulated = data.(string)
	}
	if strings.HasPrefix(line, "Parity Size") {
		data := processLine(line, valString)
		stats.ParitySize = data.(string)
	}
	if strings.HasPrefix(line, "State") {
		data := processLine(line, valString)
		stats.State = data.(string)
	}
	if strings.HasPrefix(line, "Strip Size") {
		data := processLine(line, valString)
		stats.StripSize = data.(string)
	}
	if strings.HasPrefix(line, "Number Of Drives per span") {
		data := processLine(line, valInt)
		stats.NumberOfDrivesPerSpan = data.(int)
	}
	if strings.HasPrefix(line, "Span Depth") {
		data := processLine(line, valInt)
		stats.SpanDepth = data.(int)
	}
	if strings.HasPrefix(line, "Current Cache Policy") {
		data := processLine(line, valString)
		stats.CurrentCachePolicy = data.(string)
	}
	if strings.HasPrefix(line, "Current Access Policy") {
		data := processLine(line, valString)
		stats.CurrentAccessPolicy = data.(string)
	}
	if strings.HasPrefix(line, "Disk Cache Policy") {
		data := processLine(line, valString)
		stats.DiskCachePolicy = data.(string)
	}
	if strings.HasPrefix(line, "Encryption Type") {
		data := processLine(line, valString)
		stats.EncryptionType = data.(string)
	}
	if strings.HasPrefix(line, "Bad Blocks Exist") {
		data := processLine(line, valString)
		stats.BadBlocksExist = data.(string)
	}
	if strings.HasPrefix(line, "Is VD Cached") {
		data := processLine(line, valString)
		stats.IsVDCached = data.(string)
	}

	return nil
}

func getStatsLogicalDrive(ldName string) {
	bin, err := getBin("megacli")
	if err != nil {
		if debug {
			log.Printf("collectLogicalDrive: %v", err)
		}
		os.Exit(1)
	}
	args := []string{"-LDInfo", "-Lall", "-aAll"}
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		if debug {
			log.Printf("collectLogicalDrive: %v", err)
		}
		os.Exit(1)
	}

	result := map[string]ldStats{}

	replaced := strings.ReplaceAll(string(out), "Virtual Drive", "Virtual DriveLogicalDrive")
	parts := strings.Split(replaced, "Virtual Drive")
	for _, ld := range parts {
		stats := ldStats{}
		for _, line := range strings.Split(ld, "\n") {
			stats.collect(line)
		}
		if len(stats.DeviceID) > 0 {
			result[stats.DeviceID] = stats
		}
	}
	if _, ok := result[ldName]; ok {
		//r, _ := json.MarshalIndent(result, "", " ")
		r, _ := json.Marshal(result[ldName])
		fmt.Print(string(r))
	} else {
		fmt.Printf("Logical Drive not exist %v", ldName)
		os.Exit(1)
	}
}

func discoveryLogicalDrive() {
	//megacli -LDGetNum -a0
	//megacli -LDInfo -Lall -aAll
	bin, err := getBin("megacli")
	if err != nil {
		if debug {
			log.Printf("discoveryLogicalDrive: %v", err)
		}
		os.Exit(1)
	}
	args := []string{"-LDInfo", "-Lall", "-aAll"}
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		if debug {
			log.Printf("discoveryLogicalDrive: %v", err)
		}
		os.Exit(1)
	}

	result := []discovery{}

	replaced := strings.ReplaceAll(string(out), "Virtual Drive", "Virtual DriveLogicalDrive")
	parts := strings.Split(replaced, "Virtual Drive")
	for _, ld := range parts {
		device := discovery{}
		for _, line := range strings.Split(ld, "\n") {
			device.discoveryLD(line)
		}
		if len(device.DeviceID) > 0 {
			result = append(result, device)
		}
	}
	r, _ := json.Marshal(result)
	fmt.Println(string(r))
}
