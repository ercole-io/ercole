package service

import (
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
)

// GetDefaultDatabaseTags return the default list of database tags from configuration
func (as *APIService) GetDefaultDatabaseTags() ([]string, utils.AdvancedErrorInterface) {
	return as.Config.APIService.DefaultDatabaseTags, nil
}

// GetErcoleFeatures return a map of active/inactive features
func (as *APIService) GetErcoleFeatures() (map[string]bool, utils.AdvancedErrorInterface) {
	partialList, err := as.Database.GetTechnologiesUsage("", "", utils.MAX_TIME)
	if err != nil {
		return nil, err
	}

	out := map[string]bool{}

	out["Oracle/Database"] = partialList["Oracle/Database_hostsCount"] > 0
	out["Oracle/Exadata"] = partialList["Oracle/Exadata"] > 0

	return out, nil
}

// GetTechnologyList return the list of technologies
func (as *APIService) GetTechnologyList() ([]model.TechnologyInfo, utils.AdvancedErrorInterface) {
	return as.TechnologyInfos, nil
}
