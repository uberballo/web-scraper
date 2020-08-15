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

type Stock struct {
	Name       string
	Symbol     string
	KeyFigures []string
}

type ContainerAndElement struct {
	Container []string
	Element   []string
}

//Scrape finds data from the given url
func Scrape(url string) [][][]string {
	res := scrapeData(url)
	return res
}

func scrapeData(url string) [][][]string {
	var stocksKeyFigures [][][]string

	urls := scrapeWebsite(url)

	baseURL = os.Getenv("BASE_URL")
	childContainer = os.Getenv("CHILD_CONTAINER")
	childElement = os.Getenv("CHILD_ELEMENT")
	suffix = os.Getenv("SUFFIX")

	suffixUrls := util.AppendSuffix(urls, suffix)
	workingUrls := util.PrependPrefix(suffixUrls, baseURL)
	shortList := workingUrls[:1]

	scrapedData := scrapeUrls(shortList, baseURL, childContainer, childElement)

	for _, n := range scrapedData {
		stocksKeyFigures = append(stocksKeyFigures, util.SplitList(n, 6))
	}
	return stocksKeyFigures
}

func scrapeWebsite(url string) []string {
	var res []string

	mainContainer = os.Getenv("MAIN_CONTAINER")
	mainElement = os.Getenv("MAIN_ELEMENT")
	mainObjects := ContainerAndElement{
		[]string{mainContainer},
		[]string{mainElement}}

	collectFunctions := []func([]*cdp.Node) []string{collectHrefElements}
	res = findElements(url, mainObjects, collectFunctions)
	return res
}

func scrapeUrls(list []string, baseURL, childContainer, childElement string) [][]string {
	var res [][]string

	collectFunctions := []func([]*cdp.Node) []string{collectAll}
	childrenObjects := ContainerAndElement{
		[]string{childContainer},
		[]string{childElement}}

	for _, siteURL := range list {
		stockInformation := findElements(siteURL, childrenObjects, collectFunctions)
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

func findElements(url string, elements ContainerAndElement, collectFunctions []func([]*cdp.Node) []string) []string {
	var res []string
	var nodes []*cdp.Node
	var err error

	for i := range elements.Container {
		err = visible(url, elements.Container[i])
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		err = chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Nodes(elements.Element[i], &nodes, chromedp.ByQueryAll),
		)

		if err != nil {
			log.Fatal(err)
		}
		for _, function := range collectFunctions {
			res = append(res, function(nodes)...)
		}
	}
	return res
}
