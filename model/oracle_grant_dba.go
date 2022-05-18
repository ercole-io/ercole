package model

type OracleGrantDba struct {
	Grantee     string `json:"grantee" bson:"grantee"`
	AdminOption string `json:"adminOption" bson:"adminOption"`
	DefaultRole string `json:"defaultRole" bson:"defaultRole"`
}

