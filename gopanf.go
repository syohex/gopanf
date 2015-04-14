package main

import (
	"encoding/xml"
	"fmt"
	"github.com/mgutz/ansi"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RSS struct {
	Item []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Creator     string `xml:"creator"`
}

const metacpanRSS = `https://metacpan.org/feed/recent?f=`

func main() {
	var showDesc bool
	var color bool

	limit := flag.IntP("limit", "l", 30, "limit of modules")
	flag.BoolVarP(&showDesc, "desc", "d", false, "Show description")
	flag.BoolVarP(&color, "color", "c", false, "Show description")
	flag.Parse()

	response, err := http.Get(metacpanRSS)
	if err != nil {
		log.Printf("Can't download '%s'", metacpanRSS)
		os.Exit(1)
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Can't read response: %s", err.Error())
		os.Exit(1)
	}

	var rss RSS
	if err := xml.Unmarshal(bytes, &rss); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	for _, item := range rss.Item[0:*limit] {
		if color {
			fmt.Printf("%s (%s)",
				ansi.Color(item.Title, "green+b"),
				ansi.Color(item.Creator, "yellow+b"))
		} else {
			fmt.Printf("%s (%s)", item.Title, item.Creator)
		}

		if showDesc {
			fmt.Printf(":\n    %s\n", item.Description)
		} else {
			fmt.Println("")
		}
	}
}
