package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/spiegel-im-spiegel/aozora-api"
)

type Book struct {
	ID     int    `json:"_id`
	Title  string `json:"title"`
	Author string `json:"author"`
	Link   string `json:"link"`
	Total  int    `json:"total"`
	Unique int    `json: "unique"`
	Kanji  []rune `json: "kanji"`
}

func runeInArray(c rune, array []rune) bool {

	for _, i := range array {

		if i == c {
			return true
		}
	}
	return false
}

func countKanji(text string) (int, int, []rune) {

	total := 0
	unique := 0
	var kanji []rune

	for _, c := range text {

		match, _ := regexp.MatchString("[一-龯]", string(c))

		total++

		if match && !runeInArray(c, kanji) {
			kanji = append(kanji, c)
			unique++
		}
	}
	return total, unique, kanji
}

func main() {

	var books []Book
	var text string
	count := 10577

	c := colly.NewCollector(

		colly.AllowedDomains("www.aozora.gr.jp"),
		colly.DetectCharset(),
	)

	c.OnHTML("div.main_text", func(e *colly.HTMLElement) {

		text = e.Text
	})

	// ウェストミンスター寺院 16309
	// 4629 + 60 + 290 + 5127
	// 5082-7900, 7941-12600, 12700-13200, 13300-15900, 16100-18300, 18500-24300, 24500-33100, 33300-35200, 60800
	for i := 50000; i < 60800; i++ {

		if i == 59898 {
			i = 59899
		}

		bk, err := aozora.DefaultClient().LookupBook(i)

		if err == nil {

			var temp Book
			count++
			c.Visit(bk.HTMLURL)

			fmt.Println(bk.Title, count)

			temp.ID = count
			temp.Title = bk.Title
			temp.Author = bk.Authors[0].LastName + " " + bk.Authors[0].FirstName
			temp.Link = bk.HTMLURL
			temp.Total, temp.Unique, temp.Kanji = countKanji(text)

			books = append(books, temp)
		}
	}

	jsonBooks, err := json.MarshalIndent(books, "", " ")

	if err != nil {
		fmt.Println(err)
	}

	_ = ioutil.WriteFile("aozora.json", jsonBooks, 0644)

}
