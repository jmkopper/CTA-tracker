//cta_api.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ctaReqHandler struct {
	config  *configSettings
}

type CTAPostRequest struct {
	ReqType string `json:"reqType"`
	Stpid   string `json:"stpid"`
	Rt      string `json:"rt"`
	Dir     string `json:"dir"`
}

// called by /getCTA<reqType> endpoints
// route POST requests to ctaGetRequest() which talks to the CTA API, then return the response to the frontend as JSON
func (th ctaReqHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("content-type", "application/json")

	var postBody CTAPostRequest
	if err := json.NewDecoder(r.Body).Decode(&postBody); err != nil {
		log.Println(err)
	}

	ctaResponse := make(chan map[string]interface{})
	go ctaGetRequest(ctaResponse, postBody, th.config)
	// transmit the response back to the frontend, modulo timeout error
	select {
	case resp := <-ctaResponse:
		json.NewEncoder(w).Encode(resp)
	case <-time.After(time.Second * 5):
		fmt.Fprintf(w, "timeout")
	}
}

// call the CTA API and pass the response back to the HTTP handler
func ctaGetRequest(ctaResponse chan<- map[string]interface{}, args CTAPostRequest, config *configSettings) {
	// format the request string
	var stopString, routeString, dirString string
	if args.Stpid != "" {
		stopString = fmt.Sprintf("&stpid=%s", args.Stpid)
	}
	if args.Rt != "" {
		routeString = fmt.Sprintf("&rt=%s", args.Rt)
	}
	if args.Dir != "" {
		dirString = fmt.Sprintf("&dir=%s", args.Dir)
	}
	affix := stopString + routeString + dirString
	getUrl := fmt.Sprintf("http://ctabustracker.com/bustime/api/v2/%s?key=%s%s&format=json", args.ReqType, config.apiKey, affix)
	resp, err := http.Get(getUrl)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// pass the response to the channel
	var jsonResponse map[string]interface{}
	json.Unmarshal(body, &jsonResponse)
	ctaResponse <- jsonResponse
}
