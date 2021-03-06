package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// PostInfo is a representation of a post in url including url, title, author, date
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

func (postInfo *PostInfo) handleRelatedLink() []PostInfo {
	var listPost []PostInfo
	relatedLinks := make(map[string]int)
	posts := &PostInfo{}
	links, post := posts.crawlData("https://www.thesaigontimes.vn/121624/Cuoc-cach-mang-dau-khi-da-phien.html")

	for _, link := range links {
		relatedLinks[link] = 1
	}

	listPost = append(listPost, post)

	for _, link := range links {
		newLinks, newPost := post.crawlData("https://www.thesaigontimes.vn" + link)

		for _, newLink := range newLinks {
			if relatedLinks[newLink] != 1 {
				_, newPost := post.crawlData("https://www.thesaigontimes.vn" + newLink)
				relatedLinks[newLink] = 1
				listPost = append(listPost, newPost)
			}
		}

		listPost = append(listPost, newPost)
	}
	return listPost
}

func main() {
	newCrawl := &PostInfo{}

	lists := newCrawl.handleRelatedLink()

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
