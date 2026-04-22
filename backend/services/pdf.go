package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-pdf/fpdf"
)

type PDFService struct {
	baseDir string
}

func NewPDFService(baseDir string) *PDFService {
	if baseDir == "" {
		baseDir = "/var/lib/upcycleconnect/invoices"
	}
	_ = os.MkdirAll(baseDir, 0o755)
	return &PDFService{baseDir: baseDir}
}

type InvoiceData struct {
	Number          string
	Type            string
	IssuedAt        time.Time
	CustomerName    string
	CustomerEmail   string
	PrestationTitle string
	AmountCents     int64
	TVAPercent      float64
	Currency        string
	Notes           string
}

func formatEUR(cents int64) string {
	euros := float64(cents) / 100.0
	return fmt.Sprintf("%.2f EUR", euros)
}

func (s *PDFService) Generate(data InvoiceData) (string, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	title := "FACTURE"
	if data.Type == "quote" {
		title = "DEVIS"
	}

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, title)
	pdf.Ln(14)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 7, "UpcycleConnect")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, "contact@upcycleconnect.xyz")
	pdf.Ln(5)
	pdf.Cell(0, 5, "https://upcycleconnect.xyz")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 7, "Numero :")
	pdf.Cell(0, 7, data.Number)
	pdf.Ln(6)
	pdf.Cell(40, 7, "Date :")
	pdf.Cell(0, 7, data.IssuedAt.Format("02/01/2006"))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 7, "Client")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, data.CustomerName)
	pdf.Ln(5)
	pdf.Cell(0, 5, data.CustomerEmail)
	pdf.Ln(12)

	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(110, 8, "Prestation", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Qte", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Montant", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(110, 8, data.PrestationTitle, "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 8, "1", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, formatEUR(data.AmountCents), "1", 0, "R", false, 0, "")
	pdf.Ln(12)

	tva := data.TVAPercent
	if tva <= 0 {
		tva = 20.0
	}
	amountHT := int64(float64(data.AmountCents) / (1 + tva/100))
	amountTVA := data.AmountCents - amountHT

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(140, 7, "Total HT", "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, formatEUR(amountHT), "0", 0, "R", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(140, 7, fmt.Sprintf("TVA (%.2f%%)", tva), "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, formatEUR(amountTVA), "0", 0, "R", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(140, 8, "Total TTC", "T", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, formatEUR(data.AmountCents), "T", 0, "R", false, 0, "")
	pdf.Ln(18)

	if data.Notes != "" {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, data.Notes, "", "L", false)
		pdf.Ln(4)
	}

	if data.Type == "quote" {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, "Ce devis est valable 30 jours a compter de sa date d'emission.", "", "L", false)
	} else {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, "Paiement recu - Merci de votre confiance.", "", "L", false)
	}

	filename := fmt.Sprintf("%s.pdf", data.Number)
	path := filepath.Join(s.baseDir, filename)
	if err := pdf.OutputFileAndClose(path); err != nil {
		return "", fmt.Errorf("generate pdf: %w", err)
	}
	return path, nil
}
