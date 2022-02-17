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
	"github.com/google/go-github/v28/github"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

// Index is the index of all artifact in a repository
type Index struct {
	config    config.Configuration
	log       logger.Logger
	artifacts []*ArtifactInfo
}

func readOrUpdateIndex(log logger.Logger) Index {
	index := Index{
		config: *ercoleConfig,
		log:    log,
	}

	log.Debug("Trying to read index.json...")

	indexFile, err := os.Stat(filepath.Join(index.distributedFiles(), "index.json"))
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	if os.IsNotExist(err) || indexFile.ModTime().Add(time.Duration(8)*time.Hour).Before(time.Now()) || rebuildCache {
		index.log.Debug("Scanning the repositories...")
		index.getArtifactsFromUpstreamRepositories()

		index.log.Debug("Writing the index...")
	} else {
		log.Debug("Read index.json...")
		if index.artifacts, err = readArtifactsFromFile(index.distributedFiles()); err != nil {
			log.Fatalf("Can't read artifacts from file: %s", err)
		}
	}

	index.checkInstalledArtifacts()
	index.getLocalArtifacts()
	index.saveArtifactsToFile()
	index.sortArtifactsInfo()

	return index
}

func (idx *Index) distributedFiles() string {
	return idx.config.RepoService.DistributedFiles
}

func (idx *Index) getArtifactsFromUpstreamRepositories() {
	for _, repo := range idx.config.RepoService.UpstreamRepositories {
		var err error

		switch repo.Type {
		case UpstreamTypeGitHub:
			err = idx.getArtifactsFromGithub(repo)

		case UpstreamTypeDirectory:
			err = idx.getArtifactsFromDirectory(repo)

		case UpstreamTypeErcoleRepo:
			err = idx.getArtifactsFromErcoleReposervice(repo)

		default:
			err = fmt.Errorf("Unknown repository type %q of %q, skipped", repo.Type, repo.Name)
		}

		if err != nil {
			idx.log.Warnf("%+v\n%s\n", repo, err)
		}
	}
}

func (idx *Index) getArtifactsFromGithub(upstreamRepo config.UpstreamRepository) error {
	req, err := http.NewRequest("GET", upstreamRepo.URL, nil)
	if err != nil {
		return err
	}

	if githubToken != "" {
		req.Header.Add("Authorization", "token "+githubToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Received %d from github for URL %s (body: %s)", resp.StatusCode, upstreamRepo.URL, string(bytes))
	}

	idx.log.Debugf("Fetched data from %s\n", upstreamRepo.URL)

	var releases []github.RepositoryRelease
	if err = json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return err
	}

	var artifacts []*ArtifactInfo

	for _, release := range releases {
		for _, asset := range release.Assets {
			artifactInfo := new(ArtifactInfo)

			artifactInfo.Repository = upstreamRepo.Name
			artifactInfo.Filename = asset.GetName()
			artifactInfo.Version = release.GetTagName()
			artifactInfo.ReleaseDate = asset.GetUpdatedAt().Format("2006-01-02")
			artifactInfo.UpstreamRepository = upstreamRepository{
				Type:        UpstreamTypeGitHub,
				DownloadUrl: asset.GetBrowserDownloadURL(),
			}

			if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
				idx.log.Warn(err)
				continue
			}

			if artifactInfo.Version == "latest" {
				idx.log.Debugf("Ignore latest %+v", artifactInfo.Filename)
				continue
			}

			artifacts = append(artifacts, artifactInfo)
		}
	}

	idx.artifacts = append(idx.artifacts, artifacts...)

	return nil
}

func (idx *Index) getArtifactsFromDirectory(upstreamRepo config.UpstreamRepository) error {
	files, err := ioutil.ReadDir(upstreamRepo.URL)
	if err != nil {
		return err
	}

	var artifacts []*ArtifactInfo

	for _, file := range files {
		artifactInfo := new(ArtifactInfo)

		artifactInfo.Repository = upstreamRepo.Name
		artifactInfo.Filename = filepath.Base(file.Name())
		artifactInfo.ReleaseDate = file.ModTime().Format("2006-01-02")
		artifactInfo.UpstreamRepository = upstreamRepository{
			Type:     UpstreamTypeDirectory,
			Filepath: filepath.Join(upstreamRepo.URL, file.Name()),
		}

		if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
			idx.log.Warn(err)
		}

		if artifactInfo.Version == "latest" {
			idx.log.Debug("Ignore latest %+v", artifactInfo)
			continue
		}

		artifacts = append(artifacts, artifactInfo)
	}

	idx.artifacts = append(idx.artifacts, artifacts...)

	return nil
}

func (idx *Index) getArtifactsFromErcoleReposervice(upstreamRepo config.UpstreamRepository) error {
	req, err := http.NewRequest("GET", upstreamRepo.URL+"/all", nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("Received %d from ercole reposervice for URL %s (body: %s)", resp.StatusCode, upstreamRepo.URL+"/all", string(bytes))
	}

	idx.log.Debugf("Fetched data from %s\n", upstreamRepo.URL)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	regex := regexp.MustCompile(`<a href="([^"]*)">`)

	var artifacts []*ArtifactInfo

	installedNames := make([]string, 0)

	for _, fn := range regex.FindAllStringSubmatch(string(data), -1) {
		installedNames = append(installedNames, fn[1])
	}

	for _, file := range installedNames {
		artifactInfo := new(ArtifactInfo)

		artifactInfo.Repository = upstreamRepo.Name
		artifactInfo.Filename = file
		artifactInfo.ReleaseDate = "????-??-??"
		artifactInfo.UpstreamRepository = upstreamRepository{
			Type:        UpstreamTypeErcoleRepo,
			DownloadUrl: upstreamRepo.URL + "/all/" + file,
		}

		if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
			idx.log.Warn(err)
			continue
		}

		if artifactInfo.Version == "latest" {
			idx.log.Debug("Ignore latest %+v", artifactInfo)
			continue
		}

		artifacts = append(artifacts, artifactInfo)
	}

	idx.artifacts = append(idx.artifacts, artifacts...)

	return nil
}

// getLocalArtifacts scan filesystem for installed artifacts not in index
func (idx *Index) getLocalArtifacts() {
	localFiles := getLocalFiles(idx.log, idx.artifacts, idx.distributedFiles())

	localArtifacts := make([]*ArtifactInfo, 0)

	for _, file := range localFiles {
		artifactInfo := new(ArtifactInfo)
		artifactInfo.Filename = filepath.Base(file.Name())

		if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
			idx.log.Warnf("File %q not supported: %s\n", artifactInfo.Filename, err)

			artifactInfo.Name = artifactInfo.Filename
			artifactInfo.Version = "0.0.0"
			artifactInfo.Arch = "unknown"
			artifactInfo.OperatingSystemFamily = "unknown"
			artifactInfo.OperatingSystem = "unknown"
		}

		artifactInfo.Repository = UpstreamTypeLocal
		artifactInfo.Installed = true
		artifactInfo.ReleaseDate = file.ModTime().Format("2006-01-02")
		artifactInfo.UpstreamRepository = upstreamRepository{
			Type: UpstreamTypeLocal,
		}

		localArtifacts = append(localArtifacts, artifactInfo)
	}

	idx.artifacts = append(idx.artifacts, localArtifacts...)
}

func getLocalFiles(log logger.Logger, index []*ArtifactInfo, distributedFiles string) []os.FileInfo {
	installedInIndex := make(map[string]bool)

	allDirectory := filepath.Join(distributedFiles, "all")

	for _, artifact := range index {
		if artifact.Installed {
			installedInIndex[filepath.Join(allDirectory, artifact.Filename)] = true
		}
	}

	matches, err := filepath.Glob(allDirectory + "/*")
	if err != nil {
		log.Fatal(err)
	}

	filesNotIndexed := make([]os.FileInfo, 0)

	for _, filePath := range matches {
		if installedInIndex[filePath] {
			continue
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Warnf("Something went wrong reading file: %v\n", filePath)
			continue
		}

		if fileInfo.IsDir() {
			log.Warnf("Warning! Directories in /all aren't supported, but I found: %v\n", filePath)
			continue
		}

		filesNotIndexed = append(filesNotIndexed, fileInfo)
	}

	return filesNotIndexed
}

// searchArtifactByArg get *ArtifactInfo by string
func (idx *Index) searchArtifactByArg(arg string) *ArtifactInfo {
	submatches := utils.FindNamedMatches(artifactNameRegex, arg)

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
		foundArtifact = idx.searchLatestArtifactByRepositoryAndName(repository, name)
	}

	return foundArtifact
}

func (idx *Index) searchArtifactByFilename(filename string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

	for _, f := range idx.artifacts {
		if filename != f.Filename {
			continue
		}

		if foundArtifact != nil {
			idx.log.Fatalf("Two artifact have the same filename: %v and %v", foundArtifact, f)
		}

		foundArtifact = f
	}

	return foundArtifact
}

func (idx *Index) searchArtifactByName(name string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

	for _, f := range idx.artifacts {
		if name != f.Name {
			continue
		}

		if foundArtifact == nil {
			foundArtifact = f
			continue
		}

		if foundArtifact.Repository != f.Repository {
			idx.log.Fatalf("Two artifact have the same filename: %v and %v", foundArtifact, f)
		}

		foundVersionIsLess, err := utils.IsVersionLessThan(foundArtifact.Version, f.Version)
		if err != nil {
			idx.log.Warnf("Invalid version comparing %q with %q", foundArtifact.Version, f.Version)
			continue
		}

		if foundVersionIsLess {
			foundArtifact = f
		}
	}

	return foundArtifact
}

func (idx *Index) searchArtifactByNameAndVersion(name string, version string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

	for _, f := range idx.artifacts {
		if name != f.Name || version != f.Version {
			continue
		}

		if foundArtifact != nil {
			idx.log.Fatalf("Two artifact have the same filename: %v and %v", foundArtifact, f)
		}

		foundArtifact = f
	}

	return foundArtifact
}

func (idx *Index) searchArtifactByFullname(repository, name, version string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

	for _, f := range idx.artifacts {
		if repository != f.Repository || name != f.Name {
			continue
		}

		if eq, err := utils.IsVersionEqual(version, f.Version); err != nil {
			idx.log.Warnf("Invalid version comparing %q with %q", version, f.Version)
			continue
		} else if !eq {
			continue
		}

		if foundArtifact != nil {
			idx.log.Fatalf("Two artifact have the same filename: %v and %v", foundArtifact, f)
		}

		foundArtifact = f
	}

	return foundArtifact
}

// searchLatestArtifactByRepositoryAndName get latest version of a *ArtifactInfo by repo and name
func (idx *Index) searchLatestArtifactByRepositoryAndName(repo string, name string) *ArtifactInfo {
	var foundArtifact *ArtifactInfo

	for _, f := range idx.artifacts {
		if name != f.Name || repo != f.Repository {
			continue
		}

		if foundArtifact == nil {
			foundArtifact = f
			continue
		}

		foundVersionIsLess, err := utils.IsVersionLessThan(foundArtifact.Version, f.Version)
		if err != nil {
			idx.log.Warnf("Invalid version comparing %q with %q", foundArtifact.Version, f.Version)
			continue
		}

		if foundVersionIsLess {
			foundArtifact = f
		}
	}

	return foundArtifact
}

// sortArtifactsInfo sort artifact information inside index
func (idx Index) sortArtifactsInfo() {
	artifacts := idx.artifacts

	sort.Slice(artifacts, func(i, j int) bool {
		if artifacts[i].Repository != artifacts[j].Repository {
			return artifacts[i].Repository < artifacts[j].Repository
		} else if artifacts[i].Name != artifacts[j].Name {
			return artifacts[i].Name < artifacts[j].Name
		} else {
			is, err := utils.IsVersionLessThan(artifacts[i].Version, artifacts[j].Version)
			if err != nil {
				idx.log.Warnf("Invalid version comparing %q with %q", artifacts[i].Version, artifacts[j].Version)
				return false
			}

			return is
		}
	})
}

func readArtifactsFromFile(distributedFiles string) ([]*ArtifactInfo, error) {
	file, err := os.Open(filepath.Join(distributedFiles, "index.json"))
	if err != nil {
		return nil, err
	}

	var artifacts []*ArtifactInfo
	if err = json.NewDecoder(file).Decode(&artifacts); err != nil {
		return nil, err
	}

	return artifacts, nil
}

func (idx *Index) checkInstalledArtifacts() {
	for _, art := range idx.artifacts {
		art.checkIsInstalled(idx.distributedFiles())
	}
}

func (idx *Index) saveArtifactsToFile() {
	file, err := os.Create(filepath.Join(idx.distributedFiles(), "index.json"))
	if err != nil {
		log.Fatalf("Can't create index.json file: %s", err)
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	if err := enc.Encode(idx.artifacts); err != nil {
		log.Fatalf("Can't encode artifacts: %s", err)
	}
}

func (idx *Index) Install(artifact *ArtifactInfo) {
	idx.log.Debugf("Creating the directories (if missing) %s, %s\n",
		artifact.DirectoryPath(idx.distributedFiles()),
		filepath.Join(idx.distributedFiles(), "all"),
	)

	if err := os.MkdirAll(artifact.DirectoryPath(idx.distributedFiles()), 0755); err != nil {
		idx.log.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(idx.distributedFiles(), "all"), 0755); err != nil {
		idx.log.Fatal(err)
	}

	idx.log.Debugf("Downloading the artifact %s to %s\n", artifact.Filename, artifact.FilePath(idx.distributedFiles()))

	if err := idx.Download(artifact); err != nil {
		log.Fatalf("Unable to download artifact: %s", err)
	}

	idx.log.Debugf("Linking the artifact to %s\n", filepath.Join(idx.distributedFiles(), "all", artifact.Filename))

	if err := os.Link(artifact.FilePath(idx.distributedFiles()), filepath.Join(idx.distributedFiles(), "all", artifact.Filename)); err != nil {
		log.Fatalf("Unable to link artifact to \"all\" folder: %s", err)
	}

	if strings.HasSuffix(artifact.Filename, ".rpm") {
		idx.log.Debugf("Executing createrepo %s\n", artifact.DirectoryPath(idx.distributedFiles()))
		dirPath := artifact.DirectoryPath(idx.distributedFiles())
		cmd := exec.Command("createrepo", dirPath)

		if artifact.OperatingSystem == "rhel5" {
			cmd = exec.Command("createrepo", "-s", "sha", dirPath)
		}

		if verbose {
			cmd.Stdout = os.Stdout
		}

		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Error running createrepo: %s\n", err.Error())
		}
	}

	artifact.Installed = true
	idx.log.Infof("Installed %q", artifact.FullName())
}

func (idx *Index) Remove(artifact *ArtifactInfo) {
	idx.log.Debugf("Removing the file %s\n", filepath.Join(idx.distributedFiles(), "all", artifact.Filename))

	allArtifactsPath := filepath.Join(idx.distributedFiles(), "all", artifact.Filename)
	if err := os.Remove(allArtifactsPath); err != nil {
		idx.log.Fatalf("Can't remove %s from %s: %s", artifact.Filename, allArtifactsPath, err)
	}

	artifactPath := artifact.FilePath(idx.distributedFiles())
	idx.log.Debugf("Removing the file %s\n", artifactPath)

	if err := os.Remove(artifactPath); err != nil {
		idx.log.Errorf("Can't remove %s from %s: %s", artifact.Filename, artifactPath, err)
		return
	}

	if strings.HasSuffix(artifact.Filename, ".rpm") {
		idx.log.Debugf("Executing createrepo %s\n", artifact.DirectoryPath(idx.distributedFiles()))
		cmd := exec.Command("createrepo", artifact.DirectoryPath(idx.distributedFiles())) //TODO Refactor

		if verbose {
			cmd.Stdout = os.Stdout
		}

		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	artifact.Installed = false
	idx.log.Infof("Removed %q", artifact.FullName())
}

func (idx *Index) Download(artifact *ArtifactInfo) error {
	dest := artifact.FilePath(idx.distributedFiles())

	switch artifact.UpstreamRepository.Type {
	case UpstreamTypeGitHub:
		if err := utils.DownloadFile(dest, artifact.UpstreamRepository.DownloadUrl); err != nil {
			return err
		}

	case UpstreamTypeDirectory:
		if verbose {
			fmt.Printf("Copying file from %s to %s\n", artifact.UpstreamRepository.Filepath, dest)
		}

		err := yos.CopyFile(artifact.UpstreamRepository.Filepath, dest)
		if err != nil {
			panic(err)
		}

		err = os.Chmod(dest, 0755)
		if err != nil {
			panic(err)
		}

	case UpstreamTypeErcoleRepo:
		if err := utils.DownloadFile(dest, artifact.UpstreamRepository.DownloadUrl); err != nil {
			return err
		}

	case UpstreamTypeLocal:
		fmt.Println("Nothing to do, artifact already installed")

	default:
		return fmt.Errorf("Unknown UpstreamType: %s", artifact.UpstreamRepository.Type)
	}

	return nil
}
