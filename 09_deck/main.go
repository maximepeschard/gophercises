package main

import (
	"fmt"
	"log"

	"github.com/maximepeschard/gophercises/09_deck/deck"
)

func main() {
	// removeAces := deck.Filter(func(c deck.Card) bool {
	// 	return c.Rank != deck.Ace
	// })
	// customSort := deck.Sort(func(c1, c2 deck.Card) bool {
	// 	return c1.Rank > c2.Rank
	// })
	cards, err := deck.New(deck.Jokers(2), deck.Shuffle())
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(cards)
}
