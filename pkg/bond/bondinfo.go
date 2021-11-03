package bond

import (
	"archecord/pkg/archeage"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func BondParser(serverName string) (continent [2]map[string][]string) {
	doc := archeage.BasicParser("/play/worldinfo/" + serverName)

	// Initialize map
	for i := 0; i < len(continent); i++ {
		continent[i] = make(map[string][]string)
	}

	//
	doc.Find("div.bond-info").Each(func(idx int, sel *goquery.Selection) {
		for i := 0; i < 4; i++ {
			material := sel.Find(fmt.Sprintf(archeage.SearchMaterialName, i+1)).Text()
			westRegions := sel.Find(fmt.Sprintf(archeage.SearchWestRegion, i+1)).Text()
			eastRegions := sel.Find(fmt.Sprintf(archeage.SearchEastRegion, i+1)).Text()

			if westRegions != "" {
				regions := strings.Split(strings.TrimRight(westRegions, ": 20개"), ": 20개")
				continent[archeage.West][material] = regions
			}
			if eastRegions != "" {
				regions := strings.Split(strings.TrimRight(eastRegions, ": 20개"), ": 20개")
				continent[archeage.East][material] = regions
			}

		}
	})

	return continent
}

func RecommendRoute(continent [2]map[string][]string) (route string) {
	usedMaterials := map[string]bool{
		"옷감":   false,
		"가죽":   false,
		"목재":   false,
		"철 주괴": false,
	}

	longitude := 0
	direction := 0
	if len(continent[archeage.West]) > len(continent[archeage.East]) {
		longitude = archeage.West
		direction = 1
	} else {
		longitude = archeage.East
		direction = -1
	}

	for ; longitude < 2 && longitude >= 0; longitude += direction {
		for material, region := range continent[longitude] {
			if usedMaterials[material] {
				continue
			}
			route += region[0] + " -> "
			usedMaterials[material] = true
		}
	}
	route = strings.TrimRight(route, " -> ")
	return
}
