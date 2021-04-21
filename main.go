package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aiomonitors/godiscord"
)

type Products struct {
	MainItems []MainItem `json:"products"`
}

type MainItem struct {
	Title  string  `json:"title"`
	Id     int64   `json:"id"`
	Handle string  `json:"handle"`
	Sizes  []Size  `json:"variants"`
	Images []Image `json:"images"`
}

type Size struct {
	Title   string `json:"title"`
	Sku     string `json:"sku"`
	Price   string `json:"price"`
	InStock bool   `json:"available"`
}

type Image struct {
	Src string `json:"src"`
}

var SentProductsIds = make([]int64, 0)

func ReleasesMonitor() {
	for {
		client := &http.Client{}
		link := "https://www.a-ma-maniere.com/collections/releases/products.json"

		req, reqerr := client.Get(link)
		if reqerr != nil {
			fmt.Printf("Could not send the request: %s", reqerr)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, resperr := ioutil.ReadAll(req.Body)
		if resperr != nil {
			fmt.Printf("Could not get the request response!: %s", resperr)
		}
		_ = req.Body.Close()

		var products Products
		err := json.Unmarshal(resp, &products)
		if err != nil {
			fmt.Printf("Could not unmarshal the body: %s", err)
		}
		for _, product := range products.MainItems {
			// this checks if the current product has been sent in a webhook. if true, skip the current product, and check the next one.
			if contains(SentProductsIds, product.Id) {
				continue
			} else {
				sendWebhook(product, "webhook")
				time.Sleep(time.Millisecond * 500)
				SentProductsIds = append(SentProductsIds, product.Id)
			}

		}
		time.Sleep(time.Second * 30)
		SentProductsIds = nil

	}
}

func sendWebhook(product MainItem, webhook string) {
	emb := godiscord.NewEmbed(product.Title, "", fmt.Sprintf("https://www.a-ma-maniere.com/collections/releases/products/%s", product.Handle))
	ImageString := product.Images[0].Src
	ImageUrl := strings.Replace(ImageString, `\`, "", -1)
	emb.SetThumbnail(ImageUrl)
	emb.SetColor("#ffc0cb")
	emb.SetFooter("@RixtiRobotics", "")

	for _, size := range product.Sizes {
		emb.AddField(fmt.Sprintf("Size[%s]", size.Title), fmt.Sprintln("STOCK: "+strconv.FormatBool(size.InStock)), true)

	}
	emb.Username = "A-ma-Maniere Monitor"
	emb.SendToWebhook(webhook)
}

func main() {
	ReleasesMonitor()
}

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
