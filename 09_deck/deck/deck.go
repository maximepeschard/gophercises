// Package deck provides utilities to generate decks of playing cards.
package deck

import (
	"math/rand"
	"sort"
	"time"
)

// A Deck is a collection of cards.
type Deck []Card

// An Option allows to modify a deck in place.
type Option func(*Deck) error

// New returns a generated deck of cards with the given options applied.
// Options are applied in the same order they are provided.
func New(options ...Option) (Deck, error) {
	var deck Deck

	for s := Clubs; s <= Spades; s++ {
		for r := Ace; r <= King; r++ {
			deck = append(deck, Card{Rank: r, Suit: s})
		}
	}

	for _, option := range options {
		err := option(&deck)
		if err != nil {
			return nil, err
		}
	}

	return deck, nil
}

// Jokers returns an option to add n joker cards to a deck.
func Jokers(n int) Option {
	return func(deck *Deck) error {
		jokers := make(Deck, n)
		for i := 0; i < n; i++ {
			jokers[i] = Card{Suit: Joker}
		}
		*deck = append(*deck, jokers...)
		return nil
	}
}

// Sort returns an option to sort the cards of a deck with a user-defined comparison function.
func Sort(less func(c1, c2 Card) bool) Option {
	return func(deck *Deck) error {
		sort.Slice(*deck, func(i, j int) bool {
			return less((*deck)[i], (*deck)[j])
		})
		return nil
	}
}

// Shuffle returns an option to shuffle the cards of a deck.
func Shuffle() Option {
	rand.Seed(time.Now().UnixNano())

	return func(deck *Deck) error {
		rand.Shuffle(len(*deck), func(i, j int) {
			(*deck)[i], (*deck)[j] = (*deck)[j], (*deck)[i]
		})
		return nil
	}
}

// Filter returns an option to filter out cards in a deck.
func Filter(keep func(Card) bool) Option {
	return func(deck *Deck) error {
		var newDeck Deck

		for _, card := range *deck {
			if keep(card) {
				newDeck = append(newDeck, card)
			}
		}
		*deck = newDeck

		return nil
	}
}

// Repeat returns an option to repeat the deck n times in order to build a single deck.
func Repeat(n int) Option {
	return func(deck *Deck) error {
		base := make(Deck, len(*deck))
		copy(base, *deck)
		for i := 0; i < n; i++ {
			*deck = append(*deck, base...)
		}
		return nil
	}
}
