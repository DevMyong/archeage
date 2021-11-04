package archeage

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

const (
	Token    string = ""
	BasicURL string = "https://archeage.xlgames.com"

	West = 0
	East = 1

	SearchMaterialName = "table:nth-child(1) > tbody > tr:nth-child(%d) > th"
	SearchWestRegion   = "table:nth-child(1) > tbody > tr:nth-child(%d) > td > ul > li.point"
	SearchEastRegion   = "table.table-bond.right > tbody > tr:nth-child(%d) > td > ul > li.point"

	SchemeDefault = "https"
	URLArcheage   = "archeage.xlgames.com"
)

var ServerNameMap = map[string]string{
	"누이":    	"NUI",
	"하제":    	"HAJE",
	"다후타": 	"DAHUTA",
	"모르페우스":	"MORPHEUS",
	"랑그레이": 	"RANGORA",
	"환락":    	"SEASON",
}

func BasicParser(url string) (doc *goquery.Document) {
	resp, err := http.Get(BasicURL + url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("StatusCode Error: %d %s", resp.StatusCode, resp.Status)
	}

	// Load HTML
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return
}