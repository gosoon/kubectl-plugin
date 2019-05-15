// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"os"

	"github.com/gosoon/kubectl-plugin/pkg/kubeclient"
	"github.com/gosoon/kubectl-plugin/pkg/printers"
	"github.com/gosoon/kubectl-plugin/pkg/types"
	"github.com/gosoon/kubectl-plugin/pkg/utils"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/api/core/v1"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "view-node-resource",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		client, err := kubeclient.NewClient()
		if err != nil {
			panic(err)
		}
		run(client)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.view-node-resource.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".view-node-resource" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".view-node-resource")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(client *kubeclient.Client) {
	nodeList, err := client.ListNode()
	if err != nil {
		panic(err)
	}

	podList, err := client.ListPod()
	if err != nil {
		panic(err)
	}
	nodeResourceList := NodeResouceHandler(nodeList, podList)

	printNodeResourceColumnDefinitions(nodeResourceList)
}

func printNodeResourceColumnDefinitions(nodeResourceList map[string]*types.NodeResourceList) {
	var nodeResourceColumnDefinitions []types.NodeResourceColumnDefinitions
	for _, node := range nodeResourceList {
		nodeResourceColumnDefinitions = append(nodeResourceColumnDefinitions,
			types.NodeResourceColumnDefinitions{
				Name:           node.Name,
				PodCount:       node.PodCount,
				CPURequests:    pickNodeCPURequests(node),
				MemoryRequests: pickNodeMemoryRequests(node),
				CPULimits:      pickNodeCPULimits(node),
				MemoryLimits:   pickNodeMemoryLimits(node),
			})
	}

	printers.Output(nodeResourceColumnDefinitions)
}

func NodeResouceHandler(nodeList *v1.NodeList, podList *v1.PodList) map[string]*types.NodeResourceList {
	nodeResourceList := make(map[string]*types.NodeResourceList)

	for _, node := range nodeList.Items {
		if _, existed := nodeResourceList[node.Name]; !existed {
			nodeCPU, nodeMemory := getNodeAllocatable(node.Status.Allocatable)
			nodeResourceList[node.Name] = &types.NodeResourceList{
				Name:   node.Name,
				CPU:    nodeCPU,
				Memory: nodeMemory,
			}
		}
	}

	for _, pod := range podList.Items {
		nodeName := pod.Spec.NodeName
		if _, existed := nodeResourceList[nodeName]; existed {
			for _, container := range pod.Spec.Containers {
				for name, value := range container.Resources.Requests {
					if string(name) == "cpu" {
						cr, _ := utils.ConvertCPUUnit(value.String())
						nodeResourceList[nodeName].CPURequests += cr
					}
					if string(name) == "memory" {
						mr, _ := utils.ConvertMemoryUnit(value.String())
						nodeResourceList[pod.Spec.NodeName].MemoryRequests += mr
					}
				}
				for name, value := range container.Resources.Limits {
					if string(name) == "cpu" {
						cl, _ := utils.ConvertCPUUnit(value.String())
						nodeResourceList[nodeName].CPULimits += cl
					}
					if string(name) == "memory" {
						ml, _ := utils.ConvertMemoryUnit(value.String())
						nodeResourceList[nodeName].MemoryLimits += ml
					}
				}
				nodeResourceList[nodeName].PodCount++
			}
		}
	}

	for _, node := range nodeList.Items {
		nodeName := node.Name

		nodeResourceList[nodeName].CPURequestsUsage = fmt.Sprintf("%.2f%%",
			nodeResourceList[nodeName].CPURequests/nodeResourceList[nodeName].CPU*100)

		nodeResourceList[nodeName].MemoryRequestsUsage = fmt.Sprintf("%.2f%%",
			nodeResourceList[nodeName].MemoryRequests/nodeResourceList[nodeName].Memory*100)

		nodeResourceList[nodeName].CPULimitsUsage = fmt.Sprintf("%.2f%%",
			nodeResourceList[nodeName].CPULimits/nodeResourceList[nodeName].CPU*100)

		nodeResourceList[nodeName].MemoryLimitsUsage = fmt.Sprintf("%.2f%%",
			nodeResourceList[nodeName].MemoryLimits/nodeResourceList[nodeName].Memory*100)
	}
	return nodeResourceList
}
