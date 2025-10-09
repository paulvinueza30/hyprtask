package proc

type ProcStats struct {
	cpuStats CPUStats
}
type CPUStats struct {
	uTime   uint
	sTime   uint
	cuTime  uint
	cstTime uint
	cpuTime float64 
}
