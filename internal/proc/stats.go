package proc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/prometheus/procfs"
)

var fs procfs.FS
var err error

func GetStats(pid int) {

	fs, err = procfs.NewFS("/proc")
	if err != nil {
		logger.Log.Error("cannot set fs procfs: " + err.Error())
		return
	}
	statChannel := make(chan float64)
	go getStats(pid, statChannel)
	statRes := <-statChannel
	logger.Log.Info("cpu usage for PID " + fmt.Sprintf("(%d)", pid) + ":" + fmt.Sprintf(" %f", statRes))
	
}

func getStats(pid int , statChan chan float64) {
	beforeStats := getProcStats(pid)
	beforeTime := getTotalTime()
	time.Sleep(2 * time.Second)
	afterStats := getProcStats(pid)
	afterTime := getTotalTime()
	
	totalTime := afterTime - beforeTime
	cpuUsage := calcCpuUsage(beforeStats.cpuStats, afterStats.cpuStats, totalTime)
	statChan <- cpuUsage
}
func getTotalTime() float64{
	stat , err := fs.Stat()
	if err != nil {
		logger.Log.Error("could not get fs proc start " + err.Error())
	}
	cpuTotal := stat.CPUTotal
	
	v := reflect.ValueOf(cpuTotal)

	var totalTicks float64
	for i := range v.NumField(){
		totalTicks += v.Field(i).Float()
	}
	return totalTicks
}
func getProcStats(pid int) ProcStats {
	proc, err := fs.Proc(pid)
	if err != nil {
		logger.Log.Error("could not get proc: " + err.Error())
	}
	// logger.Log.Info("proc "+fmt.Sprintf("%d", pid)+" details: ", fmt.Sprintf("%+v", proc))

	stat, err := proc.Stat()
	if err != nil {
		logger.Log.Error("could not get proc stats: " + err.Error())
	}
	// logger.Log.Info("proc "+fmt.Sprintf("%d", proc.PID), " stats: ", fmt.Sprintf("%+v", stat))

	cpuStats := CPUStats{cuTime: uint(stat.CUTime), cstTime: uint(stat.CSTime), sTime: stat.STime, uTime: stat.UTime}

	return ProcStats{cpuStats: cpuStats}
}

func calcCpuUsage(before , after CPUStats, totalTime float64)  float64{
	if totalTime == 0 {
		return 0.0
	}
	userUtil := after.uTime - before.uTime
	sysUtil := after.sTime - before.sTime
	processDelta := userUtil + sysUtil
	return 100.0 * float64(processDelta)/  totalTime
}