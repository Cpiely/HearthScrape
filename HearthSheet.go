package main

import (
    "fmt"
    "sort"
    "github.com/tealeg/xlsx"
)

func CreateFile() *xlsx.File {
    file := xlsx.NewFile()
    return file
}

func CreateSheet(file *xlsx.File, sheetName string) *xlsx.Sheet {
    sheet, err := file.AddSheet(sheetName)
    if err != nil {
        fmt.Printf(err.Error())
    }
    return sheet
}

func SaveFile(file *xlsx.File) {
    err := file.Save("MyXLSXFile.xlsx")
    if err != nil {
        fmt.Printf(err.Error())
    }
}

func PrintDeck(sheet *xlsx.Sheet, deck Deck, offset int) int{
    sort.Sort(ByCost(deck.ClassKeys))
    sort.Sort(ByCost(deck.NeutralKeys))

    for i := offset; i < offset + 9 + len(deck.ClassKeys) + len(deck.NeutralKeys); i++ {
        row := sheet.AddRow()
        for x := 0; x < 8; x++ {
            row.AddCell()
        }
    }
    sheet.Rows[offset].Cells[0].SetValue("Deck Name:")
    sheet.Rows[offset].Cells[1].SetValue(deck.Name)
    sheet.Rows[offset].Cells[3].SetValue("Class:")
    sheet.Rows[offset].Cells[4].SetValue(deck.Class)
    offset++
    sheet.Rows[offset].Cells[0].SetValue("Url:")
    sheet.Rows[offset].Cells[1].SetValue(deck.Url)
    sheet.Rows[offset].Cells[3].SetValue("Deck Rating:")
    sheet.Rows[offset].Cells[4].SetValue(deck.Rating)
    offset++
    sheet.Rows[offset].Cells[0].SetValue("Game Mode:")
    sheet.Rows[offset].Cells[1].SetValue(deck.Mode)
    if deck.Type != "Tavern Brawl" {
        sheet.Rows[offset].Cells[2].SetValue("Deck Type:")
        sheet.Rows[offset].Cells[3].SetValue(deck.Type)
    }
    sheet.Rows[offset].Cells[4].SetValue(deck.DateCreated)
    offset++
    sheet.Rows[offset].Cells[0].SetValue("Expansion:")
    sheet.Rows[offset].Cells[1].SetValue(deck.Expansion)
    sheet.Rows[offset].Cells[2].SetValue("Deck Cost:")
    sheet.Rows[offset].Cells[3].SetValue(deck.Cost)
    offset++
    sheet.Rows[offset].Cells[1].SetValue("Mana Cost:")
    sheet.Rows[offset].Cells[2].SetValue("Attack:")
    sheet.Rows[offset].Cells[3].SetValue("Health:")
    sheet.Rows[offset].Cells[4].SetValue("Count:")
    offset++
    sheet.Rows[offset].Cells[0].SetValue(deck.Class + " Cards:")
    offset++
    for idx, card := range deck.ClassKeys {
        sheet.Rows[offset + idx].Cells[0].SetValue(card.Name)
        sheet.Rows[offset + idx].Cells[1].SetValue(card.Cost)
        if card.Type == "Minion" {
            sheet.Rows[offset + idx].Cells[2].SetValue(card.Attack)
            sheet.Rows[offset + idx].Cells[3].SetValue(card.Health)
        }
        sheet.Rows[offset + idx].Cells[4].SetValue(deck.CardList[card])
        
    }
    offset += len(deck.ClassKeys)
    offset++
    sheet.Rows[offset].Cells[0].SetValue("Neutral Cards:")
    offset++
    for idx, card := range deck.NeutralKeys {
        sheet.Rows[offset + idx].Cells[0].SetValue(card.Name)
        sheet.Rows[offset + idx].Cells[1].SetValue(card.Cost)
        if card.Type == "Minion" {
            sheet.Rows[offset + idx].Cells[2].SetValue(card.Attack)
            sheet.Rows[offset + idx].Cells[3].SetValue(card.Health)
        }
        sheet.Rows[offset + idx].Cells[4].SetValue(deck.CardList[card])
    }
    offset += len(deck.NeutralKeys) + 1
    return offset
}
