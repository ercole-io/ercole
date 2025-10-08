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

package dto

import "github.com/ercole-io/ercole/v2/model"

type Database struct {
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	Version          string  `json:"version"`
	Hostname         string  `json:"hostname"`
	Environment      string  `json:"environment"`
	Location         string  `json:"location"`
	Charset          string  `json:"charset"`
	Memory           float64 `json:"memory"`       // in GB
	DatafileSize     float64 `json:"datafileSize"` // in GB
	SegmentsSize     float64 `json:"segmentSize"`  // in GB
	Archivelog       bool    `json:"archivelog"`
	HighAvailability bool    `json:"highAvailability"`
	DisasterRecovery bool    `json:"disasterRecovery"`
}

type DatabasesStatistics struct {
	TotalMemorySize   float64 `json:"total-memory-size"`   // in bytes
	TotalSegmentsSize float64 `json:"total-segments-size"` // in bytes
}

type DatabaseUsedLicense struct {
	Hostname        string  `json:"hostname" bson:"hostname"`
	DbName          string  `json:"dbName" bson:"dbName"`
	ClusterName     string  `json:"clusterName" bson:"clusterName"`
	ClusterType     string  `json:"clusterType" bson:"clusterType"`
	LicenseTypeID   string  `json:"licenseTypeID" bson:"licenseTypeID"`
	Description     string  `json:"description" bson:"description"`
	Metric          string  `json:"metric" bson:"metric"`
	UsedLicenses    float64 `json:"usedLicenses" bson:"usedLicenses"`
	ClusterLicenses float64 `json:"clusterLicenses" bson:"clusterLicenses"`
	Ignored         bool    `json:"ignored" bson:"ignored"`
	IgnoredComment  string  `json:"ignoredComment" bson:"ignoredComment"`
	OlvmCapped      bool    `json:"olvmCapped" bson:"olvmCapped"`
}

func (d *DatabaseUsedLicense) ToModel() model.LicenseUsedDatabase {
	return model.LicenseUsedDatabase{
		Hostname:        d.Hostname,
		DbName:          d.DbName,
		ClusterName:     d.ClusterName,
		ClusterType:     d.ClusterType,
		LicenseTypeID:   d.LicenseTypeID,
		Description:     d.Description,
		Metric:          d.Metric,
		UsedLicenses:    d.UsedLicenses,
		ClusterLicenses: d.ClusterLicenses,
		Ignored:         d.Ignored,
		IgnoredComment:  d.IgnoredComment,
		OlvmCapped:      d.OlvmCapped,
	}
}

type DatabaseUsedLicensePerHost struct {
	Hostname        string   `json:"hostname" bson:"hostname"`
	DatabaseNames   []string `json:"databaseNames" bson:"databaseNames"`
	LicenseTypeID   string   `json:"licenseTypeID" bson:"licenseTypeID"`
	Description     string   `json:"description" bson:"description"`
	Metric          string   `json:"metric" bson:"metric"`
	UsedLicenses    float64  `json:"usedLicenses" bson:"usedLicenses"`
	ClusterLicenses float64  `json:"clusterLicenses" bson:"clusterLicenses"`
	OlvmCapped      bool     `json:"olvmCapped" bson:"olvmCapped"`
}

func (d *DatabaseUsedLicensePerHost) ToModel() model.LicenseUsedHost {
	return model.LicenseUsedHost{
		Hostname:        d.Hostname,
		DatabaseNames:   d.DatabaseNames,
		LicenseTypeID:   d.LicenseTypeID,
		Description:     d.Description,
		Metric:          d.Metric,
		UsedLicenses:    d.UsedLicenses,
		ClusterLicenses: d.ClusterLicenses,
		OlvmCapped:      d.OlvmCapped,
	}
}

type DatabaseUsedLicensePerCluster struct {
	Cluster       string   `json:"cluster"`
	Hostnames     []string `json:"hostnames"`
	LicenseTypeID string   `json:"licenseTypeID"`
	Description   string   `json:"description"`
	Metric        string   `json:"metric"`
	UsedLicenses  float64  `json:"usedLicenses"`
}

func (d *DatabaseUsedLicensePerCluster) ToModel() model.LicenseUsedCluster {
	return model.LicenseUsedCluster{
		Cluster:       d.Cluster,
		Hostnames:     d.Hostnames,
		LicenseTypeID: d.LicenseTypeID,
		Description:   d.Description,
		Metric:        d.Metric,
		UsedLicenses:  d.UsedLicenses,
	}
}
