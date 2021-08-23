package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://www.dunyahalleri.com/haftanin-ozeti-"
const lastWeek = 250

func main() {
	file, err := os.Create("./resources.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	urls := make(map[string]uint)

	soup.SetDebug(true)

	for i := 1; i < lastWeek; i++ {
		pageLink := baseURL + strconv.Itoa(i)
		for j := 1; j < 10; j++ {
			subPageLink := pageLink + "/" + strconv.Itoa(j)
			fmt.Println(subPageLink)
			req, err := http.NewRequest("GET", subPageLink, nil)

			client := &http.Client{}
			httpResp, err := client.Do(req)

			// is request redirected
			if j != 1 && httpResp != nil && httpResp.Request.URL.String() != subPageLink {
				continue
			}

			resp, err := soup.Get(subPageLink)
			if err != nil {
				os.Exit(1)
			}

			doc := soup.HTMLParse(resp)
			div := doc.Find("div", "class", "entry-content")
			if div.Error != nil {
				os.Exit(1)
			}
			links := div.FindAll("a")
			for _, link := range links {
				l := link.Attrs()["href"]
				if strings.Contains(l, "dunyahalleri") || strings.Contains(l, "mserdark") {
					continue
				}

				u, err := url.Parse(link.Attrs()["href"])
				if err != nil {
					continue
				}

				if _, ok := urls[u.Host]; !ok {
					urls[u.Host] = 1
				} else {
					oldVal := urls[u.Host]
					urls[u.Host] = oldVal + 1
				}
			}

			time.Sleep(time.Second)
		}
	}

	var ss []kv
	for k, v := range urls {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
		file.WriteString(kv.Key + ": " + strconv.Itoa(int(kv.Value)) + "\n")
	}
}

type kv struct {
	Key   string
	Value uint
}