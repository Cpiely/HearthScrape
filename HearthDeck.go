package main

//ToDo for tomorrow:
//Impliement json or whatever to store the info

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"sort"
	"strconv"
	"sync"
	//"time"
)

func main() {

	numPages := 10

	channels := Urls(numPages)
	//List of all of the decks
	decks := Response(channels)
	//deckUrls := CreateMap(decks)
	deckList := Decks(decks)

	standard, wild, brawl := SeperateModes(deckList)
	
	fmt.Println(len(standard), "Standard Decks")
	fmt.Println(len(wild), "Wild Decks")
	fmt.Println(len(brawl), "Tavern Brawl Decks")

	file := CreateFile()

	MakeSpreadSheet(file, standard, "Standard")
	MakeSpreadSheet(file, brawl, "Brawl")
	MakeSpreadSheet(file, wild, "Wild")

	SaveFile(file)

}

//Creates the Urls for each of the pages to be scraped
func Urls(numPages int) chan []string {
	var wg sync.WaitGroup
	channels := make(chan []string, numPages)

	site := "http://www.hearthpwn.com/decks?filter-deck-tag=4&sort=-rating"

	sitePage := "http://www.hearthpwn.com/decks?filter-deck-tag=4&page="
	ext := "&sort=-rating"

	pages := make([]string, numPages)
	if numPages > 0 {
		pages[0] = site
		for i := 1; i < numPages; i++ {
			pages[i] = (sitePage + strconv.Itoa(i+1) + ext)
		}
	}

	for _, page := range pages {
		wg.Add(1)
		go func(url string, ch chan []string) {
			defer wg.Done()
			scraper := NewScraper(url)
			deckNames := scraper.UrlFinder("div span", "/")
			ch <- deckNames

		}(page, channels)
	
	}
	//time.Sleep(time.Second * 3)
	wg.Wait()

	close(channels)

	return channels
}

//Creates the unique card list and list of cards from each deck
func Decks(decks []string) []Deck {
	//usedCards := make(map[string]int)
	//var wg sync.WaitGroup
	var deckCards []Deck
	ch := make(chan Deck)
	classCards, neutralCards := CardIndexing()
	for _, deck := range decks {
	//	wg.Add(1)
		go func(url string, ch chan Deck) {
	//		defer wg.Done()
			scraper := NewScraper(url)
			deckInfo := scraper.CardsFinder(classCards, neutralCards)
			deckInfo = scraper.InfoFinder(deckInfo)
			deckInfo.Url = url

			ch <- deckInfo
			//wg.Done()
		}(deck, ch)
	}

	//wg.Wait()

	for i := 0; i < len(decks); i++ {
		deckCards = append(deckCards, <-ch)
	}
	

	return deckCards
}

//Creates each of the Sheets for Standard, Wild, and Brawl modes
func MakeSpreadSheet(file *xlsx.File, mode []Deck, sheetName string) {
	sheet := CreateSheet(file, sheetName)
	offset := 0
	for i := 0; i < len(mode); i++ {
		offset = PrintDeck(sheet, mode[i], offset)
	}
}

//Seperates the Deck lists into the different game modes
func SeperateModes(decks []Deck) ([]Deck, []Deck, []Deck) {
	
	standard := []Deck{}
	brawl := []Deck{}
	wild := []Deck{}

	for _, deck := range decks {
		if deck.Mode == "Standard" {
			standard = append(standard, deck)
		} else if deck.Mode == "Wild" {
			wild = append(wild, deck)
		} else {
			brawl = append(brawl, deck)
		}
	}
	sort.Sort(ByRating(standard))
	sort.Sort(ByRating(wild))
	sort.Sort(ByRating(brawl))
	return standard, wild, brawl
}

//Creates the list of decks
func Response(channels chan []string) []string {
	decks := make([]string, 0)
	baseUrl := "http://www.hearthpwn.com"
	for i := range channels {
		decks = append(decks, i...)
	}
	for idx, deck := range decks {
		decks[idx] = baseUrl + deck
	}
	return decks
}

func printResults(list []string) {
	for _, deck := range list {
		fmt.Println(deck)
	}
}
