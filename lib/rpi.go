package rpi

import (
	"fmt"
	"os/exec"
)

func getHostname() (result string, success bool) {
	output, err := exec.Command("hostname").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	return
}

func getUname() (result string, success bool) {
	output, err := exec.Command("uname", "-a").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	return
}

func getUptime() (result string, success bool) {
	output, err := exec.Command("uptime").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	return
}

func getFreeSpaces() (result string, success bool) {
	output, err := exec.Command("df", "-h").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	return
}

func getMemorySplit() (result string, success bool) {
	// arm memory
	output, err := exec.Command("vcgencmd", "get_mem", "arm").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	// gpu memory
	if success {
		output, err = exec.Command("vcgencmd", "get_mem", "gpu").CombinedOutput()
		result = result + string(output)
		if err == nil {
			success = true
		} else {
			success = false
		}
	}
	return
}

func getFreeMemory() (result string, success bool) {
	output, err := exec.Command("free", "-o", "-h").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	return
}

func getCpuTemperature() (result string, success bool) {
	output, err := exec.Command("vcgencmd", "measure_temp").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	return
}

func getCpuInfo() (result string, success bool) {
	output, err := exec.Command("cat", "/proc/cpuinfo").CombinedOutput()
	result = string(output)
	if err == nil {
		success = true
	} else {
		success = false
	}
	return
}

/*
	Read value for given method.
*/
func ReadValue(method string) (result string, success bool) {
	switch method {
	case "hostname": // hostname
		result, success = getHostname()
	case "uname": // uname -a
		result, success = getUname()
	case "uptime": // uptime
		result, success = getUptime()
	case "free_spaces": // df -h
		result, success = getFreeSpaces()
	case "memory_split": // vcgencmd get_mem arm && vcgencmd get_mem gpu
		result, success = getMemorySplit()
	case "free_memory": // free -o -h
		result, success = getFreeMemory()
	case "cpu_temperature": // vcgencmd measure_temp
		result, success = getCpuTemperature()
	case "cpu_info": //cat /proc/cpuinfo
		result, success = getCpuInfo()
	default:
		result = fmt.Sprintf("No such method: %s", method)
		success = false
	}

	return result, success
}
