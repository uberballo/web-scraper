package scraper

import (
	"context"
	"encoding/json"
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

//KeyFigures contains all Stocks and their key figures.
type KeyFigures struct {
	Stocks []Stock
}

//Stock struct has the symbol and all found key figures.
type Stock struct {
	Symbol     string
	KeyFigures [][]string
}

//ContainerAndElement contains strings that will point to the container element and to-be queried elements.
type ContainerAndElement struct {
	Container []string
	Element   []string
}

//Scrape finds data from the given url
func Scrape(url string) KeyFigures {
	res := scrapeData(url)
	return res
}

func ScrapeAsJSON(url string) []byte {
	res := scrapeData(url)
	jsonRes, err := toJSON(res)
	if err != nil {
		log.Print(err)
	}
	return jsonRes
}

func toJSON(stocks KeyFigures) ([]byte, error) {
	b, err := json.Marshal(stocks)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func scrapeData(url string) KeyFigures {
	urls := scrapeWebsite(url)

	baseURL = os.Getenv("BASE_URL")
	childContainer = os.Getenv("CHILD_CONTAINER")
	childElement = os.Getenv("CHILD_ELEMENT")
	suffix = os.Getenv("SUFFIX")

	suffixUrls := util.AppendSuffix(urls, suffix)
	workingUrls := util.PrependPrefix(suffixUrls, baseURL)
	shortList := workingUrls[:2]

	scrapedData := scrapeUrls(shortList, baseURL, childContainer, childElement)
	result := KeyFigures{scrapedData}

	return result
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

func scrapeUrls(list []string, baseURL, childContainer, childElement string) []Stock {
	var res []Stock
	var symbol string

	collectFunctions := []func([]*cdp.Node) []string{collectAll}
	childrenObjects := ContainerAndElement{
		[]string{childContainer},
		[]string{childElement}}

	for _, siteURL := range list {
		symbol = util.GetLastPart(siteURL)
		stockInformation := findElements(siteURL, childrenObjects, collectFunctions)
		splittedFigures := util.SplitList(stockInformation, 6)
		newStock := Stock{symbol, splittedFigures}
		res = append(res, newStock)
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
		res = append(res, collectFunctions[i](nodes)...)

	}
	return res
}
