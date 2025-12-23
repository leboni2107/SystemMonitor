package main

import (
	"fmt"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

type DiskUsage struct {
	Path              string  `json:"path"`
	FSType            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type CPUInfo struct {
	CPU        int      `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int      `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int      `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int      `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

type MemoryInfo struct {
	Total          uint64  `json:"total"`
	Available      uint64  `json:"available"`
	Used           uint64  `json:"used"`
	UsedPercent    float64 `json:"usedPercent"`
	Free           uint64  `json:"free"`
	Active         uint64  `json:"active"`
	Inactive       uint64  `json:"inactive"`
	Wired          uint64  `json:"wired"`
	Laundry        uint64  `json:"laundry"`
	Buffers        uint64  `json:"buffers"`
	Cached         uint64  `json:"cached"`
	WriteBack      uint64  `json:"writeBack"`
	Dirty          uint64  `json:"dirty"`
	WriteBackTmp   uint64  `json:"writeBackTmp"`
	Shared         uint64  `json:"shared"`
	Slab           uint64  `json:"slab"`
	Sreclaimable   uint64  `json:"sreclaimable"`
	Sunreclaim     uint64  `json:"sunreclaim"`
	PageTables     uint64  `json:"pageTables"`
	SwapCached     uint64  `json:"swapCached"`
	CommitLimit    uint64  `json:"commitLimit"`
	CommittedAS    uint64  `json:"committedAS"`
	HighTotal      uint64  `json:"highTotal"`
	HighFree       uint64  `json:"highFree"`
	LowTotal       uint64  `json:"lowTotal"`
	LowFree        uint64  `json:"lowFree"`
	SwapTotal      uint64  `json:"swapTotal"`
	SwapFree       uint64  `json:"swapFree"`
	Mapped         uint64  `json:"mapped"`
	VmallocTotal   uint64  `json:"vmallocTotal"`
	VmallocUsed    uint64  `json:"vmallocUsed"`
	VmallocChunk   uint64  `json:"vmallocChunk"`
	HugePagesTotal uint64  `json:"hugePagesTotal"`
	HugePagesFree  uint64  `json:"hugePagesFree"`
	HugePagesRsvd  uint64  `json:"hugePagesRsvd"`
	HugePagesSurp  uint64  `json:"hugePagesSurp"`
	HugePageSize   uint64  `json:"hugePageSize"`
	AnonHugePages  uint64  `json:"anonHugePages"`
}

func main() {
	diskUsage, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}

	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CPU Model:", cpuInfo[0].ModelName)
	fmt.Println("Disk Used:", diskUsage.UsedPercent, "%")
	fmt.Println("Memory Used:", memoryInfo.UsedPercent, "%")
}
