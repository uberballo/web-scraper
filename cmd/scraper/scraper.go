package scraper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var testSuffix string = "porssi/porssikurssit/osake/NOKIA/tilinpaatos"
var baseUrl string = "https://www.kauppalehti.fi"
var url string = "https://www.kauppalehti.fi/porssi/kurssit/XHEL"
var elements string = `div[class="list-striped mb-1"] > div > a`

//class["list-item-header routeless"]
func Scrape(url string) {
	scrapeData(url, "")
}

func scrapeData(url, suffix string) {
	//var res []string
	stockUrls := findElements(url, `.OMXH-list`, elements, collectHrefElements)
	t := baseUrl + stockUrls[0]
	test := findElements(t, `.list-item-wrapper`, `.stock-list-column > span, .stock-list-column > div > h5 > span`, collectAll)
	fmt.Println(test)
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
	for _, i := range nodes {
		res = append(res, i.AttributeValue("href")+"/tilinpaatos")
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
