package story

import "encoding/json"

type ArcOption struct {
	Text      string
	TargetArc string `json:"arc"`
}

type Arc struct {
	Title   string
	Story   []string
	Options []ArcOption
}

type Story struct {
	Name string
	Arcs map[string]Arc
}

func NewStory(name string) *Story {
	return &Story{Name: name}
}

func (s *Story) ParseJSON(data []byte) (*Story, error) {
	var arcs map[string]Arc
	err := json.Unmarshal(data, &arcs)
	if err != nil {
		return s, err
	}

	s.Arcs = arcs
	return s, err
}
