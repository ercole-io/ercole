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

package model

const (
	HardwareAbstractionPhysical string = "PH"
	HardwareAbstractionVirtual  string = "VIRT"
)

const (
	HardwareAbstractionTechnologyPhysical string = "PH"
	HardwareAbstractionTechnologyOvm      string = "OVM"
	HardwareAbstractionTechnologyVmware   string = "VMWARE"
	HardwareAbstractionTechnologyHyperv   string = "HYPERV"
	HardwareAbstractionTechnologyVmother  string = "VMOTHER"
	HardwareAbstractionTechnologyXen      string = "XEN"
	HardwareAbstractionTechnologyHpvirt   string = "HPVIRT"
	HardwareAbstractionTechnologyKvm      string = "KVM"
)

// Host contains info about the host
type Host struct {
	Hostname                      string  `json:"hostname" bson:"hostname"`
	CPUModel                      string  `json:"cpuModel" bson:"cpuModel"`
	CPUFrequency                  string  `json:"cpuFrequency" bson:"cpuFrequency"`
	CPUSockets                    int     `json:"cpuSockets" bson:"cpuSockets"`
	CPUCores                      int     `json:"cpuCores" bson:"cpuCores"`
	CPUThreads                    int     `json:"cpuThreads" bson:"cpuThreads"`
	ThreadsPerCore                int     `json:"threadsPerCore" bson:"threadsPerCore"`
	CoresPerSocket                int     `json:"coresPerSocket" bson:"coresPerSocket"`
	HardwareAbstraction           string  `json:"hardwareAbstraction" bson:"hardwareAbstraction"`
	HardwareAbstractionTechnology string  `json:"hardwareAbstractionTechnology" bson:"hardwareAbstractionTechnology"`
	Kernel                        string  `json:"kernel" bson:"kernel"`
	KernelVersion                 string  `json:"kernelVersion" bson:"kernelVersion"`
	OS                            string  `json:"os" bson:"os"`
	OSVersion                     string  `json:"osVersion" bson:"osVersion"`
	MemoryTotal                   float64 `json:"memoryTotal" bson:"memoryTotal"`
	SwapTotal                     float64 `json:"swapTotal" bson:"swapTotal"`
}
