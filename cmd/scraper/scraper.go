package scraper

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/uberballo/web-scraper/cmd/util"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var baseURL string
var mainContainer string
var mainElement string
var childContainer string
var childElement string
var suffix string

//Scrape finds data from the given url
func Scrape(url string) [][][]string {
	res := scrapeData(url)
	return res
}

func scrapeData(url string) [][][]string {
	var res [][][]string

	urls := scrapeWebsite(url)

	baseURL = os.Getenv("BASE_URL")
	childContainer = os.Getenv("CHILD_CONTAINER")
	childElement = os.Getenv("CHILD_ELEMENT")
	suffix = os.Getenv("SUFFIX")

	suffixUrls := util.AppendSuffix(urls, suffix)
	workingUrls := util.PrependPrefix(suffixUrls, baseURL)
	shortList := workingUrls[:5]

	scrapedData := scrapeUrls(shortList, baseURL, childContainer, childElement)

	for _, n := range scrapedData {
		res = append(res, util.SplitList(n, 6))
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
