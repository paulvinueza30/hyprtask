package procprovider

import (
	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/prometheus/procfs"
)
type Proc struct {
	PID int 
}
type ProcProvider struct {
	fs procfs.FS
}

func NewProcProvider() ProcProvider {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		logger.Log.Error("could not get procfs: " + err.Error())
		return ProcProvider{}
	}
	return ProcProvider{fs: fs}
}

func (p *ProcProvider) GetProcs() ([]Proc, error) {
	procs, err := p.fs.AllProcs()
	if err != nil {
		logger.Log.Error("could not get all procs: " + err.Error())
		return nil, err
	}
	procList:= make([]Proc, 0, len(procs))
	for _, p := range procs {
		procList = append(procList, Proc{
			PID: p.PID,
		})
	}
	return procList, nil
}