package deck

import "testing"

func TestJokers(t *testing.T) {
	basicDeck, _ := New()
	basicDeckLength := len(basicDeck)

	nbJokers := 4
	deck, _ := New(Jokers(nbJokers))
	if deckLength := len(deck); deckLength != basicDeckLength+nbJokers {
		t.Errorf("New(Jokers(%d)) result has length %d, want %d", nbJokers, deckLength, basicDeckLength+nbJokers)
	} else {
		jokerCount := 0
		for _, card := range deck {
			if card.Suit == Joker {
				jokerCount++
			}
		}
		if jokerCount != nbJokers {
			t.Errorf("New(Jokers(%d)) result has %d jokers, want %d", nbJokers, jokerCount, nbJokers)
		}
	}
}

func TestSort(t *testing.T) {
	basicDeck, _ := New()
	basicDeckLength := len(basicDeck)

	less := func(c1, c2 Card) bool { return c1.Rank >= c2.Rank }
	deck, _ := New(Sort(less))
	if deckLength := len(deck); deckLength != basicDeckLength {
		t.Errorf("New(Sort(less)) result has length %d, want %d", deckLength, basicDeckLength)
	} else {
		for i := 0; i < deckLength-1; i++ {
			if !less(deck[i], deck[i+1]) {
				t.Errorf("New(Sort(less)) result does not satisfy the less function")
				break
			}
		}
	}
}

func TestFilter(t *testing.T) {
	basicDeck, _ := New()
	basicDeckLength := len(basicDeck)
	var deck Deck

	// Test a filter that keeps all cards
	noFilter := Filter(func(c Card) bool { return true })
	deck, _ = New()
	noFilter(&deck)
	if deckLength := len(deck); deckLength != basicDeckLength {
		t.Errorf("noFilter result has length %d, want %d", deckLength, basicDeckLength)
	} else {
		for i := 0; i < deckLength; i++ {
			if deck[i] != basicDeck[i] {
				t.Errorf("noFilter result does not match original deck (first mismatch at index %d: got %v, want %v", i, deck[i], basicDeck[i])
				break
			}
		}
	}

	// Test a filter that filters out some cards
	removeThreesFilter := Filter(func(c Card) bool { return c.Rank != Three })
	deck, _ = New()
	removeThreesFilter(&deck)
	if deckLength := len(deck); deckLength != basicDeckLength-4 {
		t.Errorf("removeThreesFilter result has length %d, want %d", deckLength, basicDeckLength-4)
	} else {
		for i, card := range deck {
			if card.Rank == Three {
				t.Errorf("removeThreesFilter result contains Threes (first occurence at index %d: %v)", i, card)
				break
			}
		}
	}
}

func TestRepeat(t *testing.T) {
	basicDeck, _ := New()
	basicDeckLength := len(basicDeck)

	nbRepeat := 2
	deck, _ := New(Repeat(nbRepeat))
	if deckLength := len(deck); deckLength != (nbRepeat+1)*basicDeckLength {
		t.Errorf("New(Repeat(%d)) result has length %d, want %d", nbRepeat, deckLength, (nbRepeat+1)*basicDeckLength)
	} else {
		cardCount := make(map[Card]int)
		for _, card := range deck {
			cardCount[card]++
		}

		for card, count := range cardCount {
			if count != nbRepeat+1 {
				t.Errorf("New(Repeat(%d)) result has cards with count other than %d (first occurence: %v with count %d)", nbRepeat, nbRepeat+1, card, count)
				break
			}
		}
	}
}
