package metrics

import (
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/prometheus/procfs"
)

type SystemMonitor struct {
	fs          procfs.FS
	totalMemory uint64
	pageSize    int
	clockRate   int

	tickDuration time.Duration
}

func NewSystemMonitor(tickDuration time.Duration) (*SystemMonitor, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		logger.Log.Error("cannot set fs procfs: " + err.Error())
		return nil, err
	}

	memInfo, err := fs.Meminfo()
	if err != nil {
		logger.Log.Error("could not get fs memory info: " + err.Error())
		return nil, err
	}
	// Get the system clock rate
	clockRate := getSystemClockRate()
	pageSize := os.Getpagesize()

	return &SystemMonitor{
		fs: fs, 
		totalMemory: *memInfo.MemTotal, 
		pageSize: pageSize, 
		clockRate: clockRate,
		tickDuration: tickDuration,
	}, nil
}

func (m *SystemMonitor) GetMetrics(pid int) (*Metrics, error) {
	beforeStats, err := m.getProcStats(pid)
	if err != nil {
		return &DEFAULT_METRICS, err
	}
	beforeTime, err := m.getTotalTime()
	if err != nil {
		return nil , err
	}
	
	time.Sleep(m.tickDuration)
	
	afterStats, err := m.getProcStats(pid)
	if err != nil {
		return &DEFAULT_METRICS, err
	}
	afterTime, err := m.getTotalTime()
	if err != nil {
		return nil , err
	}
	totalTime := *afterTime - *beforeTime
	
	cpuUsage := m.calcCpuUsage(beforeStats.cpuStats, afterStats.cpuStats, totalTime)
	memUsage := m.calcMemoryUsage(afterStats.memoryStats)
	
	metrics := &Metrics{CPU: cpuUsage, MEM: memUsage, } 
	return metrics , nil
}

// GetQuickMetrics returns metrics without the sleep delay for immediate display
// CPU will be 0.0 as it requires time delta to calculate accurately
// Memory usage is calculated immediately from current stats
func (m *SystemMonitor) GetQuickMetrics(pid int) (*Metrics, error) {
	stats, err := m.getProcStats(pid)
	if err != nil {
		return &DEFAULT_METRICS, err
	}
	
	memUsage := m.calcMemoryUsage(stats.memoryStats)
	
	metrics := &Metrics{CPU: 0.0, MEM: memUsage}
	return metrics, nil
}
func (m *SystemMonitor) getTotalTime() (*float64, error) {
	stat, err := m.fs.Stat()
	if err != nil {
		return nil, err
	}
	cpuTotal := stat.CPUTotal

	v := reflect.ValueOf(cpuTotal)

	var totalTicks float64
	for i := range v.NumField() {
		totalTicks += v.Field(i).Float()
	}
	return &totalTicks, err
}
func (m *SystemMonitor) getProcStats(pid int) (*ProcStats, error) {
	proc, err := m.fs.Proc(pid)
	if err != nil {
		return nil, err
	}

	stat, err := proc.Stat()
	if err != nil {
		return nil, err
	}

	cpuStats := CPUStats{cuTime: uint(stat.CUTime), cstTime: uint(stat.CSTime), sTime: stat.STime, uTime: stat.UTime}
	memStats := MemoryStats{rss: stat.RSS}

	return &ProcStats{cpuStats: cpuStats, memoryStats: memStats}, nil
}

func getSystemClockRate() int {
	if data, err := os.ReadFile("/proc/sys/kernel/hz"); err == nil {
		if rate, err := strconv.Atoi(string(data[:len(data)-1])); err == nil {
			return rate
		}
	}
	
	// Fallback to 100 Hz (most common on modern systems)
	return 100
}

func (m *SystemMonitor) calcCpuUsage(before, after CPUStats, totalTime float64) float64 {
	if totalTime == 0 {
		return 0.0
	}
	userUtil := after.uTime - before.uTime
	sysUtil := after.sTime - before.sTime
	processDelta := userUtil + sysUtil
	
	// Convert jiffies to seconds 
	processTimeSeconds := float64(processDelta) / float64(m.clockRate)
	
	return (processTimeSeconds / totalTime) * 100.0
}

func (m *SystemMonitor) calcMemoryUsage(memStats MemoryStats) float64 {
	residentMemorySize := uint64(memStats.rss) * uint64(m.pageSize)
	memoryTotal := uint64(m.totalMemory * 1024) // kB -> b
	return float64(residentMemorySize) / float64(memoryTotal) * 100.0
}
