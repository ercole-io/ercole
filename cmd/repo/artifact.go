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
	"path/filepath"
	"regexp"

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
	ercoleRHELRegex    *regexp.Regexp = regexp.MustCompile(`^ercole-(?P<version>.*)-1\.el(?P<dist>\d+)\.(?P<arch>x86_64)\.rpm$`)
	ercoleWebRHELRegex *regexp.Regexp = regexp.MustCompile(`^ercole-web-(?P<version>.*)-1\.el(?P<dist>\d+)\.(?P<arch>noarch)\.rpm$`)

	agentRHELRegex      *regexp.Regexp = regexp.MustCompile(`^ercole-agent-(?P<version>.*)-1\.el(?P<dist>\d+)\.(?P<arch>x86_64).rpm$`)
	agentWinRegex       *regexp.Regexp = regexp.MustCompile(`^ercole-agent-setup-(?P<version>.*)\.exe$`)
	agentPerlRpmRegex   *regexp.Regexp = regexp.MustCompile(`^ercole-agent-perl-(?P<version>.*)-1\.(?P<dist>[a-z0-9.]+)\.(?P<arch>noarch)\.rpm$`)
	agentPerlTarGzRegex *regexp.Regexp = regexp.MustCompile(`^ercole-agent-perl-(?P<version>.*)-1\.(?P<dist>[a-z0-9.]+)\.(?P<arch>noarch)\.tar\.gz$`)
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
	case agentRHELRegex.MatchString(filename):
		data := utils.FindNamedMatches(agentRHELRegex, filename)
		artifact.Name = "ercole-agent-rhel" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case ercoleRHELRegex.MatchString(filename):
		data := utils.FindNamedMatches(ercoleRHELRegex, filename)
		artifact.Name = "ercole-" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case ercoleWebRHELRegex.MatchString(filename):
		data := utils.FindNamedMatches(ercoleWebRHELRegex, filename)
		artifact.Name = "ercole-web" + data["dist"]
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		artifact.OperatingSystemFamily = "rhel"
		artifact.OperatingSystem = "rhel" + data["dist"]

	case agentWinRegex.MatchString(filename):
		data := utils.FindNamedMatches(agentWinRegex, filename)
		artifact.Name = "ercole-agent-win"
		artifact.Version = data["version"]
		artifact.Arch = "x86_64"
		artifact.OperatingSystemFamily = "win"
		artifact.OperatingSystem = "win"

	case agentPerlRpmRegex.MatchString(filename):
		data := utils.FindNamedMatches(agentPerlRpmRegex, filename)
		artifact.Name = "ercole-agent-perl-" + data["dist"] + "-rpm"
		artifact.Version = data["version"]
		artifact.Arch = data["arch"]
		switch data["dist"] {
		case "aix6.1":
			artifact.OperatingSystemFamily = "aix"
			artifact.OperatingSystem = data["dist"]
		default:
			return fmt.Errorf("Unknown distribution %s", data["dist"])
		}

	case agentPerlTarGzRegex.MatchString(filename):
		data := utils.FindNamedMatches(agentPerlTarGzRegex, filename)
		artifact.Name = "ercole-agent-perl-" + data["dist"] + "-tar-gz"
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
		return fmt.Errorf("Filename %q is not supported, skipped", filename)
	}

	return nil
}
