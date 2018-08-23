package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/wybiral/hub/internal/app"
)

func main() {
	var host string
	flag.StringVar(&host, "h", "127.0.0.1", "Server Host")
	var port int
	flag.IntVar(&port, "p", 8080, "Server Port")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", host, port)
	// Create and start App server
	a := app.New(addr)
	log.Println("Serving at " + addr)
	err := a.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
