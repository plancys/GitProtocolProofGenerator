package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"git-report-generator/internal/config"
	"git-report-generator/internal/git"

	"github.com/jung-kurt/gofpdf"
)

const (
	fontDir  = "fonts"
	fontName = "DejaVu"
	fontFile = "DejaVuSans.ttf"
	fontBold = "DejaVuSans-Bold.ttf"
)

// ReportData contains all the data needed to generate a report
type ReportData struct {
	Config         *config.Config
	RepositoryName string
	RepositoryPath string // Absolute path to the repository
	BranchName     string
	AuthorEmail    string
	DateFrom       time.Time
	DateTo         time.Time
	Commits        []*git.Commit
}

// PDFGenerator handles PDF report generation
type PDFGenerator struct {
	pdf *gofpdf.Fpdf
}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (g *PDFGenerator) getAbsFontPath(fontFileName string) (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	fmt.Printf("Executable path: %s\n", executablePath)
	executableDir := filepath.Dir(executablePath)
	fmt.Printf("Executable dir: %s\n", executableDir)
	// Construct path relative to executable dir first
	relPath := filepath.Join(fontDir, fontFileName)
	// Then join with executable dir and ensure it's absolute and clean
	absPath, err := filepath.Abs(filepath.Join(executableDir, relPath))
	if err != nil {
		return "", fmt.Errorf("failed to get absolute font path for %s: %w", relPath, err)
	}
	fmt.Printf("Constructed absolute font path: %s\n", absPath)
	return absPath, nil
}

// Generate creates a PDF report based on the provided data
func (g *PDFGenerator) Generate(data *ReportData, outputPath string) error {
	executablePath, _ := os.Executable()
	executableDir := filepath.Dir(executablePath)
	absFontDir := filepath.Join(executableDir, fontDir)
	fmt.Printf("Absolute font directory for gofpdf.New: %s\n", absFontDir)

	regularFontPath, err := g.getAbsFontPath(fontFile)
	if err != nil {
		return fmt.Errorf("failed to get font path: %w", err)
	}
	boldFontPath, err := g.getAbsFontPath(fontBold)
	if err != nil {
		return fmt.Errorf("failed to get bold font path: %w", err)
	}

	g.pdf = gofpdf.New("P", "mm", "A4", absFontDir)
	fmt.Printf("Passing to AddUTF8Font (Regular): '%s'\n", regularFontPath)
	g.pdf.AddUTF8Font(fontName, "", fontFile)
	fmt.Printf("Passing to AddUTF8Font (Bold): '%s'\n", boldFontPath)
	g.pdf.AddUTF8Font(fontName, "B", fontBold)
	// Register italic font
	italicFontFile := "DejaVuSans-Oblique.ttf"
	italicFontPath, err := g.getAbsFontPath(italicFontFile)
	if err != nil {
		return fmt.Errorf("failed to get italic font path: %w", err)
	}
	fmt.Printf("Passing to AddUTF8Font (Italic): '%s'\n", italicFontPath)
	g.pdf.AddUTF8Font(fontName, "I", italicFontFile)
	g.pdf.SetFont(fontName, "", 11)
	g.pdf.AddPage()
	g.pdf.SetMargins(20, 20, 20)
	g.pdf.SetAutoPageBreak(true, 20)

	if err := g.generateHeader(data); err != nil {
		return err
	}
	if err := g.generateCommits(data); err != nil {
		return err
	}
	if err := g.pdf.OutputFileAndClose(outputPath); err != nil {
		return fmt.Errorf("failed to save PDF: %w", err)
	}
	return nil
}

// generateHeader creates the header section of the PDF
func (g *PDFGenerator) generateHeader(data *ReportData) error {
	// Parse template to extract components
	lines := strings.Split(data.Config.Header.Template, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("invalid template format: not enough lines")
	}

	// Extract key parts
	dateLocLine := lines[0] // "Kraków, {{.date_from}} - {{.date_to}}"
	titleLine := lines[1]   // "Protokół odbioru prac programistycznych"

	// Rest of the template (excluding first two lines)
	restOfTemplate := strings.Join(lines[2:], "\n")

	// Prepare template data
	templateData := map[string]interface{}{
		"executor_name":   data.Config.Header.ExecutorName,
		"executor_email":  data.Config.Header.ExecutorEmail,
		"recipient_name":  data.Config.Header.RecipientName,
		"repository_name": data.RepositoryName,
		"repository_path": data.RepositoryPath,
		"branch_name":     data.BranchName,
		"date_from":       data.DateFrom.Format("2006-01-02"),
		"date_to":         data.DateTo.Format("2006-01-02"),
	}

	// Debug: Print template parts
	fmt.Println("DEBUG: Date line:", dateLocLine)
	fmt.Println("DEBUG: Title line:", titleLine)
	fmt.Println("DEBUG: Rest of template:", restOfTemplate)

	// Render date line
	dateTmpl, err := template.New("date").Option("missingkey=zero").Parse(dateLocLine)
	if err != nil {
		return fmt.Errorf("failed to parse date template: %w", err)
	}
	var dateBuf bytes.Buffer
	if err := dateTmpl.Execute(&dateBuf, templateData); err != nil {
		return fmt.Errorf("failed to execute date template: %w", err)
	}
	dateText := dateBuf.String()

	// Render rest of template
	restTmpl, err := template.New("rest").Option("missingkey=zero").Parse(restOfTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse rest of template: %w", err)
	}
	var restBuf bytes.Buffer
	if err := restTmpl.Execute(&restBuf, templateData); err != nil {
		return fmt.Errorf("failed to execute rest of template: %w", err)
	}
	restText := restBuf.String()

	// Debug: Print rendered components
	fmt.Println("DEBUG: Rendered date:", dateText)
	fmt.Println("DEBUG: Title (not rendered):", titleLine) // Title doesn't need rendering, no variables
	fmt.Println("DEBUG: Rendered rest:", restText)

	// Generate PDF layout
	// 1. Date line at normal size
	g.pdf.SetFont(fontName, "", 11)
	g.pdf.Cell(0, 10, dateText)
	g.pdf.Ln(12)

	// 2. Title line larger, bold, and centered
	g.pdf.SetFont(fontName, "B", 16)
	titleWidth := g.pdf.GetStringWidth(titleLine)
	pageWidth, _ := g.pdf.GetPageSize()
	leftMargin, _, rightMargin, _ := g.pdf.GetMargins()
	availableWidth := pageWidth - leftMargin - rightMargin
	xPos := (availableWidth - titleWidth) / 2
	g.pdf.Cell(0, 8, titleLine)
	g.pdf.Ln(15)

	// 3. Rest of the template
	g.pdf.SetFont(fontName, "", 11)
	g.pdf.MultiCell(0, 7, restText, "", "L", false)
	g.pdf.Ln(5)
	return nil
}

// generateCommits creates the commits section of the PDF
func (g *PDFGenerator) generateCommits(data *ReportData) error {
	if len(data.Commits) == 0 {
		g.pdf.SetFont(fontName, "I", 11)
		g.pdf.Cell(0, 6, "Brak commitów w podanym okresie.")
		return nil
	}

	// Table header
	g.pdf.SetFont(fontName, "B", 10)
	g.pdf.SetFillColor(220, 220, 220)
	g.pdf.CellFormat(30, 8, "Data", "1", 0, "C", true, 0, "")
	g.pdf.CellFormat(25, 8, "SHA", "1", 0, "C", true, 0, "")
	g.pdf.CellFormat(0, 8, "Opis", "1", 1, "C", true, 0, "")

	g.pdf.SetFont(fontName, "", 10)
	for i, commit := range data.Commits {
		if i%2 == 1 {
			g.pdf.SetFillColor(245, 245, 245)
		} else {
			g.pdf.SetFillColor(255, 255, 255)
		}
		g.pdf.CellFormat(30, 7, commit.Date.Format("2006-01-02"), "1", 0, "C", true, 0, "")
		g.pdf.CellFormat(25, 7, commit.SHA, "1", 0, "C", true, 0, "")
		g.pdf.MultiCell(0, 7, fmt.Sprintf("%s\n%s", commit.Message, commit.Description), "1", "L", false)
	}

	g.pdf.Ln(8)
	g.pdf.SetFont(fontName, "B", 11)
	g.pdf.Cell(0, 8, "Podsumowanie:")
	g.pdf.Ln(8)
	g.pdf.SetFont(fontName, "", 10)
	g.pdf.Cell(0, 6, fmt.Sprintf("Łączna liczba commitów: %d", len(data.Commits)))
	g.pdf.Ln(6)
	g.pdf.Cell(0, 6, fmt.Sprintf("Autor: %s", data.AuthorEmail))
	g.pdf.Ln(6)
	g.pdf.Cell(0, 6, fmt.Sprintf("Okres: %s - %s", data.DateFrom.Format("2006-01-02"), data.DateTo.Format("2006-01-02")))
	g.pdf.Ln(10)
	g.pdf.SetFont(fontName, "I", 8)
	g.pdf.SetTextColor(120, 120, 120)
	g.pdf.Cell(0, 4, fmt.Sprintf("Raport wygenerowany: %s", time.Now().Format("2006-01-02 15:04:05")))
	return nil
}
