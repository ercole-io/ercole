package dto

import "time"

type HostCores struct {
	Date  time.Time `json:"date" bson:"date"`
	Cores int       `json:"cores" bson:"cores"`
}
