package search

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"errors"
	"golang.org/x/time/rate"
)

// ErrBlocked indicates that Google has detected that you were scraping and temporarily blocked you.
// The duration of the block is unspecified.
//
// See: https://github.com/rocketlaunchr/google-search#warning-warning
var ErrBlocked = errors.New("google block")

// RateLimit sets a global limit to how many requests to Google Search can be made in a given time interval.
// The default is unlimited (but obviously Google Search will block you temporarily if you do too many
// calls too quickly).
//
// See: https://godoc.org/golang.org/x/time/rate#NewLimiter
var RateLimit = rate.NewLimiter(rate.Inf, 0)

// Result represents a single result from Google Search.
type Result struct {

	// Rank is the order number of the search result.
	Rank int `json:"rank"`

	// URL of result.
	URL string `json:"url"`

	// Title of result.
	Title string `json:"title"`

	// Description of the result.
	Description string `json:"description"`
}

const stdGoogleBase = "https://www.google."

// GoogleDomains represents localized Google homepages. The 2 letter country code is based on ISO 3166-1 alpha-2.
//
// See: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
var GoogleDomains = map[string]string{
	"us":  "com/search?q=",
	"ac":  "ac/search?q=",
	"ad":  "ad/search?q=",
	"ae":  "ae/search?q=",
	"af":  "com.af/search?q=",
	"ag":  "com.ag/search?q=",
	"ai":  "com.ai/search?q=",
	"al":  "al/search?q=",
	"am":  "am/search?q=",
	"ao":  "co.ao/search?q=",
	"ar":  "com.ar/search?q=",
	"as":  "as/search?q=",
	"at":  "at/search?q=",
	"au":  "com.au/search?q=",
	"az":  "az/search?q=",
	"ba":  "ba/search?q=",
	"bd":  "com.bd/search?q=",
	"be":  "be/search?q=",
	"bf":  "bf/search?q=",
	"bg":  "bg/search?q=",
	"bh":  "com.bh/search?q=",
	"bi":  "bi/search?q=",
	"bj":  "bj/search?q=",
	"bn":  "com.bn/search?q=",
	"bo":  "com.bo/search?q=",
	"br":  "com.br/search?q=",
	"bs":  "bs/search?q=",
	"bt":  "bt/search?q=",
	"bw":  "co.bw/search?q=",
	"by":  "by/search?q=",
	"bz":  "com.bz/search?q=",
	"ca":  "ca/search?q=",
	"kh":  "com.kh/search?q=",
	"cc":  "cc/search?q=",
	"cd":  "cd/search?q=",
	"cf":  "cf/search?q=",
	"cat": "cat/search?q=",
	"cg":  "cg/search?q=",
	"ch":  "ch/search?q=",
	"ci":  "ci/search?q=",
	"ck":  "co.ck/search?q=",
	"cl":  "cl/search?q=",
	"cm":  "cm/search?q=",
	"cn":  "cn/search?q=",
	"co":  "com.co/search?q=",
	"cr":  "co.cr/search?q=",
	"cu":  "com.cu/search?q=",
	"cv":  "cv/search?q=",
	"cy":  "com.cy/search?q=",
	"cz":  "cz/search?q=",
	"de":  "de/search?q=",
	"dj":  "dj/search?q=",
	"dk":  "dk/search?q=",
	"dm":  "dm/search?q=",
	"do":  "com.do/search?q=",
	"dz":  "dz/search?q=",
	"ec":  "com.ec/search?q=",
	"ee":  "ee/search?q=",
	"eg":  "com.eg/search?q=",
	"es":  "es/search?q=",
	"et":  "com.et/search?q=",
	"fi":  "fi/search?q=",
	"fj":  "com.fj/search?q=",
	"fm":  "fm/search?q=",
	"fr":  "fr/search?q=",
	"ga":  "ga/search?q=",
	"gb":  "co.uk/search?q=",
	"ge":  "ge/search?q=",
	"gf":  "gf/search?q=",
	"gg":  "gg/search?q=",
	"gh":  "com.gh/search?q=",
	"gi":  "com.gi/search?q=",
	"gl":  "gl/search?q=",
	"gm":  "gm/search?q=",
	"gp":  "gp/search?q=",
	"gr":  "gr/search?q=",
	"gt":  "com.gt/search?q=",
	"gy":  "gy/search?q=",
	"hk":  "com.hk/search?q=",
	"hn":  "hn/search?q=",
	"hr":  "hr/search?q=",
	"ht":  "ht/search?q=",
	"hu":  "hu/search?q=",
	"id":  "co.id/search?q=",
	"iq":  "iq/search?q=",
	"ie":  "ie/search?q=",
	"il":  "co.il/search?q=",
	"im":  "im/search?q=",
	"in":  "co.in/search?q=",
	"io":  "io/search?q=",
	"is":  "is/search?q=",
	"it":  "it/search?q=",
	"je":  "je/search?q=",
	"jm":  "com.jm/search?q=",
	"jo":  "jo/search?q=",
	"jp":  "co.jp/search?q=",
	"ke":  "co.ke/search?q=",
	"ki":  "ki/search?q=",
	"kg":  "kg/search?q=",
	"kr":  "co.kr/search?q=",
	"kw":  "com.kw/search?q=",
	"kz":  "kz/search?q=",
	"la":  "la/search?q=",
	"lb":  "com.lb/search?q=",
	"lc":  "com.lc/search?q=",
	"li":  "li/search?q=",
	"lk":  "lk/search?q=",
	"ls":  "co.ls/search?q=",
	"lt":  "lt/search?q=",
	"lu":  "lu/search?q=",
	"lv":  "lv/search?q=",
	"ly":  "com.ly/search?q=",
	"ma":  "co.ma/search?q=",
	"md":  "md/search?q=",
	"me":  "me/search?q=",
	"mg":  "mg/search?q=",
	"mk":  "mk/search?q=",
	"ml":  "ml/search?q=",
	"mm":  "com.mm/search?q=",
	"mn":  "mn/search?q=",
	"ms":  "ms/search?q=",
	"mt":  "com.mt/search?q=",
	"mu":  "mu/search?q=",
	"mv":  "mv/search?q=",
	"mw":  "mw/search?q=",
	"mx":  "com.mx/search?q=",
	"my":  "com.my/search?q=",
	"mz":  "co.mz/search?q=",
	"na":  "com.na/search?q=",
	"ne":  "ne/search?q=",
	"nf":  "com.nf/search?q=",
	"ng":  "com.ng/search?q=",
	"ni":  "com.ni/search?q=",
	"nl":  "nl/search?q=",
	"no":  "no/search?q=",
	"np":  "com.np/search?q=",
	"nr":  "nr/search?q=",
	"nu":  "nu/search?q=",
	"nz":  "co.nz/search?q=",
	"om":  "com.om/search?q=",
	"pa":  "com.pa/search?q=",
	"pe":  "com.pe/search?q=",
	"ph":  "com.ph/search?q=",
	"pk":  "com.pk/search?q=",
	"pl":  "pl/search?q=",
	"pg":  "com.pg/search?q=",
	"pn":  "pn/search?q=",
	"pr":  "com.pr/search?q=",
	"ps":  "ps/search?q=",
	"pt":  "pt/search?q=",
	"py":  "com.py/search?q=",
	"qa":  "com.qa/search?q=",
	"ro":  "ro/search?q=",
	"rs":  "rs/search?q=",
	"ru":  "ru/search?q=",
	"rw":  "rw/search?q=",
	"sa":  "com.sa/search?q=",
	"sb":  "com.sb/search?q=",
	"sc":  "sc/search?q=",
	"se":  "se/search?q=",
	"sg":  "com.sg/search?q=",
	"sh":  "sh/search?q=",
	"si":  "si/search?q=",
	"sk":  "sk/search?q=",
	"sl":  "com.sl/search?q=",
	"sn":  "sn/search?q=",
	"sm":  "sm/search?q=",
	"so":  "so/search?q=",
	"st":  "st/search?q=",
	"sv":  "com.sv/search?q=",
	"td":  "td/search?q=",
	"tg":  "tg/search?q=",
	"th":  "co.th/search?q=",
	"tj":  "com.tj/search?q=",
	"tk":  "tk/search?q=",
	"tl":  "tl/search?q=",
	"tm":  "tm/search?q=",
	"to":  "to/search?q=",
	"tn":  "tn/search?q=",
	"tr":  "com.tr/search?q=",
	"tt":  "tt/search?q=",
	"tw":  "com.tw/search?q=",
	"tz":  "co.tz/search?q=",
	"ua":  "com.ua/search?q=",
	"ug":  "co.ug/search?q=",
	"uk":  "co.uk/search?q=",
	"uy":  "com.uy/search?q=",
	"uz":  "co.uz/search?q=",
	"vc":  "com.vc/search?q=",
	"ve":  "co.ve/search?q=",
	"vg":  "vg/search?q=",
	"vi":  "co.vi/search?q=",
	"vn":  "com.vn/search?q=",
	"vu":  "vu/search?q=",
	"ws":  "ws/search?q=",
	"za":  "co.za/search?q=",
	"zm":  "co.zm/search?q=",
	"zw":  "co.zw/search?q=",
}

// SearchOptions modifies how the Search function behaves.
type SearchOptions struct {

	// CountryCode sets the ISO 3166-1 alpha-2 code of the localized Google Search homepage to use.
	// The default is "us", which will return results from https://www.google.com.
	CountryCode string

	// LanguageCode sets the language code.
	// Default: en
	LanguageCode string

	// Limit sets how many results to fetch (at maximum).
	Limit int

	// Start sets from what rank the new result set should return.
	Start int

	// UserAgent sets the UserAgent of the http request.
	// Default: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"
	UserAgent string

	// OverLimit searches for more results than that specified by Limit.
	// It then reduces the returned results to match Limit.
	OverLimit bool

	// ProxyAddr sets a proxy address to avoid IP blocking.
	ProxyAddr string
}

// SearchGoogle returns a list of search results from Google.
func SearchGoogle(ctx context.Context, searchTerm string, opts ...SearchOptions) ([]Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := RateLimit.Wait(ctx); err != nil {
		return nil, err
	}

	// Checking options
	opt := SearchOptions{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.UserAgent == "" {
		opt.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"
	}

	client := &http.Client{}
	log.Println(getSearchURL(searchTerm, opt))
	req, err := http.NewRequest("GET", getSearchURL(searchTerm, opt), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", opt.UserAgent)
	if opt.ProxyAddr != "" {
		proxyUrl, _ := url.Parse(opt.ProxyAddr)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Received non-200 response code")
		return nil, fmt.Errorf("Received non-200 response code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	results, err := parseResults(resp)
	if err != nil {
		log.Println("Error parsing results")
		return nil, err
	}

	return results, nil
}

func containsAny(text string, values ...string) bool {
	for _, value := range values {
		if strings.Contains(text, value) {
			return true
		}
	}
	return false
}

func parseResults(resp *http.Response) ([]Result, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []Result
	s := doc.Find(".g")
	rank := 1

	log.Println("extract s", s)

	s.Each(func(i int, el *goquery.Selection) {
		result := Result{}

		result.Rank = rank
		rank++

		titleEl := el.Find(".LC20lb.DKV0Md")
		result.Title = titleEl.Text()

		link, _ := titleEl.Attr("href")
		parsedLink, _ := url.Parse(link)
		result.URL = "http://" + parsedLink.Host + parsedLink.Path

		desc := el.Find(".aCOpRe span")
		result.Description = desc.Text()

		results = append(results, result)
	})

	return results, nil
}

func getSearchURL(searchTerm string, opts SearchOptions) string {
	base := stdGoogleBase
	if val, ok := GoogleDomains[opts.CountryCode]; ok {
		base += val
	} else {
		base += GoogleDomains["us"]
	}

	query := url.QueryEscape(searchTerm)
	return base + query + "&hl=" + opts.LanguageCode + "&start=" + strconv.Itoa(opts.Start)
}

func base(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	} else {
		return stdGoogleBase + url
	}
}

//func url(searchTerm string, countryCode string, languageCode string, limit int, start int) string {
//	searchTerm = strings.Trim(searchTerm, " ")
//	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
//	countryCode = strings.ToLower(countryCode)
//
//	var url string
//
//	if googleBase, found := GoogleDomains[countryCode]; found {
//		if start == 0 {
//			url = fmt.Sprintf("%s%s&hl=%s", base(googleBase), searchTerm, languageCode)
//		} else {
//			url = fmt.Sprintf("%s%s&hl=%s&start=%d", base(googleBase), searchTerm, languageCode, start)
//		}
//	} else {
//		if start == 0 {
//			url = fmt.Sprintf("%s%s&hl=%s", stdGoogleBase+GoogleDomains["us"], searchTerm, languageCode)
//		} else {
//			url = fmt.Sprintf("%s%s&hl=%s&start=%d", stdGoogleBase+GoogleDomains["us"], searchTerm, languageCode, start)
//		}
//	}
//
//	if limit != 0 {
//		url = fmt.Sprintf("%s&num=%d", url, limit)
//	}
//
//	return url
//}
