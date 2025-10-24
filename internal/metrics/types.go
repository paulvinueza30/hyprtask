package metrics

type Metrics struct{
	CPU float64
	MEM float64
}
type ProcStats struct {
	cpuStats CPUStats
	memoryStats MemoryStats
}
type CPUStats struct {
	uTime   uint
	sTime   uint
	cuTime  uint
	cstTime uint
	cpuTime float64 
}
type MemoryStats struct {
	rss int
} 

var (
	DEFAULT_METRICS = Metrics{CPU: 0, MEM: 0}
)