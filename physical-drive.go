package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type pdStats struct {
	DeviceID                            string
	WWN                                 string
	MediaErrorCount                     int
	OtherErrorCount                     int
	PredictiveFailureCount              int
	LastPredictiveFailureEventSeqNumber int
	PDType                              string
	RawSize                             string
	SectorSize                          int
	LogicalSectorSize                   int
	PhysicalSectorSize                  int
	FirmwareState                       string
	CommissionedSpare                   string
	EmergencySpare                      string
	DeviceFirmwareLevel                 string
	ShieldCounter                       int
	InquiryData                         string
	Secured                             string
	Locked                              string
	NeedsEKMAttention                   string
	DeviceSpeed                         string
	LinkSpeed                           string
	MediaType                           string
	DriveTemperature                    string
	DriveSMARTAlert                     string
}

func (d *discovery) discoveryPD(line string) error {
	if strings.HasPrefix(line, "Drive's position") {
		data := processLine(line, valString)
		d.DeviceID = data.(string)
	}
	if strings.HasPrefix(line, "WWN") {
		data := processLine(line, valString)
		d.DeviceAlias = data.(string)
	}
	d.DeviceType = "PD"
	return nil
}

func (stats *pdStats) collect(line string) error {
	if strings.HasPrefix(line, "Drive's position") {
		data := processLine(line, valString)
		stats.DeviceID = data.(string)
	}
	if strings.HasPrefix(line, "WWN") {
		data := processLine(line, valString)
		stats.WWN = data.(string)
	}
	if strings.HasPrefix(line, "Media Error Count") {
		data := processLine(line, valInt)
		stats.MediaErrorCount = data.(int)
	}
	if strings.HasPrefix(line, "Other Error Count") {
		data := processLine(line, valInt)
		stats.OtherErrorCount = data.(int)
	}
	if strings.HasPrefix(line, "Predictive Failure Count") {
		data := processLine(line, valInt)
		stats.PredictiveFailureCount = data.(int)
	}
	if strings.HasPrefix(line, "Last Predictive Failure Event Seq Number") {
		data := processLine(line, valInt)
		stats.LastPredictiveFailureEventSeqNumber = data.(int)
	}
	if strings.HasPrefix(line, "PD Type") {
		data := processLine(line, valString)
		stats.PDType = data.(string)
	}
	if strings.HasPrefix(line, "Raw Size") {
		data := processLine(line, valString)
		stats.RawSize = data.(string)
	}
	if strings.HasPrefix(line, "Sector Size") {
		data := processLine(line, valInt)
		stats.SectorSize = data.(int)
	}
	if strings.HasPrefix(line, "Logical Sector Size") {
		data := processLine(line, valInt)
		stats.LogicalSectorSize = data.(int)
	}
	if strings.HasPrefix(line, "Physical Sector Size") {
		data := processLine(line, valInt)
		stats.PhysicalSectorSize = data.(int)
	}
	if strings.HasPrefix(line, "Firmware state") {
		data := processLine(line, valString)
		stats.FirmwareState = data.(string)
	}
	if strings.HasPrefix(line, "Commissioned Spare") {
		data := processLine(line, valString)
		stats.CommissionedSpare = data.(string)
	}
	if strings.HasPrefix(line, "Emergency Spare") {
		data := processLine(line, valString)
		stats.EmergencySpare = data.(string)
	}
	if strings.HasPrefix(line, "Device Firmware Level") {
		data := processLine(line, valString)
		stats.DeviceFirmwareLevel = data.(string)
	}
	if strings.HasPrefix(line, "Shield Counter") {
		data := processLine(line, valInt)
		stats.ShieldCounter = data.(int)
	}
	if strings.HasPrefix(line, "Inquiry Data") {
		data := processLine(line, valString)
		stats.InquiryData = data.(string)
	}
	if strings.HasPrefix(line, "Secured") {
		data := processLine(line, valString)
		stats.Secured = data.(string)
	}
	if strings.HasPrefix(line, "Locked") {
		data := processLine(line, valString)
		stats.Locked = data.(string)
	}
	if strings.HasPrefix(line, "Needs EKM Attention") {
		data := processLine(line, valString)
		stats.NeedsEKMAttention = data.(string)
	}
	if strings.HasPrefix(line, "Device Speed") {
		data := processLine(line, valString)
		stats.DeviceSpeed = data.(string)
	}
	if strings.HasPrefix(line, "Link Speed") {
		data := processLine(line, valString)
		stats.LinkSpeed = data.(string)
	}
	if strings.HasPrefix(line, "Media Type") {
		data := processLine(line, valString)
		stats.MediaType = data.(string)
	}
	if strings.HasPrefix(line, "DriveTemperature") {
		data := processLine(line, valString)
		stats.DriveTemperature = data.(string)
	}
	if strings.HasPrefix(line, "Drive has flagged a S.M.A.R.T alert") {
		data := processLine(line, valString)
		stats.DriveSMARTAlert = data.(string)
	}
	return nil
}

func getStatsPhysicalDrive(pdName string) {
	bin, err := getBin("megacli")
	if err != nil {
		if debug {
			log.Printf("getStatsPhysicalDrive: %v", err)
		}
		os.Exit(1)
	}
	args := []string{"-PDList", "-aAll"}
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		if debug {
			log.Printf("getStatsPhysicalDrive: %v", err)
		}
		os.Exit(1)
	}

	result := map[string]pdStats{}

	parts := strings.Split(string(out), "Enclosure Device ID")
	for _, ld := range parts {
		stats := pdStats{}
		for _, line := range strings.Split(ld, "\n") {
			stats.collect(line)
		}
		if len(stats.DeviceID) > 0 {
			result[stats.DeviceID] = stats
		}
	}
	if _, ok := result[pdName]; ok {
		//r, _ := json.MarshalIndent(result[pdName], "", " ")
		r, _ := json.Marshal(result[pdName])
		fmt.Print(string(r))
	} else {
		fmt.Printf("Physical Drive not exist %v", pdName)
		os.Exit(1)
	}
}

func discoveryPhysicalDrive() {
	//megacli -PDGetNum -a0
	//megacli -PDList -a0
	bin, err := getBin("megacli")
	if err != nil {
		if debug {
			log.Printf("discoveryPhisycalDrive: %v", err)
		}
		os.Exit(1)
	}
	args := []string{"-PDList", "-aAll"}
	cmd := exec.Command(bin, args...)
	out, err := cmd.Output()
	if err != nil {
		if debug {
			log.Printf("discoveryPhisycalDrive: %v", err)
		}
		os.Exit(1)
	}

	result := []discovery{}

	parts := strings.Split(string(out), "Enclosure Device ID")
	for _, pd := range parts {
		device := discovery{}
		for _, line := range strings.Split(pd, "\n") {
			device.discoveryPD(line)
		}
		if len(device.DeviceID) > 0 {
			result = append(result, device)
		}
	}
	//r, _ := json.MarshalIndent(result, "", "  ")
	r, _ := json.Marshal(result)
	fmt.Println(string(r))
}
