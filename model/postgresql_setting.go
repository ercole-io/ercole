// Copyright (c) 2022 Sorint.lab S.p.A.
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

type PostgreSQLSetting struct {
	DbVersion                  string  `json:"dbVersion" bson:"dbVersion"`
	DataDirectory              string  `json:"dataDirectory" bson:"dataDirectory"`
	WorkMem                    int     `json:"workMem" bson:"workMem"`
	ArchiveMode                bool    `json:"archiveMode" bson:"archiveMode"`
	ArchiveCommand             string  `json:"archiveCommand" bson:"archiveCommand"`
	MinWalSize                 int     `json:"minWalSize" bson:"minWalSize"`
	MaxWalSize                 int     `json:"maxWalSize" bson:"maxWalSize"`
	MaxConnections             int     `json:"maxConnections" bson:"maxConnections"`
	CheckpointCompletionTarget string  `json:"checkpointCompletionTarget" bson:"checkpointCompletionTarget"`
	DefaultStatisticsTarget    int     `json:"defaultStatisticsTarget" bson:"defaultStatisticsTarget"`
	RandomPageCost             float64 `json:"randomPageCost" bson:"randomPageCost"`
	MaintenanceWorkMem         int     `json:"maintenanceWorkMem" bson:"maintenanceWorkMem"`
	SharedBuffers              int     `json:"sharedBuffers" bson:"sharedBuffers"`
	EffectiveCacheSize         int     `json:"effectiveCacheSize" bson:"effectiveCacheSize"`
	EffectiveIoConcurrency     int     `json:"effectiveIoConcurrency" bson:"effectiveIoConcurrency"`
	MaxWorkerProcesses         int     `json:"maxWorkerProcesses" bson:"maxWorkerProcesses"`
	MaxParallelWorkers         int     `json:"maxParallelWorkers" bson:"maxParallelWorkers"`
}
