package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HostDataMap holds all informations about a host & services, in a map format
type HostDataMap map[string]interface{}

// ID getter
func (hd *HostDataMap) ID() primitive.ObjectID {
	id := (*hd)["_id"]
	return id.(primitive.ObjectID)
}

// SetID setter
func (hd *HostDataMap) SetID(id primitive.ObjectID) {
	(*hd)["_id"] = id
}

// Hostname getter
func (hd *HostDataMap) Hostname() string {
	hostname := (*hd)["Hostname"]
	return hostname.(string)
}

// CreatedAt getter
func (hd *HostDataMap) CreatedAt() time.Time {
	createdAt := (*hd)["CreatedAt"]
	return (createdAt.(primitive.DateTime)).Time().UTC()
}

// SetCreatedAt setter
func (hd *HostDataMap) SetCreatedAt(t time.Time) {
	(*hd)["CreatedAt"] = primitive.NewDateTimeFromTime(t)
}

// Extra getter
func (hd *HostDataMap) Extra() ExtraInfoMap {
	extraInfo := (*hd)["Extra"]
	return extraInfo.(map[string]interface{})
}

// Info getter
func (hd *HostDataMap) Info() HostMap {
	info := (*hd)["Info"]
	return info.(map[string]interface{})
}

// NewEmptyHostDataMap return an empty initialized HostDataMap
func NewEmptyHostDataMap() HostDataMap {
	return map[string]interface{}{
		"Hostname":  "",
		"CreatedAt": time.Time{},
		"Archived":  false,
		"Extra":     map[string]interface{}{},
		"Info": map[string]interface{}{
			"Hostname":       "",
			"Environment":    "",
			"Location":       "",
			"CPUModel":       "",
			"CPUCores":       0,
			"CPUThreads":     0,
			"Socket":         0,
			"Type":           "",
			"Virtual":        false,
			"Kernel":         "",
			"OS":             "",
			"MemoryTotal":    0,
			"SwapTotal":      0,
			"OracleCluster":  false,
			"VeritasCluster": false,
			"SunCluster":     false,
			"AixCluster":     false,
		},
	}
}
