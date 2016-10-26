package main

import (
	"errors"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"

	"golang.org/x/net/html"
)

type htmlParser struct {
	*html.Tokenizer
}

func (z *htmlParser) getItems() (entries []*Item) {
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
						if i := z.getItem(string(v)); i.Name != "" {
							entries = append(entries, i)
						}
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

func (z *htmlParser) getItem(entite string) (item *Item) {
	item = &Item{
		Entite: entite,
	}

	z.Next() // text (newline)
	z.Next() // td
	z.Next() // input
	if n, _ := z.TagName(); string(n) != "input" {
		return
	}
	for {
		k, v, more := z.TagAttr()
		if string(k) == "name" {
			item.Name = string(v)
			break
		}
		if !more {
			break
		}
	}
	z.Next() // /td

	z.Next() // td
	z.Next() // text
	date, err := time.Parse("02/01/2006", z.Token().Data)
	if err != nil {
		log.Printf("Failed to parse date %s: %v", z.Token().Data, err)
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

func (z *htmlParser) isLogged() bool {
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// end of document, done
			return false
		case html.StartTagToken:
			n, a := z.TagName()
			if string(n) == "div" && a {
				for {
					k, v, more := z.TagAttr()
					if string(k) == "id" && string(v) == "compte" {
						return true
					}
					if !more {
						break
					}
				}
			}
		}
	}
	return false
}

func (z *htmlParser) checkError() (err error) {
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			// end of document, done
			return
		case html.StartTagToken:
			n, a := z.TagName()
			if string(n) == "span" && a {
				for {
					k, v, more := z.TagAttr()
					if string(k) == "id" && string(v) == "ctl00_ContentPlaceHolder1_ctl00_ctl08_COMPTE_PRET_1_1_MSG_ERREUR" {
						z.Next()
						err = errors.New(z.Token().Data)
						return
					}
					if !more {
						break
					}
				}
			}
		}
	}
	return
}
