package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func getIntFromQuery(r *http.Request, k string, d int) int {
	t := r.URL.Query().Get(k)
	value, err := strconv.Atoi(t)
	if err != nil {
		value = d
	}
	return value
}

func timeoutHandler(w http.ResponseWriter, r *http.Request) {
	timeout := getIntFromQuery(r, "t", 60)
	time.Sleep(time.Duration(timeout) * time.Second)
	fmt.Fprint(w, "Hello world")
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	hops := getIntFromQuery(r, "h", 100)
	switch hops {
	case 1:
		fmt.Fprint(w, "Hello world")
	default:
		newURL, err := url.Parse(r.URL.String())
		if err != nil {
			fmt.Fprint(w, "Failed")
		}
		query := newURL.Query()
		query.Set("h", strconv.Itoa(hops-1))
		newURL.RawQuery = query.Encode()
		http.Redirect(w, r, newURL.String(), http.StatusFound)
	}
}

func longHeaderNameHandler(w http.ResponseWriter, r *http.Request) {
	length := getIntFromQuery(r, "l", 5000)
	w.Header().Set("f"+strings.Repeat("o", length-1), "bar")
	fmt.Fprint(w, "Hello world")
}

func longHeaderValueHandler(w http.ResponseWriter, r *http.Request) {
	length := getIntFromQuery(r, "l", 5000)
	w.Header().Set("foo", "b"+strings.Repeat("a", length-2)+"r")
	fmt.Fprint(w, "Hello world")
}

func lotOfHeadersHandler(w http.ResponseWriter, r *http.Request) {
	amount := getIntFromQuery(r, "a", 1000)
	counter := 1
	for counter < amount {
		w.Header().Set("foo-"+strconv.Itoa(counter), "bar")
		counter++
	}
	fmt.Fprint(w, "Hello world")
}

func invalidHeaderHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("foo bar", "baz")
	w.Header().Set("new\nline\n\n", "baz")
	w.Header().Set("double new\n\nline", "baz")
	w.Header().Set("🍀", "🐎")
	w.Header().Set(":", ";")
	fmt.Fprint(w, "Hello world")
}

func main() {
	http.HandleFunc("/response/timeout", timeoutHandler)

	http.HandleFunc("/response/redirects", redirectHandler)

	http.HandleFunc("/response/headers/long-name", longHeaderNameHandler)

	http.HandleFunc("/response/headers/long-value", longHeaderValueHandler)

	http.HandleFunc("/response/headers/lot-of-headers", lotOfHeadersHandler)

	http.HandleFunc("/response/headers/invalid", invalidHeaderHandler)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
