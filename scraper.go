package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func _getHtmlDoc(link string) (*goquery.Document, error) {
	res, err := http.Get(link)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func _getIdFromLink(link string) (string, string) {
	linkParts := strings.Split(link, "/")
	sluggPart := linkParts[len(linkParts)-2]
	uniqePart := linkParts[len(linkParts)-1]
	uniqeId := strings.Replace(uniqePart, ".html", "", 1)
	sluggId := sluggPart + "#" + uniqeId
	return uniqeId, sluggId
}

func scrapeBook(link string) (*Book, error) {
	doc, err := _getHtmlDoc(link)
	if err != nil {
		return nil, err
	}
	// srape id
	_, slugId := _getIdFromLink(link)

	// scrape title
	title := doc.Find("div.pr_header").Find("h1.pr_header__heading").Text()

	// scrape image
	imageUrl := doc.Find("div.pr_images__preview").Find("img#js-book-cover").AttrOr("src", "noimage")

	// scrape description
	description := doc.Find("div.pr_description").Find("span.info__text").Text()

	var pageCount int = 0
	var isbn string = ""
	attributeElms := doc.Find("div.pr_attributes").Find("tr")
	attributeElms.Each(func(i int, s *goquery.Selection) {
		valueAndKey := s.Find("td")
		attrText := valueAndKey.Text()

		// scrape page count
		if strings.Contains(attrText, "Sayfa Sayısı") {
			split := strings.Split(attrText, ":")

			pageCountStr := split[len(split)-1]
			_pageCount, err := strconv.Atoi(pageCountStr)
			if err == nil {
				pageCount = _pageCount
			}
		}

		// scrape ISBN number
		if strings.Contains(attrText, "ISBN") {
			split := strings.Split(attrText, ":")
			isbn = split[len(split)-1]
		}
	})

	// scrape author
	var author *Author = nil
	authorLink := doc.Find("div.pr_producers__manufacturer").Find("a.pr_producers__link").AttrOr("href", "")
	if authorLink != "" {
		athr, authorErr := scrapeAuthor(authorLink)

		if authorErr == nil {
			author = athr
		}
	}

	// scrape category groups
	defaultCategory := "Kitap"
	categoryGroups := [][]string{}
	categoryGroupsElm := doc.Find("ul.rel-cats__list")
	categoryGroupsElm.Each(func(i int, s *goquery.Selection) {
		// now we have selected the <ul><\ul>
		categoryGroupElm := s.Find("li.rel-cats__item")
		categoryGroupElm.Each(func(i int, s *goquery.Selection) {
			// now we have selected the <li><\li>
			categories := []string{}
			categoryElms := s.Find("span")
			categoryElms.Each(func(i int, s *goquery.Selection) {

				// now we have selected the <span><\span>
				category := s.Text()
				// remove default category because it is not
				// defining anything
				if category != defaultCategory {
					categories = append(categories, category)
				}
			})
			categoryGroups = append(categoryGroups, categories)
		})

	})

	if title != "" {
		book := &Book{
			Id:          slugId,
			Title:       title,
			ISBN:        isbn,
			ImageUrl:    imageUrl,
			Link:        link,
			Categories:  categoryGroups,
			PageCount:   pageCount,
			Description: description,
		}

		if author != nil {
			book.Author = *author
		}
		return book, nil
	}

	return nil, nil
}

func scrapeAuthor(link string) (*Author, error) {
	id, slugId := _getIdFromLink(link)

	authorGetUrl := fmt.Sprintf("https://www.kitapyurdu.com/index.php?route=product/manufacturer/manufacturer_about&manufacturer_id=%s", id)
	doc, err := _getHtmlDoc(authorGetUrl)

	if err != nil {
		return nil, err
	}

	name := doc.Find("h2").Find("span").First().Text()

	var imageUrl string = ""
	imageUrlAttr := doc.Find("a.manufacturer-image").AttrOr("style", "")
	if imageUrlAttr != "" {
		rgx, _ := regexp.Compile(`\(.+?\)`)
		imageUrl = rgx.FindString(imageUrlAttr)
		imageUrl = strings.Replace(imageUrl, "(", "", 1)
		imageUrl = strings.Replace(imageUrl, ")", "", 1)
	}
	var biography string = ""
	biographyElm := doc.Find("div#manufacturer-description").Find("p")
	if biographyElm.Length() > 0 {
		biography = biographyElm.First().Text()
	}

	if name != "" {
		return &Author{
			Id:        slugId,
			Name:      name,
			ImageUrl:  imageUrl,
			Link:      link,
			Biography: biography,
		}, nil

	}

	return nil, nil
}
