package main

import (
	"flag"
	"log"
	"net/http"

	"development.thatwebsite.xyz/utils/sysman"
)

var (
	addr = flag.String("addr", ":8000", "address to listen on")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s, err := sysman.New()
	if err != nil {
		log.Fatal(err)
	}
	s.Routes()
	if err := http.ListenAndServe(*addr, s); err != nil {
		log.Fatal(err)
	}
}
