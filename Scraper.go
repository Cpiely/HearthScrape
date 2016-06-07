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

func (s *Scraper) FindInfo(cardIndex map[string]Card) Deck {
	deck := s.InfoFinder(cardIndex)
	deck = s.ModeFinder(deck)
	deck = s.TypeFinder(deck)

	return deck
}

func FindCard (cardList []Card, name string) Card{
	for _, card := range cardList {
		if card.Name == name {
			return card
		}
	}  
	return Card{}
}

func (s *Scraper) InfoFinder(cardIndex map[string]Card) Deck {
	var err error
	keys := make([]Card, 0)
	cardList := make(map[Card]int)
	//Finds each of the cards in the deck
	s.document.Find("tr ").Each(func (i int, s *goquery.Selection) {
		info := strings.TrimSpace(strings.Replace(s.Contents().Text(), "\n", "", -1))
		splitInfo := strings.Split(info, "    ")
		//Filters out anything that isn't actually a card
		if len(splitInfo) == 2 {
			numCost := string((strings.Split(splitInfo[1], " ")[2])[:])
			//Change to have it add the card from the list of all cards 
			card := cardIndex[splitInfo[0]]		
			cardList[card], err = strconv.Atoi(string(numCost[0]))
			checkErr(err)
			keys = append(keys, card)
		}
	})

	deck := Deck {
		CardList: cardList,
		Keys: keys,
	}
	return deck
}

func cardIndexing() map[string]Card {
	url := "http://www.hearthpwn.com/cards?display=1&filter-premium=1"
	ext := "&page="

	pages := make([]string, 9)
	pages[0] = url
	for i := 1; i < 9; i++ {
		pages[i] = (url + ext + strconv.Itoa(i + 1))
	}

	cardList := make(map[string]Card, 0)
	card := Card{}
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
				} else {
					card.Class = strings.TrimSpace(info)
				}
			case 4:
				card.Cost = info
			case 5:
				card.Attack = info
			case 6:
				card.Health = info
				cardList[name] = card
				card = Card{}
				index = 0
			}
			index++
		})
	}
	return cardList
}

func (s *Scraper) ModeFinder(deck Deck) Deck {
	s.document.Find("section p").Each(func (i int, s *goquery.Selection) {
		if text := s.Contents().Text(); text == "Standard" || text == "Wild" {
		 	deck.Mode = text
		}
	})
	return deck
}

func (s *Scraper) TypeFinder(deck Deck) Deck{
	s.document.Find("ul li span").Each(func (i int, s *goquery.Selection) {
		text := s.Contents().Text()
		switch text {
		case "Combo", "Tempo", "Aggro", "Control", "Midrange", "None", "Tournament" : 
 			deck.Type = text
 		case "Tavern Brawl":
 			deck.Type = text
 			deck.Mode = text
		}
	})
	return deck	
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
