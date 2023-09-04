package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"thrive-webserver/util"

	"github.com/iancoleman/strcase"
)

// Form struct for creator form submissions
type Form struct {
	FormName string                 `json:"name"`
	SiteId   string                 `json:"site"`
	Data     map[string]interface{} `json:"data"`
}

type CreatorForm struct {
	Form
	Data struct {
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
	IsArchived bool `json:"isArchived"`
	IsDraft    bool `json:"isDraft"`
	FieldData  struct {
		Slug                   string `json:"slug"`
		Name                   string `json:"name"`
		NotAcceptingCommission bool   `json:"not-accepting-commission-products"`
		AdditionalInfo         string `json:"additional-information"`
		InstagramHandle        string `json:"instagram-handle-2"`
		Occupation             string `json:"occupation"`
		ArtistDescription      string `json:"artist-description"`
		Collaborating          string `json:"how-collaborating-works"`
		HeroImage              string `json:"artist-hero-image"`
	} `json:"fieldData"`
}

type ProductForm struct {
	Form
	Data struct {
		ProductName        string `json:"Project Name"`
		Description        string `json:"Project Description"`
		Categories         string `json:"Categories"`
		IsProductAvailable string `json:"isProjectAvailable"`
		Creator            string `json:"Creator"`
		IsChangeByToWidth  string `json:"isChangeByToWith"`
	} `json:"data"`
}

type ProductItem struct {
	PublishStatus string `json:"publishStatus"`
	ProductObject struct {
		IsArchived bool `json:"isArchived"`
		IsDraft    bool `json:"isDraft"`
		FieldData  struct {
			Slug             string `json:"slug"`
			IsShippable      bool   `json:"shippable"`
			IsChangeByToWith bool   `json:"change-by-to-with"`
			ItemName         string `json:"name"`
			ItemDescription  string `json:"description"`
			// Categories       string `json:"category"`
			// CreatorName string `json:"creator"`
		} `json:"fieldData"`
	} `json:"product"`
	Sku struct {
		FieldData struct {
			Slug  string `json:"slug"`
			Name  string `json:"name"`
			Price struct {
				Value int    `json:"value"`
				Unit  string `json:"unit"`
			} `json:"price"`
		} `json:"fieldData"`
	} `json:"sku"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// apiToken := os.Getenv("API_TOKEN")
	// creatorCollectionId := os.Getenv("COLLECTION_ID_CREATOR")
	// siteId := os.Getenv("SITE_ID")

	//For Logger
	// productID := "thrive-389702"
	// logName := "log"

	// err := logger.Init(productID, logName)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize logger: %v", err)
	// }
	// defer logger.Close()

	//For dev, get env variables from config file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	apiToken := config.APIToken
	creatorCollectionId := config.CreatorCollectionId
	siteId := config.SiteId

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
	}

	if r.Method == "POST" {
		var form Form
		var jsonData []byte
		var err error
		var url string

		//Parse body received from webflow webhook
		err = json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
			return
		}

		//Build items based on received form
		if form.FormName == "Creator Form" {
			var creatorForm CreatorForm
			jsonBlob, _ := json.Marshal(form)
			err := json.Unmarshal(jsonBlob, &creatorForm)
			if err != nil {
				log.Fatal(err)
			}
			jsonData, err = mapperCreatorToItem(creatorForm)
			if err != nil {
				log.Fatal(err)
			}
			url = "https://api.webflow.com/collections/" + creatorCollectionId + "/items"
		} else if form.FormName == "Product Form" {
			var productForm ProductForm
			jsonBlob, _ := json.Marshal(form)
			err := json.Unmarshal(jsonBlob, &productForm)
			if err != nil {
				log.Fatal(err)
			}
			jsonData, err = mapProductToItem(productForm)
			if err != nil {
				log.Fatal(err)
			}
			url = "https://api.webflow.com/beta/sites/" + siteId + "/products"
			// url = "https://webhook.site/f20fc6e6-54b7-4ebf-b027-1d466cafc692"
		} else {
			log.Fatal("Form " + form.FormName + " is not supported")
		}

		//Build POST request to CMS
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Add("accept", "application/json")
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "Bearer "+apiToken)

		//Send new item to CMS
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
func mapperCreatorToItem(form CreatorForm) ([]byte, error) {
	creatorItem := CreatorItem{
		IsArchived: false,
		IsDraft:    true,
		FieldData: struct {
			Slug                   string `json:"slug"`
			Name                   string `json:"name"`
			NotAcceptingCommission bool   `json:"not-accepting-commission-products"`
			AdditionalInfo         string `json:"additional-information"`
			InstagramHandle        string `json:"instagram-handle-2"`
			Occupation             string `json:"occupation"`
			ArtistDescription      string `json:"artist-description"`
			Collaborating          string `json:"how-collaborating-works"`
			HeroImage              string `json:"artist-hero-image"`
		}{
			Slug:                   strcase.ToKebab(form.Data.ArtistName),
			Name:                   form.Data.ArtistName,
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

func mapProductToItem(form ProductForm) ([]byte, error) {
	productItem := ProductItem{
		PublishStatus: "staging",
		ProductObject: struct {
			IsArchived bool `json:"isArchived"`
			IsDraft    bool `json:"isDraft"`
			FieldData  struct {
				Slug             string `json:"slug"`
				IsShippable      bool   `json:"shippable"`
				IsChangeByToWith bool   `json:"change-by-to-with"`
				ItemName         string `json:"name"`
				ItemDescription  string `json:"description"`
				// Categories       string `json:"category"`
				// CreatorName string `json:"creator"`
			} `json:"fieldData"`
		}{
			IsArchived: false,
			IsDraft:    true,
			FieldData: struct {
				Slug             string `json:"slug"`
				IsShippable      bool   `json:"shippable"`
				IsChangeByToWith bool   `json:"change-by-to-with"`
				ItemName         string `json:"name"`
				ItemDescription  string `json:"description"`
				// Categories       string `json:"category"`
				// CreatorName string `json:"creator"`
			}{
				Slug:             strcase.ToKebab(form.Data.ProductName),
				IsShippable:      form.Data.IsProductAvailable != "",
				IsChangeByToWith: form.Data.IsChangeByToWidth != "",
				ItemName:         form.Data.ProductName,
				ItemDescription:  form.Data.Description,
				// Categories:       form.Data.Categories,
				// CreatorName: form.Data.Creator,
			},
		},
		Sku: struct {
			FieldData struct {
				Slug  string `json:"slug"`
				Name  string `json:"name"`
				Price struct {
					Value int    `json:"value"`
					Unit  string `json:"unit"`
				} `json:"price"`
			} `json:"fieldData"`
		}{
			FieldData: struct {
				Slug  string `json:"slug"`
				Name  string `json:"name"`
				Price struct {
					Value int    `json:"value"`
					Unit  string `json:"unit"`
				} `json:"price"`
			}{
				Slug: strcase.ToKebab(form.Data.ProductName),
				Name: form.Data.ProductName,
				Price: struct {
					Value int    `json:"value"`
					Unit  string `json:"unit"`
				}{
					Value: 0,
					Unit:  "CAD",
				},
			},
		},
	}
	jsonData, err := json.Marshal(productItem)
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
