package sources

import (
	"fmt"

	"github.com/jerryluo/nettui/internal/data"
	"github.com/shirou/gopsutil/v4/process"
)

// CollectProcesses gathers all running processes.
func CollectProcesses() ([]data.Process, []data.CollectionError) {
	var errs []data.CollectionError

	procs, err := process.Processes()
	if err != nil {
		return nil, []data.CollectionError{{Source: "processes", Error: fmt.Sprintf("Processes(): %v", err)}}
	}

	result := make([]data.Process, 0, len(procs))
	for _, p := range procs {
		name, err := p.Name()
		if err != nil {
			name = ""
		}

		cmdline, err := p.Cmdline()
		if err != nil {
			cmdline = ""
		}

		user, err := p.Username()
		if err != nil {
			user = ""
		}

		result = append(result, data.Process{
			PID:     p.Pid,
			Name:    name,
			Command: cmdline,
			User:    user,
		})
	}

	return result, errs
}
