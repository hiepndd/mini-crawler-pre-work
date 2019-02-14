package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PostInfo struct {
	URL    string
	Title  string
	Author string
	Date   string
}

func (postInfo *PostInfo) crawlData(url string) {
	res, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.Item1 a").Each(func(i int, s *goquery.Selection) {
		str, exists := s.Attr("href")
		if exists {
			fmt.Println(str)
		}

	})

	doc.Find("h1 span.Title").Each(func(i int, s *goquery.Selection) {
		postInfo.Title = s.Text()
	})

	doc.Find("td span.ReferenceSourceTG").Each(func(i int, s *goquery.Selection) {
		postInfo.Author = strings.Trim(string(s.Text()), "(*)")

	})

	doc.Find("td span.Date").Each(func(i int, s *goquery.Selection) {
		postInfo.Date = s.Text()
	})
}

func main() {

	postInfo := &PostInfo{}
	postInfo.crawlData("https://www.thesaigontimes.vn/121624/Cuoc-cach-mang-dau-khi-da-phien.html")

}
