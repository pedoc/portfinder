package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

const (
	Version = "1.0.0"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <port>\n", os.Args[0])
		fmt.Printf("Version: %s\n", Version)
		os.Exit(1)
	}

	portStr := os.Args[1]

	connections, err := net.Connections("all")
	if err != nil {
		fmt.Printf("Error: Unable to get network connection information: %v\n", err)
		os.Exit(1)
	}

	// Use map to deduplicate, key is process ID
	processMap := make(map[int32]processInfo)

	for _, conn := range connections {
		connPortStr := strconv.Itoa(int(conn.Laddr.Port))
		if strings.Contains(connPortStr, portStr) {
			proc, err := process.NewProcess(conn.Pid)
			if err != nil {
				continue
			}

			info, exists := processMap[conn.Pid]
			if !exists {
				info, err = getProcessInfo(proc)
				if err != nil {
					continue
				}
				info.Protocols = make(map[string]struct{})
				info.States = make(map[string]struct{})
				info.Ports = make(map[uint32]struct{})
				info.PortIPs = make(map[string][]string)
			}
			// Record protocol type
			protoStr := ""
			switch conn.Type {
			case 1:
				protoStr = "TCP"
			case 2:
				protoStr = "UDP"
			default:
				protoStr = "OTHER"
			}
			info.Protocols[protoStr] = struct{}{}
			// Record connection status
			info.States[conn.Status] = struct{}{}
			// Record port number
			info.Ports[conn.Laddr.Port] = struct{}{}
			// Record port and IP address mapping
			portKey := strconv.Itoa(int(conn.Laddr.Port))
			ipAddr := conn.Laddr.IP

			// Handle IPv6 address display
			if ipAddr == "::" {
				ipAddr = "[::]" // IPv6 wildcard
			} else if ipAddr == "::1" {
				ipAddr = "[::1]" // IPv6 localhost
			} else if strings.Contains(ipAddr, ":") {
				// Other IPv6 addresses, wrap in brackets
				ipAddr = "[" + ipAddr + "]"
			}

			// Check if IP already exists
			existsIP := false
			for _, existingIP := range info.PortIPs[portKey] {
				if existingIP == ipAddr {
					existsIP = true
					break
				}
			}
			if !existsIP {
				info.PortIPs[portKey] = append(info.PortIPs[portKey], ipAddr)
			}
			processMap[conn.Pid] = info
		}
	}

	if len(processMap) == 0 {
		fmt.Printf("No processes found using ports containing '%s'\n", portStr)
		return
	}

	// Convert map to slice for output
	var foundProcesses []processInfo
	for _, proc := range processMap {
		foundProcesses = append(foundProcesses, proc)
	}

	for _, proc := range foundProcesses {
		printProcessInfo(proc)
	}

	// Output statistics
	fmt.Printf("Total: %d processes found\n", len(foundProcesses))
}

type processInfo struct {
	Name        string
	PID         int32
	Command     string
	WorkDir     string
	StartedTime time.Time
	Protocols   map[string]struct{}
	States      map[string]struct{}
	Ports       map[uint32]struct{}
	PortIPs     map[string][]string
}

func getProcessInfo(proc *process.Process) (processInfo, error) {
	var info processInfo

	info.PID = proc.Pid

	// Get process name
	name, err := proc.Name()
	if err != nil {
		return info, err
	}
	info.Name = name

	// Get command line
	cmdline, err := proc.Cmdline()
	if err != nil {
		return info, err
	}
	info.Command = cmdline

	// Get working directory
	cwd, err := proc.Cwd()
	if err != nil {
		// If unable to get working directory, use default value
		info.WorkDir = "~"
	} else {
		// Simplify path display, replace user home directory with ~
		homeDir, _ := os.UserHomeDir()
		if homeDir != "" && strings.HasPrefix(cwd, homeDir) {
			info.WorkDir = "~" + cwd[len(homeDir):]
		} else {
			info.WorkDir = cwd
		}
	}

	// Get start time
	createTime, err := proc.CreateTime()
	if err != nil {
		return info, err
	}
	info.StartedTime = time.Unix(createTime/1000, 0)

	return info, nil
}

func printProcessInfo(info processInfo) {
	// Output port numbers
	if len(info.Ports) > 0 {
		var ports []int
		for port := range info.Ports {
			ports = append(ports, int(port))
		}
		// Sort in ascending order
		for i := 0; i < len(ports)-1; i++ {
			for j := i + 1; j < len(ports); j++ {
				if ports[i] > ports[j] {
					ports[i], ports[j] = ports[j], ports[i]
				}
			}
		}
		// Convert to string and include IP information
		var portStrs []string
		for _, port := range ports {
			portKey := strconv.Itoa(port)
			if ips, exists := info.PortIPs[portKey]; exists && len(ips) > 0 {
				// If there are multiple IPs, sort by IP address
				for i := 0; i < len(ips)-1; i++ {
					for j := i + 1; j < len(ips); j++ {
						if ips[i] > ips[j] {
							ips[i], ips[j] = ips[j], ips[i]
						}
					}
				}
				portStrs = append(portStrs, fmt.Sprintf("%s(%s)", portKey, strings.Join(ips, ",")))
			} else {
				portStrs = append(portStrs, portKey)
			}
		}
		fmt.Printf("Port %s\n", strings.Join(portStrs, ", "))
	}

	fmt.Printf("Process %s\n", info.Name)
	fmt.Printf("PID %d\n", info.PID)
	fmt.Printf("Command %s\n", info.Command)
	fmt.Printf("WorkDirectory %s\n", info.WorkDir)
	// Output protocol type
	if len(info.Protocols) > 0 {
		var protos []string
		for proto := range info.Protocols {
			protos = append(protos, proto)
		}
		fmt.Printf("Protocol %s\n", strings.Join(protos, ", "))
	}
	// Output connection status
	if len(info.States) > 0 {
		var states []string
		for state := range info.States {
			states = append(states, state)
		}
		fmt.Printf("Status %s\n", strings.Join(states, ", "))
	}

	// Calculate running time
	duration := time.Since(info.StartedTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	var timeStr string
	if hours > 0 {
		timeStr = fmt.Sprintf("%dh", hours)
		if minutes > 0 {
			timeStr += fmt.Sprintf(" %dm", minutes)
		}
	} else {
		timeStr = fmt.Sprintf("%dm", minutes)
	}

	fmt.Printf("Started %s\n", timeStr)
	fmt.Println()
}
