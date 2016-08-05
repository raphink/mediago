package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

type account struct {
	Name     string
	Login    string
	Password string
	Items    []*Item
	Alert    bool
	Client   *http.Client
}

func (a *account) alerts(colored bool, markdown bool) (alerts string) {
	var state string
	for _, i := range a.Items {
		if markdown {
			state = i.State.MarkdownBadge(i.Date)
			alerts += fmt.Sprintf("- [ ] %s %s\n", state, i.Title)
		} else {
			if colored {
				state = i.State.ColoredString()
			} else {
				state = i.State.String()
			}
			alerts += fmt.Sprintf("[%s]\t%s\t%s\n", state, i.Date.Format("02/01/2006"), i.Title)
		}
	}
	return
}

func (a *account) report(cfg *config) {
	titleColor.Println(a.Name)
	fmt.Println(a.alerts(true, false))

	if a.Alert && cfg.Alert == "smtp" {
		SMTPAlert(cfg, a)
	}

	if cfg.Report == "gist" {
		gistReport(cfg, a)
	}
}

func (a *account) getItems() (items []*Item) {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal(err)
	}
	a.Client = &http.Client{Jar: jar}
	resp, err := a.Client.PostForm("http://www.bm-chambery.fr/opacwebaloes/index.aspx?idPage=33", url.Values{
		"ctl00$ScriptManager1":                                          {"ctl00$ContentPlaceHolder1$ctl00$ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$RadAjaxPanelConnexionPanel|ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$btnImgConnexion"},
		"ctl00_ScriptManager1_TSM":                                      {";;System.Web.Extensions, Version=3.5.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35:fr-FR:c2b5a2f3-2711-4e71-b087-b34e92289501:ea597d4b:b25378d2;Telerik.Web.UI, Version=2013.3.1324.35, Culture=neutral, PublicKeyToken=121fae78165ba3d4:fr-FR:84d93921-96f0-4f42-826e-aa3f3f71544e:16e4e7cd:ed16cbdc:874f8ea2:f7645509:24ee1bba:92fe8ea0:fa31b949:f46195d3:19620875:490a9d4e:bd8f85e4:88144a7a"},
		"ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$TextSaisie":        {a.Login},
		"ctl00$ContentPlaceHolder1$ctl00$ctl05$ctl00$TextPass":          {a.Password},
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

	//TODO: check if authentication failed

	resp, err = a.Client.Get("http://www.bm-chambery.fr/opacwebaloes/index.aspx?idPage=478")
	if err != nil {
		log.Fatal(err)
	}

	data := resp.Body
	defer data.Close()

	z := &htmlParser{
		Tokenizer: html.NewTokenizer(data),
	}

	items = z.getItems()
	return
}
