package pdfmeta

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Result struct {
	Title        string
	Description  string
	Tags         []string
	Category     string
	SourceURL    string
	Designer     string
	DesignerURL  string
	LicenseName  string
	LicenseURL   string

	PrintTimeMinutes int
	WeightGrams      float64
	Quantity         int
	NozzleMM         float64
	LayerHeightMM    float64
	Material         string
	Printer          string
}

func Extract(pdfData []byte) (*Result, error) {
	r, err := pdf.NewReader(bytes.NewReader(pdfData), int64(len(pdfData)))
	if err != nil {
		return nil, fmt.Errorf("pdf: open: %w", err)
	}

	total := r.NumPage()
	if total < 1 {
		return nil, fmt.Errorf("pdf: no pages")
	}
	_ = total

	raw := &rawExtract{}
	collectURIs(r, raw)

	plainReader, err := r.GetPlainText()
	if err == nil {
		data, _ := io.ReadAll(plainReader)
		plainLines := strings.Split(string(data), "\n")
		extractFromPlainText(plainLines, raw)
	}

	p1 := r.Page(1)
	if !p1.V.IsNull() {
		rows, _ := p1.GetTextByRow()
		var rowLines []string
		for _, row := range rows {
			var parts []string
			for _, word := range row.Content {
				parts = append(parts, word.S)
			}
			line := strings.TrimSpace(strings.Join(parts, " "))
			if line != "" {
				rowLines = append(rowLines, line)
			}
		}
		fillFromRows(rowLines, raw)
	}

	if raw.sourceURL == "" {
		return nil, fmt.Errorf("not a printables export: no printables.com model link found")
	}

	return buildResult(raw), nil
}

func extractFromPlainText(lines []string, raw *rawExtract) {
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Model files") || strings.HasPrefix(line, "Print files") ||
			strings.HasPrefix(line, "License") {
			continue
		}
		if strings.HasPrefix(line, "Summary") && len(line) > 7 {
			if raw.description == "" {
				raw.description = strings.TrimPrefix(line, "Summary")
			}
			continue
		}

		if strings.Contains(line, "HexScraper") || strings.Contains(line, "hexscraper") {
			if raw.title == "" {
				raw.title = line
				if i > 0 {
					prev := strings.TrimSpace(lines[i-1])
					if len(prev) > 1 && len(prev) < 20 && !isBoring(prev) &&
						!strings.Contains(prev, "hrs") && !strings.Contains(prev, "mm") &&
						!strings.Contains(prev, "pcs") && !strings.Contains(prev, " g ") &&
						prev != "VIEW IN BROWSER" && prev != ">" {
						raw.designer = prev
						raw.designerSet = true
					}
				}
			}
			continue
		}

		if !raw.designerSet {
			if looksLikeDesignerName(line) {
				raw.designer = line
			}
		}

		if raw.description == "" && raw.title != "" && len(line) > 30 &&
			!strings.Contains(line, "hrs") && !strings.Contains(line, "mm") &&
			!strings.Contains(line, " g ") && !strings.HasPrefix(line, "3D") &&
			!strings.HasPrefix(line, "Tags") && !strings.HasPrefix(line, "License") &&
			!strings.HasPrefix(line, "This work") && !strings.HasPrefix(line, "Creative") {
			raw.description = line
		}

		if strings.HasPrefix(line, "3D Printers") {
			raw.category = line
		}

		if strings.Contains(line, "hrs") || strings.Contains(line, "pcs") ||
			strings.Contains(line, "mm") || strings.Contains(line, " g ") {
			collectStats(line, raw)
		}

		for _, mat := range materialNamesFull {
			if strings.HasPrefix(line, mat) && (len(line) == len(mat) || line[len(mat)] == ' ') {
				raw.materials = append(raw.materials, mat)
				break
			}
		}

		for _, pr := range printerNames {
			if strings.Contains(line, pr+" ") || strings.Contains(line, " "+pr) ||
				strings.HasPrefix(line, pr) {
				raw.printers = append(raw.printers, pr)
				break
			}
		}

		if strings.HasPrefix(line, "hexagon") || strings.HasPrefix(line, "bed") ||
			strings.HasPrefix(line, "scraper") || strings.HasPrefix(line, "printbed") ||
			strings.HasPrefix(line, "hexagons") || strings.HasPrefix(line, "bedscraper") ||
			strings.HasPrefix(line, "spatula") || strings.HasPrefix(line, "hex") {
			for _, word := range strings.Fields(line) {
				if len(word) > 1 && !isPrivateUse(word) && word != "Tags:" && word != ">" && word != "|" {
					raw.tagWords = append(raw.tagWords, word)
				}
			}
		}
	}
}

func looksLikeDesignerName(line string) bool {
	if strings.EqualFold(line, "visle") {
		return true
	}
	if len(line) >= 20 || len(line) < 2 {
		return false
	}
	if isBoring(line) {
		return false
	}
	if strings.Contains(line, ">") {
		return false
	}
	for _, r := range line {
		if r >= '0' && r <= '9' {
			return false
		}
	}
	for _, kw := range []string{"hrs", "mm", "pcs", " g ", "Printers", "VIEW",
		"updated", "Attribution", "Creative", "Ironing", "This work", "Summary", "Tags",
		"hexagon", "bed", "scraper", "printbed", "hexagons", "bedscraper",
		"spatula", "Model", "Print", "License", "Commons"} {
		if strings.HasPrefix(line, kw) || strings.Contains(line, " "+kw) {
			return false
		}
	}
	return true
}

func isBoring(line string) bool {
	if len(line) == 0 {
		return true
	}
	if isPrivateUse(line) {
		return true
	}
	for _, r := range line {
		if r >= 0xF000 && r <= 0xFFFF {
			return true
		}
		if r >= 0x2700 && r <= 0x27BF {
			return true
		}
	}
	return false
}

func isPrivateUse(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := []rune(s)
	for _, c := range r {
		if c >= 0xE000 && c <= 0xF8FF {
			return true
		}
	}
	return false
}

func fillFromRows(lines []string, raw *rawExtract) {
	rowSet := make(map[string]bool)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "T ags:") || strings.HasPrefix(line, "Tags:") {
			tagPart := strings.TrimPrefix(line, "T ags: ")
			tagPart = strings.TrimPrefix(tagPart, "Tags: ")
			for _, w := range strings.Fields(tagPart) {
				if len(w) > 1 {
					raw.rowTags = append(raw.rowTags, w)
				}
			}
		}

		if !rowSet[line] && !strings.Contains(line, ">") && len(line) > 15 &&
			!strings.Contains(line, "hrs") && !strings.Contains(line, "mm") &&
			!strings.Contains(line, "pcs") && !strings.HasPrefix(line, "T ags") &&
			!strings.HasPrefix(line, "Tags") && !strings.HasPrefix(line, "updated") &&
			!strings.HasPrefix(line, "VIEW") && !strings.HasPrefix(line, "Matching") &&
			!strings.HasPrefix(line, "Summary") && !strings.HasPrefix(line, "Model") &&
			!strings.HasPrefix(line, "Print") && !strings.HasPrefix(line, "License") &&
			!strings.HasPrefix(line, "3D ") && !strings.HasPrefix(line, "f   k") &&
			!strings.HasPrefix(line, "PET") && !strings.HasSuffix(line, ".3mf") &&
			raw.description == "" {
			raw.description = line
		}
		rowSet[line] = true

		if strings.Contains(line, "hrs") || strings.Contains(line, "mm") || strings.Contains(line, " g ") {
			collectStats(line, raw)
		}
	}
}

func ExtractThumbnail(pdfData []byte) ([]byte, error) {
	var jpegChunks [][]byte

	for i := 0; i < len(pdfData)-3; i++ {
		if pdfData[i] == 0xFF && pdfData[i+1] == 0xD8 && pdfData[i+2] == 0xFF {
			if pdfData[i+3] >= 0xE0 && pdfData[i+3] <= 0xEF || pdfData[i+3] == 0xFE || pdfData[i+3] == 0xDB || pdfData[i+3] == 0xC4 || pdfData[i+3] == 0xC0 {
				for j := i + 2; j < len(pdfData)-1; j++ {
					if pdfData[j] == 0xFF && pdfData[j+1] == 0xD9 {
						jpegChunks = append(jpegChunks, append([]byte(nil), pdfData[i:j+2]...))
						i = j + 1
						break
					}
				}
			}
		}
	}

	if len(jpegChunks) == 0 {
		return nil, fmt.Errorf("no thumbnail found")
	}

	var best []byte
	for _, chunk := range jpegChunks {
		if len(chunk) > len(best) {
			best = chunk
		}
	}
	return best, nil
}

func collectURIs(r *pdf.Reader, raw *rawExtract) {
	trailer := r.Trailer()
	if trailer.Kind() != pdf.Dict {
		return
	}
	root := trailer.Key("Root")
	if root.Kind() != pdf.Dict {
		return
	}
	pages := root.Key("Pages")
	if pages.Kind() != pdf.Dict {
		return
	}
	kids := pages.Key("Kids")
	for i := 0; i < kids.Len(); i++ {
		page := kids.Index(i)
		if page.Kind() != pdf.Dict {
			continue
		}
		annots := page.Key("Annots")
		if annots.Kind() != pdf.Array {
			continue
		}
		for j := 0; j < annots.Len(); j++ {
			annot := annots.Index(j)
			if annot.Kind() != pdf.Dict {
				continue
			}
			if annot.Key("Subtype").Name() != "Link" {
				continue
			}
			a := annot.Key("A")
			if a.Kind() != pdf.Dict {
				continue
			}
			uri := a.Key("URI").Text()
			if uri == "" {
				uri = a.Key("URI").RawString()
			}
			if uri == "" {
				continue
			}

			switch {
			case raw.sourceURL == "" && strings.Contains(uri, "printables.com/model/"):
				raw.sourceURL = uri
			case raw.designerURL == "" && strings.Contains(uri, "printables.com/social/"):
				raw.designerURL = uri
			case raw.licenseURL == "" && strings.Contains(uri, "creativecommons.org/licenses/"):
				raw.licenseURL = uri
			}
		}
	}
}

var timeRE = regexp.MustCompile(`(\d+\.?\d*)\s*hrs?`)
var pcsRE = regexp.MustCompile(`(\d+)\s*pcs?`)
var distRE = regexp.MustCompile(`(\d+\.\d+)\s*mm`)
var weightRE = regexp.MustCompile(`(\d+\.?\d*)\s*g\b`)

var materialNamesFull = []string{
	"PLA", "PLA+", "PETG", "PET", "ABS", "ASA", "TPU", "TPE",
	"PC", "PA", "Nylon", "PVA", "HIPS", "PP", "POM", "PEI",
	"PEEK", "Carbon Fiber", "Wood", "Metal",
}

var printerNames = []string{
	"Prusa", "Bambu Lab", "Ender", "Creality", "Anycubic",
	"Voron", "Rat Rig", "Elegoo", "Snapmaker", "Ultimaker",
}

func collectStats(line string, raw *rawExtract) {
	if raw.printTime == "" {
		if m := timeRE.FindStringSubmatch(line); m != nil {
			raw.printTime = m[1]
		}
	}
	if raw.quantity == "" {
		if m := pcsRE.FindStringSubmatch(line); m != nil {
			raw.quantity = m[1]
		}
	}
	if raw.weight == "" {
		if m := weightRE.FindStringSubmatch(line); m != nil {
			raw.weight = m[1]
		}
	}

	matches := distRE.FindAllStringSubmatch(line, -1)
	for _, m := range matches {
		val := m[1]
		if raw.layerHeight == "" {
			v, _ := strconv.ParseFloat(val, 64)
			if v < 1.0 && v >= 0.01 {
				raw.layerHeight = val
				continue
			}
		}
		if raw.nozzle == "" {
			raw.nozzle = val
		}
	}

	upper := strings.ToUpper(line)
	for _, mat := range materialNamesFull {
		matUpper := strings.ToUpper(mat)
		if strings.HasPrefix(upper, matUpper) {
			raw.materials = append(raw.materials, mat)
			break
		}
	}
	for _, pr := range printerNames {
		if strings.Contains(line, pr+" ") || strings.Contains(line, " "+pr) ||
			strings.HasPrefix(line, pr) || strings.HasSuffix(line, pr) {
			raw.printers = append(raw.printers, pr)
			break
		}
	}
}

func buildResult(raw *rawExtract) *Result {
	r := &Result{
		Title:       raw.title,
		Designer:    raw.designer,
		SourceURL:   raw.sourceURL,
		DesignerURL: raw.designerURL,
		LicenseURL:  raw.licenseURL,
		Description: raw.description,
		Category:    raw.category,
		LicenseName: licenseNameFromURL(raw.licenseURL),
	}

	allTags := append(raw.rowTags, raw.tagWords...)
	allTags = dedup(allTags)
	r.Tags = allTags

	if v, err := strconv.ParseFloat(raw.printTime, 64); err == nil {
		r.PrintTimeMinutes = int(math.Round(v * 60))
	}
	if v, err := strconv.ParseFloat(raw.weight, 64); err == nil {
		r.WeightGrams = v
	}
	if v, err := strconv.Atoi(raw.quantity); err == nil {
		r.Quantity = v
	}
	if v, err := strconv.ParseFloat(raw.nozzle, 64); err == nil {
		r.NozzleMM = v
	}
	if v, err := strconv.ParseFloat(raw.layerHeight, 64); err == nil {
		r.LayerHeightMM = v
	}

	if len(raw.materials) > 0 {
		r.Material = strings.Join(dedup(raw.materials), " / ")
	}
	if len(raw.printers) > 0 {
		r.Printer = strings.Join(dedup(raw.printers), " / ")
	}

	if raw.title == "" && r.SourceURL != "" {
		r.Title = "Imported from Printables"
	}

	if strings.Contains(r.Title, "-") && !strings.Contains(r.Title, " - ") {
		r.Title = strings.Replace(r.Title, "-", " - ", 1)
	}
	r.Title = strings.TrimSpace(r.Title)

	return r
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

var ccLicenses = map[string]string{
	"by":       "CC BY 4.0",
	"by-sa":    "CC BY-SA 4.0",
	"by-nc":    "CC BY-NC 4.0",
	"by-nc-sa": "CC BY-NC-SA 4.0",
	"by-nd":    "CC BY-ND 4.0",
	"by-nc-nd": "CC BY-NC-ND 4.0",
	"zero":     "CC0",
}

func licenseNameFromURL(url string) string {
	if url == "" {
		return ""
	}
	for code, name := range ccLicenses {
		if strings.Contains(url, "/licenses/"+code+"/") {
			return name
		}
	}
	if strings.Contains(url, "creativecommons.org") {
		return "Creative Commons"
	}
	return ""
}

func dedup(items []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, item := range items {
		key := strings.ToLower(item)
		if !seen[key] {
			seen[key] = true
			out = append(out, item)
		}
	}
	sort.Strings(out)
	return out
}

type rawExtract struct {
	title        string
	designer     string
	designerSet  bool
	sourceURL    string
	designerURL  string
	licenseURL   string
	description  string
	category     string
	printTime    string
	weight       string
	quantity     string
	nozzle       string
	layerHeight  string
	materials    []string
	printers     []string
	rowTags      []string
	tagWords     []string
}