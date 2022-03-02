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

import (
	"strings"

	"github.com/ercole-io/ercole/v2/utils"
)

// OracleDatabase holds information about an Oracle database.
type OracleDatabase struct {
	InstanceNumber    int                               `json:"instanceNumber" bson:"instanceNumber"`
	InstanceName      string                            `json:"instanceName" bson:"instanceName"`
	Name              string                            `json:"name" bson:"name"`
	UniqueName        string                            `json:"uniqueName" bson:"uniqueName"`
	Status            string                            `json:"status" bson:"status"`
	DbID              uint                              `json:"dbID" bson:"dbID"`
	Role              string                            `json:"role" bson:"role"`
	IsCDB             bool                              `json:"isCDB" bson:"isCDB"`
	Version           string                            `json:"version" bson:"version"`
	Platform          string                            `json:"platform" bson:"platform"`
	Archivelog        bool                              `json:"archivelog" bson:"archivelog"`
	Charset           string                            `json:"charset" bson:"charset"`
	NCharset          string                            `json:"nCharset" bson:"nCharset"`
	BlockSize         int                               `json:"blockSize" bson:"blockSize"`
	CPUCount          int                               `json:"cpuCount" bson:"cpuCount"`
	SGATarget         float64                           `json:"sgaTarget" bson:"sgaTarget"`
	PGATarget         float64                           `json:"pgaTarget" bson:"pgaTarget"`
	MemoryTarget      float64                           `json:"memoryTarget" bson:"memoryTarget"`
	SGAMaxSize        float64                           `json:"sgaMaxSize" bson:"sgaMaxSize"`
	SegmentsSize      float64                           `json:"segmentsSize" bson:"segmentsSize"`
	DatafileSize      float64                           `json:"datafileSize" bson:"datafileSize"`
	Allocable         float64                           `json:"allocable" bson:"allocable"`
	Elapsed           *float64                          `json:"elapsed" bson:"elapsed"`
	DBTime            *float64                          `json:"dbTime" bson:"dbTime"`
	DailyCPUUsage     *float64                          `json:"dailyCPUUsage" bson:"dailyCPUUsage"`
	Work              *float64                          `json:"work" bson:"work"`
	ASM               bool                              `json:"asm" bson:"asm"`
	Dataguard         bool                              `json:"dataguard" bson:"dataguard"`
	Patches           []OracleDatabasePatch             `json:"patches" bson:"patches"`
	Tablespaces       []OracleDatabaseTablespace        `json:"tablespaces" bson:"tablespaces"`
	Schemas           []OracleDatabaseSchema            `json:"schemas" bson:"schemas"`
	Licenses          []OracleDatabaseLicense           `json:"licenses" bson:"licenses"`
	ADDMs             []OracleDatabaseAddm              `json:"addms" bson:"addms"`
	SegmentAdvisors   []OracleDatabaseSegmentAdvisor    `json:"segmentAdvisors" bson:"segmentAdvisors"`
	PSUs              []OracleDatabasePSU               `json:"psus" bson:"psus"`
	Backups           []OracleDatabaseBackup            `json:"backups" bson:"backups"`
	FeatureUsageStats []OracleDatabaseFeatureUsageStat  `json:"featureUsageStats" bson:"featureUsageStats"`
	PDBs              []OracleDatabasePluggableDatabase `json:"pdbs" bson:"pdbs"`
	Services          []OracleDatabaseService           `json:"services" bson:"services"`
	Changes           []Changes                         `json:"changes" bson:"changes"`
}

var (
	OracleDatabaseStatusOpen    = "OPEN"
	OracleDatabaseStatusMounted = "MOUNTED"
)

var (
	OracleDatabaseRolePrimary         = "PRIMARY"
	OracleDatabaseRoleLogicalStandby  = "LOGICAL STANDBY"
	OracleDatabaseRolePhysicalStandby = "PHYSICAL STANDBY"
	OracleDatabaseRoleSnapshotStandby = "SNAPSHOT STANDBY"
)

var (
	OracleDatabaseEditionEnterprise = "ENT"
	OracleDatabaseEditionExtreme    = "EXE"
	OracleDatabaseEditionStandard   = "STD"
	OracleDatabaseEditionExpress    = "XE"
)

func (v OracleDatabase) Edition() (dbEdition string) {
	if strings.Contains(strings.ToUpper(v.Version), "ENTERPRISE") {
		dbEdition = OracleDatabaseEditionEnterprise
	} else if strings.Contains(strings.ToUpper(v.Version), "EXTREME") {
		dbEdition = OracleDatabaseEditionExtreme
	} else if strings.Contains(strings.ToUpper(v.Version), "EXPRESS") {
		dbEdition = OracleDatabaseEditionExpress
	} else {
		dbEdition = OracleDatabaseEditionStandard
	}

	return
}

func (v OracleDatabase) CoreFactor(host Host, hostCoreFactor float64) (float64, error) {
	dbEdition := v.Edition()

	if host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyOvm ||
		host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyVmware ||
		host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyVmother ||
		host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyKvm {
		if dbEdition == OracleDatabaseEditionExtreme || dbEdition == OracleDatabaseEditionEnterprise {
			return hostCoreFactor, nil
		}

		if dbEdition == OracleDatabaseEditionStandard {
			return 0, nil
		}

		return 0, utils.NewErrorf("%q db: dbEdition %q unknown", v.Name, dbEdition)
	}

	if host.HardwareAbstractionTechnology == HardwareAbstractionTechnologyPhysical {
		if dbEdition == OracleDatabaseEditionExtreme || dbEdition == OracleDatabaseEditionEnterprise {
			return 0.5, nil
		} else if dbEdition == OracleDatabaseEditionStandard {
			return float64(host.CPUSockets), nil
		}

		return 0, utils.NewErrorf("%q db: dbEdition %q unknown", v.Name, dbEdition)
	}

	return 0, utils.NewErrorf("%q db: hardwareAbstractionTechnology %q unknown",
		v.Name, host.HardwareAbstractionTechnology)
}

// DatabaseSliceAsMap return the equivalent map of the database slice with Database.Name as Key
func DatabaseSliceAsMap(dbs []OracleDatabase) map[string]OracleDatabase {
	out := make(map[string]OracleDatabase)
	for _, db := range dbs {
		out[db.Name] = db
	}

	return out
}
