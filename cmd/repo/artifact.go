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

package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/1set/gut/yos"
	"github.com/ercole-io/ercole/v2/utils"
)

// ArtifactInfo contains info about all files in repository
type ArtifactInfo struct {
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
}

//Regex for filenames
var (
	agentRHEL5Regex              *regexp.Regexp = regexp.MustCompile(`^ercole-agent-(?P<version>.*)-1.(?P<arch>x86_64).rpm$`)
	agentRHELRegex               *regexp.Regexp = regexp.MustCompile(`^ercole-agent-(?P<version>.*)-1.el(?P<dist>\d+).(?P<arch>x86_64).rpm$`)
	agentVirtualizationRHELRegex *regexp.Regexp = regexp.MustCompile(`^ercole-agent-virtualization-(?P<version>.*)-1.el(?P<dist>\d+).(?P<arch>x86_64).rpm$`)
	agentExadataRHELRegex        *regexp.Regexp = regexp.MustCompile(`^ercole-agent-exadata-(?P<version>.*)-1.el(?P<dist>\d+).(?P<arch>x86_64).rpm$`)
	agentWinRegex                *regexp.Regexp = regexp.MustCompile(`^ercole-agent-setup-(?P<version>.*).exe$`)
	agentHpuxRegex               *regexp.Regexp = regexp.MustCompile(`^ercole-agent-hpux-(?P<version>.*).tar.gz`)
	agentAixRegexRpm             *regexp.Regexp = regexp.MustCompile(`^ercole-agent-aix-(?P<version>.*)-1.(?P<dist>.*).(?P<arch>noarch).rpm$`)
	agentAixRegexTarGz           *regexp.Regexp = regexp.MustCompile(`^ercole-agent-aix-(?P<version>.*).tar.gz$`)
	ercoleRHELRegex              *regexp.Regexp = regexp.MustCompile(`^ercole-(?P<version>.*)-1.el(?P<dist>\d+).(?P<arch>x86_64).rpm$`)
	ercoleWebRHELRegex           *regexp.Regexp = regexp.MustCompile(`^ercole-web-(?P<version>.*)-1.el(?P<dist>\d+).(?P<arch>noarch).rpm$`)
	ercoleAgentPerlRpmRegex      *regexp.Regexp = regexp.MustCompile(`^ercole-agent-perl-(?P<version>.*)-1.(?P<dist>[a-z0-9.]+).(?P<arch>noarch).rpm$`)
	ercoleAgentPerlTarGzRegex    *regexp.Regexp = regexp.MustCompile(`^ercole-agent-perl-(?P<version>.*)-1.(?P<dist>[a-z0-9.]+).(?P<arch>noarch).tar.gz$`)
)

const (
	// UpstreamTypeLocal repository upstream type
	UpstreamTypeLocal = "local"
	// UpstreamTypeDirectory repository upstream type
	UpstreamTypeDirectory = "directory"
	// UpstreamTypeGitHub repository upstream type
	UpstreamTypeGitHub = "github-release"
	// UpstreamTypeErcoleRepo repository upstream type
	UpstreamTypeErcoleRepo = "ercole-reposervice"
)

// FullName return the fullname of the file
func (artifact *ArtifactInfo) FullName() string {
	return fmt.Sprintf("%s/%s@%s", artifact.Repository, artifact.Name, artifact.Version)
}

// DirectoryPath get the path of the directory containing the artifact file
func (artifact *ArtifactInfo) DirectoryPath(distributedFiles string) string {
	return filepath.Join(distributedFiles, artifact.OperatingSystemFamily, artifact.OperatingSystem, artifact.Arch)
}

//FilePath get the path of the artifact file
func (artifact *ArtifactInfo) FilePath(distributedFiles string) string {
	return filepath.Join(artifact.DirectoryPath(distributedFiles), artifact.Filename)
}

// IsInstalled return true if file is detected in the distribution directory
func (artifact *ArtifactInfo) IsInstalled(distributedFiles string) bool {
	_, err := os.Stat(filepath.Join(distributedFiles, "all", artifact.Filename))

	return !os.IsNotExist(err)
}

// SetInfoFromFileName sets to fileInfo informations taken from filename
func (artifact *ArtifactInfo) SetInfoFromFileName(filename string) error {
	switch {
	case agentVirtualizationRHELRegex.MatchString(filename): //agent virtualization RHEL
		data := utils.FindNamedMatches(agentVirtualizationRHELRegex, filename)
		artifact.Name = "ercole-agent-virtualization-rhel" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case agentExadataRHELRegex.MatchString(filename): //agent exadata RHEL
		data := utils.FindNamedMatches(agentExadataRHELRegex, filename)
		artifact.Name = "ercole-agent-exadata-rhel" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case agentRHEL5Regex.MatchString(filename): //agent RHEL5
		data := utils.FindNamedMatches(agentRHEL5Regex, filename)
		artifact.Name = "ercole-agent-rhel5"
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel5"

	case agentRHELRegex.MatchString(filename): //agent RHEL
		data := utils.FindNamedMatches(agentRHELRegex, filename)
		artifact.Name = "ercole-agent-rhel" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case ercoleRHELRegex.MatchString(filename): //ercole RHEL
		data := utils.FindNamedMatches(ercoleRHELRegex, filename)
		artifact.Name = "ercole-" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case ercoleWebRHELRegex.MatchString(filename): //ercole-web RHEL
		data := utils.FindNamedMatches(ercoleWebRHELRegex, filename)
		artifact.Name = "ercole-web" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case agentWinRegex.MatchString(filename): //agent WIN
		data := utils.FindNamedMatches(agentWinRegex, filename)
		artifact.Name = "ercole-agent-win"
		artifact.Version = data["version"]
		artifact.Arch = "x86_64"
		artifact.OperatingSystemFamily = "win"
		artifact.OperatingSystem = "win"

	case agentHpuxRegex.MatchString(filename): //agent HPUX
		data := utils.FindNamedMatches(agentHpuxRegex, filename)
		artifact.Name = "ercole-agent-hpux"
		artifact.Version = data["version"]
		artifact.Arch = "noarch"
		artifact.OperatingSystemFamily = "hpux"
		artifact.OperatingSystem = "hpux"

	case agentAixRegexRpm.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(agentAixRegexRpm, filename)
		artifact.Name = "ercole-agent-aix"
		artifact.Version = data["version"]
		artifact.Arch = "noarch"
		artifact.OperatingSystemFamily = "aix"
		artifact.OperatingSystem = data["dist"]

	case agentAixRegexTarGz.MatchString(filename): //agent AIX
		data := utils.FindNamedMatches(agentAixRegexTarGz, filename)
		artifact.Name = "ercole-agent-aix-targz"
		artifact.Version = data["version"]
		artifact.Arch = "noarch"
		artifact.OperatingSystemFamily = "aix-tar-gz"
		artifact.OperatingSystem = "aix6.1"

	case ercoleAgentPerlRpmRegex.MatchString(filename): //agent perl rpm
		data := utils.FindNamedMatches(ercoleAgentPerlRpmRegex, filename)
		artifact.Name = "ercole-agent-" + data["dist"] + "-rpm"
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		switch data["dist"] {
		case "aix6.1":
			artifact.OperatingSystemFamily = "aix"
			artifact.OperatingSystem = data["dist"]
		default:
			return fmt.Errorf("Unknown distribution %s", data["dist"])
		}

	case ercoleAgentPerlTarGzRegex.MatchString(filename): //agent perl tar gz
		data := utils.FindNamedMatches(ercoleAgentPerlTarGzRegex, filename)
		artifact.Name = "ercole-agent-" + data["dist"] + "-tar-gz"
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		switch data["dist"] {
		case "aix6.1":
			artifact.OperatingSystemFamily = "aix"
			artifact.OperatingSystem = data["dist"]
		case "hpux":
			artifact.OperatingSystemFamily = "hpux"
			artifact.OperatingSystem = data["dist"]
		case "solaris11":
			artifact.OperatingSystemFamily = "solaris"
			artifact.OperatingSystem = data["dist"]
		default:
			return fmt.Errorf("Unknown distribution %s", data["dist"])
		}

	default:
		return fmt.Errorf("Filename %s is not supported. Please check that it is correct", filename)
	}

	return nil
}

// Download the artifact
func (artifact *ArtifactInfo) Download(verbose bool, dest string) {
	switch artifact.UpstreamType {
	case UpstreamTypeGitHub:
		utils.DownloadFile(dest, artifact.UpstreamInfo["DownloadUrl"].(string))

	case UpstreamTypeDirectory:
		if verbose {
			fmt.Printf("Copying file from %s to %s\n", artifact.UpstreamInfo["Filename"].(string), dest)
		}
		err := yos.CopyFile(artifact.UpstreamInfo["Filename"].(string), dest)
		if err != nil {
			panic(err)
		}
		err = os.Chmod(dest, 0755)
		if err != nil {
			panic(err)
		}

	case UpstreamTypeErcoleRepo:
		utils.DownloadFile(dest, artifact.UpstreamInfo["DownloadUrl"].(string))

	case UpstreamTypeLocal:
		fmt.Println("Nothing to do, artifact already installed")

	default:
		panic(artifact)
	}
}

// Install the artifact
func (artifact *ArtifactInfo) Install(verbose bool, distributedFiles string) {
	//Create missing directories
	if verbose {
		fmt.Printf("Creating the directories (if missing) %s, %s\n",
			artifact.DirectoryPath(distributedFiles),
			filepath.Join(distributedFiles, "all"),
		)
	}
	err := os.MkdirAll(artifact.DirectoryPath(distributedFiles), 0755)
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(filepath.Join(distributedFiles, "all"), 0755)
	if err != nil {
		panic(err)
	}

	//Download the file in the right location
	if verbose {
		fmt.Printf("Downloading the artifact %s to %s\n", artifact.Filename, artifact.FilePath(distributedFiles))
	}
	artifact.Download(verbose, artifact.FilePath(distributedFiles))

	//Create a link to all
	if verbose {
		fmt.Printf("Linking the artifact to %s\n", filepath.Join(distributedFiles, "all", artifact.Filename))
	}
	err = os.Link(artifact.FilePath(distributedFiles), filepath.Join(distributedFiles, "all", artifact.Filename))
	if err != nil {
		panic(err)
	}

	if strings.HasSuffix(artifact.Filename, ".rpm") {
		//Launch the createrepo command
		if verbose {
			fmt.Printf("Executing createrepo %s\n", artifact.DirectoryPath(distributedFiles))
		}
		cmd := exec.Command("createrepo", artifact.DirectoryPath(distributedFiles))
		if verbose {
			cmd.Stdout = os.Stdout
		}
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running createrepo: %s\n", err.Error())
		}
	}

	artifact.Installed = true
}

// Uninstall the artifact
func (artifact *ArtifactInfo) Uninstall(verbose bool, distributedFiles string) {
	if verbose {
		fmt.Printf("Removing the file %s\n", filepath.Join(distributedFiles, "all", artifact.Filename))
	}
	if err := os.Remove(filepath.Join(distributedFiles, "all", artifact.Filename)); err != nil {
		panic(err)
	}

	if _, errStat := os.Stat(artifact.FilePath(distributedFiles)); errStat == nil {
		if verbose {
			fmt.Printf("Removing the file %s\n", artifact.FilePath(distributedFiles))
		}
		if err := os.Remove(artifact.FilePath(distributedFiles)); err != nil {
			panic(err)
		}

		if strings.HasSuffix(artifact.Filename, ".rpm") {

			if verbose {
				fmt.Printf("Executing createrepo %s\n", artifact.DirectoryPath(distributedFiles))
			}
			cmd := exec.Command("createrepo", artifact.DirectoryPath(distributedFiles))
			if verbose {
				cmd.Stdout = os.Stdout
			}
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				panic(err)
			}
		}

		artifact.Installed = false
	}
}
