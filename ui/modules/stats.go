package modules

import (
	"fmt"
	"math"
	"syscall"
)

// SidebarStats contains compact system information for the sidebar footer.
type SidebarStats struct {
	Version     string
	DiskUsed    string
	DiskTotal   string
	DiskPercent int
}

// GetSidebarStats returns footer-ready server and disk stats.
func GetSidebarStats() SidebarStats {
	used, total, percent := rootDiskUsage()

	return SidebarStats{
		Version:     AppVersion(),
		DiskUsed:    humanBytes(used),
		DiskTotal:   humanBytes(total),
		DiskPercent: int(math.Round(percent * 100)),
	}
}

func rootDiskUsage() (uint64, uint64, float64) {
	var fs syscall.Statfs_t
	if err := syscall.Statfs("/", &fs); err != nil {
		return 0, 0, 0
	}

	total := fs.Blocks * uint64(fs.Bsize)
	available := fs.Bavail * uint64(fs.Bsize)
	used := total - available

	if total == 0 {
		return used, total, 0
	}

	return used, total, float64(used) / float64(total)
}

func humanBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	value := float64(bytes)
	idx := 0
	for value >= unit && idx < len(units)-1 {
		value /= unit
		idx++
	}

	return fmt.Sprintf("%.1f %s", value, units[idx])
}
