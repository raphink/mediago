package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"

	"github.com/0xAX/notificator"
	"github.com/BurntSushi/toml"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type Item struct {
	Entite   string
	Date     time.Time
	Location string
	Type     string
	Title    string
	Barcode  string
	RentType string
	Booked   string
}

type account struct {
	Name     string
	Login    string
	Password string
}

type config struct {
	Account     []account
	RenewBefore duration `toml:"renew_before"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return
}

var confFile = fmt.Sprintf("%s/.mediago.conf", os.Getenv("HOME"))

func main() {
	var cfg config
	if _, err := toml.DecodeFile(confFile, &cfg); err != nil {
		log.Fatal(err)
	}

	for _, a := range cfg.Account {
		i := getAccountItems(a.Name, a.Login, a.Password)

		fmt.Printf("# %s\n", a.Name)
		now := time.Now()
		tomorrow := now.Add(cfg.RenewBefore.Duration)
		var alert bool
		var output string
		for _, e := range i {
			var state string
			if now.After(e.Date) {
				state = "!!LATE!!"
				alert = true
			} else if tomorrow.After(e.Date) {
				state = "NEEDS RENEWING"
				alert = true
				// TODO: automatically renew
			} else {
				state = "OK"
			}
			//fmt.Printf("[%s] %s: %s, %s, %s, %s, %s, %s, %s\n", state, e.Entite, e.Date, e.Location, e.Type, e.Title, e.Barcode, e.RentType, e.Booked)
			output += fmt.Sprintf("[%s]\t%s\t%s\n", state, e.Date.Format("02/01/2006"), e.Title)
		}
		fmt.Println(output)

		if alert {
			notify := notificator.New(notificator.Options{
				DefaultIcon: "icon/dialog-warning.png",
				AppName:     "Mediago",
			})

			title := fmt.Sprintf("Mediathèque %s", a.Name)
			notify.Push(title, output, "/usr/share/icons/gnome/32x32/status/gtk-dialog-warning.png", notificator.UR_CRITICAL)
		}
	}
}

func getAccountItems(name, account, password string) (items []*Item) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{Jar: jar}
	resp, err := client.PostForm("http://www.bm-chambery.fr/opacwebaloes/index.aspx?idPage=33", url.Values{
		"ctl00$ScriptManager1":                                          {"ctl00$ContentPlaceHolder1$ctl00$ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$RadAjaxPanelConnexionPanel|ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$btnImgConnexion"},
		"ctl00_ScriptManager1_TSM":                                      {";;System.Web.Extensions, Version=3.5.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35:fr-FR:c2b5a2f3-2711-4e71-b087-b34e92289501:ea597d4b:b25378d2;Telerik.Web.UI, Version=2013.3.1324.35, Culture=neutral, PublicKeyToken=121fae78165ba3d4:fr-FR:84d93921-96f0-4f42-826e-aa3f3f71544e:16e4e7cd:ed16cbdc:874f8ea2:f7645509:24ee1bba:92fe8ea0:fa31b949:f46195d3:19620875:490a9d4e:bd8f85e4:88144a7a"},
		"ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$TextSaisie":        {account},
		"ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$TextPass":          {password},
		"RadAJAXControlID":                                              {"ctl00_ContentPlaceHolder1_ctl00_ctl05_ctl00_RadAjaxPanelConnexion"},
		"ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$btnImgConnexion.x": {"0"},
		"ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$btnImgConnexion.y": {"0"},
		"__VIEWSTATE":          {"/wEPDwUKLTMxNjc3NTM3NQ9kFgJmD2QWAgIFD2QWBAIFDxQrAAIUKwADDxYCHhdFbmFibGVBamF4U2tpblJlbmRlcmluZ2hkZGRkZAIHD2QWAgIBD2QWAmYPZBYGAgIPDxYCHwBoZGQCBA8WAh4Fd2lkdGgFBDEwMjQWCGYPFgIeB1Zpc2libGVnFgJmDxYGHgdjb2xzcGFuBQEzHgZoZWlnaHRkHwJnFgICBA9kFgJmD2QWAgICD2QWBgIBDxQrAAJkZGQCAw8UKwACZGRkAgcPFCsAAmQQFgFmFgEUKwACZGQPFgFmFgEFeFRlbGVyaWsuV2ViLlVJLlJhZENvbWJvQm94SXRlbSwgVGVsZXJpay5XZWIuVUksIFZlcnNpb249MjAxMy4zLjEzMjQuMzUsIEN1bHR1cmU9bmV1dHJhbCwgUHVibGljS2V5VG9rZW49MTIxZmFlNzgxNjViYTNkNBYCAgIPFCsAAWRkAgEPFgIfAmgWAmYPFgIfAmhkAgIPFgIfBAUEMTAwJRYKZg8WAh8CaGQCAQ8WAh8CaGQCAw8WAh8CaGQCBA8WAh8CaGQCBg8WAh8BBQQxMDI0FgQCAQ8PFgIfAmhkZAIDD2QWAgIBD2QWAgIBD2QWAgIHD2QWAmYPZBYEAgIPDxYCHwBoZGQCBA9kFgICEA8PFgIeEVVzZVN1Ym1pdEJlaGF2aW9yaGRkAgMPFgIfAmgWAmYPFgIfAmhkAgYPFgIfAmcWAmYPFgQfAwUBMx8EZGQYBAUeX19Db250cm9sc1JlcXVpcmVQb3N0QmFja0tleV9fFgMFD2N0bDAwJFNpbmdsZXRvbgVHY3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwNCRCb3V0b25SZWNoZXJjaGVyJEJvdXRvblJlY2hlcmNoZXIFO2N0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDUkY3RsMDAkYnRuSW1nQ29ubmV4aW9uBUxjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDAzJGN0bDAwJENyaXRlcmVfMTA0XzU0N19udW0kY2JPcGVyYXRldXJzDxQrAAJlZWQFSmN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDMkY3RsMDAkQ3JpdGVyZV8xMDRfNTQ3X251bSRjYkNyaXRlcmVzDxQrAAJlZWQFTGN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDMkY3RsMDAkQ3JpdGVyZV8xMDRfNTQ3X251bSRjYkNyb2lzZW1lbnQPFCsAAmVlZGLUfWZKFAChQXBXPiHZHjw+aMa2"},
		"__VIEWSTATEGENERATOR": {"F8F44887"},
		"__EVENTVALIDATION":    {"/wEWBgLAucfHCALo+JuXDgLn/5usBALDgOm0AwKzr8rjCgKvhOvWBzYYzUEJnaxbNve47aiHYXI9Ma41"},
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err = client.Get("http://www.bm-chambery.fr/opacwebaloes/index.aspx?idPage=478")
	if err != nil {
		log.Fatal(err)
	}

	data := resp.Body
	defer data.Close()

	z := html.NewTokenizer(data)

	items = getItems(z)
	return
}

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
