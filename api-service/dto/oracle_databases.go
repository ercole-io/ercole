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

import (
	"net/http"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OracleDatabaseResponse struct {
	Content  []OracleDatabase `json:"content" bson:"content"`
	Metadata PagingMetadata   `json:"metadata" bson:"metadata"`
}

// OracleDatabase holds information about an Oracle database.
type OracleDatabase struct {
	ID                primitive.ObjectID                      `json:"id" bson:"_id"`
	Hostname          string                                  `json:"hostname" bson:"hostname"`
	Environment       string                                  `json:"environment" bson:"environment"`
	Location          string                                  `json:"location" bson:"location"`
	InstanceNumber    int                                     `json:"instanceNumber" bson:"instanceNumber"`
	InstanceName      string                                  `json:"instanceName" bson:"instanceName"`
	Name              string                                  `json:"name" bson:"name"`
	UniqueName        string                                  `json:"uniqueName" bson:"uniqueName"`
	Status            string                                  `json:"status" bson:"status"`
	DbID              uint                                    `json:"dbID" bson:"dbID"`
	Role              string                                  `json:"role" bson:"role"`
	IsCDB             bool                                    `json:"isCDB" bson:"isCDB"`
	Version           string                                  `json:"version" bson:"version"`
	Platform          string                                  `json:"platform" bson:"platform"`
	Archivelog        bool                                    `json:"archivelog" bson:"archivelog"`
	Charset           string                                  `json:"charset" bson:"charset"`
	NCharset          string                                  `json:"nCharset" bson:"nCharset"`
	BlockSize         int                                     `json:"blockSize" bson:"blockSize"`
	CPUCount          int                                     `json:"cpuCount" bson:"cpuCount"`
	SGATarget         float64                                 `json:"sgaTarget" bson:"sgaTarget"`
	PGATarget         float64                                 `json:"pgaTarget" bson:"pgaTarget"`
	MemoryTarget      float64                                 `json:"memoryTarget" bson:"memoryTarget"`
	Memory            float64                                 `json:"memory" bson:"memory"`
	SGAMaxSize        float64                                 `json:"sgaMaxSize" bson:"sgaMaxSize"`
	SegmentsSize      float64                                 `json:"segmentsSize" bson:"segmentsSize"`
	DatafileSize      float64                                 `json:"datafileSize" bson:"datafileSize"`
	Allocable         float64                                 `json:"allocable" bson:"allocable"`
	Elapsed           *float64                                `json:"elapsed" bson:"elapsed"`
	DBTime            *float64                                `json:"dbTime" bson:"dbTime"`
	DailyCPUUsage     *float64                                `json:"dailyCPUUsage" bson:"dailyCPUUsage"`
	Work              *float64                                `json:"work" bson:"work"`
	ASM               bool                                    `json:"asm" bson:"asm"`
	Dataguard         bool                                    `json:"dataguard" bson:"dataguard"`
	Rac               bool                                    `json:"rac" bson:"rac"`
	Ha                bool                                    `json:"ha" bson:"ha"`
	Tags              []string                                `json:"tags" bson:"tags"`
	CreatedAt         time.Time                               `json:"createdAt" bson:"createdAt"`
	Patches           []model.OracleDatabasePatch             `json:"patches" bson:"patches"`
	Tablespaces       []model.OracleDatabaseTablespace        `json:"tablespaces" bson:"tablespaces"`
	Schemas           []model.OracleDatabaseSchema            `json:"schemas" bson:"schemas"`
	Licenses          []model.OracleDatabaseLicense           `json:"licenses" bson:"licenses"`
	ADDMs             []model.OracleDatabaseAddm              `json:"addms" bson:"addms"`
	SegmentAdvisors   []model.OracleDatabaseSegmentAdvisor    `json:"segmentAdvisors" bson:"segmentAdvisors"`
	PSUs              []model.OracleDatabasePSU               `json:"psus" bson:"psus"`
	Backups           []model.OracleDatabaseBackup            `json:"backups" bson:"backups"`
	FeatureUsageStats []model.OracleDatabaseFeatureUsageStat  `json:"featureUsageStats" bson:"featureUsageStats"`
	PDBs              []model.OracleDatabasePluggableDatabase `json:"pdbs" bson:"pdbs"`
	Services          []model.OracleDatabaseService           `json:"services" bson:"services"`
	Changes           []model.Changes                         `json:"changes" bson:"changes"`
	OtherInfo         map[string]interface{}                  `json:"-" bson:"-"`
}

type SearchOracleDatabasesFilter struct {
	GlobalFilter

	Search     string
	SortBy     string
	SortDesc   bool
	PageNumber int
	PageSize   int
}

func GetSearchOracleDatabasesFilter(r *http.Request) (f *SearchOracleDatabasesFilter, err error) {
	f = new(SearchOracleDatabasesFilter)

	gf, err := GetGlobalFilter(r)
	if err != nil {
		return nil, err
	}

	f.GlobalFilter = *gf

	f.Search = r.URL.Query().Get("search")
	f.SortBy = r.URL.Query().Get("sort-by")

	if f.SortDesc, err = utils.Str2bool(r.URL.Query().Get("sort-desc"), false); err != nil {
		return nil, err
	}

	if f.PageNumber, err = utils.Str2int(r.URL.Query().Get("page"), -1); err != nil {
		return nil, err
	}

	if f.PageSize, err = utils.Str2int(r.URL.Query().Get("size"), -1); err != nil {
		return nil, err
	}

	return
}

type OracleDatabasesStatistics struct {
	TotalMemorySize   float64 `json:"total-memory-size"`   // in bytes
	TotalSegmentsSize float64 `json:"total-segments-size"` // in bytes
	TotalDatafileSize float64 `json:"total-datafile-size"` // in bytes
	TotalWork         float64 `json:"total-work"`
}

type OracleDatabaseSegmentAdvisor struct {
	SegmentOwner   string    `json:"segmentOwner"`
	SegmentName    string    `json:"segmentName"`
	SegmentType    string    `json:"segmentType"`
	SegmentsSize   float64   `json:"segmentsSize"`
	PartitionName  string    `json:"partitionName"`
	Reclaimable    float64   `json:"reclaimable"`
	Retrieve       float64   `json:"retrieve"`
	Recommendation string    `json:"recommendation"`
	CreatedAt      time.Time `json:"createdAt"`
	Dbname         string    `json:"dbname"`
	Environment    string    `json:"environment"`
	Hostname       string    `json:"hostname"`
	Location       string    `json:"location"`
}
