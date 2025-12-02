package procprovider

import (
	"os/user"
	"strconv"
	"strings"
	"sync"

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
	fs          procfs.FS
	userCache   map[int]string // UID -> username cache
	userCacheMu sync.RWMutex
}

const (
	// Worker pool size for parallel process reading
	workerPoolSize = 50
)

func NewProcProvider() ProcProvider {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		logger.Log.Error("could not get procfs: " + err.Error())
		return ProcProvider{}
	}
	return ProcProvider{
		fs:        fs,
		userCache: make(map[int]string),
	}
}

func (p *ProcProvider) GetProcs() ([]Proc, error) {
	procs, err := p.fs.AllProcs()
	if err != nil {
		logger.Log.Error("could not get all procs: " + err.Error())
		return nil, err
	}

	// Use worker pool pattern for parallel processing
	procChan := make(chan procfs.Proc, len(procs))
	resultChan := make(chan Proc, len(procs))
	
	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for proc := range procChan {
				procData := p.readProcData(proc)
				resultChan <- procData
			}
		}()
	}

	// Send all processes to workers
	go func() {
		for _, proc := range procs {
			procChan <- proc
		}
		close(procChan)
	}()

	// Close result channel when all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	procList := make([]Proc, 0, len(procs))
	for procData := range resultChan {
		procList = append(procList, procData)
	}

	return procList, nil
}

func (p *ProcProvider) readProcData(proc procfs.Proc) Proc {
	procData := Proc{PID: proc.PID}

	if comm, err := proc.Comm(); err == nil {
		procData.ProgramName = comm
	}

	if cmdline, err := proc.CmdLine(); err == nil {
		procData.CommandLine = strings.Join(cmdline, " ")
	}

	if status, err := proc.NewStatus(); err == nil {
		uid := int(status.UIDs[0])
		procData.User = p.getUsername(uid)
	}

	return procData
}

func (p *ProcProvider) getUsername(uid int) string {
	// Check cache first
	p.userCacheMu.RLock()
	if username, ok := p.userCache[uid]; ok {
		p.userCacheMu.RUnlock()
		return username
	}
	p.userCacheMu.RUnlock()

	// Cache miss - lookup user
	var username string
	if u, err := user.LookupId(strconv.Itoa(uid)); err == nil {
		username = u.Username
	} else {
		username = strconv.Itoa(uid)
	}

	// Update cache
	p.userCacheMu.Lock()
	p.userCache[uid] = username
	p.userCacheMu.Unlock()

	return username
}