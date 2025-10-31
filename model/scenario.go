// Copyright (c) 2025 Sorint.lab S.p.A.
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
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Scenario struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time          `bson:"createdAt"`

	Location string          `bson:"location"`
	Hosts    []SimulatedHost `bson:"hosts"`

	LicenseCompliance LicensesComplianceScenario `bson:"licenseCompliance"`
	LicenseUsed       LicenseUsedScenario        `bson:"licenseUsed"`
}

type SimulatedHost struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	Host      HostDataBE         `bson:"host"`
	Core      int                `bson:"core"`
}

type LicensesComplianceScenario struct {
	Actual []LicenseCompliance `bson:"actual"`
	Got    []LicenseCompliance `bson:"got"`
}

type LicenseCompliance struct {
	LicenseTypeID   string  `bson:"licenseTypeID"`
	ItemDescription string  `bson:"itemDescription"`
	Metric          string  `bson:"metric"`
	Cost            float64 `bson:"cost"`
	Consumed        float64 `bson:"consumed"`
	Covered         float64 `bson:"covered"`
	Purchased       float64 `bson:"purchased"`
	Compliance      float64 `bson:"compliance"`
	Unlimited       bool    `bson:"unlimited"`
	Available       float64 `bson:"available"`
}

type LicenseUsedScenario struct {
	LicenseDatabase          LicenseUsedPerDatabase          `bson:"licenseDatabase"`
	LicenseHost              LicenseUsedPerHost              `bson:"licenseHost"`
	LicenseHypervisorCluster LicenseUsedPerHypervisorCluster `bson:"licenseHypervisorCluster"`
	LicenseClusterVeritas    LicenseUsedPerClusterVeritas    `bson:"licenseClusterVeritas"`
}

type LicenseUsedPerDatabase struct {
	Actual []LicenseUsedDatabase `bson:"actual"`
	Got    []LicenseUsedDatabase `bson:"got"`
}

type LicenseUsedPerHost struct {
	Actual []LicenseUsedHost `bson:"actual"`
	Got    []LicenseUsedHost `bson:"got"`
}

type LicenseUsedPerHypervisorCluster struct {
	Actual []LicenseUsedCluster `bson:"actual"`
	Got    []LicenseUsedCluster `bson:"got"`
}

type LicenseUsedPerClusterVeritas struct {
	Actual []LicenseUsedClusterVeritas `bson:"actual"`
	Got    []LicenseUsedClusterVeritas `bson:"got"`
}

type LicenseUsedDatabase struct {
	Hostname        string  `bson:"hostname"`
	DbName          string  `bson:"dbName"`
	ClusterName     string  `bson:"clusterName"`
	ClusterType     string  `bson:"clusterType"`
	LicenseTypeID   string  `bson:"licenseTypeID"`
	Description     string  `bson:"description"`
	Metric          string  `bson:"metric"`
	UsedLicenses    float64 `bson:"usedLicenses"`
	ClusterLicenses float64 `bson:"clusterLicenses"`
	Ignored         bool    `bson:"ignored"`
	IgnoredComment  string  `bson:"ignoredComment"`
	OlvmCapped      bool    `bson:"olvmCapped"`
}

type LicenseUsedHost struct {
	Hostname        string   `bson:"hostname"`
	DatabaseNames   []string `bson:"databaseNames"`
	LicenseTypeID   string   `bson:"licenseTypeID"`
	Description     string   `bson:"description"`
	Metric          string   `bson:"metric"`
	UsedLicenses    float64  `bson:"usedLicenses"`
	ClusterLicenses float64  `bson:"clusterLicenses"`
	OlvmCapped      bool     `bson:"olvmCapped"`
}

type LicenseUsedCluster struct {
	Cluster       string   `bson:"cluster"`
	Hostnames     []string `bson:"hostnames"`
	LicenseTypeID string   `bson:"licenseTypeID"`
	Description   string   `bson:"description"`
	Metric        string   `bson:"metric"`
	UsedLicenses  float64  `bson:"usedLicenses"`
}

type LicenseUsedClusterVeritas struct {
	ID            string   `bson:"id"`
	Hostnames     []string `bson:"hostnames"`
	LicenseTypeID string   `bson:"licenseTypeID"`
	Description   string   `bson:"description"`
	Metric        string   `bson:"metric"`
	Count         float64  `bson:"count"`
}
