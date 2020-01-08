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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/amreo/ercole-services/config"
	"github.com/amreo/ercole-services/utils"
	"github.com/google/go-github/v28/github"
	"github.com/spf13/cobra"
)

// githubToken contains the token used to avoid rate limit
var githubToken string

// fileInfo contains info about all files in repository
type fileInfo struct {
	repository      string
	installed       bool
	version         string
	latest          bool
	releaseDate     string
	filename        string
	name            string
	operatingSystem string
	arch            string
	download        func(string) error
	install         func(string) error
}

//Regex for filenames
var AgentRHEL5Regex *regexp.Regexp = regexp.MustCompile("^ercole-agent-(?P<version>.*)-1.(?P<arch>x86_64).rpm$")
var AgentRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")
var AgentVirtualizationRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-virtualization-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")
var AgentExadataRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-exadata-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")
var AgentWinRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-setup-(?P<version>.*).exe$")
var AgentHpuxRegex *regexp.Regexp = regexp.MustCompile("^ercole-agent-hpux-(?P<version>.*).tar.gz")
var AgentAixRegexRpm *regexp.Regexp = regexp.MustCompile("^ercole-agent-aix-(?P<version>.*)-1.(?P<dist>.*).(?P<arch>noarch).rpm$")
var AgentAixRegexTarGz *regexp.Regexp = regexp.MustCompile("^ercole-agent-aix-(?P<version>.*).tar.gz$")

// setInfoFromFileName sets to fileInfo informations taken from filename
func setInfoFromFileName(filename string, fileInfo *fileInfo) error {
	switch {
	case AgentVirtualizationRHELRegex.MatchString(filename): //agent virtualization RHEL
		data := utils.FindNamedMatches(AgentVirtualizationRHELRegex, filename)
		fileInfo.name = "ercole-agent-virtualization-rhel" + data["dist"]
		fileInfo.version = data["version"]
		fileInfo.arch = data["arch"]
		fileInfo.operatingSystem = "RHEL" + data["dist"]
	case AgentExadataRHELRegex.MatchString(filename): //agent exadata RHEL
		data := utils.FindNamedMatches(AgentExadataRHELRegex, filename)
		fileInfo.name = "ercole-agent-exadata-rhel" + data["dist"]
		fileInfo.version = data["version"]
		fileInfo.arch = data["arch"]
		fileInfo.operatingSystem = "RHEL" + data["dist"]
	case AgentRHEL5Regex.MatchString(filename): //agent RHEL5
		data := utils.FindNamedMatches(AgentRHEL5Regex, filename)
		fileInfo.name = "ercole-agent-rhel5"
		fileInfo.version = data["version"]
		fileInfo.arch = data["arch"]
		fileInfo.operatingSystem = "RHEL5"
	case AgentRHELRegex.MatchString(filename): //agent RHEL
		data := utils.FindNamedMatches(AgentRHELRegex, filename)
		fileInfo.name = "ercole-agent-rhel" + data["dist"]
		fileInfo.version = data["version"]
		fileInfo.arch = data["arch"]
		fileInfo.operatingSystem = "RHEL" + data["dist"]
	case AgentWinRegex.MatchString(filename): //agent WIN
		data := utils.FindNamedMatches(AgentWinRegex, filename)
		fileInfo.name = "ercole-agent-win"
		fileInfo.version = data["version"]
		fileInfo.arch = "x86_64"
		fileInfo.operatingSystem = "WIN"
	case AgentHpuxRegex.MatchString(filename): //agent HPUX
		data := utils.FindNamedMatches(AgentHpuxRegex, filename)
		fileInfo.name = "ercole-agent-hpux"
		fileInfo.version = data["version"]
		fileInfo.arch = "noarch"
		fileInfo.operatingSystem = "HPUX"
	case AgentAixRegexRpm.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(AgentAixRegexRpm, filename)
		fileInfo.name = "ercole-agent-aix"
		fileInfo.version = data["version"]
		fileInfo.arch = "noarch"
		fileInfo.operatingSystem = data["dist"]
	case AgentAixRegexTarGz.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(AgentAixRegexTarGz, filename)
		fileInfo.name = "ercole-agent-aix-targz"
		fileInfo.version = data["version"]
		fileInfo.arch = "noarch"
		fileInfo.operatingSystem = "aix6.1"
	default:
		return fmt.Errorf("Filename %s is not supported. Please check that is correct", filename)
	}

	return nil
}

// scanGithubReleaseRepository scan a github releases repository and return detected files
func scanGithubReleaseRepository(repo config.UpstreamRepository, verbose bool) ([]fileInfo, error) {
	//Fetch the data
	req, _ := http.NewRequest("GET", repo.URL, nil)
	if githubToken != "" {
		req.Header.Add("Authorization", "token "+githubToken)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Received %d from github for URL %s (body: %s)", resp.StatusCode, repo.URL, string(bytes))
	}

	if verbose {
		fmt.Printf("Fetched data from %s\n", repo.URL)
	}

	//Decode data
	var data []github.RepositoryRelease
	json.NewDecoder(resp.Body).Decode(&data)

	//Add data to out
	var out []fileInfo
	for _, release := range data {
		for _, asset := range release.Assets {
			fileInfo := fileInfo{}

			fileInfo.repository = repo.Name
			fileInfo.filename = asset.GetName()
			fileInfo.version = release.GetTagName()
			fileInfo.latest = (fileInfo.version == "latest")
			fileInfo.releaseDate = asset.GetUpdatedAt().Format("2006-01-02")
			if err := setInfoFromFileName(asset.GetName(), &fileInfo); err != nil {
				return nil, err
			}
			fileInfo.installed = fileInfo.checkInstalled()
			fileInfo.download = func(path string) error {
				return utils.DownloadFile(path, asset.GetBrowserDownloadURL())
			}

			out = append(out, fileInfo)
		}
	}

	return out, nil
}

// scanRepository scan a single repository and return detected files
func scanRepository(repo config.UpstreamRepository, verbose bool) ([]fileInfo, error) {
	switch repo.Type {
	case "github-release":
		return scanGithubReleaseRepository(repo, verbose)
	default:
		return nil, fmt.Errorf("Unknown repository type %q of %q", repo.Type, repo.Name)
	}

}

// getFullName return the fullname of the file
func (f *fileInfo) getFullName() string {
	return fmt.Sprintf("%s/%s@%s", f.repository, f.name, f.version)
}

// checkInstalled return true if file is detected in the distribution directory
func (f *fileInfo) checkInstalled() bool {
	if _, err := os.Stat(filepath.Join(ercoleConfig.RepoService.DistributedFiles, f.filename)); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// scanRepositories scan all configured repositories and return a map of file names to info
func scanRepositories(verbose bool) ([]fileInfo, error) {
	out := make([]fileInfo, 0)

	for _, repo := range ercoleConfig.RepoService.UpstreamRepositories {
		//Get repository files
		files, err := scanRepository(repo, verbose)
		if err != nil {
			return nil, err
		}

		out = append(out, files...)
	}

	return out, nil
}

// parseNameOfFile parse a string and return the relative complete filename
func parseNameOfFile(arg string, list []fileInfo) (*fileInfo, error) {
	//valid formats
	//	- <filename> 					<-- supported
	//	- <name>						<-- NOT supported
	//	- <name>@<version>				<-- NOT supported
	//	- <repository>/<name>@<version> <-- supported
	//	- <repository>/<name>			<-- NOT supported
	var regex *regexp.Regexp = regexp.MustCompile("^(?:(?P<repository>[a-z-0-9]+)/)?(?P<name>[a-z-.0-9]+)(?:@(?P<version>[a-z0-9.0-9]+))?$")
	submatches := utils.FindNamedMatches(regex, arg)

	var repository string = submatches["repository"]
	var name string = submatches["name"]
	var version string = submatches["version"]

	for _, f := range list {
		//TODO: add support for missing format
		if name == f.filename {
			return &f, nil
		} else if repository == "" && version == "" && name == f.name {
			return &f, nil
		} else if repository == f.repository && name == f.name && version == f.version {
			return &f, nil
		}
	}

	return nil, fmt.Errorf("Unable to understand the fullname of %q", arg)
}

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "manage the internal repository",
	Long:  `manage the internal repository. It requires to be run where the repository is installed`,
}

//Commands to be added
//	- list
//	- fetch
//	- update
//	- remove
//	- clean
//	- hide
//	- unhide
//	- info
func init() {
	rootCmd.AddCommand(repoCmd)
	rootCmd.PersistentFlags().StringVarP(&githubToken, "github-token", "g", "", "Github token used to perform requests")
}
