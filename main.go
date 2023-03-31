package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"strconv"
)

type ApiResponse struct {
	Title     string `json:"title"`
	ThumbNail string `json:"thumbnailUrl"`
	Uploader  string `json:"uploader"`
	Duration  int    `json:"duration"`
	Views     int    `json:"views"`
}

func main() {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("text/html")

		if string(ctx.Path()) == "/" {
			_, err := ctx.WriteString("<!DOCTYPE html><head><meta http-equiv=\"Refresh\" content=\"0; url='https://piped.kavin.rocks/'\"></head>")
			if err != nil {
				log.Println("An error occurred when trying to redirect to Piped: ", err)
				return
			}
		} else if string(ctx.Path()) != "/" {
			var embedInfo ApiResponse
			path := "https://pipedapi.kavin.rocks/streams" + string(ctx.Path())
			statusCode, body, err := fasthttp.Get(nil, path)
			if err != nil || statusCode != fasthttp.StatusOK {
				log.Println("An error occurred when trying to get video info: ", err, statusCode)
				return
			}

			if err := json.Unmarshal(body, &embedInfo); err != nil {
				log.Println("An error occurred when trying to parse JSON response: ", err)
				return
			}

			_, err = ctx.WriteString(fmt.Sprintf("<!DOCTYPE html><head><meta content=\"%s\" property=\"og:title\"><meta content=\"Channel: %s | Views: %s | Duration: %s\" property=\"og:description\"><meta content=\"https://piped.kavin.rocks%s\" property=\"og:url\"><meta content=\"%s\" property=\"og:image\"><meta http-equiv=\"Refresh\" content=\"0; url='https://piped.kavin.rocks%s'\"></head>", embedInfo.Title, embedInfo.Uploader, sortViews(embedInfo.Views), sortTime(embedInfo.Duration), string(ctx.Path()), embedInfo.ThumbNail, string(ctx.Path())))
			if err != nil {
				log.Println("An error occurred when trying to send video info: ", err)
			}
		}
	}

	if err := fasthttp.ListenAndServe("127.0.0.1:9072", requestHandler); err != nil {
		log.Fatal(err)
	}
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
