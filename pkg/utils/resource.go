package utils

import (
	"strconv"
	"strings"
)

// ConvertCPUUnit is convert cpu unit to C.
func ConvertCPUUnit(cpu string) (float64, error) {
	cpufloat := float64(0)
	if strings.Contains(cpu, "m") {
		rcpu, err := strconv.ParseFloat(strings.Split(cpu, "m")[0], 64)
		if err != nil {
			return float64(0), err
		}
		cpufloat = rcpu / 1000
	} else {
		rcpu, err := strconv.ParseFloat(cpu, 64)
		if err != nil {
			return float64(0), err
		}
		cpufloat = rcpu
	}
	return cpufloat, nil
}

// ConvertMemoryUnit is convert memory unit to Gi.
func ConvertMemoryUnit(mem string) (float64, error) {
	memfloat := float64(0)
	if strings.Contains(mem, "Gi") {
		rmem, err := strconv.ParseFloat(strings.Split(mem, "Gi")[0], 64)
		if err != nil {
			return float64(0), err
		}
		memfloat = rmem
	} else if strings.Contains(mem, "Mi") {
		rmem, err := strconv.ParseFloat(strings.Split(mem, "Mi")[0], 64)
		if err != nil {
			return float64(0), err
		}
		memfloat = rmem / 1024
	} else if strings.Contains(mem, "Ki") {
		rmem, err := strconv.ParseFloat(strings.Split(mem, "Ki")[0], 64)
		if err != nil {
			return float64(0), err
		}
		memfloat = rmem / (1024 * 1024)
	} else if strings.Contains(mem, "m") {
		rmem, err := strconv.ParseFloat(strings.Split(mem, "m")[0], 64)
		if err != nil {
			return float64(0), err
		}
		memfloat = rmem / (1024 * 1024 * 1024 * 1000)
	} else {
		rmem, err := strconv.ParseFloat(mem, 64)
		if err != nil {
			return float64(0), err
		}
		memfloat = rmem / (1024 * 1024 * 1024)
	}
	return memfloat, nil
}
