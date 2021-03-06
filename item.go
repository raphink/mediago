package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

type Item struct {
	Name     string
	Entite   string
	Date     time.Time
	Location string
	Type     string
	Title    string
	Barcode  string
	RentType string
	Booked   string
	State    state
}

func (i *Item) processState(renewBefore time.Duration) (alert bool) {
	now := time.Now()
	renewDate := now.Add(renewBefore)

	if now.After(i.Date) {
		i.State = stateLate
		alert = true
	} else if renewDate.After(i.Date) {
		i.State = stateNeedsRenewing
		alert = true
	} else {
		i.State = stateOK
	}
	return
}

func (i *Item) renew(c *http.Client) (err error) {
	if i.RentType == "R1" || i.Booked != "0" {
		// Already renewed once or booked by someone else
		i.State = stateCannotRenew
		return
	}

	fmt.Println("Autorenewing is not implemented yet")
	return

	resp, err := c.PostForm("http://www.bm-chambery.fr/opacwebaloes/index.aspx?idPage=478", url.Values{
		"ctl00$ScriptManager1":     {"ctl00$ContentPlaceHolder1$ctl00$ctl00$ContentPlaceHolder1$ctl00$RadAjaxPanelPODPanel|ctl00$ContentPlaceHolder1$ctl00$ctl08$COMPTE_PRET_1_1$GrillePrets$ctl00$ctl02$ctl00$btnProlonger"},
		"ctl00_ScriptManager1_TSM": {";;System.Web.Extensions, Version=3.5.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35:fr-FR:c2b5a2f3-2711-4e71-b087-b34e92289501:ea597d4b:b25378d2;Telerik.Web.UI, Version=2013.3.1324.35, Culture=neutral, PublicKeyToken=121fae78165ba3d4:fr-FR:84d93921-96f0-4f42-826e-aa3f3f71544e:16e4e7cd:ed16cbdc:874f8ea2:f7645509:24ee1bba:92fe8ea0:fa31b949:f46195d3:19620875:490a9d4e:bd8f85e4:88144a7a:58366029:2003d0b8:1e771326:aa288e2d"},
		"__VIEWSTATE":              {"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a2lucy9Ta2luX2NoYW1iZXJ5X25ldy9pbWFnZXMvYm91dG9ucy9vZmYvcHJvbG9uZ2VyLmdpZiceC29uTW91c2VPdmVyBVB0aGlzLnNyYz0nL29wYWN3ZWJhbG9lcy9za2lucy9Ta2luX2NoYW1iZXJ5X25ldy9pbWFnZXMvYm91dG9ucy9vbi9wcm9sb25nZXIuZ2lmJ2QCBQ8PFgIfAmhkZAIHDw8WBB8RaB8CaGRkAgEPDxYGHghDc3NDbGFzcwUIIHJnUGFnZXIeBF8hU0ICAh8CaGQWAmYPDxYCHxICCGQWAmYPZBYCAgEPZBYCAgEPZBYIZg9kFgRmDw8WAh8RaGRkAgIPDxYCHxFoZGQCAQ9kFgJmDw8WBB8aBQ1yZ0N1cnJlbnRQYWdlHxsCAmRkAgIPZBYEZg8PFgIfEWhkZAIDDw8WAh8RaGRkAgMPDxYEHxoFEHJnV3JhcCByZ0FkdlBhcnQfGwICZBYEZg8PFgIfD2VkZAIBDxQrAAIPFhoeDFRhYmxlU3VtbWFyeWUeFUVuYWJsZUVtYmVkZGVkU2NyaXB0c2ceCklucHV0VGl0bGVlHhlSZWdpc3RlcldpdGhTY3JpcHRNYW5hZ2VyZx4cRW5hYmxlRW1iZWRkZWRCYXNlU3R5bGVzaGVldGceDFRhYmxlQ2FwdGlvbgUQUGFnZVNpemVDb21ib0JveB8GZx4TRW5hYmxlRW1iZWRkZWRTa2luc2geHE9uQ2xpZW50U2VsZWN0ZWRJbmRleENoYW5nZWQFLlRlbGVyaWsuV2ViLlVJLkdyaWQuQ2hhbmdlUGFnZVNpemVDb21ib0hhbmRsZXIfGwKAAh4TY2FjaGVkU2VsZWN0ZWRWYWx1ZWQeEUVuYWJsZUFyaWFTdXBwb3J0aB4FV2lkdGgbAAAAAAAAR0ABAAAAZA8UKwAEFCsAAg8WBh8PBQIxMB4FVmFsdWUFAjEwHghTZWxlY3RlZGgWAh4Qb3duZXJUYWJsZVZpZXdJZAVHY3RsMDBfQ29udGVudFBsYWNlSG9sZGVyMV9jdGwwMF9jdGwwOF9DT01QVEVfUFJFVF8xXzFfR3JpbGxlUHJldHNfY3RsMDBkFCsAAg8WBh8PBQIyMB8nBQIyMB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZBQrAAIPFgYfDwUCNTAfJwUCNTAfKGgWAh8pBUdjdGwwMF9Db250ZW50UGxhY2VIb2xkZXIxX2N0bDAwX2N0bDA4X0NPTVBURV9QUkVUXzFfMV9HcmlsbGVQcmV0c19jdGwwMGQUKwACDxYGHw8FAjgwHycFAjgwHyhnFgIfKQVHY3RsMDBfQ29udGVudFBsYWNlSG9sZGVyMV9jdGwwMF9jdGwwOF9DT01QVEVfUFJFVF8xXzFfR3JpbGxlUHJldHNfY3RsMDBkDxQrAQRmZmZmFgEFeFRlbGVyaWsuV2ViLlVJLlJhZENvbWJvQm94SXRlbSwgVGVsZXJpay5XZWIuVUksIFZlcnNpb249MjAxMy4zLjEzMjQuMzUsIEN1bHR1cmU9bmV1dHJhbCwgUHVibGljS2V5VG9rZW49MTIxZmFlNzgxNjViYTNkNBYMZg8PFgQfGgUJcmNiSGVhZGVyHxsCAmRkAgEPDxYEHxoFCXJjYkZvb3Rlch8bAgJkZAICDw8WBh8PBQIxMB8nBQIxMB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAIDDw8WBh8PBQIyMB8nBQIyMB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAIEDw8WBh8PBQI1MB8nBQI1MB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAIFDw8WBh8PBQI4MB8nBQI4MB8oZxYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAICD2QWCGYPDxYEHw8FBiZuYnNwOx8CaGRkAgEPDxYEHw8FBiZuYnNwOx8CaGRkAgIPDxYCHw8FBiZuYnNwO2RkAgMPDxYCHwJoZGQCAQ9kFgRmD2QWFmYPDxYCHw8FBiZuYnNwO2RkAgEPDxYCHw8FBiZuYnNwO2RkAgIPDxYCHw8FBiZuYnNwO2RkAgMPDxYCHw8FBiZuYnNwO2RkAgQPDxYCHw8FBiZuYnNwO2RkAgUPDxYCHw8FBiZuYnNwO2RkAgYPDxYCHw8FBiZuYnNwO2RkAgcPDxYCHw8FBiZuYnNwO2RkAggPDxYCHw8FBiZuYnNwO2RkAgkPDxYCHw8FBiZuYnNwO2RkAgoPDxYCHw8FBiZuYnNwO2RkAgEPDxYEHxoFCCByZ1BhZ2VyHxsCAmQWAmYPDxYCHxICCGQWAmYPZBYCAgEPZBYCAgEPZBYIZg9kFgRmDw8WAh8RaGRkAgIPDxYCHxFoZGQCAQ9kFgJmDw8WBB8aBQ1yZ0N1cnJlbnRQYWdlHxsCAmRkAgIPZBYEZg8PFgIfEWhkZAIDDw8WAh8RaGRkAgMPDxYEHxoFEHJnV3JhcCByZ0FkdlBhcnQfGwICZBYEZg8PFgIfD2VkZAIBDxQrAAIPFiAfHGUfHWcfHmUfH2cfDwUCODAfIGcfIQUQUGFnZVNpemVDb21ib0JveB4EU2tpbgUHRGVmYXVsdB8bAoACHwZnHyJoHyMFLlRlbGVyaWsuV2ViLlVJLkdyaWQuQ2hhbmdlUGFnZVNpemVDb21ib0hhbmRsZXIfAGgfJGQfJWgfJhsAAAAAAABHQAEAAABkDxQrAAQUKwACDxYGHw8FAjEwHycFAjEwHyhoFgIfKQVHY3RsMDBfQ29udGVudFBsYWNlSG9sZGVyMV9jdGwwMF9jdGwwOF9DT01QVEVfUFJFVF8xXzFfR3JpbGxlUHJldHNfY3RsMDBkFCsAAg8WBh8PBQIyMB8nBQIyMB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZBQrAAIPFgYfDwUCNTAfJwUCNTAfKGgWAh8pBUdjdGwwMF9Db250ZW50UGxhY2VIb2xkZXIxX2N0bDAwX2N0bDA4X0NPTVBURV9QUkVUXzFfMV9HcmlsbGVQcmV0c19jdGwwMGQUKwACDxYGHw8FAjgwHycFAjgwHyhnFgIfKQVHY3RsMDBfQ29udGVudFBsYWNlSG9sZGVyMV9jdGwwMF9jdGwwOF9DT01QVEVfUFJFVF8xXzFfR3JpbGxlUHJldHNfY3RsMDBkDxQrAQRmZmZmFgEFeFRlbGVyaWsuV2ViLlVJLlJhZENvbWJvQm94SXRlbSwgVGVsZXJpay5XZWIuVUksIFZlcnNpb249MjAxMy4zLjEzMjQuMzUsIEN1bHR1cmU9bmV1dHJhbCwgUHVibGljS2V5VG9rZW49MTIxZmFlNzgxNjViYTNkNBYMZg8PFgQfGgUJcmNiSGVhZGVyHxsCAmRkAgEPDxYEHxoFCXJjYkZvb3Rlch8bAgJkZAICDw8WBh8PBQIxMB8nBQIxMB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAIDDw8WBh8PBQIyMB8nBQIyMB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAIEDw8WBh8PBQI1MB8nBQI1MB8oaBYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAIFDw8WBh8PBQI4MB8nBQI4MB8oZxYCHykFR2N0bDAwX0NvbnRlbnRQbGFjZUhvbGRlcjFfY3RsMDBfY3RsMDhfQ09NUFRFX1BSRVRfMV8xX0dyaWxsZVByZXRzX2N0bDAwZAICDw8WAh8QBQEwFgIeBmVudGl0ZQUIMTAxOTQzODgWFGYPDxYCHwJoZBYCZg8PFgIfEWhkZAIBDw8WBB8PBQYmbmJzcDsfAmhkZAIDDw8WBB8PBQYmbmJzcDsfAmhkZAIEDw8WAh8PBQoxMi8wOC8yMDE2ZGQCBQ8PFgIfDwUeQmlibGlvdGjDqHF1ZSBHZW9yZ2VzIEJyYXNzZW5zZGQCBg8PFgIfDwUFTGl2cmVkZAIHDw8WAh8PBUc5IG1vaXMgZGFucyBsZSB2ZW50cmUgZGUgTWFtYW4gLyBOYXRhY2hhIEZyYWRpbi4gLSBNaWxhbiBqZXVuZXNzZSwgMjAwNmRkAggPDxYCHw8FBzE0NDcwNTZkZAIJDw8WAh8PBQJSMWRkAgoPDxYCHw8FATBkZAIDD2QWAmYPDxYCHwJoZGQCBA8PFgIfEAUBMRYCHysFCDEwMTk0Mzg5FhRmDw8WAh8CaGQWAmYPDxYCHxFoZGQCAQ8PFgQfDwUGJm5ic3A7HwJoZGQCAw8PFgQfDwUGJm5ic3A7HwJoZGQCBA8PFgIfDwUKMTIvMDgvMjAxNmRkAgUPDxYCHw8FHkJpYmxpb3Row6hxdWUgR2VvcmdlcyBCcmFzc2Vuc2RkAgYPDxYCHw8FBUxpdnJlZGQCBw8PFgIfDwVDRG9yYSBjaGV6IGxlIGRvY3RldXIgLyBOaWNrZWxvZGVvbiBwcm9kdWN0aW9uLiAtIEFsYmluIE1pY2hlbCwgMjAxNGRkAggPDxYCHw8FBzE1NjY4NTFkZAIJDw8WAh8PBQJSMWRkAgoPDxYCHw8FATBkZAIFD2QWAmYPDxYCHwJoZGQCBg8PFgIfEAUBMhYCHysFCDEwMTk0MzkwFhRmDw8WAh8CaGQWAmYPDxYCHxFoZGQCAQ8PFgQfDwUGJm5ic3A7HwJoZGQCAw8PFgQfDwUGJm5ic3A7HwJoZGQCBA8PFgIfDwUKMTIvMDgvMjAxNmRkAgUPDxYCHw8FHkJpYmxpb3Row6hxdWUgR2VvcmdlcyBCcmFzc2Vuc2RkAgYPDxYCHw8FBUxpdnJlZGQCBw8PFgIfDwVGVHUgY3JpZXMgY29tbWVudCwgQ2hpZW4gbGUgY2hpZW4gPyAvIE1vIFdpbGxlbXMuIC0gS2Fsw6lpZG9zY29wZSwgMjAxM2RkAggPDxYCHw8FBzE1NTIyNjFkZAIJDw8WAh8PBQJSMWRkAgoPDxYCHw8FATBkZAIHD2QWAmYPDxYCHwJoZGQCCg9kFgICBw9kFgJmD2QWAmYPZBYCAgMPPCsADgIAFCsAAg8WBh8AaB8FAgEfBmdkFwMFD1NlbGVjdGVkSW5kZXhlcxYABQtFZGl0SW5kZXhlcxYABQhQYWdlU2l6ZQJQARYCFgtkZGUUKwAACysEATwrAAcACysFAWRkFggfC2gfDBYAHwZnHwVmZGYWBGYPFCsAA2RkZGQCAQ8WBRQrAAIPFggfC2gfDBYAHwZnHwVmZBcGBRBDdXJyZW50UGFnZUluZGV4ZgUGXyFEU0lDZgUDX2ZlZQUIXyFQQ291bnQCAQULXyFJdGVtQ291bnRmBQhQYWdlU2l6ZQJQFgIfDRYCHw5mZBYAFgJnZxYCZg9kFgJmD2QWBmYPDxYCHwJoZGQCAQ8PFgIfAmhkZAICDw8WAh8SZhYCHxMFEHRleHQtYWxpZ246bGVmdDtkAgMPFgIfAmgWAmYPFgIfAmhkAgYPFgIfAmcWAmYPFgQfAwUBMx8EZGQYBgVkY3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwOCRDT01QVEVfUFJFVF8xXzEkR3JpbGxlUHJldHMkY3RsMDAkY3RsMDMkY3RsMDEkUGFnZVNpemVDb21ib0JveA8UKwACBQI4MAUCODBkBR5fX0NvbnRyb2xzUmVxdWlyZVBvc3RCYWNrS2V5X18WDwUPY3RsMDAkU2luZ2xldG9uBUdjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDA0JEJvdXRvblJlY2hlcmNoZXIkQm91dG9uUmVjaGVyY2hlcgU/Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwNiRDT01QVEVfSU5GT1NfJEdyaWxsZUluZm9zBTtjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDA3JGN0bDAwJEJvdXRvblJTU0NvbXB0ZQVFY3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwNyRDT01QVEVfUkVUQVJEXzFfMSRHcmlsbGVSZXRhcmRzBTtjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDA4JGN0bDAwJEJvdXRvblJTU0NvbXB0ZQVBY3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwOCRDT01QVEVfUFJFVF8xXzEkR3JpbGxlUHJldHMFYGN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDgkQ09NUFRFX1BSRVRfMV8xJEdyaWxsZVByZXRzJGN0bDAwJGN0bDAyJGN0bDAwJGJ0blByb2xvbmdlcgV0Y3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwOCRDT01QVEVfUFJFVF8xXzEkR3JpbGxlUHJldHMkY3RsMDAkY3RsMDIkY3RsMDIkQ2xpZW50U2VsZWN0Q29sdW1uU2VsZWN0Q2hlY2tCb3gFZGN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDgkQ09NUFRFX1BSRVRfMV8xJEdyaWxsZVByZXRzJGN0bDAwJGN0bDAzJGN0bDAxJFBhZ2VTaXplQ29tYm9Cb3gFbmN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDgkQ09NUFRFX1BSRVRfMV8xJEdyaWxsZVByZXRzJGN0bDAwJGN0bDA0JENsaWVudFNlbGVjdENvbHVtblNlbGVjdENoZWNrQm94BW5jdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDA4JENPTVBURV9QUkVUXzFfMSRHcmlsbGVQcmV0cyRjdGwwMCRjdGwwNiRDbGllbnRTZWxlY3RDb2x1bW5TZWxlY3RDaGVja0JveAVuY3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwOCRDT01QVEVfUFJFVF8xXzEkR3JpbGxlUHJldHMkY3RsMDAkY3RsMDgkQ2xpZW50U2VsZWN0Q29sdW1uU2VsZWN0Q2hlY2tCb3gFO2N0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDkkY3RsMDAkQm91dG9uUlNTQ29tcHRlBUxjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDA5JENPTVBURV9QUkVUQVRUXzBfMSRHcmlsbGVQcmV0c0F0dGVuZHVzBUxjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDAzJGN0bDAwJENyaXRlcmVfMTA0XzU0N19udW0kY2JDcm9pc2VtZW50DxQrAAJlZWQFZGN0bDAwJENvbnRlbnRQbGFjZUhvbGRlcjEkY3RsMDAkY3RsMDgkQ09NUFRFX1BSRVRfMV8xJEdyaWxsZVByZXRzJGN0bDAwJGN0bDAyJGN0bDAxJFBhZ2VTaXplQ29tYm9Cb3gPFCsAAmUFAjgwZAVMY3RsMDAkQ29udGVudFBsYWNlSG9sZGVyMSRjdGwwMCRjdGwwMyRjdGwwMCRDcml0ZXJlXzEwNF81NDdfbnVtJGNiT3BlcmF0ZXVycw8UKwACZWVkBUpjdGwwMCRDb250ZW50UGxhY2VIb2xkZXIxJGN0bDAwJGN0bDAzJGN0bDAwJENyaXRlcmVfMTA0XzU0N19udW0kY2JDcml0ZXJlcw8UKwACZWVkgU1xBDtk9OeKVaaL9a/m1aqIKVs="},
		"__VIEWSTATEGENERATOR":     {"F8F44887"},
		"__EVENTVALIDATION":        {"/wEWFwLEjY7SCQLo+JuXDgLn/5usBAL/y8uHDgLeqoPSDQLrifOeCALbzvW5BgLh1IqgAQLi1IqgAQLj1IqgAQLc1IqgAQLd1IqgAQLe1IqgAQLf1IqgAQLMhd7WAgLLhd7WAgLRhd7WAgLWhd7WAgLVhd7WAgKzqdOdAwL9lb6DBwKzkae3AQK+6+PtDmFn5vCGOO3+g8PxRzCqt+J+g3Ct"},
		"ctl00$ContentPlaceHolder1$ctl00$ctl08$COMPTE_PRET_1_1$GrillePrets$ctl00$ctl03$ctl01$PageSizeComboBox": {"80"},
		i.Name:        {"on"},
		"__ASYNCPOST": {"true"},
		"ctl00$ContentPlaceHolder1$ctl00$ctl08$COMPTE_PRET_1_1$GrillePrets$ctl00$ctl02$ctl00$btnProlonger.x": {"57"},
		"ctl00$ContentPlaceHolder1$ctl00$ctl08$COMPTE_PRET_1_1$GrillePrets$ctl00$ctl02$ctl00$btnProlonger.y": {"2"},
		"RadAJAXControlID": {"ctl00_ContentPlaceHolder1_ctl00_RadAjaxPanelPOD"},
	})
	if err != nil {
		err = fmt.Errorf("failed to commit renewal request: %v", err)
		return
	}

	data := resp.Body
	defer data.Close()

	z := &htmlParser{
		Tokenizer: html.NewTokenizer(data),
	}

	err = z.checkError()

	if err != nil {
		i.State = stateFailedRenewing
	} else {
		i.State = stateRenewed
	}

	return
}
