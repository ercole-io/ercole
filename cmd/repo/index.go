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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/ercole-io/ercole/utils"
)

// Index is the index of all artifact in a repository
type Index []*ArtifactInfo

// SearchArtifactByArg get *ArtifactInfo by string
func (idx *Index) SearchArtifactByArg(arg string) *ArtifactInfo {
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

	var foundArtifact *ArtifactInfo

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
		foundArtifact = idx.SearchLatestArtifactByRepositoryAndName(repository, name)
	}

	return foundArtifact
}

func (idx *Index) searchArtifactByFilename(filename string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

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

func (idx *Index) searchArtifactByName(name string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

	//Find the artifact
	for _, f := range *idx {
		//TODO: add support for missing format
		if name == f.Name {
			if foundArtifact == nil {
				foundArtifact = f
			} else if foundArtifact.Repository == f.Repository {
				if utils.IsVersionLessThan(foundArtifact.Version, f.Version) {
					foundArtifact = f
				}
			} else {
				panic(fmt.Errorf("Two artifact have the same name: %v and %v", foundArtifact, f))
			}
		}
	}

	return foundArtifact
}

func (idx *Index) searchArtifactByNameAndVersion(name string, version string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

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

func (idx *Index) searchArtifactByFullname(repository, name string, version string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

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

// SearchLatestArtifactByRepositoryAndName get latest version of a *ArtifactInfo by repo and name
func (idx *Index) SearchLatestArtifactByRepositoryAndName(repo string, name string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

	//Find the artifact
	for _, f := range *idx {
		//TODO: add support for missing format
		if name == f.Name && repo == f.Repository {
			if foundArtifact == nil {
				foundArtifact = f
			} else if utils.IsVersionLessThan(foundArtifact.Version, f.Version) {
				foundArtifact = f
			}
		}
	}

	return foundArtifact
}

// SortArtifactInfo sort artifact information inside index
func (idx Index) SortArtifactInfo() {
	sort.Slice(idx, func(i, j int) bool {
		if idx[i].Repository != idx[j].Repository {
			return idx[i].Repository < idx[j].Repository
		} else if idx[i].Name != idx[j].Name {
			return idx[i].Name < idx[j].Name
		} else {
			return utils.IsVersionLessThan(idx[i].Version, idx[j].Version)
		}
	})
}

// ReadIndexFromFile read index from the filesystem
func ReadIndexFromFile(distributedFiles string) Index {

	file, err := os.Open(filepath.Join(distributedFiles, "index.json"))
	if err != nil {
		panic(err)
	}

	var index Index
	if err = json.NewDecoder(file).Decode(&index); err != nil {
		panic(err)
	}

	return index
}

// SaveOnFile save index on the filesystem
func (idx Index) SaveOnFile(distributedFiles string) {
	file, err := os.Create(filepath.Join(distributedFiles, "index.json"))
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(idx)
}
