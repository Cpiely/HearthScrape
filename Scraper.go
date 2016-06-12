package main

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
	"strconv"
	"strings"
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
	s.document.Find("tr ").Each(func(i int, s *goquery.Selection) {
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

	deck := Deck{
		CardList:    cardList,
		ClassKeys:   classKeys,
		NeutralKeys: neutralKeys,
	}
	return deck
}

//Gets the Lists of Neutral and Class cards
func CardIndexing() (map[string]Card, map[string]Card) {
	url := "http://www.hearthpwn.com/cards?display=1&filter-premium=1"
	ext := "&page="

	pages := make([]string, 9)
	pages[0] = url
	for i := 1; i < 9; i++ {
		pages[i] = (url + ext + strconv.Itoa(i+1))
	}

	neutralList := make(map[string]Card, 0)
	classList := make(map[string]Card, 0)
	card := Card{}
	neutral := false
	var err error
	for _, page := range pages {
		index := 0
		name := ""
		scraper := NewScraper(page)
		scraper.document.Find("td ").Each(func(i int, s *goquery.Selection) {
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
				card.Cost, err = strconv.Atoi(info)
				checkErr(err)
			case 5:
				card.Attack, err = strconv.Atoi(info)
				checkErr(err)
			case 6:
				card.Health, err = strconv.Atoi(info)
				checkErr(err)
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
	return neutralList, classList
}

//Gets the Moditified, Rating Mode, Class, Name, Type, Expansion, Cost, and Creation Date
func (s *Scraper) InfoFinder(deck Deck) Deck {
	//Last time the deck was modified
	date := (s.document.Find("li abbr").First().Text())
	deck.DateModified = Date(date)

	notFirst := true
	var err error
	s.document.Find("div").Each(func(i int, s *goquery.Selection) {
		text := s.Contents().Text()
		if strings.HasPrefix(text, "+") && notFirst {
			text = strings.TrimLeft(text, "+")
			deck.Rating, err = strconv.Atoi(text)
			checkErr(err)
			notFirst = false
		}
	})

	s.document.Find("section p").Each(func(i int, s *goquery.Selection) {
		if text := s.Contents().Text(); text == "Standard" || text == "Wild" {
			deck.Mode = text
		}
	})

	s.document.Find("ul li span").Each(func(i int, s *goquery.Selection) {
		text := s.Contents().Text()
		switch i {
		case 63:
			deck.Class = text
		case 64:
			deck.Name = text
		case 68:
			if text == "Tavern Brawl" {
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

//Takes in a Date in the form of Jan 1, 2001 and returns 1/1/2001
func Date(date string) string {
	
	var newDate string
	splitDate := strings.Split(date, " ")
	months := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",}
	
	for i, month := range(months){
		if month == splitDate[0] {
			newDate = strconv.Itoa(i)
			
		}
	}

	newDate += "/"
	newDate += string(splitDate[1][:len(splitDate[1])-1]) + "/"
	newDate += splitDate[2]
	return newDate
}

//Finds the Urls ont the list of deck pages
func (s *Scraper) UrlFinder(tags string, prefix string) []string {
	urls := make([]string, 0)
	s.document.Find(tags).Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Contents().Attr("href"); ok && strings.HasPrefix(val, prefix) {
			urls = append(urls, val)
		}
	})
	return urls
}

//Gets the document from the site
func (s *Scraper) getDocument() *goquery.Document {
	resp := s.getResponse()
	doc, err := goquery.NewDocumentFromResponse(resp)
	checkErr(err)
	er := resp.Body.Close()
	checkErr(er)
	return doc
}

//Gets the Webpage from the URl
func (s *Scraper) getResponse() *http.Response {
	resp, err := http.Get(s.url)
	checkErr(err)
	return resp
}

//Checks for errors and prints it if there is one
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func fuckImports() {
	fmt.Println("fuck you")
}
