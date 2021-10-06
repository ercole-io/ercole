package patch

import (
	"encoding/json"

	"github.com/robertkrimen/otto"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// PatchHostdata patch a single hostdata using the pf PatchingFunction.
// It doesn't check if pf.Hostname equals hostdata["Hostname"]
func PatchHostdata(pf model.PatchingFunction, hostdata model.HostDataBE) (model.HostDataBE, error) {
	//FIXME: avoid repeated marshalling/unmarshalling...

	//Initialize the vm
	vm := otto.New()

	//Convert hostdata om map[string]interface{}
	var tempHD map[string]interface{}
	tempRaw, err := json.Marshal(hostdata)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DATA_PATCHING")
	}
	err = json.Unmarshal(tempRaw, &tempHD)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DATA_PATCHING")
	}

	//Set the global variables
	err = vm.Set("hostdata", tempHD)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DATA_PATCHING")
	}
	err = vm.Set("vars", pf.Vars)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DATA_PATCHING")
	}

	//Run the code
	_, err = vm.Run(pf.Code)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DATA_PATCHING")
	}

	//Convert tempHD to hostdata
	tempRaw, err = json.Marshal(tempHD)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DATA_PATCHING")
	}
	err = json.Unmarshal(tempRaw, &hostdata)
	if err != nil {
		return model.HostDataBE{}, utils.NewError(err, "DATA_PATCHING")
	}

	return hostdata, nil
}
