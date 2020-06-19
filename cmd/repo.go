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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/1set/gut/yos"
	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	"github.com/google/go-github/v28/github"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

// githubToken contains the token used to avoid rate limit
var githubToken string
var rebuildCache bool

type index []*artifactInfo

// artifactInfo contains info about all files in repository
type artifactInfo struct {
	Repository            string
	Installed             bool `json:"-"`
	Version               string
	ReleaseDate           string
	Filename              string
	Name                  string
	OperatingSystemFamily string
	OperatingSystem       string
	Arch                  string
	UpstreamType          string
	UpstreamInfo          map[string]interface{}
	Install               func(ai *artifactInfo)              `json:"-"`
	Uninstall             func(ai *artifactInfo)              `json:"-"`
	Download              func(ai *artifactInfo, dest string) `json:"-"`
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
var ErcoleRHELRegex *regexp.Regexp = regexp.MustCompile("^ercole-(?P<version>.*)-1.el(?P<dist>\\d+).(?P<arch>x86_64).rpm$")

func cmpVersion(a, b string) bool {
	va, err := version.NewVersion(a)
	if err != nil {
		panic(err)
	}
	vb, err := version.NewVersion(b)
	if err != nil {
		panic(err)
	}
	return va.LessThan(vb)
}

// setInfoFromFileName sets to fileInfo informations taken from filename
func setInfoFromFileName(filename string, artifactInfo *artifactInfo) {
	switch {
	case AgentVirtualizationRHELRegex.MatchString(filename): //agent virtualization RHEL
		data := utils.FindNamedMatches(AgentVirtualizationRHELRegex, filename)
		artifactInfo.Name = "ercole-agent-virtualization-rhel" + data["dist"]
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = data["arch"]
		artifactInfo.OperatingSystemFamily = "rhel"
		artifactInfo.OperatingSystem = "rhel" + data["dist"]
	case AgentExadataRHELRegex.MatchString(filename): //agent exadata RHEL
		data := utils.FindNamedMatches(AgentExadataRHELRegex, filename)
		artifactInfo.Name = "ercole-agent-exadata-rhel" + data["dist"]
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = data["arch"]
		artifactInfo.OperatingSystemFamily = "rhel"
		artifactInfo.OperatingSystem = "rhel" + data["dist"]
	case AgentRHEL5Regex.MatchString(filename): //agent RHEL5
		data := utils.FindNamedMatches(AgentRHEL5Regex, filename)
		artifactInfo.Name = "ercole-agent-rhel5"
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = data["arch"]
		artifactInfo.OperatingSystemFamily = "rhel"
		artifactInfo.OperatingSystem = "rhel5"
	case AgentRHELRegex.MatchString(filename): //agent RHEL
		data := utils.FindNamedMatches(AgentRHELRegex, filename)
		artifactInfo.Name = "ercole-agent-rhel" + data["dist"]
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = data["arch"]
		artifactInfo.OperatingSystemFamily = "rhel"
		artifactInfo.OperatingSystem = "rhel" + data["dist"]
	case ErcoleRHELRegex.MatchString(filename): //ercole RHEL
		data := utils.FindNamedMatches(ErcoleRHELRegex, filename)
		artifactInfo.Name = "ercole-" + data["dist"]
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = data["arch"]
		artifactInfo.OperatingSystemFamily = "rhel"
		artifactInfo.OperatingSystem = "rhel" + data["dist"]
	case AgentWinRegex.MatchString(filename): //agent WIN
		data := utils.FindNamedMatches(AgentWinRegex, filename)
		artifactInfo.Name = "ercole-agent-win"
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = "x86_64"
		artifactInfo.OperatingSystemFamily = "win"
		artifactInfo.OperatingSystem = "win"
	case AgentHpuxRegex.MatchString(filename): //agent HPUX
		data := utils.FindNamedMatches(AgentHpuxRegex, filename)
		artifactInfo.Name = "ercole-agent-hpux"
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = "noarch"
		artifactInfo.OperatingSystemFamily = "hpux"
		artifactInfo.OperatingSystem = "hpux"
	case AgentAixRegexRpm.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(AgentAixRegexRpm, filename)
		artifactInfo.Name = "ercole-agent-aix"
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = "noarch"
		artifactInfo.OperatingSystemFamily = "aix"
		artifactInfo.OperatingSystem = data["dist"]
	case AgentAixRegexTarGz.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(AgentAixRegexTarGz, filename)
		artifactInfo.Name = "ercole-agent-aix-targz"
		artifactInfo.Version = data["version"]
		artifactInfo.Arch = "noarch"
		artifactInfo.OperatingSystemFamily = "aix-tar-gz"
		artifactInfo.OperatingSystem = "aix6.1"
	default:
		panic(fmt.Errorf("Filename %s is not supported. Please check that is correct", filename))
	}
}

// setDownloader set the downloader of the artifact
func setDownloader(artifact *artifactInfo) {
	switch artifact.UpstreamType {
	case "github-release":
		artifact.Download = func(ai *artifactInfo, dest string) {
			utils.DownloadFile(dest, ai.UpstreamInfo["DownloadUrl"].(string))
		}
	case "directory":
		artifact.Download = func(ai *artifactInfo, dest string) {
			if verbose {
				fmt.Printf("Copying file from %s to %s\n", ai.UpstreamInfo["Filename"].(string), dest)
			}
			err := yos.CopyFile(ai.UpstreamInfo["Filename"].(string), dest)
			if err != nil {
				panic(err)
			}
			err = os.Chmod(dest, 0755)
			if err != nil {
				panic(err)
			}
		}
	case "ercole-reposervice":
		artifact.Download = func(ai *artifactInfo, dest string) {
			utils.DownloadFile(dest, ai.UpstreamInfo["DownloadUrl"].(string))
		}
	default:
		panic(artifact)
	}
}

// setInstaller set the installer of the artifact
func setInstaller(artifact *artifactInfo) {
	switch {
	case strings.HasSuffix(artifact.Filename, ".rpm"):
		artifact.Install = func(ai *artifactInfo) {
			//Create missing directories
			if verbose {
				fmt.Printf("Creating the directories (if missing) %s, %s\n",
					filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch),
					filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all"),
				)
			}
			err := os.MkdirAll(filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch), 0755)
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all"), 0755)
			if err != nil {
				panic(err)
			}

			//Download the file in the right location
			if verbose {
				fmt.Printf("Downloading the artifact %s to %s\n", ai.Filename, filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			ai.Download(ai, filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))

			//Create a link to all
			if verbose {
				fmt.Printf("Linking the artifact to %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			}
			err = os.Link(filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename), filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Launch the createrepo command
			if verbose {
				fmt.Printf("Executing createrepo %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			}
			cmd := exec.Command("createrepo", filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			if verbose {
				cmd.Stdout = os.Stdout
			}
			cmd.Stderr = os.Stderr
			cmd.Run()

			//Settint it to installed
			ai.Installed = true
		}
	default:
		artifact.Install = func(ai *artifactInfo) {
			//Create missing directories
			if verbose {
				fmt.Printf("Creating the directories (if missing) %s, %s\n",
					filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch),
					filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all"),
				)
			}
			err := os.MkdirAll(filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch), 0755)
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all"), 0755)
			if err != nil {
				panic(err)
			}

			//Download the file in the right location
			if verbose {
				fmt.Printf("Downloading the artifact %s to %s\n", ai.Filename, filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			ai.Download(ai, filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, "/", ai.OperatingSystem, ai.Arch, ai.Filename))

			//Create a link to all
			if verbose {
				fmt.Printf("Linking the artifact to %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			}
			err = os.Link(filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, "/", ai.OperatingSystem, ai.Arch, ai.Filename), filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Settint it to installed
			ai.Installed = true
		}
	}
}

// setInstaller set the installer of the artifact
func setUninstaller(artifact *artifactInfo) {
	switch {
	case strings.HasSuffix(artifact.Filename, ".rpm"):
		artifact.Uninstall = func(ai *artifactInfo) {
			//Removing the link to all
			if verbose {
				fmt.Printf("Removing Linking the artifact to %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			}
			err := os.Remove(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Removing the file
			if verbose {
				fmt.Printf("Removing the file %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			err = os.Remove(filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			if err != nil {
				panic(err)
			}

			//Launch the createrepo command
			if verbose {
				fmt.Printf("Executing createrepo %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			}
			cmd := exec.Command("createrepo", filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch))
			if verbose {
				cmd.Stdout = os.Stdout
			}
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				panic(err)
			}

			//Set it to not installed
			ai.Installed = false
		}
	default:
		artifact.Uninstall = func(ai *artifactInfo) {
			//Removing the link to all
			if verbose {
				fmt.Printf("Removing Linking the artifact to %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			}
			err := os.Remove(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", ai.Filename))
			if err != nil {
				panic(err)
			}

			//Removing the file
			if verbose {
				fmt.Printf("Removing the file %s\n", filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			}
			err = os.Remove(filepath.Join(ercoleConfig.RepoService.DistributedFiles, ai.OperatingSystemFamily, ai.OperatingSystem, ai.Arch, ai.Filename))
			if err != nil {
				panic(err)
			}

			//Set it to not installed
			ai.Installed = false
		}
	}
}

// scanGithubReleaseRepository scan a github releases repository and return detected files
func scanGithubReleaseRepository(repo config.UpstreamRepository) (index, error) {
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
	var out index
	for _, release := range data {
		for _, asset := range release.Assets {
			artifactInfo := new(artifactInfo)

			artifactInfo.Repository = repo.Name
			artifactInfo.Filename = asset.GetName()
			artifactInfo.Version = release.GetTagName()
			artifactInfo.ReleaseDate = asset.GetUpdatedAt().Format("2006-01-02")
			artifactInfo.UpstreamType = "github-release"
			artifactInfo.UpstreamInfo = map[string]interface{}{
				"DownloadUrl": asset.GetBrowserDownloadURL(),
			}
			setInfoFromFileName(artifactInfo.Filename, artifactInfo)
			if artifactInfo.Version == "latest" {
				continue
			}
			out = append(out, artifactInfo)
		}
	}

	return out, nil
}

// scanDirectoryRepository scan the local directory and return detected files
func scanDirectoryRepository(repo config.UpstreamRepository) (index, error) {
	//Fetch the data
	files, err := ioutil.ReadDir(repo.URL)
	if err != nil {
		log.Fatal(err)
	}

	//Add data to out
	var out index
	for _, file := range files {
		artifactInfo := new(artifactInfo)

		artifactInfo.Repository = repo.Name
		artifactInfo.Filename = filepath.Base(file.Name())
		artifactInfo.ReleaseDate = file.ModTime().Format("2006-01-02")
		artifactInfo.UpstreamType = "directory"
		artifactInfo.UpstreamInfo = map[string]interface{}{
			"Filename": filepath.Join(repo.URL, file.Name()),
		}
		setInfoFromFileName(artifactInfo.Filename, artifactInfo)
		if artifactInfo.Version == "latest" {
			continue
		}
		out = append(out, artifactInfo)
	}

	return out, nil
}

// scanGithubReleaseRepository scan a github releases repository and return detected files
func scanErcoleReposerviceRepository(repo config.UpstreamRepository) (index, error) {
	//Fetch the data
	req, _ := http.NewRequest("GET", repo.URL+"/all", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Received %d from ercole reposervice for URL %s (body: %s)", resp.StatusCode, repo.URL+"/all", string(bytes))
	}

	if verbose {
		fmt.Printf("Fetched data from %s\n", repo.URL)
	}

	//Extract the filenames
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile("<a href=\"([^\"]*)\">")

	// //Add data to out
	var out index
	installedNames := make([]string, 0)

	for _, fn := range regex.FindAllStringSubmatch(string(data), -1) {
		installedNames = append(installedNames, fn[1])
	}

	//Fetch the data
	req, _ = http.NewRequest("GET", repo.URL+"/index.json", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Received %d from ercole reposervice for URL %s (body: %s)", resp.StatusCode, repo.URL+"/index.json", string(bytes))
	}

	if verbose {
		fmt.Printf("Fetched data from %s\n", repo.URL)
	}

	//Decode data
	var data2 []map[string]string
	json.NewDecoder(resp.Body).Decode(&data2)

	for _, d := range data2 {
		if !utils.Contains(installedNames, d["Filename"]) {
			continue
		}

		artifactInfo := new(artifactInfo)

		artifactInfo.Repository = repo.Name
		artifactInfo.Filename = d["Filename"]
		artifactInfo.ReleaseDate = d["ReleaseDate"]
		artifactInfo.UpstreamType = "ercole-reposervice"
		artifactInfo.UpstreamInfo = map[string]interface{}{
			"DownloadUrl": repo.URL + "/all/" + d["Filename"],
		}
		setInfoFromFileName(artifactInfo.Filename, artifactInfo)
		if artifactInfo.Version == "latest" {
			continue
		}
		out = append(out, artifactInfo)
	}
	return out, nil
}

// scanRepository scan a single repository and return detected files
func scanRepository(repo config.UpstreamRepository) (index, error) {
	switch repo.Type {
	case "github-release":
		return scanGithubReleaseRepository(repo)
	case "directory":
		return scanDirectoryRepository(repo)
	case "ercole-reposervice":
		return scanErcoleReposerviceRepository(repo)
	default:
		return nil, fmt.Errorf("Unknown repository type %q of %q", repo.Type, repo.Name)
	}

}

// getFullName return the fullname of the file
func (f *artifactInfo) getFullName() string {
	return fmt.Sprintf("%s/%s@%s", f.Repository, f.Name, f.Version)
}

// checkInstalled return true if file is detected in the distribution directory
func (f *artifactInfo) checkInstalled() bool {
	if _, err := os.Stat(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "all", f.Filename)); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// scanRepositories scan all configured repositories and return a map of file names to info
func scanRepositories() index {
	out := make(index, 0)

	for _, repo := range ercoleConfig.RepoService.UpstreamRepositories {
		//Get repository files
		files, err := scanRepository(repo)
		if err != nil {
			panic(err)
		}

		out = append(out, files...)
	}

	return out
}

// parseNameOfFile parse a string and return the relative complete filename
func (idx *index) searchArtifactByArg(arg string) *artifactInfo {
	//valid formats
	//	- <filename>
	//	- <name>
	//	- <name>@<version>
	//	- <repository>/<name>@<version>
	//	- <repository>/<name>
	var regex *regexp.Regexp = regexp.MustCompile("^(?:(?P<repository>[a-z-0-9]+)/)?(?P<name>[a-z-.0-9]+)(?:@(?P<version>[a-z0-9.0-9]+))?$")
	submatches := utils.FindNamedMatches(regex, arg)

	var repository string = submatches["repository"]
	var name string = submatches["name"]
	var version string = submatches["version"]

	var foundArtifact *artifactInfo

	foundArtifact = idx.searchArtifactByFilename(arg)
	if foundArtifact != nil {
		return foundArtifact
	}

	switch {
	case repository == "" && version == "": //	<name> case
		foundArtifact = idx.searchArtifactByName(name)
	case repository == "" && version != "": //	<name>@<version> case
		foundArtifact = idx.searchArtifactByNameAndVersion(name, version)
	case repository != "" && version != "": //	<repository>/<name>@<version> case
		foundArtifact = idx.searchArtifactByFullname(repository, name, version)
	case repository != "" && version == "": //	<repository>/<name> case
		foundArtifact = idx.searchArtifactByRepositoryAndName(repository, name)
	}

	return foundArtifact
}

func (idx *index) searchArtifactByFilename(filename string) *artifactInfo {
	var foundArtifact *artifactInfo

	//Find the artifact
	for _, f := range *idx {
		//TODO: add support for missing format
		if filename == f.Filename {
			if foundArtifact == nil {
				foundArtifact = f
			} else {
				panic(fmt.Errorf("Two artifact have the same filename: %v and %v", foundArtifact, f))
			}

		}
	}

	return foundArtifact
}

func (idx *index) searchArtifactByName(name string) *artifactInfo {
	var foundArtifact *artifactInfo

	//Find the artifact
	for _, f := range *idx {
		//TODO: add support for missing format
		if name == f.Name {
			if foundArtifact == nil {
				foundArtifact = f
			} else if foundArtifact.Repository == f.Repository {
				if cmpVersion(foundArtifact.Version, f.Version) {
					foundArtifact = f
				}
			} else {
				panic(fmt.Errorf("Two artifact have the same name: %v and %v", foundArtifact, f))
			}
		}
	}

	return foundArtifact
}

func (idx *index) searchArtifactByNameAndVersion(name string, version string) *artifactInfo {
	var foundArtifact *artifactInfo

	//Find the artifact
	for _, f := range *idx {
		//TODO: add support for missing format
		if name == f.Name && version == f.Version {
			if foundArtifact == nil {
				foundArtifact = f
			} else {
				panic(fmt.Errorf("Two artifact have the same name and version: %v and %v", foundArtifact, f))
			}
		}
	}

	return foundArtifact
}

func (idx *index) searchArtifactByFullname(repository, name string, version string) *artifactInfo {
	var foundArtifact *artifactInfo

	//Find the artifact
	for _, f := range *idx {
		//TODO: add support for missing format
		if repository == f.Repository && name == f.Name && version == f.Version {
			if foundArtifact == nil {
				foundArtifact = f
			} else {
				panic(fmt.Errorf("Two artifact have the same fullname: %v and %v", foundArtifact, f))
			}
		}
	}

	return foundArtifact
}

func (idx *index) searchArtifactByRepositoryAndName(repo string, name string) *artifactInfo {
	var foundArtifact *artifactInfo

	//Find the artifact
	for _, f := range *idx {
		//TODO: add support for missing format
		if name == f.Name && repo == f.Repository {
			if foundArtifact == nil {
				foundArtifact = f
			} else if cmpVersion(foundArtifact.Version, f.Version) {
				foundArtifact = f
			}
		}
	}

	return foundArtifact
}

// getOrUpdateIndex return a index of available artifacts
func readOrUpdateIndex() index {
	// Get stat about index.json
	var index index

	// Check file status
	if verbose {
		fmt.Fprintln(os.Stderr, "Trying to read index.json...")
	}
	fi, err := os.Stat(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "index.json"))
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	} else if os.IsNotExist(err) || fi.ModTime().Add(time.Duration(8)*time.Hour).Before(time.Now()) || rebuildCache {
		// Rebuild the index
		if verbose {
			fmt.Fprintln(os.Stderr, "Scanning the repositories...")
		}
		index = scanRepositories()

		// Save the index
		if verbose {
			fmt.Fprintln(os.Stderr, "Writing the index...")
		}
		file, err := os.Create(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "index.json"))
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		enc.Encode(index)
	} else {
		//Read the index
		if verbose {
			fmt.Fprintln(os.Stderr, "Read index.json...")
		}
		file, err := os.Open(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "index.json"))
		if err != nil {
			panic(err)
		}
		json.NewDecoder(file).Decode(&index)
	}

	//Sort the index
	sort.Slice(index, func(i, j int) bool {
		if index[i].Repository != index[j].Repository {
			return index[i].Repository < index[j].Repository
		} else if index[i].Name != index[j].Name {
			return index[i].Name < index[j].Name
		} else {
			return cmpVersion(index[i].Version, index[j].Version)
		}
	})
	// Set flag and handlers
	for _, art := range index {
		art.Installed = art.checkInstalled()
		setDownloader(art)
		setInstaller(art)
		setUninstaller(art)
	}

	return index
}

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage the internal repository",
	Long:  `Manage the internal repository. It requires to be run where the repository is installed`,
}

//Commands to be added
//	- list
//	- install
//	- update
//	- remove
//	- info
func init() {
	rootCmd.AddCommand(repoCmd)
	rootCmd.PersistentFlags().StringVarP(&githubToken, "github-token", "g", "", "Github token used to perform requests")
	rootCmd.PersistentFlags().BoolVar(&rebuildCache, "rebuild-cache", false, "Force the rebuild the cache")
}
