package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func importDataRequest(filename string, content []byte, path string) {
	client := &http.Client{}
	tr := http.DefaultTransport.(*http.Transport).Clone()

	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client.Transport = tr

	req, err := http.NewRequest("POST", ercoleConfig.DataService.RemoteEndpoint+path, bytes.NewReader(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %s", err)
		os.Exit(1)
	}

	src := ercoleConfig.DataService.AgentUsername + ":" + ercoleConfig.DataService.AgentPassword
	bearer := "Basic " + base64.StdEncoding.EncodeToString([]byte(src))
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send data from %s: %v\n", filename, err)
		os.Exit(1)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		fmt.Fprintf(os.Stderr, "File: %s Status: %d Cause: %s\n", filename, resp.StatusCode, string(out))
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("File: %s Status: %d\n", filename, resp.StatusCode)
	}
}
