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
	"time"

	"github.com/gosoon/kubectl-plugin/pkg/kubeclient"
	"github.com/gosoon/kubectl-plugin/pkg/printers"
	"github.com/gosoon/kubectl-plugin/pkg/types"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "view-node-taints",
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
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.view-node-taints.yaml)")
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
		// Search config in home directory with name ".view-node-taints" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".view-node-taints")
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
	taintsColumnDefinitions := getTaintsColumnDefinitions(nodeList)
	printTaintsColumnDefinitions(taintsColumnDefinitions)
}

func getTaintsColumnDefinitions(nodeList *v1.NodeList) []types.TaintsColumnDefinitions {
	var taintsColumnDefinitions []types.TaintsColumnDefinitions
	for _, node := range nodeList.Items {
		status := NodeNotReady
		for _, condition := range node.Status.Conditions {
			if condition.Type == NodeReady && condition.Status == "True" {
				status = NodeReady
			}
		}
		taintsList := "<none>"
		if len(node.Spec.Taints) != 0 {
			taintsList = visitTaints(node.Spec.Taints)
		}

		if node.Spec.Unschedulable {
			status += "," + NodeSchedulingDisabled
		}

		nodeAge := convertToAge(node.CreationTimestamp)
		taintsColumnDefinitions = append(taintsColumnDefinitions, types.TaintsColumnDefinitions{
			Name:    node.Name,
			Status:  status,
			Age:     nodeAge,
			Version: node.Status.NodeInfo.KubeletVersion,
			Taints:  taintsList,
		})
	}
	return taintsColumnDefinitions
}

func printTaintsColumnDefinitions(taintsColumnDefinitions []types.TaintsColumnDefinitions) {
	printers.Output(taintsColumnDefinitions)
}

func visitTaints(taints []v1.Taint) string {
	var taintsStr string
	var ct int
	taintsLen := len(taints)
	for _, taint := range taints {
		taintsStr += fmt.Sprintf("%v=%v:%v", taint.Key, taint.Value, taint.Effect)
		ct += 1
		if ct < taintsLen {
			taintsStr += ","
		}
	}
	return taintsStr
}

func convertToAge(creationTimestamp metav1.Time) string {
	s := creationTimestamp.Time.Format("2006-01-02T15:04:05Z")
	t, _ := time.Parse("2006-01-02T15:04:05Z", s)

	now := time.Now()
	sub := now.Sub(t)
	subHours := sub.Hours()

	if subHours < float64(24) {
		return fmt.Sprintf("%vh", subHours)
	}

	hours := int(subHours)
	days := hours / 24
	if (hours % 24) > 0 {
		days += 1
	}
	return fmt.Sprintf("%vd", days)
}
