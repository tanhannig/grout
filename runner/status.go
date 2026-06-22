package runner

import (
	"fmt"
	"time"
)

// ServiceStatus represents the current state of a managed service.
type ServiceStatus struct {
	Name    string
	Running bool
	PID     int
	Uptime  time.Duration
	Error   string
}

// Status returns the current status of all managed services.
func (r *Runner) Status() []ServiceStatus {
	r.mu.Lock()
	defer r.mu.Unlock()

	statuses := make([]ServiceStatus, 0, len(r.services))
	for _, svc := range r.services {
		s := ServiceStatus{
			Name: svc.name,
		}
		if svc.cmd != nil && svc.cmd.Process != nil {
			s.PID = svc.cmd.Process.Pid
			s.Running = svc.running
			if svc.startedAt != (time.Time{}) {
				s.Uptime = time.Since(svc.startedAt).Round(time.Second)
			}
		}
		if svc.lastErr != nil {
			s.Error = svc.lastErr.Error()
		}
		statuses = append(statuses, s)
	}
	return statuses
}

// FormatStatus returns a human-readable table of service statuses.
func FormatStatus(statuses []ServiceStatus) string {
	if len(statuses) == 0 {
		return "No services defined.\n"
	}
	out := fmt.Sprintf("%-20s %-10s %-8s %-12s %s\n", "NAME", "STATUS", "PID", "UPTIME", "ERROR")
	out += fmt.Sprintf("%-20s %-10s %-8s %-12s %s\n",
		"----", "------", "---", "------", "-----")
	for _, s := range statuses {
		status := "stopped"
		if s.Running {
			status = "running"
		}
		pid := "-"
		if s.PID != 0 {
			pid = fmt.Sprintf("%d", s.PID)
		}
		uptime := "-"
		if s.Running && s.Uptime > 0 {
			uptime = s.Uptime.String()
		}
		errStr := "-"
		if s.Error != "" {
			errStr = s.Error
		}
		out += fmt.Sprintf("%-20s %-10s %-8s %-12s %s\n",
			s.Name, status, pid, uptime, errStr)
	}
	return out
}
