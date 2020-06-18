// Copyright (c) 2019 Sorint.lab S.p.A.
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

package utils

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/ercole-io/ercole/model"
	"github.com/robertkrimen/otto"
	"go.mongodb.org/mongo-driver/bson"
)

var MIN_TIME time.Time = time.Unix(0, 0)
var MAX_TIME time.Time = time.Now().AddDate(1000, 0, 0)

//ToJSON convert v to a string containing the equivalent json rappresentation
func ToJSON(v interface{}) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

//ToMongoJSON convert v to a string containing the equivalent json rappresentation
func ToMongoJSON(v interface{}) string {
	raw, err := bson.MarshalExtJSON(v, false, false)
	if err != nil {
		panic(err)
	}
	return string(raw)
}

//FromJSON convert a json str to interface containing the equivalent json rappresentation
func FromJSON(str []byte) interface{} {
	var out map[string]interface{}
	json.Unmarshal(str, &out)
	return out
}

//Intptr return a point to the int passed in the argument
func Intptr(v int64) *int64 {
	return &v
}

// Contains return true if a contains x, otherwise false.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Remove return slice without element at position i, mantaining order
func Remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}

// Str2bool parse a string to a boolean
func Str2bool(in string, defaultValue bool) (bool, AdvancedErrorInterface) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseBool(in); err != nil {
		return false, NewAdvancedErrorPtr(err, "Unable to parse string to bool")
	} else {
		return val, nil
	}
}

// Str2int parse a string to a int
func Str2int(in string, defaultValue int) (int, AdvancedErrorInterface) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseInt(in, 10, 32); err != nil {
		return -1, NewAdvancedErrorPtr(err, "Unable to parse string to int")
	} else {
		return int(val), nil
	}
}

// Str2float32 parse a string to a float32
func Str2float32(in string, defaultValue float32) (float32, AdvancedErrorInterface) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := strconv.ParseFloat(in, 32); err != nil {
		return -1, NewAdvancedErrorPtr(err, "Unable to parse string to float")
	} else {
		return float32(val), nil
	}
}

// Str2time parse a string to a time
func Str2time(in string, defaultValue time.Time) (time.Time, AdvancedErrorInterface) {
	if in == "" {
		return defaultValue, nil
	} else if val, err := time.Parse(time.RFC3339, in); err != nil {
		return time.Time{}, NewAdvancedErrorPtr(err, "Unable to parse string to time.Time")
	} else {
		return val, nil
	}
}

// NewAPIUrl return a new url crafted using the parameters
func NewAPIUrl(baseURL string, username string, password string, path string, params url.Values) *url.URL {
	u := NewAPIUrlNoParams(baseURL, username, password, path)
	u.RawQuery = params.Encode()

	return u
}

// NewAPIUrlNoParams return a new url crafted using the parameters
func NewAPIUrlNoParams(baseURL string, username string, password string, path string) *url.URL {
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	u.User = url.UserPassword(username, password)
	u.Path += path

	return u
}

// FindNamedMatches return the map of the groups of str
func FindNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

// DownloadFile download the file from url into the filepath
func DownloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// PatchHostdata patch a single hostdata using the pf PatchingFunction.
// It doesn't check if pf.Hostname equals hostdata["Hostname"]
func PatchHostdata(pf model.PatchingFunction, hostdata model.HostDataBE) (model.HostDataBE, AdvancedErrorInterface) {
	//FIXME: avoid repeated marshalling/unmarshalling...

	//Initialize the vm
	vm := otto.New()

	//Convert hostdata om map[string]interface{}
	var tempHD map[string]interface{}
	tempRaw, err := json.Marshal(hostdata)
	if err != nil {
		return model.HostDataBE{}, NewAdvancedErrorPtr(err, "DATA_PATCHING")
	}
	err = json.Unmarshal(tempRaw, &tempHD)
	if err != nil {
		return model.HostDataBE{}, NewAdvancedErrorPtr(err, "DATA_PATCHING")
	}

	//Set the global variables
	err = vm.Set("hostdata", tempHD)
	if err != nil {
		return model.HostDataBE{}, NewAdvancedErrorPtr(err, "DATA_PATCHING")
	}
	err = vm.Set("vars", pf.Vars)
	if err != nil {
		return model.HostDataBE{}, NewAdvancedErrorPtr(err, "DATA_PATCHING")
	}

	//Run the code
	_, err = vm.Run(pf.Code)
	if err != nil {
		return model.HostDataBE{}, NewAdvancedErrorPtr(err, "DATA_PATCHING")
	}

	//Convert tempHD to hostdata
	tempRaw, err = json.Marshal(tempHD)
	if err != nil {
		return model.HostDataBE{}, NewAdvancedErrorPtr(err, "DATA_PATCHING")
	}
	err = json.Unmarshal(tempRaw, &hostdata)
	if err != nil {
		return model.HostDataBE{}, NewAdvancedErrorPtr(err, "DATA_PATCHING")
	}

	return hostdata, nil
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// ParsePrivateKey converts a private key expressed as []byte to interface{}
func ParsePrivateKey(raw []byte) (interface{}, interface{}, AdvancedErrorInterface) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, nil, NewAdvancedErrorPtr(errors.New("Unable to parse the private key"), "PARSE_PRIVATE_KEY")
	}

	privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, NewAdvancedErrorPtr(err, "PARSE_PRIVATE_KEY")
	}
	return privatekey, &privatekey.PublicKey, nil
}

// ParsePublicKey converts a private key expressed as []byte to interface{}
func ParsePublicKey(raw []byte) (interface{}, AdvancedErrorInterface) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, NewAdvancedErrorPtr(errors.New("Unable to parse the public key"), "PARSE_PUBLIC_KEY")
	}

	publickey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, NewAdvancedErrorPtr(err, "PARSE_PUBLIC_KEY")
	}
	return publickey, nil
}
