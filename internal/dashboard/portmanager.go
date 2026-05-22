package dashboard

import (
	"fmt"
	"net"
	"sync"

	"github.com/pocketbase/pocketbase/core"
)

// PortManager manages port allocation across the system
type PortManager struct {
	mu   sync.Mutex
	used map[int]bool
}

func NewPortManager() *PortManager {
	return &PortManager{
		used: make(map[int]bool),
	}
}

func (pm *PortManager) Next(app core.App) (int, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Get all used ports from database
	dbPorts := make(map[int]bool)

	// Sites ports
	if siteRecords, err := app.FindAllRecords("sites"); err == nil {
		for _, rec := range siteRecords {
			port := rec.GetInt("port")
			if port > 0 {
				dbPorts[port] = true
			}
		}
	}

	// Databases ports
	if dbRecords, err := app.FindAllRecords("databases"); err == nil {
		for _, rec := range dbRecords {
			port := rec.GetInt("port")
			if port > 0 {
				dbPorts[port] = true
			}
		}
	}

	// Scan for available ports in real-time
	for port := PortRangeStart; port <= PortRangeEnd; port++ {
		// Skip if used in DB
		if dbPorts[port] {
			continue
		}
		// Skip if reserved in-memory during this session
		if pm.used[port] {
			continue
		}
		// Skip if actually occupied on the host system
		if !pm.isPortInUse(port) {
			pm.used[port] = true
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports in range %d-%d", PortRangeStart, PortRangeEnd)
}

func (pm *PortManager) Release(port int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.used, port)
}

// isPortInUse checks if a port is actually occupied on the system
func (pm *PortManager) isPortInUse(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true // port is in use
	}
	ln.Close()
	return false
}
