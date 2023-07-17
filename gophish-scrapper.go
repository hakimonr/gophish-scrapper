package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	apiKey := "******" ///Enter your Gophish Server's API key!!!
	campaignID := "***" ///ChangeCampaignID!!!

	url := fmt.Sprintf("https://YOUR-GOPHISP-SERVER-IP:3333/api/campaigns/%s/results?api_key=%s", campaignID, apiKey)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var results []map[string]interface{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		// Try unmarshaling as a single object instead
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			panic(err)
		}
		results = append(results, result)
	}

	clickedLinkFile, err := os.Create("clicked_links.txt")
	if err != nil {
		panic(err)
	}
	defer clickedLinkFile.Close()

	submittedDataFile, err := os.Create("submitted_data.txt")
	if err != nil {
		panic(err)
	}
	defer submittedDataFile.Close()

	for _, result := range results {
		if result["results"] != nil {
			results := result["results"].([]interface{})
			for _, r := range results {
				result := r.(map[string]interface{})
				email := result["email"].(string)
				if result["status"].(string) == "Clicked Link" {
					clickedLinkFile.WriteString(email + "\n")
				} else if result["status"].(string) == "Submitted Data" {
					submittedDataFile.WriteString(email + "\n")
				}
			}
		}
	}
}
