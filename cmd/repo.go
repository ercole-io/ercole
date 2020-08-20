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
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ercole-io/ercole/cmd/repo"
	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/utils"
	"github.com/google/go-github/v28/github"
	"github.com/spf13/cobra"
)

// githubToken contains the token used to avoid rate limit
var githubToken string
var rebuildCache bool

// scanGithubReleaseRepository scan a github releases repository and return detected files
func scanGithubReleaseRepository(upstreamRepo config.UpstreamRepository) (repo.Index, error) {
	//Fetch the data
	req, _ := http.NewRequest("GET", upstreamRepo.URL, nil)
	if githubToken != "" {
		req.Header.Add("Authorization", "token "+githubToken)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Received %d from github for URL %s (body: %s)", resp.StatusCode, upstreamRepo.URL, string(bytes))
	}

	if verbose {
		fmt.Printf("Fetched data from %s\n", upstreamRepo.URL)
	}

	//Decode data
	var data []github.RepositoryRelease
	json.NewDecoder(resp.Body).Decode(&data)

	//Add data to out
	var out repo.Index
	for _, release := range data {
		for _, asset := range release.Assets {
			artifactInfo := new(repo.ArtifactInfo)

			artifactInfo.Repository = upstreamRepo.Name
			artifactInfo.Filename = asset.GetName()
			artifactInfo.Version = release.GetTagName()
			artifactInfo.ReleaseDate = asset.GetUpdatedAt().Format("2006-01-02")
			artifactInfo.UpstreamType = "github-release"
			artifactInfo.UpstreamInfo = map[string]interface{}{
				"DownloadUrl": asset.GetBrowserDownloadURL(),
			}
			if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
				panic(err)
			}
			if artifactInfo.Version == "latest" {
				continue
			}
			out = append(out, artifactInfo)
		}
	}

	return out, nil
}

// scanDirectoryRepository scan the local directory and return detected files
func scanDirectoryRepository(upstreamRepo config.UpstreamRepository) (repo.Index, error) {
	//Fetch the data
	files, err := ioutil.ReadDir(upstreamRepo.URL)
	if err != nil {
		log.Fatal(err)
	}

	//Add data to out
	var out repo.Index
	for _, file := range files {
		artifactInfo := new(repo.ArtifactInfo)

		artifactInfo.Repository = upstreamRepo.Name
		artifactInfo.Filename = filepath.Base(file.Name())
		artifactInfo.ReleaseDate = file.ModTime().Format("2006-01-02")
		artifactInfo.UpstreamType = repo.UpstreamTypeDirectory
		artifactInfo.UpstreamInfo = map[string]interface{}{
			"Filename": filepath.Join(upstreamRepo.URL, file.Name()),
		}

		if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
			panic(err)
		}

		if artifactInfo.Version == "latest" {
			continue
		}

		out = append(out, artifactInfo)
	}

	return out, nil
}

// scanGithubReleaseRepository scan a github releases repository and return detected files
func scanErcoleReposerviceRepository(upstreamRepo config.UpstreamRepository) (repo.Index, error) {
	//Fetch the data
	req, _ := http.NewRequest("GET", upstreamRepo.URL+"/all", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Received %d from ercole reposervice for URL %s (body: %s)", resp.StatusCode, upstreamRepo.URL+"/all", string(bytes))
	}

	if verbose {
		fmt.Printf("Fetched data from %s\n", upstreamRepo.URL)
	}

	//Extract the filenames
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile("<a href=\"([^\"]*)\">")

	// //Add data to out
	var out repo.Index
	installedNames := make([]string, 0)

	for _, fn := range regex.FindAllStringSubmatch(string(data), -1) {
		installedNames = append(installedNames, fn[1])
	}

	//Fetch the data
	req, _ = http.NewRequest("GET", upstreamRepo.URL+"/index.json", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Received %d from ercole reposervice for URL %s (body: %s)", resp.StatusCode, upstreamRepo.URL+"/index.json", string(bytes))
	}

	if verbose {
		fmt.Printf("Fetched data from %s\n", upstreamRepo.URL)
	}

	//Decode data
	var data2 []map[string]string
	json.NewDecoder(resp.Body).Decode(&data2)

	for _, d := range data2 {
		if !utils.Contains(installedNames, d["Filename"]) {
			continue
		}

		artifactInfo := new(repo.ArtifactInfo)

		artifactInfo.Repository = upstreamRepo.Name
		artifactInfo.Filename = d["Filename"]
		artifactInfo.ReleaseDate = d["ReleaseDate"]
		artifactInfo.UpstreamType = "ercole-reposervice"
		artifactInfo.UpstreamInfo = map[string]interface{}{
			"DownloadUrl": upstreamRepo.URL + "/all/" + d["Filename"],
		}
		if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
			panic(err)
		}
		if artifactInfo.Version == "latest" {
			continue
		}
		out = append(out, artifactInfo)
	}
	return out, nil
}

// scanRepository scan a single repository and return detected files
func scanRepository(upstreamRepo config.UpstreamRepository) (repo.Index, error) {
	if strings.TrimSpace(upstreamRepo.Name) == "local" {
		return nil,
			fmt.Errorf("\"local\" isn't a valid name for an upstream repository")
	}

	switch upstreamRepo.Type {

	case repo.UpstreamTypeGitHub:
		return scanGithubReleaseRepository(upstreamRepo)

	case repo.UpstreamTypeDirectory:
		return scanDirectoryRepository(upstreamRepo)

	case repo.UpstreamTypeErcoleRepo:
		return scanErcoleReposerviceRepository(upstreamRepo)

	default:
		return nil, fmt.Errorf("Unknown repository type %q of %q", upstreamRepo.Type, upstreamRepo.Name)
	}
}

// scanRepositories scan all configured repositories and return an index
func scanRepositories() repo.Index {
	out := make(repo.Index, 0)

	for _, repo := range ercoleConfig.RepoService.UpstreamRepositories {
		//Get repository files
		files, err := scanRepository(repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "I can't get files from: %+v\n%s\n", repo, err)
			continue
		}

		out = append(out, files...)
	}

	return out
}

// readOrUpdateIndex return an index of available artifacts
func readOrUpdateIndex() repo.Index {
	var index repo.Index

	if verbose {
		fmt.Fprintln(os.Stderr, "Trying to read index.json...")
	}
	indexFile, err := os.Stat(filepath.Join(ercoleConfig.RepoService.DistributedFiles, "index.json"))
	if err != nil && !os.IsNotExist(err) {
		panic(err)

	} else if os.IsNotExist(err) || indexFile.ModTime().Add(time.Duration(8)*time.Hour).Before(time.Now()) || rebuildCache {
		if verbose {
			fmt.Fprintln(os.Stderr, "Scanning the repositories...")
		}
		index = scanRepositories()

		if verbose {
			fmt.Fprintln(os.Stderr, "Writing the index...")
		}
		index.SaveOnFile(ercoleConfig.RepoService.DistributedFiles)

	} else {
		if verbose {
			fmt.Fprintln(os.Stderr, "Read index.json...")
		}
		index = repo.ReadIndexFromFile(ercoleConfig.RepoService.DistributedFiles)
	}

	for _, art := range index {
		art.Installed = art.IsInstalled(ercoleConfig.RepoService.DistributedFiles)
	}

	index = append(index, getArtifactsNotIndexed(index, ercoleConfig.RepoService.DistributedFiles)...)

	for _, art := range index {
		art.SetDownloader(verbose)
		art.SetInstaller(verbose, ercoleConfig.RepoService.DistributedFiles)
		art.SetUninstaller(verbose, ercoleConfig.RepoService.DistributedFiles)
	}

	index.SortArtifactInfo()

	return index
}

// getArtifactsNotIndexed scan filesystem for installed artifacts not in index
func getArtifactsNotIndexed(index repo.Index, distributedFiles string) repo.Index {
	filesNotIndexed := getFilesNotIndexed(index, distributedFiles)

	artifactsNotIndexed := make(repo.Index, 0)

	for _, file := range filesNotIndexed {
		artifactInfo := new(repo.ArtifactInfo)
		artifactInfo.Filename = filepath.Base(file.Name())

		if err := artifactInfo.SetInfoFromFileName(artifactInfo.Filename); err != nil {
			fmt.Fprintf(os.Stderr, "Warning! File %v is not a supported filename\n", artifactInfo.Filename)

			artifactInfo.Name = artifactInfo.Filename
			artifactInfo.Version = "0.0.0"
			artifactInfo.Arch = "unknown"
			artifactInfo.OperatingSystemFamily = "unknown"
			artifactInfo.OperatingSystem = "unknown"
		}

		artifactInfo.Repository = repo.UpstreamTypeLocal
		artifactInfo.Installed = true
		artifactInfo.ReleaseDate = file.ModTime().Format("2006-01-02")
		artifactInfo.UpstreamType = repo.UpstreamTypeLocal

		artifactsNotIndexed = append(artifactsNotIndexed, artifactInfo)
	}

	return artifactsNotIndexed
}

func getFilesNotIndexed(index repo.Index, distributedFiles string) []os.FileInfo {
	installedInIndex := make(map[string]bool)

	allDirectory := filepath.Join(distributedFiles, "all")

	for _, artifact := range index {
		if artifact.Installed {
			installedInIndex[filepath.Join(allDirectory, artifact.Filename)] = true
		}
	}

	matches, err := filepath.Glob(allDirectory + "/*")
	if err != nil {
		panic(err)
	}

	filesNotIndexed := make([]os.FileInfo, 0)

	for _, filePath := range matches {
		if installedInIndex[filePath] {
			continue
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Something went wrong reading file: %v\n", filePath)
			continue
		}

		if fileInfo.IsDir() {
			fmt.Fprintf(os.Stderr, "Warning! Directories in /all aren't supported, but I found: %v\n", filePath)
			continue
		}

		filesNotIndexed = append(filesNotIndexed, fileInfo)
	}

	return filesNotIndexed
}

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage the internal repository",
	Long:  `Manage the internal repository. It requires to be run where the repository is installed`,
}

func init() {
	rootCmd.AddCommand(repoCmd)

	repoCmd.PersistentFlags().StringVarP(&githubToken, "github-token", "g", "", "Github token used to perform requests")
	repoCmd.PersistentFlags().BoolVar(&rebuildCache, "rebuild-cache", false, "Force the rebuild the cache")
}
