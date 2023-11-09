package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

var today *finaljeop
var random *finaljeop

type finaljeop struct {
	category string
	clue     string
	answer   string
}

func newFinalJeop(cat string, clue string, answer string) *finaljeop {

	fj := finaljeop{category: cat}
	fj.clue = clue
	fj.answer = answer
	return &fj
}

func didYouGetHere() {
	fmt.Println("You made it this far")
}

func scrape() {
	c := colly.NewCollector(
	// colly.AllowedDomains("jeopardy.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		// could set headers here
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Response Code", r.StatusCode)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("error", err)
	})

	c.OnHTML(".final_round", func(h *colly.HTMLElement) {
		fullTable := h.DOM
		catName := fullTable.Find(".category_name").Text()
		// fmt.Println(catName)

		// This might look a little weird to you in the future, but the clue_text div applies to the clue itself and the hidden answer underneath
		// In order to separate them into two separate entities, we just remove the answer from the clue string itself
		clueSlashAnswer := fullTable.Find(".clue_text").Text()
		answer := fullTable.Find(".correct_response").Text()

		clueText := strings.ReplaceAll(clueSlashAnswer, answer, "")

		fj := newFinalJeop(catName, clueText, answer)

		if today == nil {
			today = fj
		} else {
			random = fj
		}
	})

	c.Visit("https://j-archive.com/")

	// fmt.Println("Today's", today)
	// fmt.Println("Random:", random)
}
