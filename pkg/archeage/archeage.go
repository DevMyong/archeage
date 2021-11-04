package archeage

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Archeage struct {
	client *http.Client
}

func New(c *http.Client) *Archeage {
	return &Archeage{c}
}

func (a *Archeage) Get(url string) (*goquery.Document, error) {
	return a.do("GET", url, nil)
}

func (a *Archeage) post(url string, body io.Reader) (*goquery.Document, error) {
	return a.do("POST", url, body)
}

func (a *Archeage) do(method, url string, body io.Reader) (*goquery.Document, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc, err
}

func SetURI(host string, path string, params *url.Values) url.URL {
	return url.URL{Scheme: SchemeDefault, Host: host, Path: path, RawQuery: params.Encode()}
}
