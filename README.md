# Web scraper  

A simple web scraper made with GO. Designed to scrape this one Finnish stock site. Uses Chromedp.  

## How it works  
First the application scrapes the given URL to find contents `MAIN_CONTAINER` and from there grabs all the `MAIN_ELEMENT`'s. From the elements, it collects all the href's that point to the stocks page. With the href's, application will repeat the same process again. Find `CHILD_CONTAINER` and collect `CHILD_ELEMENT`'s. With the scraped data, application will build a list of lists, where the inner list is \[name, value1, value2, value3, value4, value5\].  

As the website uses JS to show the data, we wait for the page to load before scraping it.  

Currently we collect data from 5 stocks, as the website will otherwise time us out and I'd rather avoid bombarding their site.   

## How to use   
Required following variables in `.env` file: 
- BASE_URL  
- SUFFIX  
- MAIN_URL  
- MAIN_ELEMENT  
- MAIN_CONTAINER  
- CHILD_CONTAINER  
- CHILD_ELEMENT  
  
Run the program with: `go run main.go`  
  
Example:
```
<div class="list-1">
  <div class="element-1">
    <span href="/site"> hello </span>
  </div>
</div>
```
MAIN_CONTAINER=.list-1  
MAIN_ELEMENT=.element-1 > span  
And so on. I can explain it better if wanted.  

