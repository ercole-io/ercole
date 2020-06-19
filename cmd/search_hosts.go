// Copyright (c) 2020 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var searchHostsHostname string
var searchHostsDatabase string
var searchHostsTechnology string
var searchHostsHardwareAbstractionTechnology string
var searchHostsCluster string
var searchHostsPhysicalHost string
var searchHostsOperatingSystem string
var searchHostsKernel string
var searchHostsLTEMemoryTotal float64
var searchHostsGTEMemoryTotal float64
var searchHostsLTESwapTotal float64
var searchHostsGTESwapTotal float64
var searchHostsIsMemberOfCluster string
var searchHostsCPUModel string
var searchHostsLTECPUCores int
var searchHostsGTECPUCores int
var searchHostsLTECPUThreads int
var searchHostsGTECPUThreads int

var searchHostsOption apiOption = apiOption{
	addOption: func(cmd *cobra.Command) {
		cmd.Flags().StringVar(&searchHostsHostname, "hostname", "", "Filter by hostname")
		cmd.Flags().StringVar(&searchHostsDatabase, "database", "", "Filter by a database")
		cmd.Flags().StringVar(&searchHostsTechnology, "technology", "", "Filter by a technology")
		cmd.Flags().StringVar(&searchHostsHardwareAbstractionTechnology, "hardware-abstraction-technology", "", "Filter by the HA technology")
		cmd.Flags().StringVar(&searchHostsCluster, "cluster", "", "Filter by the name of virtualization cluster. use 'NULL' for searching host that aren't part of a virtualization cluster")
		cmd.Flags().StringVar(&searchHostsPhysicalHost, "physical-host", "", "Filter by the name of physical host on which the host runs")
		cmd.Flags().StringVar(&searchHostsOperatingSystem, "operating-system", "", "Filter by the operating system")
		cmd.Flags().StringVar(&searchHostsKernel, "kernel", "", "Filter by the kernel")
		cmd.Flags().Float64Var(&searchHostsLTEMemoryTotal, "memory-total-lte", -1, "Filter the hosts with the Info.MemoryTotal less than memory-total-lte value")
		cmd.Flags().Float64Var(&searchHostsGTEMemoryTotal, "memory-total-gte", -1, "Filter the hosts with the Info.MemoryTotal less than memory-total-gte value")
		cmd.Flags().Float64Var(&searchHostsLTESwapTotal, "swap-total-lte", -1, "Filter the hosts with the Info.SwapTotal less than swap-total-lte value")
		cmd.Flags().Float64Var(&searchHostsGTESwapTotal, "swap-total-gte", -1, "Filter the hosts with the Info.SwapTotal less than swap-total-gte value")
		cmd.Flags().StringVar(&searchHostsIsMemberOfCluster, "is-member-of-cluster", "", "Filter the host by operating systems cluster membership")
		cmd.Flags().StringVar(&searchHostsCPUModel, "cpu-model", "", "Filter by CPU model")
		cmd.Flags().IntVar(&searchHostsLTECPUCores, "cpu-cores-lte", -1, "Filter the hosts with the Info.CPUCores less than cpu-cores-lte value")
		cmd.Flags().IntVar(&searchHostsGTECPUCores, "cpu-cores-gte", -1, "Filter the hosts with the Info.CPUCores less than cpu-cores-gte value")
		cmd.Flags().IntVar(&searchHostsLTECPUThreads, "cpu-threads-lte", -1, "Filter the hosts with the Info.CPUThreads less than cpu-threads-lte value")
		cmd.Flags().IntVar(&searchHostsGTECPUThreads, "cpu-threads-gte", -1, "Filter the hosts with the Info.CPUThreads less than cpu-threads-gte value")
	},
	addParam: func(params url.Values) {
		params.Set("hostname", searchHostsHostname)
		params.Set("database", searchHostsDatabase)
		params.Set("technology", searchHostsTechnology)
		params.Set("hardware-abstraction-technology", searchHostsHardwareAbstractionTechnology)
		params.Set("cluster", searchHostsCluster)
		params.Set("physical-host", searchHostsPhysicalHost)
		params.Set("operating-system", searchHostsOperatingSystem)
		params.Set("kernel", searchHostsKernel)
		params.Set("memory-total-lte", fmt.Sprintf("%f", searchHostsLTEMemoryTotal))
		params.Set("memory-total-gte", fmt.Sprintf("%f", searchHostsGTEMemoryTotal))
		params.Set("swap-total-lte", fmt.Sprintf("%f", searchHostsLTESwapTotal))
		params.Set("swap-total-gte", fmt.Sprintf("%f", searchHostsGTESwapTotal))
		params.Set("is-member-of-cluster", searchHostsIsMemberOfCluster)
		params.Set("cpu-model", searchHostsCPUModel)
		params.Set("cpu-cores-lte", fmt.Sprintf("%d", searchHostsLTECPUCores))
		params.Set("cpu-cores-gte", fmt.Sprintf("%d", searchHostsGTECPUCores))
		params.Set("cpu-threads-lte", fmt.Sprintf("%d", searchHostsLTECPUThreads))
		params.Set("cpu-threads-gte", fmt.Sprintf("%d", searchHostsGTECPUThreads))
	},
}

func init() {
	searchHostsCmd := simpleAPIRequestCommand("search-hosts",
		"Search current hosts",
		`search-hosts search the most matching hosts to the arguments`,
		true, []apiOption{searchHostsOption, modeOption, locationOption, environmentOption, sortingOptions, olderThanOptions}, true,
		"/hosts",
		"Failed to search hosts data: %v\n",
		"Failed to search hosts data(Status: %d): %s\n",
	)

	apiCmd.AddCommand(searchHostsCmd)
}
