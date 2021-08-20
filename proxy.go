package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

var rEverything = regexp.MustCompile(`.*`) // Route these to the backend ES cluster
var rDoc = regexp.MustCompile(`/_doc`)     //any singular indexing operation

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

// Given a request send it to the appropriate url
func proxy(res http.ResponseWriter, req *http.Request) {
	// backing cluster URL
	url := "http://localhost:9200"
	serveReverseProxy(url, res, req)
}

func home(w http.ResponseWriter, r *http.Request) {
	//this could be proxied through if needed too
	w.Header().Set("Content-Type", "application/json")
	//return string as json
	fmt.Fprintf(w, "{  \"name\" : \"mocked\",  \"cluster_name\" : \"mocked\",  \"cluster_uuid\" : \"mo-sA\",  \"version\" : {    \"number\" : \"7.13.2\",    \"build_flavor\" : \"default\",    \"build_type\" : \"docker\",    \"build_hash\" : \"4d960a0733be83dd2543ca018aa4ddc42e956800\",    \"build_date\" : \"2021-06-10T21:01:55.251515791Z\",    \"build_snapshot\" : false,    \"lucene_version\" : \"8.8.2\",    \"minimum_wire_compatibility_version\" : \"6.8.0\",    \"minimum_index_compatibility_version\" : \"6.0.0-beta1\"  },  \"tagline\" : \"You Know, for Search\"}")
}

func bulk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	//mock bulk API response
	h := json.RawMessage(`{
		"_index" : "mocked",
		"_type" : "_doc",
		"_id" : "mocked",
		"_version" : 1,
		"result" : "created",
		"_shards" : {
		  "total" : 2,
		  "successful" : 1,
		  "failed" : 0
		},
		"_seq_no" : 0,
		"_primary_term" : 1
	  }`)

	jsonResp, err := json.Marshal(h)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func route(w http.ResponseWriter, r *http.Request) {

	//authenticate supplied credentials
	switch {
	case r.URL.Path == "/":
		home(w, r)
	case (r.URL.Path == "/_bulk" && r.Method == "POST"):
		bulk(w, r)
	case (rDoc.MatchString(r.URL.Path) && r.Method == "POST"):
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
