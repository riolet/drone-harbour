package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/riolet/drone-go/drone"
	"github.com/riolet/drone-go/plugin"
)

const (
	respFormat      = "Webhook %d\n  URL: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
	debugRespFormat = "Webhook %d\n  URL: %s\n  METHOD: %s\n  HEADERS: %s\n  REQUEST BODY: %s\n  RESPONSE STATUS: %s\n  RESPONSE BODY: %s\n"
)

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone Webhook Plugin built from %s\n", buildCommit)

	system := drone.System{}
	repo := drone.Repo{}
	build := drone.Build{}
	vargs := Params{}

	plugin.Param("system", &system)
	plugin.Param("repo", &repo)
	plugin.Param("build", &build)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if vargs.Method == "" {
		vargs.Method = "POST"
	}

	if vargs.ContentType == "" {
		vargs.ContentType = "application/json"
	}

	// Creates the payload, by default the payload
	// is the build details in json format, but a custom
	// template may also be used.

	var buf bytes.Buffer


	data := struct {
		System 		drone.System 		`json:"system"`
		Repo   		drone.Repo   		`json:"repo"`
		Build  		drone.Build  		`json:"build"`
		Registry  	string 			`json:"registry"`
		Image 		string 			`json:"image"`
		Name  		string 			`json:"name"`
		Tag 		string 			`json:"tag"`
		Ports 		[]int 			`json:"ports"`
		PortBindings 	map[string]string 	`json:"port_bindings"`
		Env 		[]string 		`json:"env"`
		Links 		map[string]string 	`json:"links"`
		Volumes 	map[string]string 	`json:"volumes"`
		PublishAllPorts bool `json:"publish_all_ports"`
	}{system, repo, build, vargs.Registry, vargs.Repo, vargs.Name, vargs.Tag,
		vargs.Ports, vargs.PortBindings, vargs.Env,
		vargs.Links, vargs.Volumes, vargs.PublishAllPorts}

	if err := json.NewEncoder(&buf).Encode(&data); err != nil {
		fmt.Printf("Error: Failed to encode JSON payload. %s\n", err)
		os.Exit(1)
	}


	// build and execute a request for each url.
	// all auth, headers, method, template (payload),
	// and content_type values will be applied to
	// every webhook request.

	for i, rawurl := range vargs.URLs {
		uri, err := url.Parse(rawurl+"/drone-harbour-run")

		if err != nil {
			fmt.Printf("Error: Failed to parse the hook URL. %s\n", err)
			os.Exit(1)
		}

		b := buf.Bytes()
		r := bytes.NewReader(b)

		req, err := http.NewRequest(vargs.Method, uri.String(), r)

		if err != nil {
			fmt.Printf("Error: Failed to create the HTTP request. %s\n", err)
			os.Exit(1)
		}

		req.Header.Set("Content-Type", vargs.ContentType)

		for key, value := range vargs.Headers {
			req.Header.Set(key, value)
		}

		if vargs.Auth.Username != "" {
			req.SetBasicAuth(vargs.Auth.Username, vargs.Auth.Password)
		}

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Printf("Error: Failed to execute the HTTP request. %s\n", err)
			os.Exit(1)
		}

		defer resp.Body.Close()


		if vargs.Debug || resp.StatusCode >= http.StatusBadRequest {
			body, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				fmt.Printf("Error: Failed to read the HTTP response body. %s\n", err)
			}

			if vargs.Debug {
				fmt.Printf(
					debugRespFormat,
					i+1,
					req.URL,
					req.Method,
					req.Header,
					string(b),
					resp.Status,
					string(body),
				)
			} else {
				fmt.Printf(
					respFormat,
					i+1,
					req.URL,
					resp.Status,
					string(body),
				)
			}
			if resp.StatusCode >= http.StatusBadRequest {
				fmt.Printf("Error: Deployment Failed.\n")
				os.Exit(1)
			}
		}
	}
}
