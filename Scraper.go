package main
import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
	"strings"
	"strconv"
)

type Scraper struct {
	url      string
	document *goquery.Document
}

func NewScraper(url string) *Scraper {
	s := new(Scraper)
	s.url = url
	s.document = s.getDocument()
	return s
}

func (s *Scraper) CardsFinder(classIndex, neutralIndex map[string]Card) Deck {
	var err error
	classKeys := make([]Card, 0)
	neutralKeys := make([]Card, 0)
	cardList := make(map[Card]int)
	card := Card{}
	ok := true
	//Finds each of the cards in the deck
	s.document.Find("tr ").Each(func (i int, s *goquery.Selection) {
		info := strings.TrimSpace(strings.Replace(s.Contents().Text(), "\n", "", -1))
		splitInfo := strings.Split(info, "    ")
		//Filters out anything that isn't actually a card
		if len(splitInfo) == 2 {
			numCost := string((strings.Split(splitInfo[1], " ")[2])[:])
			//Change to have it add the card from the list of all cards 
			if card, ok = classIndex[splitInfo[0]]; !ok {
				card = neutralIndex[splitInfo[0]]
				neutralKeys = append(neutralKeys, card)
			} else {
				classKeys = append(classKeys, card)
			}

			cardList[card], err = strconv.Atoi(string(numCost[0]))
			checkErr(err)
		}
	})

	deck := Deck {
		CardList: cardList,
		ClassKeys: classKeys,
		NeutralKeys: neutralKeys,
	}
	return deck
}

func cardIndexing() (map[string]Card ,map[string]Card) {
	url := "http://www.hearthpwn.com/cards?display=1&filter-premium=1"
	ext := "&page="

	pages := make([]string, 9)
	pages[0] = url
	for i := 1; i < 9; i++ {
		pages[i] = (url + ext + strconv.Itoa(i + 1))
	}

	neutralList := make(map[string]Card, 0)
	classList := make(map[string]Card, 0)

	card := Card{}
	neutral := false
	for _, page := range pages {
		index := 0
		name := ""
		scraper := NewScraper(page)
		scraper.document.Find("td ").Each(func (i int, s *goquery.Selection) {
			info := s.Contents().Text()
			switch index {
			case 1:
				name = strings.Replace(info, "\n", "", -1)
				card.Name = name
			case 2:
				card.Type = info
			case 3:
				if info == "" {
					card.Class = "Neutral"
					neutral = true
				} else {
					card.Class = strings.TrimSpace(info)
					neutral = false
				}
			case 4:
				card.Cost = info
			case 5:
				card.Attack = info
			case 6:
				card.Health = info
				if neutral {
					neutralList[name] = card
				} else { 	
					classList[name] = card
				}
				card = Card{}
				index = 0
			}
			index++
		})
	}
	return classList, neutralList
}

func (s *Scraper) ModeFinder(deck Deck) Deck {
	s.document.Find("section p").Each(func (i int, s *goquery.Selection) {
		if text := s.Contents().Text(); text == "Standard" || text == "Wild" {
		 	deck.Mode = text
		}
	})
	return deck
}

func (s *Scraper) InfoFinder(deck Deck) Deck{
	date := (s.document.Find("li abbr").First().Text())
	deck.DateModified = Date(date)

	notFirst := true
	s.document.Find("div").Each(func (i int, s *goquery.Selection) {
		text := s.Contents().Text()
		if strings.HasPrefix(text, "+") && notFirst{
			deck.Rating = text
			notFirst = false
		}
	})

	s.document.Find("ul li span").Each(func (i int, s *goquery.Selection) {
		text := s.Contents().Text()
		
		switch i {
		case 63:
			deck.Class = text
		case 64:
			deck.Name = text
		case 68:
			if text == "Tavern Brawl"{
				deck.Mode = text
			}
			deck.Type = text
		case 67:
			deck.Expansion = text
		case 69:
			deck.Cost = text
		case 72:
			created := strings.Split(text, " ")
			deck.DateCreated = created[0]
		}
		
	})
	return deck
}

func Date (date string) string {
	var newDate string
	splitDate := strings.Split(date, " ")
	switch splitDate[0]{
		case "Jan":
			newDate = strconv.Itoa(1) 
		case "Feb":
			newDate = strconv.Itoa(2) 
		case "Mar":
			newDate = strconv.Itoa(3) 
		case "Apr":
			newDate = strconv.Itoa(4) 
		case "May":
			newDate = strconv.Itoa(5) 
		case "Jun":
			newDate = strconv.Itoa(6) 
		case "Jul":
			newDate = strconv.Itoa(7) 
		case "Aug":
			newDate = strconv.Itoa(8) 
		case "Sep":
			newDate = strconv.Itoa(9) 
		case "Oct":
			newDate = strconv.Itoa(10) 
		case "Nov":
			newDate = strconv.Itoa(11) 
		case "Dec": 
			newDate = strconv.Itoa(12) 
	} 
	newDate += "/"
	newDate += string(splitDate[1][:len(splitDate[1])-1]) + "/"
	newDate += splitDate[2]
	return newDate
}

func (s *Scraper) UrlFinder(tags string, prefix string) [] string {
	urls := make([]string, 0)
	s.document.Find(tags).Each(func (i int, s *goquery.Selection) {
		if val, ok := s.Contents().Attr("href"); ok && strings.HasPrefix(val, prefix){
			urls = append(urls, val)
		}
	})
	return urls
}

func (s *Scraper) getDocument() *goquery.Document {
	resp := s.getResponse()
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	checkErr(err)
	return doc
}

func (s *Scraper) getResponse() *http.Response {
	//fmt.Println(s.url)
	resp, err := http.Get(s.url)
	checkErr(err)
	return resp
}

func checkErr(err error){
	if err != nil {
    	fmt.Println(err)
   	}
}
