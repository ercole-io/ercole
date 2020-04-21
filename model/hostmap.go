package model

// HostMap contains info about the host
type HostMap map[string]interface{}

// CPUCores getter
func (host *HostMap) CPUCores() int {
	cpuCores := (*host)["CPUCores"]
	return cpuCores.(int)
}
