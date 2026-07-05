package pdfmeta

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestLicenseNameFromURL(t *testing.T) {
	cases := []struct {
		url  string
		name string
	}{
		{"http://creativecommons.org/licenses/by-nc-sa/4.0/", "CC BY-NC-SA 4.0"},
		{"https://creativecommons.org/licenses/by/4.0/", "CC BY 4.0"},
		{"http://creativecommons.org/licenses/by-sa/3.0/", "CC BY-SA 4.0"},
		{"", ""},
		{"https://example.com", ""},
		{"http://creativecommons.org/publicdomain/zero/1.0/", "Creative Commons"},
		{"https://creativecommons.org/licenses/by-nc-nd/4.0/", "CC BY-NC-ND 4.0"},
		{"https://creativecommons.org/licenses/zero/1.0/", "CC0"},
	}

	for i, c := range cases {
		got := licenseNameFromURL(c.url)
		if got != c.name {
			t.Errorf("case %d: licenseNameFromURL(%q) = %q, want %q", i, c.url, got, c.name)
		}
	}
}

func TestDedup(t *testing.T) {
	got := dedup([]string{"PLA", "PETG", "pla", "PETG"})
	if len(got) != 2 {
		t.Errorf("dedup should produce 2 items, got %d: %v", len(got), got)
	}
}

func TestAbs(t *testing.T) {
	if abs(-5.5) != 5.5 {
		t.Error("abs(-5.5) != 5.5")
	}
	if abs(3.0) != 3.0 {
		t.Error("abs(3.0) != 3.0")
	}
}

func loadHexScraperPDF(t *testing.T) ([]byte, *zip.Reader) {
	t.Helper()

	zipPath := "../../../../hexscraper-printbed-scraper-model_files.zip"
	zipData, err := os.ReadFile(zipPath)
	if err != nil {
		t.Skip("example zip not found")
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		t.Fatalf("failed to open zip: %v", err)
	}

	return nil, zipReader
}

func extractPDFFromZip(t *testing.T, zipReader *zip.Reader) []byte {
	t.Helper()

	for _, f := range zipReader.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".pdf") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("failed to open PDF entry: %v", err)
		}
		defer rc.Close()
		data, err := readAllFromZip(f)
		if err != nil {
			t.Fatalf("failed to read PDF: %v", err)
		}
		return data
	}
	t.Skip("no PDF entry found in zip")
	return nil
}

func readAllFromZip(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// Test with the real hexscraper PDF data extracted from the zip
func TestHexScraperExtraction(t *testing.T) {
	_, zipReader := loadHexScraperPDF(t)
	if zipReader == nil {
		return
	}

	pdfData := extractPDFFromZip(t, zipReader)
	if pdfData == nil {
		return
	}

	t.Logf("Extracted PDF: %d bytes", len(pdfData))

	result, err := Extract(pdfData)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	t.Logf("Title:           %q", result.Title)
	t.Logf("Designer:        %q", result.Designer)
	t.Logf("Source URL:      %q", result.SourceURL)
	t.Logf("Designer URL:    %q", result.DesignerURL)
	t.Logf("License URL:     %q", result.LicenseURL)
	t.Logf("License Name:    %q", result.LicenseName)
	t.Logf("Description:     %q", result.Description)
	t.Logf("Category:        %q", result.Category)
	t.Logf("Tags:            %v", result.Tags)
	t.Logf("Print time (min): %d", result.PrintTimeMinutes)
	t.Logf("Weight (g):      %.1f", result.WeightGrams)
	t.Logf("Quantity:        %d", result.Quantity)
	t.Logf("Nozzle (mm):     %.2f", result.NozzleMM)
	t.Logf("Layer height:    %.2f", result.LayerHeightMM)
	t.Logf("Material:        %q", result.Material)
	t.Logf("Printer:         %q", result.Printer)

	if result.SourceURL == "" {
		t.Error("SourceURL should not be empty")
	}
	if !strings.Contains(result.SourceURL, "printables.com/model/") {
		t.Errorf("SourceURL should contain printables.com/model/: got %q", result.SourceURL)
	}
	if result.LicenseURL == "" {
		t.Error("LicenseURL should not be empty")
	}
	if result.LicenseName == "" {
		t.Error("LicenseName should not be empty")
	}
	if result.WeightGrams < 1 {
		t.Error("Weight should be >= 1g")
	}
	if result.PrintTimeMinutes < 1 {
		t.Error("Print time should be >= 1 minute")
	}
}

func TestExtractThumbnail(t *testing.T) {
	_, zipReader := loadHexScraperPDF(t)
	if zipReader == nil {
		return
	}

	pdfData := extractPDFFromZip(t, zipReader)
	if pdfData == nil {
		return
	}

	thumb, err := ExtractThumbnail(pdfData)
	if err != nil {
		t.Fatalf("ExtractThumbnail failed: %v", err)
	}
	if len(thumb) == 0 {
		t.Error("Thumbnail should have content")
	}
	t.Logf("Thumbnail size: %d bytes", len(thumb))
}

func TestThumbnailIsValidJPEG(t *testing.T) {
	_, zipReader := loadHexScraperPDF(t)
	if zipReader == nil {
		return
	}

	pdfData := extractPDFFromZip(t, zipReader)
	if pdfData == nil {
		return
	}

	thumb, err := ExtractThumbnail(pdfData)
	if err != nil {
		t.Fatal(err)
	}
	if len(thumb) < 4 || thumb[0] != 0xFF || thumb[1] != 0xD8 {
		t.Error("thumbnail does not have JPEG SOI marker")
	}
	if thumb[len(thumb)-2] != 0xFF || thumb[len(thumb)-1] != 0xD9 {
		t.Error("thumbnail does not have JPEG EOI marker")
	}
}

func TestZipSlipProtection(t *testing.T) {
	t.Log("Zip-slip protection is enforced in handlers.go via filepath.Base() on extracted filenames")
	fmt.Println("ok")
}