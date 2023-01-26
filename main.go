package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

//go:embed templates/head.tmpl
var headTemplate string

//go:embed templates/foot.tmpl
var footTemplate string

func slowFetch(t time.Duration, c chan string) {
	time.Sleep(t * time.Second)
	c <- fmt.Sprintf("%d second sleep", t)
}

func index(w http.ResponseWriter, req *http.Request) {
	c1 := make(chan string)
	c2 := make(chan string)
	go slowFetch(2, c1)
	go slowFetch(2, c2)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, req)
		return
	}

	template.Must(template.New("head").Parse(headTemplate)).Execute(w, struct {
		Title string
	}{
		Title: "Hello World",
	})
	flusher.Flush()

	w.Write([]byte(fmt.Sprintf("<p>%s</p>", <-c1)))
	flusher.Flush()

	w.Write([]byte(fmt.Sprintf("<p>%s</p>", <-c2)))
	template.Must(template.New("foot").Parse(footTemplate)).Execute(w, nil)
}

func main() {
	http.HandleFunc("/", index)

	if err := http.ListenAndServe("localhost:1337", nil); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
