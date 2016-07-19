package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
)

var (
	buildstamp string
	githash    string

	releaseAPI  = "https://api.github.com/repos/%s/releases/%s"
	artefactAPI = "https://api.github.com/repos/%s/releases/assets/%d"
	flags       = Flags{}
)

// Flags ...
type Flags struct {
	Debug    bool
	Version  bool
	GitToken string
	Repo     string
	Release  string
	Artefact string
}

// ReleaseResponse the response from GET /repos/:owner/:repo/releases/:id
type ReleaseResponse struct {
	Assets []struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	} `json:"assets"`
}

func init() {
	flag.BoolVar(&flags.Debug, "d", false, "Debug logging")
	flag.BoolVar(&flags.Version, "v", false, "Print current version")
	flag.StringVar(&flags.GitToken, "token", "", "Git personal access token, to access private repos")
	flag.StringVar(&flags.Repo, "repo", "", "user/repo")
	flag.StringVar(&flags.Release, "release", "", "Release tag")
	flag.StringVar(&flags.Artefact, "artefact", "", "Artefacts to download, comma separated")
}

func main() {
	flag.Parse()

	if flags.Version {
		fmt.Printf("Commit hash: %s Time: %s\n", githash, buildstamp)
		os.Exit(0)
	}

	if flags.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// validate flags
	artefacts := strings.Split(flags.Artefact, ",")

	var endpoint string
	if flags.Release != "latest" {
		endpoint = fmt.Sprintf(releaseAPI, flags.Repo, "tags/"+flags.Release)
	} else {
		endpoint = fmt.Sprintf(releaseAPI, flags.Repo, flags.Release)
	}
	result, err := get(endpoint)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	releaseResponse := ReleaseResponse{}
	err = json.Unmarshal(result, &releaseResponse)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	sort.Strings(artefacts)
	for _, r := range releaseResponse.Assets {
		i := sort.SearchStrings(artefacts, r.Name)
		if i < len(artefacts) && artefacts[i] == r.Name {

			wg.Add(1)
			go func(ID uint64, name string) {
				log.Debugf("Downloading %s", name)
				err = download(ID, name, &wg)
				if err != nil {
					log.Error(err)
					os.Exit(1)
				}
			}(r.ID, r.Name)
		}
	}
	wg.Wait()
}

func get(endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if flags.GitToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", flags.GitToken))
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func download(ID uint64, name string, wg *sync.WaitGroup) error {
	defer wg.Done()
	fh, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		return err
	}
	defer fh.Close()

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf(artefactAPI, flags.Repo, ID), nil)
	if err != nil {
		return err
	}

	if flags.GitToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", flags.GitToken))
	}

	req.Header.Add("Accept", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(fh, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
