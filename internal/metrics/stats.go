package metrics

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/prometheus/procfs"
)

type SystemMonitor struct {
	fs          procfs.FS
	totalMemory uint64
	pageSize    int

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
	pageSize := os.Getpagesize()

	return &SystemMonitor{fs: fs, totalMemory: *memInfo.MemTotal, pageSize: pageSize, tickDuration: tickDuration}, nil
}

func (m *SystemMonitor) GetMetrics(pid int) (*Metrics, error) {
	beforeStats, err := m.getProcStats(pid)
	if err != nil {
		return nil , err
	}
	
	beforeTime, err := m.getTotalTime()
	if err != nil {
		return nil , err
	}
	
	time.Sleep(m.tickDuration)
	
	afterStats, err := m.getProcStats(pid)
	if err != nil {
		return nil , err
	}
	afterTime, err := m.getTotalTime()
	if err != nil {
		return nil , err
	}
	totalTime := *afterTime - *beforeTime
	
	cpuUsage := m.calcCpuUsage(beforeStats.cpuStats, afterStats.cpuStats, totalTime)
	memUsage := m.calcMemoryUsage(afterStats.memoryStats)
	
	metrics := &Metrics{CPU: cpuUsage, MEM: memUsage, } 
	logger.Log.Info(fmt.Sprintf("usage stats for proc(%d) :", pid), "metrics " , metrics)
	return metrics , nil
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
		logger.Log.Error("could not get proc: " + err.Error())
		return nil, err
	}

	stat, err := proc.Stat()
	if err != nil {
		logger.Log.Error("could not get proc stats: " + err.Error())
		return nil, err
	}

	cpuStats := CPUStats{cuTime: uint(stat.CUTime), cstTime: uint(stat.CSTime), sTime: stat.STime, uTime: stat.UTime}
	memStats := MemoryStats{rss: stat.RSS}

	return &ProcStats{cpuStats: cpuStats, memoryStats: memStats}, nil
}

func (m *SystemMonitor) calcCpuUsage(before, after CPUStats, totalTime float64) float64 {
	if totalTime == 0 {
		return 0.0
	}
	userUtil := after.uTime - before.uTime
	sysUtil := after.sTime - before.sTime
	processDelta := userUtil + sysUtil
	return float64(processDelta) / totalTime * 100.0
}

func (m *SystemMonitor) calcMemoryUsage(memStats MemoryStats) float64 {
	residentMemorySize := uint64(memStats.rss) * uint64(m.pageSize)
	memoryTotal := uint64(m.totalMemory * 1024) // kB -> b
	return float64(residentMemorySize) / float64(memoryTotal) * 100.0
}
