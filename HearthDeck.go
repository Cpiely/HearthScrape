package main

//ToDo for tomorrow:
//More info on decks: save the URLs, Authors, date modified, rating, dust cost
//Impliement json or whatever to store the info
//

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"github.com/tealeg/xlsx"
)

func main() {
	
	numPages := 5
	
	channels := Urls(numPages) 
	//List of all of the decks
	decks := Response(channels)
	//deckUrls := CreateMap(decks)
	deckList := Decks(decks)

	standard, wild, brawl := SeperateModes(deckList)
	
	file := CreateFile()

	MakeSpreadSheet(file, standard, "Standard")
	MakeSpreadSheet(file, brawl, "Brawl")
	MakeSpreadSheet(file, wild, "Wild")

	SaveFile(file)

}

func Urls(numPages int) []chan []string {
	var wg sync.WaitGroup
	channels := Channels(numPages)

	site := "http://www.hearthpwn.com/decks?filter-deck-tag=4&sort=-rating"

	sitePage := "http://www.hearthpwn.com/decks?filter-deck-tag=4&page="
	ext := "&sort=-rating"

	pages := make([]string, numPages)
	if numPages > 0{
		pages[0] = site
		for i := 1; i < numPages; i++ {
			pages[i] = (sitePage + strconv.Itoa(i + 1) + ext)
		}
	}

	for idx, page := range pages{
		wg.Add(1)
		go func (url string, ch chan []string) {
			scraper := NewScraper(url)
			deckNames := scraper.UrlFinder("div span", "/")
			ch <- deckNames
			
		}(page, channels[idx])
		wg.Done()
	}
	wg.Wait()
	
	return channels
}

//Creates the unique card list and list of cards from each deck
func Decks(decks []string) []Deck {
	//usedCards := make(map[string]int)
	var wg sync.WaitGroup
	var deckCards []Deck
	ch := make(chan Deck)
	classCards, neutralCards  := cardIndexing()
	for _,deck := range decks {
		wg.Add(1)
		go func (url string, ch chan Deck) {
			scraper := NewScraper(url)
			deckInfo := scraper.CardsFinder(classCards, neutralCards)
			deckInfo = scraper.ModeFinder(deckInfo)
			deckInfo = scraper.InfoFinder(deckInfo)
			deckInfo.Url = url

			ch <- deckInfo
			wg.Done()
		}(deck, ch)
	}
	for i := 0; i < len(decks); i++ {
		deckCards = append(deckCards, <-ch)
	}
	wg.Wait()
	close(ch)
	return deckCards
}

func MakeSpreadSheet(file *xlsx.File, mode []Deck, sheetName string) {
	sheet := CreateSheet(file, sheetName)
	offset := 0
	for i := 0; i < len(mode); i++ {
		offset = PrintDeck(sheet, mode[i], offset)
	}
}

//Seperates the Deck lists into the different game modes
func SeperateModes(decks []Deck) ([]Deck, []Deck, []Deck) {
	sort.Sort(ByMode(decks))
	mode := "Standard"
	start := 0
	standard := []Deck{}
	brawl := []Deck{}
	wild := []Deck{}
	for i := 0; i < len(decks); i++ {
		if decks[i].Mode != mode {
			if mode == "Standard" {
				standard = decks[start:i -1]
			} else {
				brawl = decks[start:i - 1]
			}
			mode = decks[i].Mode
			start = i
		}
	}
	wild = decks[start:len(decks)]
	return standard, wild, brawl
}

//Creates the Channels
func Channels(num int) []chan []string {
	channels := make([]chan []string, num)
	for i := 0; i < num; i++ {
		channels[i] = make(chan []string)
	}
	return channels
}

//Creates the list of decks
func Response(channels []chan []string) []string{
	decks := make([]string, 0)
	baseUrl := "http://www.hearthpwn.com"
	for i := 0; i < len(channels); i++ {	
		decks = append(decks, <-channels[i]...)		
	}
	for idx, deck := range decks{
		decks[idx] = baseUrl + deck
	}
	return decks
}

func printResults(list []string) {
	for _, deck := range list {
		fmt.Println(deck)
	}
}
