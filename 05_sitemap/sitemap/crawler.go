package sitemap

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/maximepeschard/gophercises/04_link/link"
)

var namespace = "http://www.sitemaps.org/schemas/sitemap/0.9"

// Sitemap represents a sitemap built from a source URL
type Sitemap struct {
	sm     sitemap
	source *url.URL
}

type sitemap struct {
	XMLName   xml.Name `xml:"urlset"`
	Namespace string   `xml:"xmlns,attr"`
	Pages     []page   `xml:"url"`
}

type page struct {
	XMLName xml.Name `xml:"url"`
	URL     string   `xml:"loc"`
}

// New returns a empty sitemap with the given source URL
func New(u *url.URL) *Sitemap {
	return &Sitemap{sm: sitemap{Namespace: namespace}, source: u}
}

func (s Sitemap) String() string {
	return fmt.Sprintf("Sitemap(source: %s)", s.source)
}

// Build builds the sitemap starting from its source URL and following
// links up to maxDepth levels (or infitenely if maxDepth is -1)
func (s *Sitemap) Build(maxDepth int) *Sitemap {
	visited := make(map[string]*url.URL)
	queue := []*url.URL{s.source}
	depth := 0
	for len(queue) > 0 && (maxDepth == -1 || depth <= maxDepth) {
		depth++
		var wg sync.WaitGroup
		urls := make(chan *url.URL)

		// launch one goroutine per page to parse its links
		for _, u := range queue {
			_, prs := visited[u.Path]
			if prs {
				continue
			}
			visited[u.Path] = u

			wg.Add(1)
			go parsePage(u, urls, &wg)
		}

		// goroutine to close the results channel once every page is parsed
		go func() {
			wg.Wait()
			close(urls)
		}()

		// add clean URLs with the same domain as the source to the queue
		queue = nil
		for u := range urls {
			if u := cleanURL(*u); u.Host == s.source.Host {
				queue = append(queue, &u)
			}
		}
	}

	s.sm.Pages = nil
	for _, u := range visited {
		s.sm.Pages = append(s.sm.Pages, page{URL: u.String()})
	}

	return s
}

// WriteXML writes the XML encoding of the sitemap to w
func (s Sitemap) WriteXML(w io.Writer) error {
	enc := xml.NewEncoder(w)
	enc.Indent("", "    ")
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	return enc.Encode(s.sm)
}

func cleanURL(u url.URL) url.URL {
	u.RawQuery = ""
	u.Fragment = ""
	return u
}

func parsePage(u *url.URL, urls chan<- *url.URL, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Get(u.String())
	if err != nil || response.StatusCode != 200 {
		return
	}

	links, err := link.Parse(response.Body)
	if err != nil {
		return
	}

	for _, l := range links {
		href, err := url.Parse(l.Href)
		if err != nil {
			continue
		}

		urls <- u.ResolveReference(href)
	}
}
