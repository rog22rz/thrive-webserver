package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"thrive-webserver/logger"

	"github.com/iancoleman/strcase"
)

// Form struct for creator form submissions
type Form struct {
	FormName string `json:"name"`
	SiteId   string `json:"site"`
	Data     struct {
		ArtistName            string `json:"Artist Name"`
		Email                 string `json:"Email"`
		PhoneNumber           string `json:"Phone Number"`
		Occupation            string `json:"Occupation"`
		ShortDescription      string `json:"Short Description"`
		LongDescription       string `json:"Long Description"`
		Website               string `json:"Website"`
		Socials               string `json:"Socials"`
		HowCollaboartingWorks string `json:"How Collaboration Works"`
		MainImage             string `json:"Main Image"`
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
	apiToken := os.Getenv("_API_TOKEN")
	creatorCollectionId := os.Getenv("_COLLECTION_ID_CREATOR")

	projectID := "thrive-webserver"
	logName := "log"

	err := logger.Init(projectID, logName)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	// config, err := util.LoadConfig(".")
	// if err != nil {
	// 	log.Fatal("cannot load config:", err)
	// }

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
		jsonData, err := mapperCreatorToItem(form)
		if err != nil {
			log.Fatal(err)
		}

		logger.Log(creatorCollectionId)
		logger.Log(apiToken)

		//Build new POST request to send to Creator colleciton in CMS
		url := "https://api.webflow.com/collections/" + creatorCollectionId + "/items"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "Bearer "+apiToken)

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
func mapperCreatorToItem(form Form) ([]byte, error) {
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
			Slug:                   strcase.ToKebab(form.Data.ArtistName),
			Name:                   form.Data.ArtistName,
			IsArchived:             false,
			IsDraft:                true,
			NotAcceptingCommission: form.Data.HowCollaboartingWorks == "",
			AdditionalInfo:         form.Data.LongDescription,
			InstagramHandle:        form.Data.Socials,
			Occupation:             form.Data.Occupation,
			ArtistDescription:      form.Data.ShortDescription,
			Collaborating:          form.Data.HowCollaboartingWorks,
			HeroImage:              form.Data.MainImage,
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
