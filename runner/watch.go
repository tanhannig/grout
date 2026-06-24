package runner

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchConfig holds configuration for file watching.
type WatchConfig struct {
	Paths    []string
	Exts     []string
	Debounce time.Duration
}

// DefaultWatchConfig returns a sensible default watch configuration.
func DefaultWatchConfig(paths []string) WatchConfig {
	return WatchConfig{
		Paths:    paths,
		Exts:     []string{".go", ".py", ".js", ".ts", ".rb", ".env"},
		Debounce: 500 * time.Millisecond,
	}
}

// WatchAndRestart watches the given paths and restarts the named service
// on relevant file changes. It blocks until the done channel is closed.
func (r *Runner) WatchAndRestart(serviceName string, cfg WatchConfig, done <-chan struct{}) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for _, p := range cfg.Paths {
		if err := watcher.Add(p); err != nil {
			log.Printf("watch: could not watch %s: %v", p, err)
		}
	}

	var timer *time.Timer

	for {
		select {
		case <-done:
			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if !r.relevantEvent(event, cfg.Exts) {
				continue
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(cfg.Debounce, func() {
				log.Printf("watch: change detected, restarting %s", serviceName)
				if err := r.Restart(serviceName, DefaultRestartPolicy()); err != nil {
					log.Printf("watch: restart failed for %s: %v", serviceName, err)
				}
			})
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("watch: watcher error: %v", err)
		}
	}
}

// relevantEvent returns true if the file system event concerns a file
// whose extension is in the allowed list.
func (r *Runner) relevantEvent(event fsnotify.Event, exts []string) bool {
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
		return false
	}
	ext := strings.ToLower(filepath.Ext(event.Name))
	for _, allowed := range exts {
		if ext == allowed {
			return true
		}
	}
	return false
}
