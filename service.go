package main

import (
	"log"
	"net/http"
)

func Hello(rw http.ResponseWriter, r *http.Request) {
	/*
	 * The response writer is what the server will respond with for any request. That's how you write text, json or whatever you want to return to the client
	 * The request includes things like headers and cookies and anything else that the client (web browser) sends to the server
	 */
	rw.Write([]byte("Hello"))
}

// I want the answer to be returned in link markdown that takes you to a simple page that displays the answer
func takeItAwayKen(rw http.ResponseWriter, r *http.Request) {
	scrape()
	rw.Write([]byte(today.clue))
}

func hitTheTrebekVault(rw http.ResponseWriter, r *http.Request) {
	scrape()
	rw.Write([]byte(random.clue))
}

func main() {
	http.HandleFunc("/todayjeop", takeItAwayKen)
	http.HandleFunc("/jeop", hitTheTrebekVault)

	didYouGetHere()

	log.Fatal(http.ListenAndServe(":1964", nil))
}
