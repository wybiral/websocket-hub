package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/wybiral/websocket-hub/internal/app"
)

// Current hub version
const version = "0.1.0"

func main() {
	flag.Usage = func() {
		usage()
		os.Exit(0)
	}
	var host string
	flag.StringVar(&host, "h", "127.0.0.1", "Server Host")
	var port int
	flag.IntVar(&port, "p", 8080, "Server Port")
	var dir string
	flag.StringVar(&dir, "d", "", "Public directory (optional)")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", host, port)
	// Create and start App server
	a := app.New(addr, dir)
	log.Println("Serving at " + addr)
	err := a.Server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Print("NAME:\n")
	fmt.Print("   hub\n\n")
	fmt.Print("USAGE:\n")
	fmt.Print("   hub -h 127.0.0.1 -p 8080 [-d /public/directory]\n\n")
	fmt.Print("VERSION:\n")
	fmt.Printf("   %s\n\n", version)
}
