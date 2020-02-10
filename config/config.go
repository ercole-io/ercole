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

// Package config contains configuration utilities, like readConfig()
package config

// Configuration contains Ercole DataService configuration
type Configuration struct {
	// DataService contains configuration about the data service
	DataService DataService
	// DataService contains configuration about the alert service
	AlertService AlertService
	// APIService contains configuration about the api service
	APIService APIService
	// RepoService contains configuration about the repo service
	RepoService RepoService
	// Mongodb contains configuration about database connection, some data logic and migration
	Mongodb Mongodb
	// Version contains the version of the server
	Version string
}

// DataService contains configuration about the data service
type DataService struct {
	// RemoteEndpoint contains the endpoint used to connect to the DataService
	RemoteEndpoint string
	// BindIP contains the bind ip
	BindIP string
	// Port contains the port of the internal http server
	Port uint16
	// LogHTTPRequest enable the logging of the internal http serverl
	LogHTTPRequest bool
	// LogInsertingHostdata enable the logging of the inserting hostdata
	LogInsertingHostdata bool
	// AgentUsername contains the username of the agent
	AgentUsername string
	// AgentPassword contains the password of the agent
	AgentPassword string
	// CurrentHostCleaningJob contains the parameters of the current host cleaning
	CurrentHostCleaningJob CurrentHostCleaningJob
	// ArchivedCleaningJob contains the parameters of the archived host cleaning
	ArchivedHostCleaningJob ArchivedHostCleaningJob
}

// AlertService contains configuration about the alert service
type AlertService struct {
	// RemoteEndpoint contains the endpoint used to connect to the AlertService
	RemoteEndpoint string
	// BindIP contains the bind ip
	BindIP string
	// Port contains the port of the internal http server
	Port uint16
	// LogHTTPRequest enable the logging of the internal http serverl
	LogHTTPRequest bool
	// LogHTTPRequest enable the logging of the received messages
	LogMessages bool
	// LogThrows enable the logging of alert throws
	LogAlertThrows bool
	// PublisherUsername contains the username of the agent
	PublisherUsername string
	// PublisherPassword contains the password of the agent
	PublisherPassword string
	// FreshnessCheckJob contains the parameters of the freshness check
	FreshnessCheckJob FreshnessCheckJob
}

// APIService contains configuration about the api service
type APIService struct {
	// RemoteEndpoint contains the endpoint used to connect to the APIService
	RemoteEndpoint string
	// BindIP contains the bind ip
	BindIP string
	// Port contains the port of the internal http server
	Port uint16
	// LogHTTPRequest enable the logging of the internal http serverl
	LogHTTPRequest bool
	// UserUsername contains the username of the user
	UserUsername string
	// UserPassword contains the password of the user
	UserPassword string
	// OperatingSystemAggregationRules contains rules used to aggregate various operating systems
	OperatingSystemAggregationRules []AggregationRule
}

// RepoService contains configuration about the repo service
type RepoService struct {
	// UpstreamRepository contains the list of upstream repositories
	UpstreamRepositories []UpstreamRepository
	// DistributedFiles contains the path to the files to be served
	DistributedFiles string
	// HTTP contains the configuration about the HTTP server
	HTTP HTTPRepoService
	// SFTP contains the configuration about the SFTP server
	SFTP SFTPRepoService
}

// Mongodb contains configuration about the database connection, some data logic and migration
type Mongodb struct {
	// URI contains MongoDB connection string/URI like 'mongodb://localhost:27017/ercole'
	URI string
	// DBName contains the name of the database
	DBName string
	// Migrate is true when mongodb should update/migrate data/schema during the initializazion
	Migrate bool
	// LicensesList contains the filename of the file that contains the list of licenses
	LicensesList string
}

// FreshnessCheckJob contains parameters for the freshness check
type FreshnessCheckJob struct {
	// Crontab contains the crontab string used to schedule the freshness check
	Crontab string
	// DaysThreshold contains the threshdold of the freshness check
	DaysThreshold int
	// RunAtStartup contains true if the job should run when the service start, otherwise false
	RunAtStartup bool
}

// CurrentHostCleaningJob contains parameters for the current host cleaning
type CurrentHostCleaningJob struct {
	// Crontab contains the crontab string used to schedule the cleaning
	Crontab string
	// DaysThreshold contains the threshdold of the current host cleaning
	HourThreshold int
	// RunAtStartup contains true if the job should run when the service start, otherwise false
	RunAtStartup bool
}

// ArchivedHostCleaningJob contains parameters for the archived host cleaning
type ArchivedHostCleaningJob struct {
	// Crontab contains the crontab string used to schedule the cleaning
	Crontab string
	// DaysThreshold contains the threshdold of the archived host cleaning
	HourThreshold int
	// RunAtStartup contains true if the job should run when the service start, otherwise false
	RunAtStartup bool
}

// HTTPRepoService contains parameters for a single serving service
type HTTPRepoService struct {
	// Enable contains true it the service is enabled, otherwise false
	Enable bool
	// RemoteEndpoint contains the endpoint used to connect to the HTTPRepoService
	RemoteEndpoint string
	// BindIP contains the bind ip
	BindIP string
	// Port contains the port of the internal http server
	Port uint16
	// LogHTTPRequest enable the logging of the internal http serverl
	LogHTTPRequest bool
}

// SFTPRepoService contains parameters for a single serving service
type SFTPRepoService struct {
	// Enable contains true it the service is enabled, otherwise false
	Enable bool
	// RemoteEndpoint contains the endpoint used to connect to the SFTPRepoService
	RemoteEndpoint string
	// BindIP contains the bind ip
	BindIP string
	// Port contains the port of the sftp server
	Port uint16
	// PrivateKey contains the path to the private key
	PrivateKey string
	// LogConnections contains true if log connections, otherwise false
	LogConnections bool
	// DebugConnections contains true if debug connections, otherwise false
	DebugConnections bool
}

// AggregationRule contains a rule used to aggregate string per group
type AggregationRule struct {
	// Regex contains the regular expression used for matching the aggregation group
	Regex string
	// Group contains the name of the group
	Group string
}

// UpstreamRepository contains info about a upstream repository
type UpstreamRepository struct {
	// Name of the repository
	Name string
	// Type of the repository
	Type string
	// URL of the repository where to find files
	URL string
}
