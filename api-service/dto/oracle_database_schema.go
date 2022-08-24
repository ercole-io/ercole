package dto

type OracleDatabaseSchema struct {
	Hostname     string `json:"hostname" bson:"hostname"`
	DatabaseName string `json:"databaseName" bson:"databaseName"`
	Indexes      int    `json:"indexes" bson:"indexes"`
	LOB          int    `json:"lob" bson:"lob"`
	Tables       int    `json:"tables" bson:"tables"`
	Total        int    `json:"total" bson:"total"`
	User         string `json:"user" bson:"user"`
}
