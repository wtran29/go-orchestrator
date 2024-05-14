package stats

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type Stats struct {
	MemStats  *mem.VirtualMemoryStat
	DiskStats *disk.UsageStat
	CpuStats  []cpu.TimesStat
	LoadStats *load.AvgStat
	TaskCount int
}

// type LoadAvg struct {
// 	Load1  float64 `json:"load1"`
// 	Load5  float64 `json:"load5"`
// 	Load15 float64 `json:"load15"`
// }

func (s *Stats) MemTotalKb() uint64 {
	return s.MemStats.Total
}

func (s *Stats) MemAvailableKb() uint64 {
	return s.MemStats.Available
}

func (s *Stats) MemUsedKb() uint64 {
	return s.MemStats.Used
}

func (s *Stats) MemUsedPercent() uint64 {
	return uint64(s.MemStats.UsedPercent)
}

func (s *Stats) DiskTotal() uint64 {
	return s.DiskStats.Total
}

func (s *Stats) DiskFree() uint64 {
	return s.DiskStats.Free
}

func (s *Stats) DiskUsed() uint64 {
	return s.DiskStats.Used
}

type CPUStat struct {
	Id        string  `json:"id"`
	User      float64 `json:"user"`
	Nice      float64 `json:"nice"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	IOWait    float64 `json:"iowait"`
	IRQ       float64 `json:"irq"`
	SoftIRQ   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guest_nice"`
}

func (s *Stats) CpuUsage() []CPUStat {
	result := make([]CPUStat, len(s.CpuStats))
	for i, stat := range s.CpuStats {
		result[i] = CPUStat{
			Id:        stat.CPU,
			User:      stat.User,
			System:    stat.System,
			Idle:      stat.Idle,
			Nice:      stat.Nice,
			IOWait:    stat.Iowait,
			IRQ:       stat.Irq,
			SoftIRQ:   stat.Softirq,
			Steal:     stat.Steal,
			Guest:     stat.Guest,
			GuestNice: stat.GuestNice,
		}
	}

	return result
}

func GetStats() *Stats {
	memStats, _ := mem.VirtualMemory()
	diskStats, _ := disk.Usage("/")
	cpuStats, _ := cpu.Times(true)
	loadAvgStats, _ := load.Avg()
	return &Stats{
		MemStats:  memStats,
		DiskStats: diskStats,
		CpuStats:  cpuStats,
		LoadStats: loadAvgStats,
	}
}

// GetMemoryInfo See https://godoc.org/github.com/c9s/goprocinfo/linux#MemInfo
// func GetMemoryInfo() *linux.MemInfo {
// 	memstats, err := linux.ReadMemInfo("/proc/meminfo")
// 	if err != nil {
// 		log.Printf("Error reading from /proc/meminfo")
// 		return &linux.MemInfo{}
// 	}
// 	return memstats
// }

// GetDiskInfo See https://godoc.org/github.com/c9s/goprocinfo/linux#Disk
// func GetDiskInfo() *linux.Disk {
// 	diskstats, err := linux.ReadDisk("/")
// 	if err != nil {
// 		log.Printf("Error reading from /")
// 		return &linux.Disk{}
// 	}
// 	return diskstats
// }

// GetCpuInfo See https://godoc.org/github.com/c9s/goprocinfo/linux#CPUStat
// func GetCpuStats() *linux.CPUStat {
// 	stats, err := linux.ReadStat("/proc/stat")
// 	if err != nil {
// 		log.Printf("Error reading from /proc/stat")
// 		return &linux.CPUStat{}
// 	}
// 	return &stats.CPUStatAll
// }

// GetLoadAvg See https://godoc.org/github.com/c9s/goprocinfo/linux#LoadAvg
// func GetLoadAvg() *linux.LoadAvg {
// 	loadavg, err := linux.ReadLoadAvg("/proc/loadavg")
// 	if err != nil {
// 		log.Printf("Error reading from /proc/loadavg")
// 		return &linux.LoadAvg{}
// 	}

// 	return loadavg
// }
