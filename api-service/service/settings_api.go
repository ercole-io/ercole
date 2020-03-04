package service

import "github.com/amreo/ercole-services/utils"

// GetDefaultDatabaseTags return the default list of database tags from configuration
func (as *APIService) GetDefaultDatabaseTags() ([]string, utils.AdvancedErrorInterface) {
	return as.Config.APIService.DefaultDatabaseTags, nil
}
