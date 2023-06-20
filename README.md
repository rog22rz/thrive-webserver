# Thrive Webserver

[![Go Version](https://img.shields.io/badge/go-v1.16-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-brightgreen)](LICENSE)

TLDR: This is a webserver that has been designed to post a new item to a Webflow CMS when it receives a POST from a weflow webhook. 

## About

This project was started because a friend of mine approach me with a problem she was having with her Webflow website. She wanted to have a form that would directly add items to their CMS when the form was submitted. The problems was that Webflow doesn't allow that functionality, which forces users to go through 3rd party softwares such as Zapier to link a form to a CMS, and costs a monthly subscription. 

After a bit of research, I saw that Webflow allow their users to have webhooks which send POST requests on form submit with the details of the form, and also an API for adding items to their CMS. They simply lacked the middle layer, a webserver, which could connect the two functionalities. I thought to myself that such a webserver, which would simply handle the data from the webhook, format it and then POST it to the CMS would be quite an easy task and would allow me to play around with Golang. I would also save my friend around 40$ a month, which would allow more funds for her grassroot project. 

Therefore, this repo contains the code that I have written for the webserver to handle the task that I have described above. I then deployed the webserver using Docker and Google Cloud Run. 

Please feel free to reach out to me if you have any questions or if you just want to connect!

PS: Also go check out my friend's awesome project which showcases the new generation of Toronto Creatives: https://thrive-artwork-afefc46898-d9f19e4ce600d.webflow.io/

## Installation

To install, simply clone the project and make sure you have all the packages installed. There is an app.env file missing for obvious security reasons which would contain the API key and CMS collection ID. However, if you replace the url on line 72 of main.go with your own link form https://webhook.site/, and send a POST request to your localhost:8080 with an item such as the JSON below, you should see the request appear on your https://webhook.site/ with the new format of the CMS. 

```json
{
  "name": "Email Form",
  "site": "dsa888dsa888dsa",
  "data": {
    "Name of Artist": "Test",
    "Name": "Test2",
    "Pronouns": "",
    "Preferred Method of Contact": "Phone Call",
    "Phone/Email": "321-321-3211",
    "Budget": "",
    "Delivery Date": "",
    "field": "Test"
  },
  "d": "2023-06-09T20:50:58.746Z",
  "_id": "dsa888dsa888das888dsa"
}
```

### Prerequisites

- Go (version 1.20.5 or higher): [Installation Guide](https://golang.org/doc/install)

### Clone the repository

```bash
git clone https://github.com/rog22rz/thrive-webserver
