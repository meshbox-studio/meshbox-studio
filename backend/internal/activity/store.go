package activity

import "sync"

var store = NewStore()

type Store struct {
	mu    sync.RWMutex
	items []Item
}

func NewStore() *Store {
	s := &Store{}
	s.seed()
	return s
}

func (s *Store) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Item, len(s.items))
	copy(result, s.items)
	return result
}

func (s *Store) seed() {
	s.items = []Item{
		{
			ID:          "a1",
			At:          "Today, 09:42",
			Title:       "Queued iteration 07",
			Description: "Hex Lamp V2 moved to Bambu Lab P1S queue.",
		},
		{
			ID:          "a2",
			At:          "Yesterday, 20:14",
			Title:       "Updated drawer geometry",
			Description: "Tool Grid Drawers got new label clearance tolerances.",
		},
		{
			ID:          "a3",
			At:          "2 days ago",
			Title:       "Archived project",
			Description: "Coffee Dosing Cup marked as finalized and archived.",
		},
	}
}