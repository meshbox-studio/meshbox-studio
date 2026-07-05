package projects

type State string

const (
	StateActive   State = "active"
	StatePrinting State = "printing"
	StateArchived State = "archived"
	StateTrash    State = "trash"
	StateDraft    State = "draft"
)

type Outcome string

const (
	OutcomeSuccess Outcome = "success"
	OutcomeFailed  Outcome = "failed"
	OutcomeQueued  Outcome = "queued"
)

type FileType string

const (
	FileTypeSTL   FileType = "stl"
	FileType3MF   FileType = "3mf"
	FileTypeGCode FileType = "gcode"
	FileTypeStep  FileType = "step"
	FileTypePDF   FileType = "pdf"
)

type File struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Type      FileType `json:"type" enum:"stl,3mf,gcode,step,pdf"`
	Size      string   `json:"size"`
	UpdatedAt string   `json:"updatedAt"`
}

type Iteration struct {
	ID        string  `json:"id"`
	Label     string  `json:"label"`
	Printer   string  `json:"printer"`
	Material  string  `json:"material"`
	StartedAt string  `json:"startedAt"`
	Duration  string  `json:"duration"`
	Outcome   Outcome `json:"outcome" enum:"success,failed,queued"`
	Notes     string  `json:"notes"`
}

type Source struct {
	URL         string `json:"url"`
	Designer    string `json:"designer"`
	DesignerURL string `json:"designerUrl"`
	Category    string `json:"category"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PrintProfile struct {
	PrintTimeMinutes int     `json:"printTimeMinutes"`
	WeightGrams      float64 `json:"weightGrams"`
	Quantity         int     `json:"quantity"`
	NozzleMM         float64 `json:"nozzleMm"`
	LayerHeightMM    float64 `json:"layerHeightMm"`
	Material         string  `json:"material"`
	Printer          string  `json:"printer"`
}

type Project struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	State        State         `json:"state" enum:"active,printing,archived,trash,draft"`
	LastActivity string        `json:"lastActivity"`
	Tags         []string      `json:"tags"`
	Files        []File        `json:"files"`
	Iterations   []Iteration   `json:"iterations"`
	Notes        []string      `json:"notes"`
	Source       *Source       `json:"source"`
	License      *License      `json:"license"`
	PrintProfile *PrintProfile `json:"printProfile"`
	HasThumbnail bool          `json:"hasThumbnail"`
}

type StateUpdate struct {
	State State `json:"state" enum:"active,printing,archived,trash,draft"`
}

type EditPayload struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Source      *Source  `json:"source"`
	License     *License `json:"license"`
}