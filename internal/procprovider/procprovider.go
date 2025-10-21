package procprovider

import (
	"os/user"
	"strconv"
	"strings"

	"github.com/paulvinueza30/hyprtask/internal/logger"
	"github.com/prometheus/procfs"
)
type Proc struct {
	PID         int
	ProgramName string
	User        string
	CommandLine string
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
	procList := make([]Proc, 0, len(procs))
	for _, proc := range procs {
		procData := Proc{PID: proc.PID}
		
		if comm, err := proc.Comm(); err == nil {
			procData.ProgramName = comm
		}
		
		if cmdline, err := proc.CmdLine(); err == nil {
			procData.CommandLine = strings.Join(cmdline, " ")
		}
		
		if status, err := proc.NewStatus(); err == nil {
			uid := status.UIDs[0]
			if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
				procData.User = u.Username
			} else {
				procData.User = strconv.Itoa(int(uid))
			}
		}
		
		procList = append(procList, procData)
	}
	return procList, nil
}