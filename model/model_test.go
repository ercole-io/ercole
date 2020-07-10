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

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestClusterMembershipStatus(t *testing.T) {

	cms := ClusterMembershipStatus{
		OracleClusterware:    false,
		VeritasClusterServer: false,
		SunCluster:           false,
		HACMP:                false,
		OtherInfo: map[string]interface{}{
			"pippo":    "pluto",
			"topolino": float64(42),
		},
	}

	out, err := cms.MarshalJSON()

	if err != nil {
		t.Fail()
	}

	fmt.Println(string(out))

	var newCMS ClusterMembershipStatus
	newCMS.UnmarshalJSON(out)

	assert.Equal(t, cms, newCMS)

	//out, err := cms.MarshalBSON()

	//if err != nil {

	//	t.Fail()
	//}

	//fmt.Println(string(out))

	//var newCMS ClusterMembershipStatus

	//newCMS.UnmarshalBSON(out)

	//assert.Equal(t, cms, newCMS)
}

func TestLicenseCount(t *testing.T) {
	l := LicenseCount{
		Name:             "pippo",
		CostPerProcessor: 42.42,
		Count:            12,
		Unlimited:        false,
	}

	out, err := bson.Marshal(l)

	if err != nil {
		t.Fail()
	}

	fmt.Println("#####################", string(out), "##########################")

	var l2 LicenseCount
	bson.Unmarshal(out, l2)

	assert.Equal(t, l, l2)
}
