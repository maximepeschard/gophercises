package story

import "encoding/json"

// An ArcOption represents a choice of story arc
type ArcOption struct {
	Text      string
	TargetArc string `json:"arc"`
}

// An Arc represents a portion of a story
type Arc struct {
	Title   string
	Text    []string `json:"story"`
	Options []ArcOption
}

// Story represents a whole story with its arcs
type Story struct {
	Name string
	Arcs map[string]Arc
}

// NewStory returns a new story (with no arcs) given a name
func NewStory(name string) *Story {
	return &Story{Name: name}
}

// ParseJSON uses the provided JSON data to populate the arcs of the story
func (s *Story) ParseJSON(data []byte) (*Story, error) {
	var arcs map[string]Arc
	err := json.Unmarshal(data, &arcs)
	if err != nil {
		return s, err
	}

	s.Arcs = arcs
	return s, err
}
