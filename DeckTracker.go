package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func startTracking(deckLists AllDecks) {
	fmt.Println("Would you like to start a game?")
	text := Prompt("(y | n)")

	if text == "y\n" {
		posDecks := Mode(deckLists)
		first := true
        cont := true
		enemyDeck := Deck{
			Name:     "Enemy",
			CardList: make(map[string]int),
		}
		for cont {
			if first {
				fmt.Println()
				fmt.Println("Options:")
				fmt.Println("add")
				fmt.Println("view")
				fmt.Println("quit")
				fmt.Println("help")
				fmt.Println()
				first = false
			}
			text := Prompt("Please enter choice: ")
			switch text {
			case "add\n":
				card := Prompt("Enter Card Played")
				enemyDeck, posDecks = CardPlayed(enemyDeck, posDecks, card)
				fmt.Println(strconv.Itoa(len(posDecks)) + " possible decks remaining.")
                //Make a spreadsheet of all the possible decks and highlight the cards that have been played
			case "view\n":
				for i, posDeck := range posDecks {
					fmt.Println("(" + strconv.Itoa(i) + ")" + posDeck.Name)
				}
				fmt.Println()
				deckNum := Prompt("Enter Number of deck to view or \"n\" to continue:")
				deck, err := strconv.Atoi(deckNum[0 : len(deckNum)-1])
				checkErr(err)
				if (deck <= (len(posDecks) - 1)) && err == nil{
					for card, count := range posDecks[deck].CardList {
						fmt.Println(card, count)
					}
					fmt.Println()
				} else {
                    fmt.Println("Not valid deck number, exited view.")
                }
			case "quit\n":
				cont = false
			case "help\n":
				fmt.Println("Options:")
				fmt.Println("add")
				fmt.Println("view decks")
				fmt.Println("quit")
				fmt.Println("help")
			default:
				fmt.Println("Please enter correct choice.")
			}
		}
		fmt.Println()
		fmt.Println(enemyDeck.Name)
		for card, count := range enemyDeck.CardList {
			for i := 0; i < count; i++ {
				fmt.Println(card)
			}
		}
	}

}

func Mode(decks AllDecks) []Deck {
	fmt.Println("What mode are you playing? ")
	text := Prompt("(Standard | Wild | Brawl)")
	text = strings.ToLower(text)
	if text == "standard\n" {
		return decks.StandardDecks
	} else if text == "wild\n" {
		return decks.WildDecks
	} else if text == "brawl\n" {
		return decks.BrawlDecks
	} else {
		fmt.Println("Please enter standard, wild, or brawl:")
		return Mode(decks)
	}
}

func CardPlayed(enemyDeck Deck, decks []Deck, s string) (Deck, []Deck) {
	s = s[:len(s)-1]
	newDecks := make([]Deck, 0)
	found := false
	for _, deck := range decks {
		for _, card := range deck.ClassKeys {
			if strings.ToLower(card.Name) == strings.ToLower(s) {
				newDecks = append(newDecks, deck)
				if !found {
					enemyDeck.CardList[card.Name] += 1
					enemyDeck.ClassKeys = append(enemyDeck.ClassKeys, card)
					found = true
				}
				break
			}
		}
		for _, card := range deck.NeutralKeys {
			if strings.ToLower(card.Name) == strings.ToLower(s) {
				newDecks = append(newDecks, deck)
				if !found {
					enemyDeck.CardList[card.Name] += 1
					enemyDeck.NeutralKeys = append(enemyDeck.NeutralKeys, card)
					found = true
				}
				break
			}
		}
	}
	return enemyDeck, newDecks
}

func Prompt(text string) string {
	fmt.Println(text)
	text, err := reader.ReadString('\n')
	checkErr(err)
	return text
}
