package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {

	type ApiResponse struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ThumbNail   string `json:"thumbnailUrl"`
		Uploader    string `json:"uploader"`
		Duration    int    `json:"duration"`
		Views       int    `json:"views"`
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if r.URL.Path == "/" {
			_, err := fmt.Fprintf(w, "<!DOCTYPE html><head><meta http-equiv=\"Refresh\" content=\"0; url='https://piped.kavin.rocks/'\"></head>")
			if err != nil {
				return
			}
		} else if r.URL.Path != "/" {
			path := "https://pipedapi.kavin.rocks/streams/" + r.URL.Path
			resp, err := http.Get(path)
			if err != nil {
				log.Fatal(err)
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}(resp.Body)
			body, err := ioutil.ReadAll(resp.Body)

			var embedInfo ApiResponse
			if err := json.Unmarshal(body, &embedInfo); err != nil {
				log.Fatal(err)
			}

			_, err = fmt.Fprintf(w, "<!DOCTYPE html><head><meta content='%s' property='og:title'><meta content='Channel: %s | Views: %d | Duration: %s' property='og:description'><meta content='https://piped.kavin.rocks%s' property='og:url'><meta content='%s' property='og:image'><meta http-equiv=\"Refresh\" content=\"0; url='https://piped.kavin.rocks%s'\"></head>", embedInfo.Title, embedInfo.Uploader, embedInfo.Views, sortTime(embedInfo.Duration), r.URL.Path, embedInfo.ThumbNail, r.URL.Path)
			if err != nil {
				return
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

