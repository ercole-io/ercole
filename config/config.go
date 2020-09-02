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

package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/ercole-io/ercole/utils"
	"github.com/goraz/onion"
	"github.com/goraz/onion/layers/directorylayer"
	_ "github.com/goraz/onion/loaders/toml-0.5.0" // Needed to load toml files
	"github.com/goraz/onion/onionwriter"
	"github.com/sirupsen/logrus"
)

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
	// ChartService contains configuration about the chart service
	ChartService ChartService
	// Mongodb contains configuration about database connection, some data logic and migration
	Mongodb Mongodb
	// Version contains the version of the server
	Version string
	// ResourceFilePath contains the directory of the resources
	ResourceFilePath string
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
	// LogDataPatching enable the logging of the data patching
	LogDataPatching bool
	// EnablePatching enable the patching of the arrived hostdata
	EnablePatching bool
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
	// QueueBufferSize contains the size of the buffer of the queue
	QueueBufferSize int
	// FreshnessCheckJob contains the parameters of the freshness check
	FreshnessCheckJob FreshnessCheckJob
	// Emailer contains the settings about the emailer
	Emailer Emailer
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
	// ReadOnly disable modifing APIs
	ReadOnly bool
	// EnableInsertingCustomPatchingFunction enable the API for inserting custom patching functions
	EnableInsertingCustomPatchingFunction bool
	// DebugOracleDatabaseAgreementsAssignmentAlgorithm enable the debugging of the Oracle/Database agreements assignment algorithm
	DebugOracleDatabaseAgreementsAssignmentAlgorithm bool
	// AuthenticationProvider contains info about how the users are authenticated
	AuthenticationProvider AuthenticationProviderConfig
	// OperatingSystemAggregationRules contains rules used to aggregate various operating systems
	OperatingSystemAggregationRules []AggregationRule
	// DefaultDatabaseTags contains the default list of database tags
	DefaultDatabaseTags []string
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

// ChartService contains configuration about the chart service
type ChartService struct {
	// RemoteEndpoint contains the endpoint used to connect to the ChartService
	RemoteEndpoint string
	// BindIP contains the bind ip
	BindIP string
	// Port contains the port of the internal http server
	Port uint16
	// LogHTTPRequest enable the logging of the internal http serverl
	LogHTTPRequest bool
}

// Mongodb contains configuration about the database connection, some data logic and migration
type Mongodb struct {
	// URI contains MongoDB connection string/URI like 'mongodb://localhost:27017/ercole'
	URI string
	// DBName contains the name of the database
	DBName string
	// Migrate is true when mongodb should update/migrate data/schema during the initializazion
	Migrate bool
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
	// Product contains vendor/name of the operating system
	Product string
}

// UpstreamRepository contains info about a upstream repository
type UpstreamRepository struct {
	// Name of the repository
	Name string
	// Type of the repository
	// Supported types are:
	//		- github-release
	//		- directory
	//		- ercole-reposervice
	Type string
	// URL of the repository where to find files
	URL string
}

// Emailer contains settings used to send emails
type Emailer struct {
	// Enabled contains true if the emailer is enabled, otherwise false
	Enabled bool
	// From contains the source email
	From string
	// To contains the destinations
	To []string
	// SMTPServer contains the address or hostname of the server
	SMTPServer string
	// SMTPPort contains the port of the server
	SMTPPort int
	// SMTPUsername contains the username used to connect to the server
	SMTPUsername string
	// SMTPPassword contains the password used to connect to the server
	SMTPPassword string
	// DisableSSLCertificateValidation contains true if disable the certification validation, otherwise false
	DisableSSLCertificateValidation bool
}

// AuthenticationProviderConfig contains the settings used to authenticate the users
type AuthenticationProviderConfig struct {
	// Type contains the type of the source. Supported types are:
	//	- basic
	// 	- ldap
	Type string
	// Username is the username of the user if type == "basic"
	Username string
	// Password is the password of the user if type == "basic"
	Password string
	// PrivateKey is the filename of the key
	PrivateKey string
	// PublicKey is the filename of the key
	PublicKey string
	// TokenValidityTimeout contains the number of seconds in which the token is still valid
	TokenValidityTimeout int
	Host                 string
	Port                 int
	LDAPBase             string
	LDAPUseSSL           bool
	LDAPBindDN           string
	LDAPBindPassword     string
	LDAPUserFilter       string
	LDAPGroupFilter      string
}

// ReadConfig read, parse and return a Configuration from the configuration file
func ReadConfig(log *logrus.Logger, extraConfigFile string) (configuration Configuration) {
	layers := make([]onion.Layer, 0)

	layers = addFileLayers(log, layers, "/opt/ercole/config.toml")

	dataDirs := xdg.DataDirs()
	for i := 0; i < len(dataDirs); i++ {
		dataDirs[i] = filepath.Join(dataDirs[i], "ercole/config.toml")
	}
	layers = addFileLayers(log, layers, dataDirs...)

	layers = addFileLayers(log, layers, "/etc/ercole/ercole.toml")

	etcErcoleDirectory := "/etc/ercole/conf.d/"
	directoryLayer, err := directorylayer.NewDirectoryLayer(etcErcoleDirectory, "toml")
	if err == nil {
		layers = append(layers, directoryLayer)
	} else if !strings.Contains(err.Error(), "no such file or directory") {
		log.Warnf("error reading file [%s]: [%s]", etcErcoleDirectory, err)
	}

	layers = addFileLayers(log, layers,
		xdg.ConfigHome()+"/ercole.toml",
		"config.toml",
		extraConfigFile,
	)

	configOnion := onion.New(layers...)

	err = onionwriter.DecodeOnion(configOnion, &configuration)
	if err != nil {
		log.Fatal("something went wrong while reading your configuration files")
	}

	patchConfiguration(&configuration)

	return configuration
}

func addFileLayers(log *logrus.Logger, layers []onion.Layer, configFiles ...string) []onion.Layer {

	for _, file := range configFiles {
		layer, err := onion.NewFileLayer(file, nil)

		var pathErr *os.PathError

		if err == nil {
			layers = append(layers, layer)
		} else if !errors.As(err, &pathErr) {
			log.Warnf("error reading file [%s]: [%s]", file, err)
		}
	}

	return layers
}

// PatchConfiguration change the value of the fields for meeting some requirements(?)
func patchConfiguration(config *Configuration) {
	cwd, _ := os.Readlink("/proc/self/exe")
	cwd = filepath.Dir(cwd)

	if config.RepoService.DistributedFiles == "" {
		config.RepoService.DistributedFiles = "/var/lib/ercole/distributed_files"
	} else if filepath.IsAbs(config.RepoService.DistributedFiles) && !strings.HasSuffix(config.RepoService.DistributedFiles, "/") {
		config.RepoService.DistributedFiles = config.RepoService.DistributedFiles + "/"
	} else if !filepath.IsAbs(config.RepoService.DistributedFiles) {
		config.RepoService.DistributedFiles = cwd + filepath.Join("/", config.RepoService.DistributedFiles) + "/"
	}

	if config.ResourceFilePath == "" {
		if utils.FileExists(filepath.Join(cwd, "resources")) {
			config.ResourceFilePath = filepath.Join(cwd, "resources")
		} else if utils.FileExists("/usr/share/ercole") {
			config.ResourceFilePath = "/usr/share/ercole"
		} else {
			config.ResourceFilePath = "RESOURCES_NOT_FOUND"
		}
	} else if !filepath.IsAbs(config.ResourceFilePath) {
		config.ResourceFilePath = cwd + filepath.Join("/", config.ResourceFilePath)
	}
}
