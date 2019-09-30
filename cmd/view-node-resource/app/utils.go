package app

import (
	"fmt"
	"strconv"

	"github.com/gosoon/kubectl-plugin/pkg/types"
	"github.com/gosoon/kubectl-plugin/pkg/utils"

	v1 "k8s.io/api/core/v1"
)

func getNodeAllocatable(allocatable v1.ResourceList) (float64, float64) {
	nodeCPU := float64(0)
	nodeMemory := float64(0)
	for name, value := range allocatable {
		if string(name) == "cpu" {
			cpu, _ := strconv.ParseFloat(value.String(), 64)
			nodeCPU = cpu
		} else if string(name) == "memory" {
			memory, _ := utils.ConvertMemoryUnit(value.String())
			nodeMemory += memory
		}
	}
	return nodeCPU, nodeMemory
}

func pickNodeCPURequests(node *types.NodeResourceList) string {
	return fmt.Sprintf("%.2f (%v)", node.CPURequests, node.CPURequestsUsage)
}

func pickNodeMemoryRequests(node *types.NodeResourceList) string {
	return fmt.Sprintf("%.1f (%v)", node.MemoryRequests, node.MemoryRequestsUsage)
}

func pickNodeCPULimits(node *types.NodeResourceList) string {
	return fmt.Sprintf("%.2f (%v)", node.CPULimits, node.CPULimitsUsage)
}

func pickNodeMemoryLimits(node *types.NodeResourceList) string {
	return fmt.Sprintf("%.1f (%v)", node.MemoryLimits, node.MemoryLimitsUsage)
}
