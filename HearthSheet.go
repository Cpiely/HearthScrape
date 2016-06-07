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

    for i := offset; i < offset + 5 + len(deck.ClassKeys) + len(deck.NeutralKeys); i++ {
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
    sheet.Rows[offset].Cells[0].SetValue("Class Card Name:")
    sheet.Rows[offset].Cells[2].SetValue("Card Cost:")
    sheet.Rows[offset].Cells[3].SetValue("Card Count:")
    sheet.Rows[offset].Cells[4].SetValue("Card Class:")
    sheet.Rows[offset].Cells[5].SetValue("Card Type:")
    sheet.Rows[offset].Cells[6].SetValue("Attack:")
    sheet.Rows[offset].Cells[7].SetValue("Health:")
    offset++
   
    for cIdx, cCard := range deck.ClassKeys {
        sheet.Rows[offset + cIdx].Cells[0].SetValue(cCard.Name)
        sheet.Rows[offset + cIdx].Cells[2].SetValue(cCard.Cost)
        sheet.Rows[offset + cIdx].Cells[3].SetValue(deck.CardList[cCard])
        sheet.Rows[offset + cIdx].Cells[4].SetValue(cCard.Class)
        sheet.Rows[offset + cIdx].Cells[5].SetValue(cCard.Type)
        sheet.Rows[offset + cIdx].Cells[6].SetValue(cCard.Attack)
        sheet.Rows[offset + cIdx].Cells[7].SetValue(cCard.Health)
    }
    offset += len(deck.ClassKeys)
    sheet.Rows[offset].Cells[0].SetValue("Neutral Card Name:")
    offset++
    for nIdx, nCard := range deck.NeutralKeys {
        sheet.Rows[offset + nIdx].Cells[0].SetValue(nCard.Name)
        sheet.Rows[offset + nIdx].Cells[2].SetValue(nCard.Cost)
        sheet.Rows[offset + nIdx].Cells[3].SetValue(deck.CardList[nCard])
        sheet.Rows[offset + nIdx].Cells[4].SetValue(nCard.Class)
        sheet.Rows[offset + nIdx].Cells[5].SetValue(nCard.Type)
        sheet.Rows[offset + nIdx].Cells[6].SetValue(nCard.Attack)
        sheet.Rows[offset + nIdx].Cells[7].SetValue(nCard.Health)
    }
    offset += len(deck.NeutralKeys) + 1
    return offset

}