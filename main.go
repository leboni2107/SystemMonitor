package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/sensors"
)

const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Underline = "\033[4m"

	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	ScreenWidth = 36
) // ANSI Formatting

func GetData() (*disk.UsageStat, []cpu.InfoStat, []float64, *mem.VirtualMemoryStat, []net.IOCountersStat, []sensors.TemperatureStat, error) {
	var diskUsage *disk.UsageStat
	var cpuInfo []cpu.InfoStat
	var cpuPercent []float64
	var memoryInfo *mem.VirtualMemoryStat
	var netInfo []net.IOCountersStat
	var temperatures []sensors.TemperatureStat
	var err error

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		diskUsage, _ = disk.Usage("/")
	}()

	go func() {
		defer wg.Done()
		cpuInfo, _ = cpu.Info()
	}()

	go func() {
		defer wg.Done()
		cpuPercent, _ = cpu.Percent(time.Second, false)
	}()

	go func() {
		defer wg.Done()
		memoryInfo, _ = mem.VirtualMemory()
	}()

	go func() {
		defer wg.Done()
		netInfo, _ = net.IOCounters(true)
	}()

	go func() {
		defer wg.Done()
		temperatures, _ = sensors.TemperaturesWithContext(context.Background())
	}()

	wg.Wait()
	return diskUsage, cpuInfo, cpuPercent, memoryInfo, netInfo, temperatures, err
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			return
		}
	}
}

func printValue(value string, cursorTop int, cursorLeft int, screenWidth int) {
	fmt.Printf("\033[%d;%dH%s", cursorTop, cursorLeft, value)
	fmt.Printf("\033[%d;%dH┃\n", cursorTop, screenWidth)
}

func PrintMenu(diskUsage *disk.UsageStat, cpuInfo []cpu.InfoStat, cpuPercent []float64, memoryInfo *mem.VirtualMemoryStat, netInfo []net.IOCountersStat, tempInfo []sensors.TemperatureStat, err error, BytesRecvDelta float64) {
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓\n")
	fmt.Printf("┃          %s%sSystem Monitor%s          ┃\n", Cyan, Bold, Reset)
	fmt.Printf("┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫\n")

	padding := 0
	modelLine := ""
	if len(cpuInfo[0].ModelName) > ScreenWidth/2+1 {
		padding = 1
		modelLine = fmt.Sprintf("┃ %sCPU Model:   %s%s%s%s\n┃              %s",
			Blue, Reset, cpuInfo[0].ModelName[:ScreenWidth/2], "-", Reset, cpuInfo[0].ModelName[18:])

		printValue(modelLine, 4, 0, ScreenWidth)
		fmt.Printf("\033[5;%sH", ScreenWidth) // Set cursor
		fmt.Printf("┃")
	} else {
		modelLine = fmt.Sprintf("┃ %s%sCPU Model:  %s %s", Blue, Bold, Reset, cpuInfo[0].ModelName)
		printValue(modelLine, 4, 0, ScreenWidth)
	}

	cpuLine := fmt.Sprintf("┃ %s%sCPU Used:%s    %s [%.2f%%]", Blue, Bold, Reset, GetProgressBar(int(cpuPercent[0]), 10), cpuPercent[0])
	printValue(cpuLine, 5+padding, 0, ScreenWidth)

	memLine := fmt.Sprintf("┃ %s%sMemory Used:%s %s [%.2f%%]", Yellow, Bold, Reset, GetProgressBar(int(memoryInfo.UsedPercent), 10), memoryInfo.UsedPercent)
	printValue(memLine, 7+padding, 0, ScreenWidth)

	diskLine := fmt.Sprintf("┃ %s%sDisk Used:%s   %s [%.2f%%]", Green, Bold, Reset, GetProgressBar(int(diskUsage.UsedPercent), 10), diskUsage.UsedPercent)
	printValue(diskLine, 6+padding, 0, ScreenWidth)

	tempLine := fmt.Sprintf("┃ %s%sTemps:%s       [%.2f°C]", Red, Bold, Reset, tempInfo[0].Temperature)
	printValue(tempLine, 8+padding, 0, ScreenWidth)

	prefixes := [6]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	i := 0
	for i = 0; BytesRecvDelta >= 1024; i++ {
		BytesRecvDelta /= 1024
	}
	netLine := fmt.Sprintf("┃ %s%sNetwork:%s     [%.2f%s]", Magenta, Bold, Reset, BytesRecvDelta, prefixes[i])
	printValue(netLine, 9+padding, 0, ScreenWidth)

	netLine2 := fmt.Sprintf("┃ %s%sConnections:%s [%d]", Magenta, Bold, Reset, len(netInfo))
	printValue(netLine2, 10+padding, 0, ScreenWidth)

	fmt.Printf("\033[%d;0H", 11+padding)
	fmt.Printf("┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛\n")
	fmt.Printf("\033[%d;0H", 12+padding)
}

func GetProgressBar(progress int, base int) string {
	p2 := float64(progress) / 100.0
	p3 := p2 * float64(base)
	progress = int(p3)

	bar := ""
	for i := 0; i < base; i++ {
		if i < progress {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}

func main() {
	/*
	 * TODO: Net Graph 0-Max
	 * TODO: Parallel Network Menu
	 */
	diskUsage, cpuInfo, cpuPercent, memoryInfo, netInfo, temperatureInfo, err := GetData()
	for {
		BytesRecvLastIt := netInfo[0].BytesRecv
		diskUsage, cpuInfo, cpuPercent, memoryInfo, netInfo, temperatureInfo, err = GetData()
		BytesRecvDelta := netInfo[0].BytesRecv - BytesRecvLastIt

		ClearScreen()
		PrintMenu(diskUsage, cpuInfo, cpuPercent, memoryInfo, netInfo, temperatureInfo, err, float64(BytesRecvDelta))

		time.Sleep(2000 * time.Millisecond)
	}
}
