package services

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/go-pdf/fpdf"
)

type PDFService struct {
	baseDir string
}

func NewPDFService(baseDir string) (*PDFService, error) {
	if baseDir == "" {
		baseDir = "/var/lib/upcycleconnect/invoices"
	}
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("pdf service: create base dir %q: %w", baseDir, err)
	}
	return &PDFService{baseDir: baseDir}, nil
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
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	title := "FACTURE"
	if data.Type == "quote" {
		title = "DEVIS"
	}

	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(0, 10, tr(title))
	pdf.Ln(14)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 7, tr("UpcycleConnect"))
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, tr("contact@upcycleconnect.xyz"))
	pdf.Ln(5)
	pdf.Cell(0, 5, tr("https://upcycleconnect.xyz"))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(40, 7, tr("Numéro :"))
	pdf.Cell(0, 7, tr(data.Number))
	pdf.Ln(6)
	pdf.Cell(40, 7, tr("Date :"))
	pdf.Cell(0, 7, tr(data.IssuedAt.Format("02/01/2006")))
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 7, tr("Client"))
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, tr(data.CustomerName))
	pdf.Ln(5)
	pdf.Cell(0, 5, tr(data.CustomerEmail))
	pdf.Ln(12)

	pdf.SetFillColor(240, 240, 240)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(110, 8, tr("Prestation"), "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, tr("Qté"), "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, tr("Montant"), "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(110, 8, tr(data.PrestationTitle), "1", 0, "L", false, 0, "")
	pdf.CellFormat(30, 8, tr("1"), "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, tr(formatEUR(data.AmountCents)), "1", 0, "R", false, 0, "")
	pdf.Ln(12)

	tva := data.TVAPercent
	if tva <= 0 {
		tva = 20.0
	}
	amountHT := int64(math.Round(float64(data.AmountCents) / (1 + tva/100)))
	amountTVA := data.AmountCents - amountHT

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(140, 7, tr("Total HT"), "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, tr(formatEUR(amountHT)), "0", 0, "R", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(140, 7, tr(fmt.Sprintf("TVA (%.2f%%)", tva)), "0", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, tr(formatEUR(amountTVA)), "0", 0, "R", false, 0, "")
	pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(140, 8, tr("Total TTC"), "T", 0, "R", false, 0, "")
	pdf.CellFormat(30, 8, tr(formatEUR(data.AmountCents)), "T", 0, "R", false, 0, "")
	pdf.Ln(18)

	if data.Notes != "" {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, tr(data.Notes), "", "L", false)
		pdf.Ln(4)
	}

	if data.Type == "quote" {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, tr("Ce devis est valable 30 jours à compter de sa date d'émission."), "", "L", false)
	} else {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, tr("Paiement reçu — Merci de votre confiance."), "", "L", false)
	}

	filename := fmt.Sprintf("%s.pdf", data.Number)
	path := filepath.Join(s.baseDir, filename)
	if err := pdf.OutputFileAndClose(path); err != nil {
		return "", fmt.Errorf("generate pdf: %w", err)
	}
	return path, nil
}
