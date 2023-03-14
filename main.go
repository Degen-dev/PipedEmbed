package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {

	type ApiResponse struct {
		Title     string `json:"title"`
		ThumbNail string `json:"thumbnailUrl"`
		Uploader  string `json:"uploader"`
		Duration  int    `json:"duration"`
		Views     int    `json:"views"`
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		if r.URL.Path == "/" {
			_, err := fmt.Fprintf(w, "<!DOCTYPE html><head><meta http-equiv=\"Refresh\" content=\"0; url='https://piped.kavin.rocks/'\"></head>")
			if err != nil {
				log.Println("An error occured when trying to redirect user: ", err)
			}
		} else if r.URL.Path != "/" {
			var embedInfo ApiResponse
			path := "https://pipedapi.kavin.rocks/streams/" + r.URL.Path
			resp, err := http.Get(path)
			if err != nil {
				log.Println("An error occured when trying to get video info: ", err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println("An error occured when trying to get JSON response ", err)
				}
			}(resp.Body)
			body, err := ioutil.ReadAll(resp.Body)
			if err := json.Unmarshal(body, &embedInfo); err != nil {
				log.Println("An error occured when trying to parse JSON response: ", err)
			}

			_, err = fmt.Fprintf(w, "<!DOCTYPE html><head><meta content=\"%s\" property=\"og:title\"><meta content=\"Channel: %s | Views: %s | Duration: %s\" property=\"og:description\"><meta content=\"https://piped.kavin.rocks%s\" property=\"og:url\"><meta content=\"%s\" property=\"og:image\"><meta http-equiv=\"Refresh\" content=\"0; url='https://piped.kavin.rocks%s'\"></head>", embedInfo.Title, embedInfo.Uploader, sortViews(embedInfo.Views), sortTime(embedInfo.Duration), r.URL.Path, embedInfo.ThumbNail, r.URL.Path)
			if err != nil {
				log.Println("An error occured when trying to send video data: ", err)
			}
		}
	})

	log.Fatal(http.ListenAndServe(":9072", nil))
}

func sortTime(num int) string {
	var vidLength string
	if num < 60 {
		vidLength = strconv.Itoa(num) + "s"
	} else if num >= 60 {
		vidLength = strconv.Itoa(num/60) + "m " + strconv.Itoa(num%60) + "s"
	}

	return vidLength
}

func sortViews(num int) string {
	p := message.NewPrinter(language.English)
	sorted := p.Sprintf("%d", num)

	return sorted
}
