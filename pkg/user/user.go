package user

import (
	"fmt"
	"github.com/DevMyong/archeage/pkg/archeage"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
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

	PathSearch = "/search"

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

func getParams(name, server string) (params *url.Values) {
	params = &url.Values{}
	params.Set("dt", "characters")
	params.Add("keyword", name)
	params.Add("server", server)
	return
}
func NewUserInfo(args Args, flagMap FlagMap) (userInfo User, err error) {
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

	docUser, err := getUserData(userInfo.UUID)
	if err != nil {
		return
	}

	userInfo.Character, err = parseUser(docUser)
	if err != nil {
		return
	}

	for opt := range flagMap {
		if opt == "-stat" {
			userInfo.Stat, err = parseUserDetail(docUser)
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
	urlSearchUser := archeage.SetURI(archeage.URLArcheage, PathSearch, params)

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
	urlUsers := archeage.SetURI(archeage.URLArcheage, PathSearch, params)
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

func getUserData(uuid string) (docUser *goquery.Document, err error) {
	aa := archeage.New(http.DefaultClient)
	path := fmt.Sprintf(archeage.PathCharacterFormat, uuid)
	urlUser := archeage.SetURI(archeage.URLArcheage, path, nil)
	docUser, err = aa.Get(urlUser.String())
	if err != nil {
		return
	}
	return
}
func parseUser(docUser *goquery.Document) (character *Character, err error) {
	character = &Character{}

	character.Name = strings.Trim(docUser.Find(archeage.CharacterName).Text(), "\t\n\t ")
	character.Server = docUser.Find(archeage.CharacterServer).Text()[1:]
	character.Union = docUser.Find(archeage.CharacterUnion).Text()
	character.Expedition = docUser.Find(archeage.CharacterExpedition).Find("span").Text()
	character.Score = docUser.Find(archeage.CharacterEquipScore).Text()
	character.SavedAt = time.Now().Format("2006-01-02 15:04")

	if character.Expedition == "" {
		character.Expedition = "-"
	}
	return
}
func parseUserDetail(doc *goquery.Document) (stats *Status, err error) {
	stats = &Status{}
	n := reflect.TypeOf(Status{}).NumField()
	for i := 0; i < n; i++ {
		statQuery := fmt.Sprintf(archeage.CharacterBasicDPS, i+1)
		stat := doc.Find(statQuery).Text()

		idx := strings.Index(stat, "\n")
		if idx == -1 {
			idx = len(stat) - 1
		}

		stat = stat[:idx]
		reflect.ValueOf(stats).Elem().Field(i).SetString(stat)
	}

	return
}
