package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var today *finaljeop
var random *finaljeop

/*
cat: The category of the question
clue: The question itself
answer: The correct answer to the question
id: the randomly generated ID that will be used to look up the answer
url: the generated url that the correct answer to this particular question will be saved at
*/
type finaljeop struct {
	category string
	clue     string
	answer   string
	id       int
	url      string
}

/*
Creates a new Final Jeopardy Question and saves every answer
*/
func newFinalJeop(cat string, clue string, answer string) *finaljeop {

	fj := finaljeop{category: cat}
	fj.clue = clue
	fj.answer = answer

	// random 6 digit ID for answer URL TKTKTK clear this every x minutes
	rand.Seed(time.Now().UnixNano())
	min := 99999
	max := 1000000
	randomID := rand.Intn(max-min+1) + min
	fj.id = randomID
	fj.url = "localhost:1964/answer?param=" + strconv.Itoa(randomID)

	rememberThat(randomID, answer, clue)

	return &fj
}

func rememberThat(id int, answer string, clue string) {
	f, err := os.Create(strconv.Itoa(id) + "answer.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.WriteString(answer + "~!~" + clue)
	if err != nil {
		fmt.Println(err, l)
		f.Close()
		return
	}

	// fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Deprecated
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

	// fmt.Println("Today's: ", today)
	// fmt.Println("Random:", random)
}
