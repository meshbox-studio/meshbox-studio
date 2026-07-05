package projects

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Store struct {
	mu       sync.RWMutex
	dataDir  string
	projects []Project
}

func NewStore(dataDir string) *Store {
	s := &Store{dataDir: dataDir}
	s.load()
	return s
}

func (s *Store) projectsDir() string {
	return filepath.Join(s.dataDir, "projects")
}

func (s *Store) projectDir(id string) string {
	return filepath.Join(s.projectsDir(), id)
}

func (s *Store) projectPath(id string) string {
	return filepath.Join(s.projectDir(id), "project.json")
}

func (s *Store) filesDir(id string) string {
	return filepath.Join(s.projectDir(id), "files")
}

func (s *Store) thumbnailPath(id string) string {
	return filepath.Join(s.projectDir(id), "thumbnail.jpg")
}

func (s *Store) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	dir := s.projectsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			s.seed()
			s.persistAll()
			return
		}
		s.seed()
		return
	}

	var loaded []Project
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pj := s.projectPath(entry.Name())
		data, err := os.ReadFile(pj)
		if err != nil {
			continue
		}
		var p Project
		if err := json.Unmarshal(data, &p); err != nil {
			continue
		}
		if p.ID == "" {
			p.ID = entry.Name()
		}
		loaded = append(loaded, p)
	}

	if len(loaded) == 0 {
		s.seed()
		s.persistAll()
		return
	}

	s.projects = loaded
}

func (s *Store) persist(p Project) {
	dir := s.projectDir(p.ID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return
	}
	if err := os.MkdirAll(s.filesDir(p.ID), 0755); err != nil {
		return
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(s.projectPath(p.ID), data, 0644)
}

func (s *Store) persistAll() {
	for _, p := range s.projects {
		s.persist(p)
	}
}

func (s *Store) List() []Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Project, len(s.projects))
	copy(result, s.projects)
	return result
}

func (s *Store) Get(id string) (Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, p := range s.projects {
		if p.ID == id {
			return p, true
		}
	}
	return Project{}, false
}

func (s *Store) Create(input DraftCreate) Project {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := generateID(input.Title)

	project := Project{
		ID:           id,
		Title:        input.Title,
		Description:  input.Description,
		State:        StateDraft,
		LastActivity: "just now",
		Tags:         []string{},
		Files:        []File{},
		Iterations:   []Iteration{},
		Notes:        []string{},
	}

	s.projects = append([]Project{project}, s.projects...)
	s.persist(project)
	return project
}

func (s *Store) CreateFromImport(project Project) Project {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.projects = append([]Project{project}, s.projects...)
	s.persist(project)
	return project
}

func (s *Store) Update(id string, update StateUpdate) (Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.projects {
		if p.ID == id {
			p.State = update.State
			p.LastActivity = "just now"
			s.projects[i] = p
			s.persist(p)
			return p, true
		}
	}
	return Project{}, false
}

func (s *Store) Edit(id string, input EditPayload) (Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.projects {
		if p.ID == id {
			if input.Title != "" {
				p.Title = input.Title
			}
			p.Description = input.Description
			if input.Tags != nil {
				p.Tags = input.Tags
			}
			if input.Source != nil {
				p.Source = input.Source
			}
			if input.License != nil {
				p.License = input.License
			}
			p.LastActivity = "just now"
			s.projects[i] = p
			s.persist(p)
			return p, true
		}
	}
	return Project{}, false
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.projects {
		if p.ID == id {
			s.projects = append(s.projects[:i], s.projects[i+1:]...)
			os.RemoveAll(s.projectDir(id))
			return true
		}
	}
	return false
}

func generateID(title string) string {
	slug := strings.TrimSpace(strings.ToLower(title))
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, slug)
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "new-project"
	}

	return fmt.Sprintf("%s-%s", slug, time.Now().Format("060102150405"))
}

type DraftCreate struct {
	Title       string
	Description string
}

func (s *Store) seed() {
	s.projects = []Project{
		{
			ID: "hex-lamp-v2", Title: "Hex Lamp V2",
			Description: "Modular lamp shades with print-in-place clip joints and alternate diffuser thickness tests.",
			State: StatePrinting,
			LastActivity: "2 hours ago",
			Tags: []string{"lighting", "functional", "v2"},
			Files: []File{
				{ID: "f-hex-1", Name: "hex-lamp-cap-v2.3mf", Type: FileType3MF, Size: "14.3 MB", UpdatedAt: "2026-07-03"},
				{ID: "f-hex-2", Name: "hex-lamp-body-v2.stl", Type: FileTypeSTL, Size: "6.8 MB", UpdatedAt: "2026-07-02"},
				{ID: "f-hex-3", Name: "hex-lamp-cap-v2.gcode", Type: FileTypeGCode, Size: "21.7 MB", UpdatedAt: "2026-07-03"},
			},
			Iterations: []Iteration{
				{ID: "it-hex-1", Label: "Iteration 07", Printer: "Bambu Lab P1S", Material: "PETG - Warm Sand", StartedAt: "2026-07-03 20:12", Duration: "03h 41m", Outcome: OutcomeQueued, Notes: "Running bridge and overhang checks with 0.16 layer height."},
				{ID: "it-hex-2", Label: "Iteration 06", Printer: "Prusa MK4S", Material: "PLA - Bone White", StartedAt: "2026-07-01 18:33", Duration: "04h 02m", Outcome: OutcomeFailed, Notes: "Warp on panel corner. Need higher bed temp and 8 mm brim."},
				{ID: "it-hex-3", Label: "Iteration 05", Printer: "Prusa MK4S", Material: "PLA - Bone White", StartedAt: "2026-06-29 09:17", Duration: "03h 35m", Outcome: OutcomeSuccess, Notes: "Clips fit tight. Diffuser thickness changed to 0.9 mm."},
			},
			Notes: []string{"Target: keep each panel under 4h print time.", "Need to test heat-soak near LED strip after 2h runtime.", "Move cable channel 1.2 mm left in next revision."},
		},
		{
			ID: "tool-grid-drawers", Title: "Tool Grid Drawers",
			Description: "Parametric workshop drawer bins with snap labels and TPU anti-slip inserts.",
			State: StateActive,
			LastActivity: "Yesterday",
			Tags: []string{"organization", "workshop", "parametric"},
			Files: []File{
				{ID: "f-grid-1", Name: "drawer-40x60-v8.step", Type: FileTypeStep, Size: "3.1 MB", UpdatedAt: "2026-07-02"},
				{ID: "f-grid-2", Name: "drawer-40x60-v8.3mf", Type: FileType3MF, Size: "11.2 MB", UpdatedAt: "2026-07-02"},
			},
			Iterations: []Iteration{
				{ID: "it-grid-1", Label: "Iteration 03", Printer: "Prusa MK4S", Material: "PLA - Graphite", StartedAt: "2026-07-02 19:05", Duration: "02h 19m", Outcome: OutcomeSuccess, Notes: "Label tabs fit perfectly; no elephant foot after first layer tweak."},
				{ID: "it-grid-2", Label: "Iteration 02", Printer: "Prusa MK4S", Material: "PLA - Graphite", StartedAt: "2026-06-25 17:44", Duration: "02h 08m", Outcome: OutcomeFailed, Notes: "Front edge lifted. Added mouse ears for next run."},
			},
			Notes: []string{"Label slots need a looser tolerance for matte filaments.", "TPU insert profile now uses 12% gyroid infill."},
		},
		{
			ID: "coffee-dosing-cup", Title: "Coffee Dosing Cup",
			Description: "Single-dose espresso cup with anti-static ribs and magnetic base ring.",
			State: StateArchived,
			LastActivity: "2 weeks ago",
			Tags: []string{"kitchen", "daily", "done"},
			Files: []File{
				{ID: "f-coffee-1", Name: "dosing-cup-v4.stl", Type: FileTypeSTL, Size: "2.6 MB", UpdatedAt: "2026-06-18"},
				{ID: "f-coffee-2", Name: "dosing-cup-v4.gcode", Type: FileTypeGCode, Size: "9.3 MB", UpdatedAt: "2026-06-18"},
			},
			Iterations: []Iteration{
				{ID: "it-coffee-1", Label: "Iteration 04", Printer: "Bambu Lab P1S", Material: "PETG - Smoke", StartedAt: "2026-06-18 07:20", Duration: "01h 11m", Outcome: OutcomeSuccess, Notes: "Final release candidate."},
			},
			Notes: []string{"Archived after final fit check with 58.5 mm basket."},
		},
		{
			ID: "camera-cable-clip", Title: "Camera Cable Clip",
			Description: "Small desk cable guide for webcam arm and USB-C power routing.",
			State: StateTrash,
			LastActivity: "3 weeks ago",
			Tags: []string{"desk", "quick-fix"},
			Files: []File{
				{ID: "f-clip-1", Name: "clip-v1.stl", Type: FileTypeSTL, Size: "0.9 MB", UpdatedAt: "2026-06-10"},
			},
			Iterations: []Iteration{
				{ID: "it-clip-1", Label: "Iteration 01", Printer: "Prusa MK4S", Material: "PLA - Carbon Black", StartedAt: "2026-06-10 13:42", Duration: "00h 24m", Outcome: OutcomeFailed, Notes: "Too brittle. Replaced by TPU version in another project."},
			},
			Notes: []string{"Moved to trash after switching to snap-on arm clip."},
		},
	}
}