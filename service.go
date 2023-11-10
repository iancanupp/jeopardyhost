package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
Function to route to either the final jeopardy question du jour, or grab a randomized
one from the archives
*/
func placeYourWagers(rw http.ResponseWriter, r *http.Request) {
	values := r.URL.Query().Get("param")
	// for k, v := range values {
	// 	fmt.Println(k, " => ", v)
	// }
	if strings.Contains(values, "today") {
		takeItAwayKen(rw)
	} else {
		hitTheTrebekVault(rw)
	}
}

/*
Opens up the corresponding file from the baked-in answerID in the url,
then calls the helper to format it it for display
*/
func displayAnswer(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	answerID := r.URL.Query().Get("param")

	contents, err := os.ReadFile(answerID + "answer.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	rw.Write([]byte(formatAnswerAsHTML(string(contents))))
}

/*
Helper to the reading file function that takes the raw file contents
and splits it into answer/clue in the html for display
*/
func formatAnswerAsHTML(contents string) string {
	splitOnComma := strings.Split(contents, "~!~")
	answer := splitOnComma[0]
	clue := splitOnComma[1]

	htmlStr := "<h1 style='color: blue'>What is... " + answer + "</h1>\n" + "<p>was the correct answer to the clue..." + clue + "</p>"

	return htmlStr
}

/*
Route taken when the request is for today's final jeopardy question.
Otherwise, does the same thing as the randomized one
*/
func takeItAwayKen(w http.ResponseWriter) {
	scrape()
	strbuffer := "In the category " + today.category + ", your clue is:\n" + today.clue + "\n[Time's up! Reveal answer](" + today.url + ")"
	w.Write([]byte(strbuffer))
}

/*
Route taken when the request is missing an indication of 'today', this will
bring up a random final jeopardy
*/
func hitTheTrebekVault(w http.ResponseWriter) {
	// definitely do a separate file of lead ins to category and clue here"
	scrape()

	strbuffer := "In the category " + random.category + ", your clue is:\n" + random.clue + "\n[Time's up! Reveal answer](" + random.url + ")"
	w.Write([]byte(strbuffer))

	// This was for testing the load answer function in html instead of markdown
	// w.Header().Set("Content-Type", "text/html")
	// urlhtml := createLinkHTML(random.url)
	// w.Write([]byte("<p>" + random.category + "</p><p>" + random.clue + "</p><p>" + urlhtml + "</p>"))

}

// This was for testing the load answer function in html instead of markdown
func createLinkHTML(url string) string {
	// print string with double quotes
	fmt.Printf(`"%s"
`, url)

	// create string with double quotes
	quotedURL := fmt.Sprintf(`"%s"`, url)

	finalStr := "<a href=" + quotedURL + ">Answer</a>"
	return finalStr
}

func main() {
	http.HandleFunc("/jeop", placeYourWagers)
	http.HandleFunc("/answer", displayAnswer)
	// TKTKTKTK maybe a help endpoint later

	// didYouGetHere()

	log.Fatal(http.ListenAndServe(":1964", nil))
}
