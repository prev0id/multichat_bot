package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
)

var (
	port string
)

func init() {
	flag.StringVar(&port, "port", "8080", "sets port for serving website, by default 8080")
}

func main() {
	http.HandleFunc("/", rootHandler)

	fs := http.FileServer(http.Dir("website/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	slog.Info(
		"starting server",
		slog.String("port", port),
	)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err.Error())
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "website/index.html")
}
