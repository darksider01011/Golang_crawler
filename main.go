package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

func js(url string) []string {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var links []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second), // Wait for JS to execute
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &links),
	)
	if err != nil {
		log.Fatal(err)
	}
	keys := make(map[string]bool)
	link := []string{}
	for _, li := range links {
		if !keys[li] {
			keys[li] = true
			link = append(link, li)
		}
	}
	for _, a := range link {
		fmt.Println(a)
	}
	return link
}

func parse(urll string) []string {
	rand.Seed(time.Now().UnixNano())
	var parse_link []string

	// User-agents
	user_agent := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:124.0) Gecko/20100101 Firefox/124.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.2420.81",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 OPR/109.0.0.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.4; rv:124.0) Gecko/20100101 Firefox/124.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.4.1 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_4_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 OPR/109.0.0.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux i686; rv:124.0) Gecko/20100101 Firefox/124.0",
	}

	// Random user-agent generator
	rand := rand.Intn(len(user_agent))
	//fmt.Println(user_agent[rand])

	// new collector
	c := colly.NewCollector(
		colly.UserAgent(user_agent[rand]),
	)

	// onError callback
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error:", err)
	})

	// onHtml callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// relative path to absolute path
		if strings.HasPrefix(link, "/") {
			link = urll + link
		}

		// Remove not valid links
		if strings.HasPrefix(link, "https://") {
			parse_link = append(parse_link, link)
		}

		// Remove not valid links
		if strings.HasPrefix(link, "http://") {
			parse_link = append(parse_link, link)
		}

	})
	// Function delay
	time.Sleep(1 * time.Second)

	c.Visit(urll)

	keys := make(map[string]bool)

	//unique links
	link := []string{}

	// Append unique link to slice array
	for _, li := range parse_link {
		if !keys[li] {
			keys[li] = true
			link = append(link, li)
		}
	}

	// Remove root url from slice array
	s_url := urll + "/"
	for i, a := range link {
		if a == urll {
			link = append(link[:i], link[i+1:]...)
		}
		if a == s_url {
			link = append(link[:i], link[i+1:]...)
		}
	}

	u, _ := url.Parse(urll)
	host := u.Hostname()
	domain := strings.TrimPrefix(host, "www.")
	alter_domain := "www." + domain

	// Filter only target domain
	var linkk []string
	for _, a := range link {
		u, err := url.Parse(a)
		if err != nil {
			continue
		}
		host := u.Hostname()

		if host == domain || host == alter_domain {
			linkk = append(linkk, a)
		}
	}

	// Iterate over linkk slice array
	for _, a := range linkk {
		fmt.Println(a)
	}

	return linkk
}

// Depth one crawl
func depth_one(uri string) {
	fmt.Println("")
	fmt.Println("Depth one crawl")

	fmt.Println("")
	word := "Crawling => " + uri
	fmt.Println(word)
	for i := 0; i < len(word); i++ {
		fmt.Print("=")
	}
	fmt.Println("")
	fmt.Println("")

	a := parse(uri)

	keys := make(map[string]bool)
	var uniq []string

	for i, b := range a {

		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		word := "Crawling => " + a[i]
		fmt.Println(word)
		len := len(word)
		for i := 0; i < len; i++ {
			fmt.Print("=")
		}
		fmt.Println("")
		fmt.Println("")

		c := parse(b)
		for _, li := range c {
			if !keys[li] {
				keys[li] = true
				uniq = append(uniq, li)
			}
		}
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Results:")
	fmt.Println("")

	// Print results
	for _, li := range uniq {
		fmt.Println(li)
	}
	u, _ := url.Parse(uri)
	host := u.Hostname()
	domain := strings.TrimPrefix(host, "www.")
	file_name := "depth_one_" + domain + ".txt"
	file, _ := os.Create(string(file_name))

	defer file.Close()
	// Write to file
	fmt.Println("")
	fmt.Println("Saved to " + file_name)
	fmt.Println("")
	for _, li := range uniq {
		file.WriteString(li + "\n")
	}
}

func depth_two(uri string) {
	fmt.Println("")
	fmt.Println("Depth one crawl")

	fmt.Println("")
	word := "Crawling => " + uri
	fmt.Println(word)
	for i := 0; i < len(word); i++ {
		fmt.Print("=")
	}
	fmt.Println("")
	fmt.Println("")

	a := parse(uri)

	key := make(map[string]bool)
	var uniq []string

	keyy := make(map[string]bool)
	var uniqq []string

	var all []string
	var alll []string
	keyyy := make(map[string]bool)

	for i, b := range a {

		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		word := "Crawling => " + a[i]
		fmt.Println(word)
		len := len(word)
		for i := 0; i < len; i++ {
			fmt.Print("=")
		}
		fmt.Println("")
		fmt.Println("")

		one := parse(b)
		for _, li := range one {
			if !key[li] {
				key[li] = true
				uniq = append(uniq, li)
			}
		}
	}
	fmt.Println("")
	fmt.Println("")
	fmt.Println("")
	fmt.Println("Depth one results:")
	fmt.Println("")

	for _, li := range uniq {
		all = append(all, li)
		fmt.Println(li)
	}

	fmt.Println("")
	fmt.Println("Depth two crawl")
	fmt.Println("")

	for i, li := range uniq {
		fmt.Println("")
		fmt.Println("")
		word := "Crawling => " + uniq[i]
		fmt.Println(word)
		len := len(word)
		for i := 0; i < len; i++ {
			fmt.Print("=")
		}
		fmt.Println("")
		fmt.Println("")

		two := parse(li)
		for _, li := range two {
			if !keyy[li] {
				keyy[li] = true
				uniqq = append(uniqq, li)
			}
		}
	}
	fmt.Println("")
	fmt.Println("Depth two results:")
	fmt.Println("")
	for _, li := range uniqq {
		all = append(all, li)
		fmt.Println(li)
	}

	u, _ := url.Parse(uri)
	host := u.Hostname()
	domain := strings.TrimPrefix(host, "www.")
	file_name := "depth_two_" + domain + ".txt"
	file, _ := os.Create(string(file_name))

	defer file.Close()
	// unduplicate urls
	for _, li := range all {
		if !keyyy[li] {
			keyyy[li] = true
			alll = append(alll, li)
		}
	}
	// Write to file
	fmt.Println("")
	fmt.Println("Saved to " + file_name)
	fmt.Println("")
	for _, li := range alll {
		file.WriteString(li + "\n")
	}
}

func main() {
	uri := "https://shamsipour.nus.ac.ir"
	depth_two(uri)
}
