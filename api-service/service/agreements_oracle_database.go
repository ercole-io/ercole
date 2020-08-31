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

package service

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/api-service/database"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoadOracleDatabaseAgreementPartsList loads the list of Oracle/Database agreement parts and store it to as.OracleDatabaseAgreementParts.
func (as *APIService) LoadOracleDatabaseAgreementPartsList() {
	// read the list content
	listContentRaw, err := ioutil.ReadFile(as.Config.ResourceFilePath + "/oracle_database_agreement_parts_list.json")
	if err != nil {
		as.Log.Warnf("Unable to read %s: %v\n", as.Config.ResourceFilePath+"/oracle_database_agreement_parts_list.json", err)
		return
	}

	// unmarshal to OracleDatabaseAgreementParts
	err = json.Unmarshal(listContentRaw, &as.OracleDatabaseAgreementParts)
	if err != nil {
		as.Log.Warnf("Unable to unmarshal %s: %v\n", as.Config.ResourceFilePath+"/oracle_database_agreement_parts_list.json", err)
		return
	}
}

// GetOracleDatabaseAgreementPartsList return the list of Oracle/Database agreement parts
func (as *APIService) GetOracleDatabaseAgreementPartsList() ([]model.OracleDatabaseAgreementPart, utils.AdvancedErrorInterface) {
	return as.OracleDatabaseAgreementParts, nil
}

// AddOracleDatabaseAgreements return the list of Oracle/Database agreement parts
func (as *APIService) AddOracleDatabaseAgreements(req apimodel.OracleDatabaseAgreementsAddRequest) (interface{}, utils.AdvancedErrorInterface) {
	//Check and resolve every part id
	var parts []*model.OracleDatabaseAgreementPart = make([]*model.OracleDatabaseAgreementPart, len(req.PartsID))
	for i, pid := range req.PartsID {
		found := false
		for j, vpid := range as.OracleDatabaseAgreementParts {
			if pid == vpid.PartID {
				found = true
				parts[i] = &as.OracleDatabaseAgreementParts[j]
				break
			}
		}
		if !found {
			return nil, utils.AerrOracleDatabaseAgreementInvalidPartID
		}
	}

	//Get the list of hosts not in cluster
	notInClusterHosts, aerr := as.SearchHosts("hostnames", "", database.SearchHostsFilters{
		GTECPUCores:    -1,
		LTECPUCores:    -1,
		LTECPUThreads:  -1,
		LTEMemoryTotal: -1,
		GTECPUThreads:  -1,
		GTESwapTotal:   -1,
		GTEMemoryTotal: -1,
		LTESwapTotal:   -1,
	}, "", false, -1, -1, "", "", utils.MAX_TIME)
	if aerr != nil {
		return nil, aerr
	}
	notInClusterHostnames := make([]string, len(notInClusterHosts))
	for i, h := range notInClusterHosts {
		notInClusterHostnames[i] = h["hostname"].(string)
	}

	//Check every host in req.Hosts
	for _, host := range req.Hosts {
		found := false
		for _, vhost := range notInClusterHostnames {
			if host == vhost {
				found = true
				break
			}
		}
		if !found {
			return nil, utils.AerrHostNotFound
		}
	}

	//expode req in multple agreement
	aggs := make([]model.OracleDatabaseAgreement, len(req.PartsID))
	for i, part := range parts {
		aggs[i].AgreementID = req.AgreementID
		aggs[i].CSI = req.CSI
		aggs[i].CatchAll = req.CatchAll
		aggs[i].Count = req.Count
		aggs[i].Hosts = req.Hosts
		aggs[i].ID = primitive.NewObjectIDFromTimestamp(as.TimeNow())
		aggs[i].ItemDescription = part.ItemDescription
		aggs[i].Metrics = part.Metrics
		aggs[i].PartID = part.PartID
		aggs[i].ReferenceNumber = req.ReferenceNumber
		aggs[i].Unlimited = req.Unlimited
	}

	//insert it to the database
	res := make([]interface{}, len(aggs))
	for i, agg := range aggs {
		if res[i], aerr = as.Database.InsertOracleDatabaseAgreement(agg); aerr != nil {
			return nil, aerr
		}
	}

	return res, nil
}
