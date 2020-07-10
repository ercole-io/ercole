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
		VeritasClusterServer: true,
		SunCluster:           false,
		HACMP:                true,
		OtherInfo: map[string]interface{}{
			"pippo":    "pluto",
			"topolino": float64(42),
		},
	}

	//out, err := cms.MarshalJSON()

	//if err != nil {
	//	t.Fail()
	//}

	//fmt.Println(string(out))

	//var newCMS ClusterMembershipStatus
	//newCMS.UnmarshalJSON(out)

	//assert.Equal(t, cms, newCMS)

	// BSON /////////////////////////////
	out, err := cms.MarshalBSON()

	if err != nil {
		t.Fail()
	}

	fmt.Println(string(out))

	newCMS2 := new(ClusterMembershipStatus)

	err = newCMS2.UnmarshalBSON(out)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, cms, *newCMS2)
}

func TestLicenseCount(t *testing.T) {
	l := LicenseCount{
		Name:             "pippo",
		CostPerProcessor: 42.42,
		Count:            12,
		Unlimited:        true,
	}

	out, err := bson.Marshal(l)

	if err != nil {
		t.Fail()
	}

	fmt.Println("#####################", string(out), "##########################")

	l2 := new(LicenseCount)
	err = bson.Unmarshal(out, l2)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	assert.Equal(t, l, *l2)
}
