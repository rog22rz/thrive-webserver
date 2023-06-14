package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Form struct for creator form submissions
type Form struct {
	SiteId string `json:"site"`
	Data   struct {
		ArtistName         string `json:"Name of Artist"`
		Name               string `json:"Name"`
		Pronouns           string `json:"Pronouns"`
		ContactMethod      string `json:"Preferred Method of Contact"`
		Contact            string `json:"Phone/Email"`
		Budget             string `json:"Budget"`
		Date               string `json:"Delivery Date"`
		ProjectDescription string `json:"Field"`
	} `json:"data"`
}

// CreatorItem struct for CreatorItem to add to the CMS
type CreatorItem struct {
	Fields struct {
		Slug                   string `json:"slug"`
		Name                   string `json:"name"`
		IsArchived             bool   `json:"_archived"`
		IsDraft                bool   `json:"_draft"`
		NotAcceptingCommission bool   `json:"not-accepting-commission-projects"`
		AdditionalInfo         string `json:"additional-information"`
		InstagramHandle        string `json:"instagram-handle-2"`
		Occupation             string `json:"occupation"`
		ArtistDescription      string `json:"artist-description"`
		Collaborating          string `json:"how-collaborating-works"`
		HeroImage              string `json:"artist-hero-image"`
	} `json:"fields"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "This is the Thrive webserver")
	}

	if r.Method == "POST" {
		var form Form

		//Parse body received from webflow webhook
		err := json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
			return
		}

		//Build CreatorItem based on received form
		jsonData, err := buildCreatorItem(form)
		if err != nil {
			log.Fatal(err)
		}

		//Build new POST request to send to Creator colleciton in CMS
		url := "https://api.webflow.com/collections//items"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "Bearer ")

		//Send new CreatorItem to CMS
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		log.Println("Response Status:", resp.Status)
	}
}

// Build a CreatorItem based on Form
func buildCreatorItem(form Form) ([]byte, error) {
	creatorItem := CreatorItem{
		Fields: struct {
			Slug                   string `json:"slug"`
			Name                   string `json:"name"`
			IsArchived             bool   `json:"_archived"`
			IsDraft                bool   `json:"_draft"`
			NotAcceptingCommission bool   `json:"not-accepting-commission-projects"`
			AdditionalInfo         string `json:"additional-information"`
			InstagramHandle        string `json:"instagram-handle-2"`
			Occupation             string `json:"occupation"`
			ArtistDescription      string `json:"artist-description"`
			Collaborating          string `json:"how-collaborating-works"`
			HeroImage              string `json:"artist-hero-image"`
		}{
			Slug:                   "testSlug",
			Name:                   "testName",
			IsArchived:             false,
			IsDraft:                true,
			NotAcceptingCommission: false,
			AdditionalInfo:         "testAdditionalInfo",
		},
	}

	jsonData, err := json.Marshal(creatorItem)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func main() {
	fmt.Println("Thrive Webserver Started")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
