package main

type Deck struct {
	Name string
	Type string
	Url string
	Rating string
	Cost string
	Class string
	Expansion string
	DateCreated string
	DateModified string
	Mode string
	CardList map[Card]int
	ClassKeys []Card
	NeutralKeys []Card
}

type Card struct {
	Name string
	Type string
	Class string
	Cost string
	Attack string
	Health string
}


//Used to sort the Structs by different parameters
type ByName []Card
type ByCost []Card
type ByMode []Deck
type ByType []Deck
type ByDeckName[]Deck
//Sort Decks By name
func (slice ByDeckName) Len() int {return len(slice)}

func (slice ByDeckName) Less(i, j int) bool {return slice[i].Name < slice[j].Name}

func (slice ByDeckName) Swap(i, j int) {slice[i], slice[j] = slice[j], slice[i]}
//Sort Cards By Name
func (slice ByName) Len() int {return len(slice)}

func (slice ByName) Less(i, j int) bool {return slice[i].Name < slice[j].Name}

func (slice ByName) Swap(i, j int) {slice[i], slice[j] = slice[j], slice[i]}
//Sort Cards By Cost
func (slice ByCost) Len() int {return len(slice)}

func (slice ByCost) Less(i, j int) bool {return slice[i].Cost < slice[j].Cost}

func (slice ByCost) Swap(i, j int) {slice[i], slice[j] = slice[j], slice[i]}
//Sort Deck by deck type
func (slice ByType) Len() int {return len(slice)}

func (slice ByType) Less(i, j int) bool {return slice[i].Type < slice[j].Type}

func (slice ByType) Swap(i, j int) {slice[i], slice[j] = slice[j], slice[i]}
//Sort Deck by game mode
func (slice ByMode) Len() int {return len(slice)}

func (slice ByMode) Less(i, j int) bool {return slice[i].Mode < slice[j].Mode}

func (slice ByMode) Swap(i, j int) {slice[i], slice[j] = slice[j], slice[i]}

