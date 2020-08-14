package scraper

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var testSuffix string = "porssi/porssikurssit/osake/NOKIA/tilinpaatos"
var baseURL string
var mainContainer string
var mainElement string
var childContainer string
var childElement string

func Scrape(url string) {
	fmt.Println(scrapeData(url, ""))
}

func scrapeData(url, suffix string) [][][]string {
	var res [][][]string

	urls := scrapeWebsite(url)

	baseURL = os.Getenv("BASE_URL")
	childContainer = os.Getenv("CHILD_CONTAINER")
	childElement = os.Getenv("CHILD_ELEMENT")

	suffixUrls := appendSuffix(urls, "/tilinpaatos")
	workingUrls := prependPrefix(suffixUrls, baseURL)
	shortList := workingUrls[:5]

	scrapedData := scrapeUrls(shortList, baseURL, childContainer, childElement)

	for _, n := range scrapedData {
		res = append(res, splitList(n))
	}
	return res
}

func scrapeWebsite(url string) []string {
	var res []string
	mainContainer = os.Getenv("MAIN_CONTAINER")
	mainElement = os.Getenv("MAIN_ELEMENT")

	res = findElements(url, mainContainer, mainElement, collectHrefElements)
	return res
}

func scrapeUrls(list []string, baseURL, childContainer, childElement string) [][]string {
	var res [][]string
	for _, siteURL := range list {
		stockInformation := findElements(siteURL, childContainer, childElement, collectAll)
		res = append(res, stockInformation)
	}
	return res
}

func splitList(list []string) [][]string {
	var res [][]string
	var tempSlice []string
	for i := 0; i < len(list); i++ {
		tempSlice = append(tempSlice, list[i])
		if (i+1)%6 == 0 {
			res = append(res, tempSlice)
			tempSlice = nil
		}

	}
	return res
}

func visible(url, element string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	timeoutContext, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	err := chromedp.Run(timeoutContext,
		chromedp.Navigate(url),
		chromedp.WaitVisible(element),
	)

	if err == nil {
		timeoutContext.Done()
		return nil
	}

	return err
}

func appendSuffix(list []string, suffix string) []string {
	var res []string
	for _, n := range list {
		res = append(res, n+suffix)
	}
	return res
}

func prependPrefix(list []string, prefix string) []string {
	var res []string
	for _, n := range list {
		res = append(res, prefix+n)
	}
	return res

}

func collectHrefElements(nodes []*cdp.Node) []string {
	var res []string
	for _, n := range nodes {
		res = append(res, n.AttributeValue("href"))
	}
	return res
}

func collectAll(nodes []*cdp.Node) []string {
	var res []string
	for _, i := range nodes {
		for _, j := range i.Children {
			res = append(res, j.NodeValue)
		}
	}
	return res
}

func findElements(url, containerElement, element string, collectFunction func([]*cdp.Node) []string) []string {
	var nodes []*cdp.Node
	var err error

	err = visible(url, containerElement)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Nodes(element, &nodes, chromedp.ByQueryAll),
	)

	if err != nil {
		log.Fatal(err)
	}

	res := collectFunction(nodes)
	return res
}
