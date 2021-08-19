package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
)

var rEverything = regexp.MustCompile(`.*`) // Route these to the backend ES cluster
var rDoc = regexp.MustCompile(`/_doc`)     //any singular indexing operation

func proxy(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	//proxy this request through to the backend ES cluster and return the response (choose cluster based on config file)
}

func home(w http.ResponseWriter, r *http.Request) {
	//return string as json
	fmt.Fprintf(w, "{  \"name\" : \"instance-0000000045\",  \"cluster_name\" : \"ddc7994029894f8d8ed1847318505939\",  \"cluster_uuid\" : \"VJM2yDPqTLaaXEK0wAO-sA\",  \"version\" : {    \"number\" : \"7.13.2\",    \"build_flavor\" : \"default\",    \"build_type\" : \"docker\",    \"build_hash\" : \"4d960a0733be83dd2543ca018aa4ddc42e956800\",    \"build_date\" : \"2021-06-10T21:01:55.251515791Z\",    \"build_snapshot\" : false,    \"lucene_version\" : \"8.8.2\",    \"minimum_wire_compatibility_version\" : \"6.8.0\",    \"minimum_index_compatibility_version\" : \"6.0.0-beta1\"  },  \"tagline\" : \"You Know, for Search\"}")
}

func bulk(w http.ResponseWriter, r *http.Request) {
	//mock bulk API response
	resp := make(map[string]string)
	resp["message"] = "Bulked"

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}

	resp["request"] = string(requestDump)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func index(w http.ResponseWriter, r *http.Request) {
	//mock POST index/_doc API response
	resp := make(map[string]string)
	resp["message"] = "Indexed"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func route(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//authenticate supplied credentials
	switch {
	case r.URL.Path == "/":
		home(w, r)
	case r.URL.Path == "/_bulk":
		bulk(w, r)
	case rDoc.MatchString(r.URL.Path):
		index(w, r)
	case rEverything.MatchString(r.URL.Path):
		proxy(w, r)
	default:
		fmt.Println("Unknown pattern")
	}
}

func main() {
	http.HandleFunc("/", route)

	http.ListenAndServe(":9243", nil)
}
