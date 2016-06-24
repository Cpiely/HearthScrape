package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"sort"
	"strconv"
	"sync"
    "encoding/json"
    "io/ioutil"
    "bytes"
)

func main() {
    
	numPages := 10

	channels := Urls(numPages)
	//List of all of the deck names
	decks := Response(channels)
    //List of decks, List of cards
	deckList, cardList := Decks(decks)

    deckList = RemoveDupes(deckList)
    //AllDecks
	finalDecks := FinalDeckInfo(deckList)
    toJson(finalDecks, cardList)
	printResults(finalDecks)

    //startTracking(finalDecks)

	file := CreateFile()

	MakeSpreadSheet(file, finalDecks.StandardDecks, "Standard")
	MakeSpreadSheet(file, finalDecks.WildDecks, "Wild")
	MakeSpreadSheet(file, finalDecks.BrawlDecks, "Brawl")

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
	wg.Wait()

	close(channels)

	return channels
}

//Creates the unique card list and list of cards from each deck
func Decks(decks []string) ([]Deck, Cards) {
	//usedCards := make(map[string]int)
	//var wg sync.WaitGroup
	var deckCards []Deck
	ch := make(chan Deck)

	classCards, neutralCards := CardIndexing()
    cards := Cards {
        Neutral: neutralCards,
        Class: classCards,
    }
	for _, deck := range decks {
		//	wg.Add(1)
		go func(url string, ch chan Deck) {
			//		defer wg.Done()
			scraper := NewScraper(url)
			deckInfo := scraper.CardsFinder(cards)
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

	return deckCards, cards
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
func FinalDeckInfo(decks []Deck) AllDecks {

	standard := []Deck{}
	brawl := []Deck{}
	wild := []Deck{}

	sort.Sort(ByType(decks))

	standardTypes := make(map[string]int)
	wildTypes := make(map[string]int)

	standardTypeKey := make([]string, 0)
	wildTypeKey := make([]string, 0)

	for _, deck := range decks {
		if deck.Mode == "Tavern Brawl" {
			brawl = append(brawl, deck)
		} else if deck.Mode == "Standard" {
			standardTypeKey = AddKey(standardTypeKey, deck.Type)
			standardTypes[deck.Type] += 1
			standard = append(standard, deck)
		} else {
			wildTypeKey = AddKey(wildTypeKey, deck.Type)
			wildTypes[deck.Type] += 1
			wild = append(wild, deck)
		}
	}

	sort.Sort(ByRating(standard))
	sort.Sort(ByRating(wild))
	sort.Sort(ByRating(brawl))

	finalDecks := AllDecks{
		NumDecks:        len(decks),
		NumStandard:     len(standard),
		NumWild:         len(wild),
		NumBrawl:        len(brawl),
		StandardTypes:   standardTypes,
		WildTypes:       wildTypes,
		WildTypeKey:     wildTypeKey,
		StandardTypeKey: standardTypeKey,
		StandardDecks:   standard,
		BrawlDecks:      brawl,
		WildDecks:       wild,
		AllDecks:        decks,
	}

	return finalDecks
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

func AddKey(keys []string, key string) []string {
	contains := false
	for _, i := range keys {
		if i == key {
			contains = true
		}
	}
	if !contains {
		keys = append(keys, key)
	}
	return keys
}

func RemoveDupes(deckList []Deck) []Deck {
    toRemove := make([]int, 0)
    for i, deck := range(deckList) {
        for x := i+1; x < len(deckList); x++ {
            if deck.DeckID == deckList[x].DeckID {
                toRemove = append(toRemove, x)
            }
        }
    }

    for _,i := range toRemove {
        if i < len(deckList) {
            deckList = append(deckList[:i], deckList[i+1:]...)
        }
    }
    return deckList
}

func printResults(finalDecks AllDecks) {
	fmt.Println(finalDecks.NumStandard, "Standard Decks:")
	for _, typ := range finalDecks.StandardTypeKey {
		fmt.Println(finalDecks.StandardTypes[typ], typ)
	}
	fmt.Println()
	fmt.Println(finalDecks.NumWild, "Wild Decks:")
	for _, typ := range finalDecks.WildTypeKey {
		fmt.Println(finalDecks.WildTypes[typ], typ)
	}
	fmt.Println()
	fmt.Println(finalDecks.NumBrawl, "Tavern Brawl Decks")
	fmt.Println()
	fmt.Println(finalDecks.NumDecks, "Total Decks")
    fmt.Println()
}

func toJson(decks AllDecks, cards Cards) {
    var outDeck bytes.Buffer
    var outCards bytes.Buffer
    deckJson, err := json.Marshal(decks)
    checkErr(err)
    cardsJson, err := json.Marshal(cards)
    checkErr(err)
    json.Indent(&outDeck, deckJson, "", "\t")
    json.Indent(&outCards, cardsJson, "", "\t")

    err = ioutil.WriteFile("./data/deck.json", outDeck.Bytes(), 0644)
    checkErr(err)
    err = ioutil.WriteFile("./data/cards.json", outCards.Bytes(), 0644)
    checkErr(err)
}
