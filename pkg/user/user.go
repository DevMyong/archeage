package user

import (
	"fmt"
	"github.com/DevMyong/archeage/pkg/archeage"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type User struct {
	UUID      string
	Character *Character
	Stat      *Status
	Old       []Character
}
type Character struct {
	Name       string
	Server     string
	Union      string
	Expedition string
	Score      string
	SavedAt    string
}
type Status struct {
	RangeDPS        string
	MeleeDPS        string
	MagicDPS        string
	HealingDPS      string
	PhysicalDefense string
	MagicDefense    string
}

type FlagMap map[string]string
type Args []string

const (
	sep           = "@"
	DefaultServer = "MORPHEUS"

	SchemeDefault = "https"
	URLArcheage   = "archeage.xlgames.com"
	PathSearch    = "/search"

	PageBar             = ".pg_inner"
	CharacterList       = ".s_char_lst"
	CharacterCard       = ".character_card"
	CharacterNameInList = ".character_name"
)

var ServerNameKorToEng = map[string]string{
	"누이":    "NUI",
	"하제":    "HAJE",
	"다후타":   "DAHUTA",
	"모르페우스": "MORPHEUS",
	"랑그레이":  "RANGORA",
	"환락":    "SEASON",
}
var ServerNameEngToKor = map[string]string{
	"NUI":      "누이",
	"HAJE":     "하제",
	"DAHUTA":   "다후타",
	"MORPHEUS": "모르페우스",
	"RANGORA":  "랑그레이",
	"SEASON":   "환락",
}
func getParams(name, server string)(params *url.Values){
	params = &url.Values{}
	params.Set("dt", "characters")
	params.Add("keyword", name)
	params.Add("server", server)
	return
}
func GetUserInfo(args Args, flagMap FlagMap) (userInfo User, err error) {
	name, server, err := getNameAndServer(args[0])
	if err != nil {
		return
	}

	params := getParams(name, server)

	docUsers, err := searchUser(params)
	if err != nil {
		return
	}

	totalPages := getTotalPage(docUsers)

	for i := 0; i < totalPages && userInfo.UUID == ""; i++ {
		params.Set("page", strconv.Itoa(i+1))

		userInfo.UUID, err = findUUID(params)
		if err != nil {
			return
		}
	}

	if userInfo.UUID == "" {
		err = fmt.Errorf("Can't find %s@%s ", name, ServerNameEngToKor[server])
		return
	}

	for opt, _ := range flagMap {
		switch opt {
		case "-stat":
			//getUserStatus(userInfo.UUID)
		case "-save":
			//getUserSummary(userInfo.UUID)
			//saveUserSummary()
		case "-load":
			//getUserSummary(userInfo.UUID)
		case "-diff":
			//getUserSummary(userInfo.UUID)
			//getUserHistory(userInfo.UUID)
		case "-history":
			//getUserHistory(userInfo.UUID)
		}
	}
	return
}

func getNameAndServer(args string) (name, server string, err error) {
	userInfo := strings.Split(args, sep)
	l := len(userInfo)

	name = userInfo[0]

	if l > 2 {
		err = fmt.Errorf("Too many '@' separator in name field. Please use this format: name@server ")
		return
	} else if l == 2 {
		server = userInfo[1]
	}

	if serverEng, ok := ServerNameKorToEng[server]; ok {
		server = serverEng
	} else {
		server = DefaultServer
	}

	return
}
func searchUser(params *url.Values) (docUsers *goquery.Document, err error) {
	urlSearchUser := archeage.SetURI(URLArcheage, PathSearch, params)

	aa := archeage.New(http.DefaultClient)

	docUsers, err = aa.Get(urlSearchUser.String())
	if err != nil {
		return
	}

	return
}
func getTotalPage(doc *goquery.Document) int {
	totalPages := 0

	doc.Find(PageBar).Each(func(i int, sel *goquery.Selection) {
		totalPages = sel.Find("a").Length() - 1
	})

	if totalPages == 0 {
		return 1
	}
	return totalPages
}

func findUUID(params *url.Values) (uuid string, err error) {
	aa := archeage.New(http.DefaultClient)
	urlUsers := archeage.SetURI(URLArcheage, PathSearch, params)
	docUsers, err := aa.Get(urlUsers.String())
	if err != nil {
		return
	}

	charList := docUsers.Find(CharacterList)

	charList.Find(CharacterCard).Each(func(idx int, card *goquery.Selection) {
		searchedName := strings.TrimSpace(card.Find(CharacterNameInList).Text())

		if searchedName == params.Get("keyword") {
			link, _ := card.Find("a").Attr("href")
			path := strings.Split(link, "/")
			uuid = path[len(path)-1]
		}
	})

	return
}
