package projects

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/meshbox-studio/meshbox-studio/internal/auth"
	"github.com/meshbox-studio/meshbox-studio/internal/projects/pdfmeta"
)

var store *Store

func RegisterRoutes(api huma.API, mux *http.ServeMux, dataDir string) {
	store = NewStore(dataDir)

	mw := auth.RequireSession(api)

	huma.Register(api, huma.Operation{
		OperationID: "list-projects",
		Method:      http.MethodGet,
		Path:        "/api/projects",
		Summary:     "List projects",
		Description: "Returns all projects.",
		Tags:        []string{"Projects"},
		Middlewares: huma.Middlewares{mw},
	}, listHandler)

	huma.Register(api, huma.Operation{
		OperationID: "create-project",
		Method:      http.MethodPost,
		Path:        "/api/projects",
		Summary:     "Create a draft project",
		Description: "Creates a new project in draft state.",
		Tags:        []string{"Projects"},
		Middlewares: huma.Middlewares{mw},
		DefaultStatus: http.StatusCreated,
	}, createHandler)

	huma.Register(api, huma.Operation{
		OperationID: "import-project",
		Method:      http.MethodPost,
		Path:        "/api/projects/import",
		Summary:     "Import project from Printables zip",
		Description: "Upload a Printables model zip. Extracts files, PDF metadata, and thumbnail. Returns 422 if no valid Printables PDF found.",
		Tags:        []string{"Projects"},
		Middlewares: huma.Middlewares{mw},
		DefaultStatus: http.StatusCreated,
	}, importHandler)

	huma.Register(api, huma.Operation{
		OperationID: "update-project",
		Method:      http.MethodPatch,
		Path:        "/api/projects/{id}",
		Summary:     "Update project state",
		Description: "Updates a project's state (e.g. active, trash, archived).",
		Tags:        []string{"Projects"},
		Middlewares: huma.Middlewares{mw},
	}, updateHandler)

	huma.Register(api, huma.Operation{
		OperationID: "edit-project",
		Method:      http.MethodPut,
		Path:        "/api/projects/{id}",
		Summary:     "Edit project metadata",
		Description: "Updates a project's title, description, tags, source, and license.",
		Tags:        []string{"Projects"},
		Middlewares: huma.Middlewares{mw},
	}, editHandler)

	huma.Register(api, huma.Operation{
		OperationID: "delete-project",
		Method:      http.MethodDelete,
		Path:        "/api/projects/{id}",
		Summary:     "Delete project forever",
		Description: "Permanently removes the project.",
		Tags:        []string{"Projects"},
		Middlewares: huma.Middlewares{mw},
	}, deleteHandler)

	mux.HandleFunc("GET /api/projects/{id}/files/{fileId}", fileDownloadHandler)
	mux.HandleFunc("GET /api/projects/{id}/thumbnail", thumbnailHandler)
}

type ListOutput struct {
	Body []Project
}

type GetInput struct {
	ID string `path:"id"`
}

type GetOutput struct {
	Body Project
}

type CreateInput struct {
	Body struct {
		Title       string `json:"title" minLength:"1" maxLength:"120"`
		Description string `json:"description"`
	}
}

type ImportInput struct {
	RawBody multipart.Form
}

type ImportOutput struct {
	Body Project
}

type EditInputBody struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Source      *Source  `json:"source"`
	License     *License `json:"license"`
}

type EditInput struct {
	ID   string        `path:"id"`
	Body EditInputBody
}

type EditOutput struct {
	Body Project
}

type UpdateInput struct {
	ID   string      `path:"id"`
	Body StateUpdate
}

type UpdateOutput struct {
	Body Project
}

func listHandler(ctx context.Context, input *struct{}) (*ListOutput, error) {
	return &ListOutput{Body: store.List()}, nil
}

func createHandler(ctx context.Context, input *CreateInput) (*GetOutput, error) {
	project := store.Create(DraftCreate{
		Title:       input.Body.Title,
		Description: input.Body.Description,
	})
	return &GetOutput{Body: project}, nil
}

func importHandler(ctx context.Context, input *ImportInput) (*ImportOutput, error) {
	files := input.RawBody.File["file"]
	if len(files) == 0 {
		return nil, huma.Error422UnprocessableEntity("missing file", &huma.ErrorDetail{Message: "No file provided in 'file' field"})
	}

	fileHeader := files[0]
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".zip") {
		return nil, huma.Error422UnprocessableEntity("invalid file type", &huma.ErrorDetail{Message: "Only .zip files are accepted"})
	}

	zipData, err := readFileHeader(fileHeader)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("cannot read file", &huma.ErrorDetail{Message: err.Error()})
	}

	var pdfData []byte
	var extractedFiles []extractedFile

	zipReader, err := zip.NewReader(strings.NewReader(string(zipData)), int64(len(zipData)))
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("invalid zip", &huma.ErrorDetail{Message: "Cannot read zip: " + err.Error()})
	}

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(f.Name), ".pdf") {
			pdfData, err = readZipFile(f)
			if err != nil {
				return nil, huma.Error422UnprocessableEntity("zip read error", &huma.ErrorDetail{Message: "Cannot read PDF from zip: " + err.Error()})
			}
			continue
		}
		ef := extractedFile{
			name: filepath.Base(f.Name),
			data: mustReadZipFile(f),
		}
		if ef.data != nil {
			extractedFiles = append(extractedFiles, ef)
		}
	}

	if pdfData == nil {
		return nil, huma.Error422UnprocessableEntity("no pdf", &huma.ErrorDetail{Message: "No PDF found in the zip file. This doesn't look like a Printables export."})
	}

	meta, err := pdfmeta.Extract(pdfData)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity("pdf parse error", &huma.ErrorDetail{Message: err.Error()})
	}

	thumbData, _ := pdfmeta.ExtractThumbnail(pdfData)

	id := generateID(meta.Title)
	now := time.Now().Format("2006-01-02")

	var source *Source
	if meta.SourceURL != "" {
		source = &Source{
			URL:         meta.SourceURL,
			Designer:    meta.Designer,
			DesignerURL: meta.DesignerURL,
			Category:    meta.Category,
		}
	}

	var license *License
	if meta.LicenseURL != "" || meta.LicenseName != "" {
		license = &License{
			Name: meta.LicenseName,
			URL:  meta.LicenseURL,
		}
	}

	var profile *PrintProfile
	if meta.PrintTimeMinutes > 0 || meta.WeightGrams > 0 || meta.NozzleMM > 0 || meta.LayerHeightMM > 0 {
		profile = &PrintProfile{
			PrintTimeMinutes: meta.PrintTimeMinutes,
			WeightGrams:      meta.WeightGrams,
			Quantity:         meta.Quantity,
			NozzleMM:         meta.NozzleMM,
			LayerHeightMM:    meta.LayerHeightMM,
			Material:         meta.Material,
			Printer:          meta.Printer,
		}
	}

	var filesList []File
	pdfFileID := fmt.Sprintf("f-pdf-%s", id)
	filesList = append(filesList, File{
		ID:        pdfFileID,
		Name:      filepath.Base(fileHeader.Filename + ".pdf"), 
		Type:      FileTypePDF,
		Size:      humanBytes(int64(len(pdfData))),
		UpdatedAt: now,
	})

	for i, ef := range extractedFiles {
		fileID := fmt.Sprintf("f-%d-%s", i, id)
		filesList = append(filesList, File{
			ID:        fileID,
			Name:      ef.name,
			Type:      detectFileType(ef.name),
			Size:      humanBytes(int64(len(ef.data))),
			UpdatedAt: now,
		})
	}

	if meta.Title == "" && meta.SourceURL != "" {
		meta.Title = "Imported from Printables"
	}
	if meta.Title == "" {
		meta.Title = strings.TrimSuffix(fileHeader.Filename, ".zip")
	}
	if meta.Description == "" && meta.Title != "" {
		meta.Description = "Imported from " + meta.Title
	}

	project := Project{
		ID:           id,
		Title:        meta.Title,
		Description:  meta.Description,
		State:        StateDraft,
		LastActivity: "just now",
		Tags:         meta.Tags,
		Files:        filesList,
		Iterations:   []Iteration{},
		Notes:        []string{},
		Source:       source,
		License:      license,
		PrintProfile: profile,
		HasThumbnail: thumbData != nil,
	}

	store.CreateFromImport(project)

	storeDir := store.filesDir(id)
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return nil, huma.Error422UnprocessableEntity("storage error", &huma.ErrorDetail{Message: "Cannot create storage directory"})
	}

	pdfStorePath := filepath.Join(storeDir, pdfFileID+".pdf")
	if err := os.WriteFile(pdfStorePath, pdfData, 0644); err != nil {
		return nil, huma.Error422UnprocessableEntity("storage error", &huma.ErrorDetail{Message: "Cannot store PDF file"})
	}

	for i, ef := range extractedFiles {
		fileID := fmt.Sprintf("f-%d-%s", i, id)
		storePath := filepath.Join(storeDir, fileID+filepath.Ext(ef.name))
		if err := os.WriteFile(storePath, ef.data, 0644); err != nil {
			continue
		}
	}

	if thumbData != nil {
		thumbPath := store.thumbnailPath(id)
		os.WriteFile(thumbPath, thumbData, 0644)
	}

	return &ImportOutput{Body: project}, nil
}

func updateHandler(ctx context.Context, input *UpdateInput) (*UpdateOutput, error) {
	project, ok := store.Update(input.ID, input.Body)
	if !ok {
		return nil, huma.Error404NotFound("project not found")
	}
	return &UpdateOutput{Body: project}, nil
}

func editHandler(ctx context.Context, input *EditInput) (*EditOutput, error) {
	project, ok := store.Edit(input.ID, EditPayload{
		Title:       input.Body.Title,
		Description: input.Body.Description,
		Tags:        input.Body.Tags,
		Source:      input.Body.Source,
		License:     input.Body.License,
	})
	if !ok {
		return nil, huma.Error404NotFound("project not found")
	}
	return &EditOutput{Body: project}, nil
}

func deleteHandler(ctx context.Context, input *GetInput) (*struct{}, error) {
	if !store.Delete(input.ID) {
		return nil, huma.Error404NotFound("project not found")
	}
	return nil, nil
}

func fileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fileID := r.PathValue("fileId")

	_, ok := store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	storeDir := store.filesDir(id)

	entries, err := os.ReadDir(storeDir)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		entryID := strings.TrimSuffix(name, filepath.Ext(name))
		if entryID == fileID || strings.HasPrefix(name, fileID) {
			http.ServeFile(w, r, filepath.Join(storeDir, name))
			return
		}
	}

	http.NotFound(w, r)
}

func thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	project, ok := store.Get(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	thumbPath := store.thumbnailPath(id)
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		if project.HasThumbnail {
			http.Error(w, "thumbnail missing", http.StatusNotFound)
		} else {
			http.NotFound(w, r)
		}
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, thumbPath)
}

type extractedFile struct {
	name string
	data []byte
}

func readFileHeader(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func readZipFile(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func mustReadZipFile(f *zip.File) []byte {
	data, err := readZipFile(f)
	if err != nil {
		return nil
	}
	return data
}

func detectFileType(name string) FileType {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".stl":
		return FileTypeSTL
	case ".3mf":
		return FileType3MF
	case ".gcode", ".gco":
		return FileTypeGCode
	case ".step", ".stp":
		return FileTypeStep
	case ".pdf":
		return FileTypePDF
	default:
		return FileTypeSTL
	}
}

func humanBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	value := float64(bytes)
	idx := 0
	for value >= unit && idx < len(units)-1 {
		value /= unit
		idx++
	}
	return fmt.Sprintf("%.1f %s", value, units[idx])
}