package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PostInfo struct {
	URL    string
	Title  string
	Author string
	Date   string
}

func (postInfo *PostInfo) crawlData(url string) ([]string, PostInfo) {
	post := PostInfo{}

	var links []string

	post.URL = url

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

	doc.Find("#ctl00_cphContent_Article_LienQuan div.Item1 a").Each(func(i int, s *goquery.Selection) {
		str, exists := s.Attr("href")
		if exists {
			links = append(links, str)
		}

	})

	doc.Find("h1 span.Title").Each(func(i int, s *goquery.Selection) {
		post.Title = s.Text()
	})

	doc.Find("td span.ReferenceSourceTG").Each(func(i int, s *goquery.Selection) {
		post.Author = strings.Trim(string(s.Text()), "(*)")

	})

	doc.Find("td span.Date").Each(func(i int, s *goquery.Selection) {
		post.Date = s.Text()
	})

	return links, post
}

func (postInfo *PostInfo) crawlRelatedLink(url string) PostInfo {
	query := "https://www.thesaigontimes.vn" + url

	post := PostInfo{}

	post.URL = query

	res, err := http.Get(query)

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

	doc.Find("h1 span.Title").Each(func(i int, s *goquery.Selection) {
		post.Title = s.Text()
	})

	doc.Find("td span.ReferenceSourceTG").Each(func(i int, s *goquery.Selection) {
		post.Author = strings.Trim(string(s.Text()), "(*)")

	})

	doc.Find("td span.Date").Each(func(i int, s *goquery.Selection) {
		post.Date = s.Text()
	})

	return post
}

func (posts *PostInfo) final() []PostInfo {
	var listPost []PostInfo
	postInfo := &PostInfo{}
	links, post := postInfo.crawlData("https://www.thesaigontimes.vn/121624/Cuoc-cach-mang-dau-khi-da-phien.html")
	listPost = append(listPost, post)
	for _, link := range links {
		newPost := post.crawlRelatedLink(link)
		listPost = append(listPost, newPost)
	}
	return listPost
}

func main() {
	example := &PostInfo{}

	lists := example.final()

	var content [][]string

	for _, list := range lists {
		post := []string{list.URL, list.Title, list.Author, list.Date}
		content = append(content, post)
	}

	file, err := os.OpenFile("result.csv", os.O_CREATE|os.O_WRONLY, 0777)

	defer file.Close()

	if err != nil {
		os.Exit(1)
	}

	csvWriter := csv.NewWriter(file)

	csvWriter.WriteAll(content)

	csvWriter.Flush()

}
