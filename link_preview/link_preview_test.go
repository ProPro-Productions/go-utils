package link_preview

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createMockServer(handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)

	return server
}

func TestGetLinkPreviewItems(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title>Test Page</title></head><body><h1>Welcome to the Test Page</h1></body></html>`))
	})
	defer server.Close()

	link := server.URL
	doc, err := GetLinkPreviewItems(link, 10)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, "Test Page", doc.Preview.Title)
}

func TestGetLinkPreviewItemsFail(t *testing.T) {
	link := "invalidurl"
	_, err := GetLinkPreviewItems(link, 10)

	assert.Error(t, err)
}

func TestToFragmentUrl(t *testing.T) {
	t.Run("fragmented url", func(t *testing.T) {
		link := "http://test.com/#!data"
		url, _ := url.Parse(link)
		scraper := &Scraper{Url: url, MaxRedirect: 10}
		err := scraper.toFragmentUrl()

		assert.NoError(t, err)
		assert.NotNil(t, scraper.EscapedFragmentUrl)
		assert.Equal(t, "http://test.com/?_escaped_fragment_=data", scraper.EscapedFragmentUrl.String())
	})

	t.Run("non-fragmented url", func(t *testing.T) {
		link := "http://test.com/data"
		url, _ := url.Parse(link)
		scraper := &Scraper{Url: url, MaxRedirect: 10}
		err := scraper.toFragmentUrl()

		assert.NoError(t, err)
		assert.Nil(t, scraper.EscapedFragmentUrl)
	})
}

func TestGetDocument(t *testing.T) {
	server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html><head><title>Test Page</title></head><body><h1>Welcome to the Test Page</h1></body></html>`))
	})
	defer server.Close()

	link := server.URL
	url, _ := url.Parse(link)
	scraper := &Scraper{Url: url, MaxRedirect: 10}
	doc, err := scraper.getDocument()

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	//assert.Equal(t, "Welcome to the Test Page", doc.Body.Text())
}

func TestGetDocumentFail(t *testing.T) {
	link := "http://invalidurl.com"
	url, _ := url.Parse(link)
	scraper := &Scraper{Url: url, MaxRedirect: 10}
	_, err := scraper.getDocument()

	assert.Error(t, err)
}

func TestParseDocument(t *testing.T) {
	//link := "http://test.com"
	//url, _ := url.Parse(link)
	//scraper := &Scraper{Url: url, MaxRedirect: 10}

	// Prepare a test document for parsing
	goquery.NewDocumentFromReader(strings.NewReader(`<html><head><title>Test Page</title></head><body><h1>Welcome to the Test Page</h1></body></html>`))

	//document := scraper.parseDocument(doc)
	//
	//assert.Equal(t, "Test Page", document.PageInfo.PageTitle)
}
