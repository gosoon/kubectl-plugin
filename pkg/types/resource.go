package types

// NodeResourceList xxx
type NodeResourceList struct {
	Name     string
	PodCount int

	CPU    float64
	Memory float64

	// request resource
	CPURequests    float64
	MemoryRequests float64

	// request resource usage
	CPURequestsUsage    string
	MemoryRequestsUsage string

	// limit resource
	CPULimits    float64
	MemoryLimits float64

	// limit resource usage
	CPULimitsUsage    string
	MemoryLimitsUsage string
}

// NodeResourceColumnDefinitions xxx
type NodeResourceColumnDefinitions struct {
	Name     string
	PodCount int

	// requests resource
	CPURequests    string
	MemoryRequests string

	// limits resource
	CPULimits    string
	MemoryLimits string
}
