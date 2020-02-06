package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type config struct {
	HTTPClient       *http.Client
	oauthURL         string
	headers          map[string]string
	requestBodyForm  map[string]string
	dumpRequest      bool
	removeAuthHeader bool
}

func main() {

	// flag.Var(&headersFlags, "list1", "Some description for this param.")
	var c config
	c.initConfig()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", c.proxyHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))

}

//set the config values to be used in the request
func (c *config) initConfig() {

	tr := &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	c.HTTPClient = &http.Client{Transport: tr}
	flag.StringVar(&c.oauthURL, "oauthURL", "", "The url to proxy the call to")
	flag.BoolVar(&c.dumpRequest, "dumpRequest", false, "will output the request to the logs")
	flag.BoolVar(&c.removeAuthHeader, "removeAuthHeader", true, "will remove any authentication header set")
	flag.StringToStringVar(&c.headers, "headers", map[string]string{}, "Header values passed as: h1=v1,h2=v2")
	flag.StringToStringVar(&c.requestBodyForm, "requestBodyForm", map[string]string{}, "Request Body form parameters passed as: h1=v1,h2=v2")
	flag.Parse()

	log.Info("Proxy configuration has been set with the values...")
	log.Infof("oauthURL: %s", c.oauthURL)
	log.Infof("headers: %v", c.headers)
	log.Infof("requestBodyForm: %v", c.requestBodyForm)

}

func (c *config) proxyHandler(w http.ResponseWriter, req *http.Request) {

	log.Info("Proxing oauth request...")

	form := url.Values{}
	for key, val := range c.requestBodyForm {
		log.Info(key + "::::" + val)
		form.Add(key, val)
	}

	proxyReq, err := http.NewRequest("POST", c.oauthURL, strings.NewReader(form.Encode()))

	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	proxyReq.Header = req.Header
	for key, val := range c.headers {
		proxyReq.Header.Set(key, val)
	}

	//will unset any Authorization values
	if c.removeAuthHeader {
		proxyReq.Header.Set("Authorization", "")
	}

	// DumpRequest
	if c.dumpRequest {
		requestDump, err := httputil.DumpRequest(proxyReq, true)
		if err != nil {
			log.Error(err.Error())
		}
		log.Infof("Request: %v", string(requestDump))
	}

	resp, err := c.HTTPClient.Do(proxyReq)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if resp.StatusCode != 200 {
		log.Infof("Request status code: %d", resp.StatusCode)
		log.Infof("Request response: %s", []byte(body))
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Encoding", resp.Header.Get("Content-Encoding"))
	w.WriteHeader(resp.StatusCode)
	w.Write([]byte(body))
}
