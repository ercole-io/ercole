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
	Installed             bool
	Version               string
	ReleaseDate           string
	Filename              string
	Name                  string
	OperatingSystemFamily string
	OperatingSystem       string
	Arch                  string
	UpstreamRepository    upstreamRepository
}

type upstreamRepository struct {
	Type        string
	DownloadUrl string
	Filepath    string
}

//valid formats
//	- <filename>
//	- <name>
//	- <name>@<version>
//	- <repository>/<name>@<version>
//	- <repository>/<name>
var artifactNameRegex *regexp.Regexp = regexp.MustCompile(`^(?:(?P<repository>[a-z-0-9]+)/)?(?P<name>[a-z-.0-9]+)(?:@(?P<version>[a-z0-9.0-9]+))?$`)

// https://semver.org/
var semVer = `(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?`

//Regex for filenames
var (
	ercoleRHELRegex    *regexp.Regexp = regexp.MustCompile(`^ercole-(?P<version>` + semVer + `)-1\.el(?P<dist>\d+)\.(?P<arch>x86_64)\.rpm$`)
	ercoleWebRHELRegex *regexp.Regexp = regexp.MustCompile(`^ercole-web-(?P<version>` + semVer + `)-1\.el(?P<dist>\d+)\.(?P<arch>noarch)\.rpm$`)

	agentRHELRegex      *regexp.Regexp = regexp.MustCompile(`^ercole-agent-(?P<version>` + semVer + `)-1\.el(?P<dist>\d+)\.(?P<arch>x86_64).rpm$`)
	agentWinRegex       *regexp.Regexp = regexp.MustCompile(`^ercole-agent-setup-(?P<version>` + semVer + `)\.exe$`)
	agentPerlRpmRegex   *regexp.Regexp = regexp.MustCompile(`^ercole-agent-perl-(?P<version>` + semVer + `)-1\.(?P<dist>[a-z0-9.]+)\.(?P<arch>noarch)\.rpm$`)
	agentPerlTarGzRegex *regexp.Regexp = regexp.MustCompile(`^ercole-agent-perl-(?P<version>` + semVer + `)-1\.(?P<dist>[a-z0-9.]+)\.(?P<arch>noarch)\.tar\.gz$`)
)

const (
	UpstreamTypeLocal      = "local"
	UpstreamTypeDirectory  = "directory"
	UpstreamTypeGitHub     = "github-release"
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

// checkIsInstalled return true if file is detected in the distribution directory
func (artifact *ArtifactInfo) checkIsInstalled(distributedFiles string) {
	_, err := os.Stat(filepath.Join(distributedFiles, "all", artifact.Filename))

	artifact.Installed = !os.IsNotExist(err)
}

// SetInfoFromFileName sets to fileInfo informations taken from filename
func (artifact *ArtifactInfo) SetInfoFromFileName(filename string) error {
	switch {
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

	case agentRHELRegex.MatchString(filename):
		data := utils.FindNamedMatches(agentRHELRegex, filename)
		artifact.Name = "ercole-agent-rhel" + data["dist"]
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
		case "solaris10", "solaris11":
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
