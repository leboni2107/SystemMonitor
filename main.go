package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func getData() (*disk.UsageStat, []cpu.InfoStat, []float64, *mem.VirtualMemoryStat, error) {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return diskUsage, cpuInfo, cpuPercent, memoryInfo, nil
}

func print(diskUsage *disk.UsageStat, cpuInfo []cpu.InfoStat, cpuPercent []float64, memoryInfo *mem.VirtualMemoryStat, err error) {
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("---System Monitor---")
	fmt.Println("CPU Model:", cpuInfo[0].ModelName)
	fmt.Printf("CPU Usage: %.2f%%\n", cpuPercent[0])
	fmt.Printf("Disk Used: %.2f%%\n", diskUsage.UsedPercent)
	fmt.Printf("Memory Used: %.2f%%\n", memoryInfo.UsedPercent)

	// Wait 2 seconds before next update
	time.Sleep(2 * time.Second)
}

func main() {
	diskUsage, cpuInfo, cpuPercent, memoryInfo, err := getData()
	print(diskUsage, cpuInfo, cpuPercent, memoryInfo, err)
}
