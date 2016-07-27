package main

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/html"
)

func getItems(z *html.Tokenizer) (entries []*Item) {
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// end of document, done
			return
		case html.StartTagToken:
			n, a := z.TagName()
			if string(n) == "tr" && a {
				for {
					k, v, more := z.TagAttr()
					if string(k) == "entite" {
						entries = append(entries, getItem(z, string(v)))
						break
					}
					if !more {
						break
					}
				}
			}
		}
	}
}

func getItem(z *html.Tokenizer, entite string) (item *Item) {
	item = &Item{
		Entite: entite,
	}

	z.Next() // text (newline)
	z.Next() // td
	z.Next() // input
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	date, err := time.Parse("02/01/2006", z.Token().Data)
	if err != nil {
		log.Printf("Failed to parse date %s", z.Token().Data)
		os.Exit(1)
	}
	item.Date = date
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	item.Location = z.Token().Data
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	item.Type = z.Token().Data
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	item.Title = z.Token().Data
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	item.Barcode = z.Token().Data
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	item.RentType = z.Token().Data
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	item.Booked = z.Token().Data
	z.Next() // /td

	z.Next() // text (newline)
	z.Next() // /tr

	return
}
