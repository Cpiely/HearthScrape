package main

//ToDo for tomorrow:
//Impliement json or whatever to store the info


import (
	"fmt"
	"sort"
	"strconv"
	"time"
	"github.com/tealeg/xlsx"
)

func main() {
	
	numPages := 30
	
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
	//var wg sync.WaitGroup
	channels := make(chan []string, numPages)

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

	for _, page := range pages{
		//wg.Add(1)
		go func (url string, ch chan []string) {
			scraper := NewScraper(url)
			deckNames := scraper.UrlFinder("div span", "/")
			ch <- deckNames	
		
		}(page, channels)
		//wg.Done()
	}
	time.Sleep(time.Second * 10)
	close(channels)
	
	
	return channels
}

//Creates the unique card list and list of cards from each deck
func Decks(decks []string) []Deck {
	//usedCards := make(map[string]int)
	//var wg sync.WaitGroup
	var deckCards []Deck
	ch := make(chan Deck)
	classCards, neutralCards  := CardIndexing()
	for _,deck := range decks {
		//wg.Add(1)
		go func (url string, ch chan Deck) {
			scraper := NewScraper(url)
			deckInfo := scraper.CardsFinder(classCards, neutralCards)
			deckInfo = scraper.InfoFinder(deckInfo)
			deckInfo.Url = url

			ch <- deckInfo
			//wg.Done()
		}(deck, ch)
	}
	for i := 0; i < len(decks); i++ {
		deckCards = append(deckCards, <-ch)
	}
	//wg.Wait()
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

//Creates the list of decks
func Response(channels chan []string) []string{
	decks := make([]string, 0)
	baseUrl := "http://www.hearthpwn.com"
	for i := range channels {	
		decks = append(decks, i...)	
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
