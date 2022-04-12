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
)

type OracleDatabaseResponse struct {
	Content  []OracleDatabase `json:"content" bson:"content"`
	Metadata PagingMetadata   `json:"metadata" bson:"metadata"`
}

// OracleDatabase holds information about an Oracle database.
type OracleDatabase struct {
	Hostname     string                        `json:"hostname" bson:"hostname"`
	Environment  string                        `json:"environment" bson:"environment"`
	Location     string                        `json:"location" bson:"location"`
	Name         string                        `json:"name" bson:"name"`
	UniqueName   string                        `json:"uniqueName" bson:"uniqueName"`
	Status       string                        `json:"status" bson:"status"`
	IsCDB        bool                          `json:"isCDB" bson:"isCDB"`
	Version      string                        `json:"version" bson:"version"`
	Archivelog   bool                          `json:"archivelog" bson:"archivelog"`
	Charset      string                        `json:"charset" bson:"charset"`
	BlockSize    int                           `json:"blockSize" bson:"blockSize"`
	CPUCount     int                           `json:"cpuCount" bson:"cpuCount"`
	MemoryTarget float64                       `json:"memoryTarget" bson:"memoryTarget"`
	Memory       float64                       `json:"memory" bson:"memory"`
	SegmentsSize float64                       `json:"segmentsSize" bson:"segmentsSize"`
	DatafileSize float64                       `json:"datafileSize" bson:"datafileSize"`
	Work         *float64                      `json:"work" bson:"work"`
	Dataguard    bool                          `json:"dataguard" bson:"dataguard"`
	Rac          bool                          `json:"rac" bson:"rac"`
	Ha           bool                          `json:"ha" bson:"ha"`
	DbID         uint                          `json:"dbID" bson:"dbID"`
	Role         string                        `json:"role" bson:"role"`
	Services     []model.OracleDatabaseService `json:"services" bson:"services"`
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
