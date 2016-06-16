package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func startTracking(deckLists AllDecks) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Would you like to start a game? (y/n) ")
    text, err := reader.ReadString('\n')
    checkErr(err)
    if text == "y\n" {
        fmt.Println("What mode are you plating? ")
        posDecks := Mode(deckLists)
        cont := true
        for cont {
            fmt.Println("Please enter card enemy played: ")
            text, err := reader.ReadString('\n')
            checkErr(err)
            posDecks = CardPlayed(posDecks, text)

            fmt.Println()
            for _ ,deck := range(posDecks){
                fmt.Println(deck.Name)
            }
            fmt.Println()
            fmt.Println("Possible Decks Updated.")
            fmt.Println("Continue? (y/n) ")
            text, err = reader.ReadString('\n')
            checkErr(err)
            if text != "y\n"{
                cont = false
            }
        }
    } 
}

func Mode(decks AllDecks) []Deck{
    reader := bufio.NewReader(os.Stdin)
    text, err := reader.ReadString('\n')
    checkErr(err)
    text = strings.ToLower(text)
    if text == "standard\n"{
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

func CardPlayed(decks []Deck, s string) []Deck {
    s = s[:len(s)-1]
    newDecks := make([]Deck, 0)
    for _,deck := range(decks){
        for _,card := range(deck.ClassKeys){
            if strings.ToLower(card.Name) == strings.ToLower(s) {
                newDecks = append(newDecks, deck)
            }
        }
        for _,card := range(deck.NeutralKeys){
            if strings.ToLower(card.Name) == strings.ToLower(s) {
                newDecks = append(newDecks, deck)
            }
        }
    }
    return newDecks
}
