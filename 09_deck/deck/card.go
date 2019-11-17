package deck

import (
	"fmt"
	"strconv"
)

// Rank represents the rank of a playing card.
type Rank int

// Playing cards ranks.
const (
	_ Rank = iota // skip first so that Ace == 1
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "1"
	case Ten:
		return "T"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	}

	return strconv.Itoa(int(r))
}

// Suit represents the suit of a playing card.
type Suit int

// Playing cards suits.
const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
	Joker
)

func (s Suit) String() string {
	switch s {
	case Clubs:
		return "c"
	case Diamonds:
		return "d"
	case Hearts:
		return "h"
	case Spades:
		return "s"
	}

	return ""
}

// A Card represents a playing card.
type Card struct {
	Rank Rank
	Suit Suit
}

// String returns the string form of a playing card.
// All cards except jokers are represented by two
// characters : a first uppercase character for the
// rank, and a second lowercase character for the suit.
// Jokers are represented by the string "Joker".
func (c Card) String() string {
	if c.Suit == Joker {
		return "Joker"
	}
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}
