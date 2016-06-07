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

    for i := offset; i < offset + 4 + len(deck.Keys); i++ {
        row := sheet.AddRow()
        for x := 0; x < 8; x++ {
            row.AddCell()
        }
    }

    sheet.Rows[offset].Cells[0].SetValue("Deck Name:")
    sheet.Rows[offset].Cells[1].SetValue(deck.Name)
    offset++
    sheet.Rows[offset].Cells[0].SetValue("Game Mode:")
    sheet.Rows[offset].Cells[1].SetValue(deck.Mode)
    sheet.Rows[offset].Cells[2].SetValue("Deck Type:")
    sheet.Rows[offset].Cells[3].SetValue(deck.Type)
    offset++
    sheet.Rows[offset].Cells[0].SetValue("Card Name:")
    sheet.Rows[offset].Cells[2].SetValue("Card Cost:")
    sheet.Rows[offset].Cells[3].SetValue("Card Count:")
    sheet.Rows[offset].Cells[4].SetValue("Card Class:")
    sheet.Rows[offset].Cells[5].SetValue("Card Type:")
    sheet.Rows[offset].Cells[6].SetValue("Attack:")
    sheet.Rows[offset].Cells[7].SetValue("Health:")
    offset++
    sort.Sort(ByCost(deck.Keys))
    for idx, card := range deck.Keys {
        sheet.Rows[offset + idx].Cells[0].SetValue(card.Name)
        sheet.Rows[offset + idx].Cells[2].SetValue(card.Cost)
        sheet.Rows[offset + idx].Cells[3].SetValue(deck.CardList[card])
        sheet.Rows[offset + idx].Cells[4].SetValue(card.Class)
        sheet.Rows[offset + idx].Cells[5].SetValue(card.Type)
        sheet.Rows[offset + idx].Cells[6].SetValue(card.Attack)
        sheet.Rows[offset + idx].Cells[7].SetValue(card.Health)
    }
    offset += len(deck.Keys) + 1
    return offset

}