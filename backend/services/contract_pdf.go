package services

import (
	"bytes"
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-pdf/fpdf"
)

type ContractData struct {
	Number          string
	SignedAt        time.Time
	CustomerName    string
	CustomerEmail   string
	CustomerAddress string
	CustomerPhone   string
	PrestationTitle string
	PrestationDesc  string
	ProviderName    string
	ProviderEmail   string
	AmountCents     int64
	Currency        string
	SignedIP        string
	SignaturePNG    []byte
}

func (s *PDFService) GenerateContract(data ContractData) (string, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "B", 22)
	pdf.SetTextColor(0, 29, 50)
	pdf.Cell(0, 10, tr("CONTRAT DE PRESTATION"))
	pdf.Ln(11)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 116, 139)
	pdf.Cell(0, 5, tr("UpcycleConnect — La plateforme de l'upcycling intelligent"))
	pdf.Ln(8)

	pdf.SetDrawColor(0, 109, 53)
	pdf.SetLineWidth(0.8)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(6)

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(30, 41, 59)
	pdf.CellFormat(45, 6, tr("Numéro de contrat :"), "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, tr(data.Number), "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(45, 6, tr("Date de signature :"), "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, tr(data.SignedAt.Format("02/01/2006 à 15:04")), "", 1, "L", false, 0, "")
	pdf.Ln(6)

	sectionTitle(pdf, tr, "Article 1 — Entre les soussignés")

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 5, tr("La société :"))
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"UpcycleConnect SAS\n"+
			"174, rue La Fayette, 75010 Paris\n"+
			"Email : contact@upcycleconnect.xyz\n"+
			"Site : https://upcycleconnect.xyz\n"+
			"Ci-après désignée « le Prestataire »"), "", "L", false)
	pdf.Ln(2)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 5, tr("Et :"))
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 9)
	clientLines := data.CustomerName + "\n" + data.CustomerEmail
	if data.CustomerPhone != "" {
		clientLines += "\nTéléphone : " + data.CustomerPhone
	}
	if data.CustomerAddress != "" {
		clientLines += "\nAdresse : " + data.CustomerAddress
	}
	clientLines += "\nCi-après désigné(e) « le Client »"
	pdf.MultiCell(0, 4.5, tr(clientLines), "", "L", false)
	pdf.Ln(4)

	sectionTitle(pdf, tr, "Article 2 — Objet du contrat")

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"Le présent contrat a pour objet la prestation de service suivante, "+
			"commandée par le Client auprès du Prestataire via la plateforme UpcycleConnect :"), "", "J", false)
	pdf.Ln(3)

	desc := data.PrestationDesc
	if desc == "" {
		desc = "(Voir descriptif détaillé sur la fiche prestation en ligne.)"
	}
	startY := pdf.GetY()
	pdf.SetX(25)
	pdf.SetFont("Arial", "B", 10)
	pdf.MultiCell(163, 5.5, tr(data.PrestationTitle), "", "L", false)
	pdf.SetX(25)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(163, 4.5, tr(desc), "", "J", false)
	endY := pdf.GetY()
	pdf.SetFillColor(0, 109, 53)
	pdf.Rect(20, startY, 1.2, endY-startY-1, "F")
	pdf.Ln(3)

	if data.ProviderName != "" {
		pdf.SetFont("Arial", "", 9)
		pdf.Cell(0, 4.5, tr("Prestation fournie par : "+data.ProviderName))
		pdf.Ln(5)
	}
	pdf.Ln(2)

	sectionTitle(pdf, tr, "Article 3 — Prix et modalités de paiement")

	pdf.SetFont("Arial", "", 9)
	pdf.Cell(0, 5, tr("Montant total TTC :"))
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 109, 53)
	pdf.Cell(0, 8, tr(formatEUR(data.AmountCents)))
	pdf.SetTextColor(30, 41, 59)
	pdf.Ln(9)

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"Le paiement est effectué en une seule fois, en ligne, via le prestataire de paiement Stripe au "+
			"moment de la signature du présent contrat. La TVA française au taux normal de 20 % est incluse "+
			"dans le montant total. Une facture définitive est émise après réception du paiement et "+
			"transmise au Client par email ainsi que dans son espace personnel."), "", "J", false)
	pdf.Ln(3)

	sectionTitle(pdf, tr, "Article 4 — Exécution de la prestation")
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"Le Prestataire s'engage à exécuter la prestation avec soin et dans le respect des règles de l'art. "+
			"Les modalités précises de réalisation (date, lieu, durée) sont convenues entre les parties via "+
			"la messagerie et l'agenda intégrés à la plateforme UpcycleConnect. Le Client s'engage à fournir "+
			"au Prestataire toutes les informations nécessaires à la bonne exécution de la prestation."), "", "J", false)
	pdf.Ln(3)

	sectionTitle(pdf, tr, "Article 5 — Droit de rétractation")
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"Conformément aux articles L.221-18 et suivants du Code de la consommation, le Client dispose d'un "+
			"délai de quatorze (14) jours à compter de la conclusion du présent contrat pour exercer son droit "+
			"de rétractation, sans avoir à motiver sa décision. Toutefois, si l'exécution de la prestation "+
			"commence, à la demande expresse du Client, avant l'expiration de ce délai, le Client renonce "+
			"expressément à son droit de rétractation pour la part de prestation déjà exécutée."), "", "J", false)
	pdf.Ln(3)

	sectionTitle(pdf, tr, "Article 6 — Protection des données personnelles")
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"Les données personnelles du Client sont traitées conformément au Règlement (UE) 2016/679 (RGPD) "+
			"et à la loi Informatique et Libertés. Elles sont collectées pour l'exécution du présent contrat "+
			"et conservées pendant la durée légale de conservation des documents contractuels. Le Client "+
			"dispose d'un droit d'accès, de rectification, d'effacement et de portabilité de ses données, "+
			"exerçable par email à dpo@upcycleconnect.xyz."), "", "J", false)
	pdf.Ln(3)

	sectionTitle(pdf, tr, "Article 7 — Loi applicable et règlement des litiges")
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"Le présent contrat est soumis au droit français. En cas de litige, les parties s'engagent à "+
			"rechercher en priorité une solution amiable. À défaut, le litige sera soumis aux tribunaux "+
			"compétents du ressort de Paris. Le Client peut également recourir gratuitement au service de "+
			"médiation MEDICYS (www.medicys.fr)."), "", "J", false)
	pdf.Ln(6)

	if pdf.GetY() > 220 {
		pdf.AddPage()
	}

	sectionTitle(pdf, tr, "Signatures")

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 4.5, tr(
		"Les soussignés reconnaissent avoir pris connaissance de l'intégralité des conditions exposées "+
			"ci-dessus et les acceptent sans réserve. Le présent contrat est conclu en deux exemplaires, "+
			"un pour chaque partie, et fait foi entre elles."), "", "J", false)
	pdf.Ln(5)

	sigY := pdf.GetY()

	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(20, sigY)
	pdf.Cell(80, 6, tr("Pour le Prestataire :"))
	pdf.SetXY(20, sigY+6)
	pdf.SetFont("Arial", "", 9)
	pdf.Cell(80, 5, tr("UpcycleConnect SAS"))
	pdf.SetXY(20, sigY+11)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.Cell(80, 5, tr("(Signature électronique de la plateforme)"))
	pdf.SetTextColor(30, 41, 59)

	pdf.SetDrawColor(0, 109, 53)
	pdf.SetLineWidth(0.4)
	pdf.Rect(20, sigY+18, 80, 24, "D")
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 109, 53)
	pdf.SetXY(20, sigY+26)
	pdf.CellFormat(80, 8, tr("✓ UpcycleConnect"), "", 0, "C", false, 0, "")
	pdf.SetTextColor(30, 41, 59)

	pdf.SetFont("Arial", "B", 10)
	pdf.SetXY(110, sigY)
	pdf.Cell(80, 6, tr("Pour le Client :"))
	pdf.SetXY(110, sigY+6)
	pdf.SetFont("Arial", "", 9)
	pdf.Cell(80, 5, tr(data.CustomerName))
	pdf.SetXY(110, sigY+11)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.Cell(80, 5, tr("(Signature manuscrite ci-dessous)"))
	pdf.SetTextColor(30, 41, 59)
	pdf.SetDrawColor(200, 210, 220)
	pdf.SetLineWidth(0.3)
	pdf.Rect(110, sigY+18, 80, 24, "D")
	if len(data.SignaturePNG) > 0 {
		opt := fpdf.ImageOptions{ImageType: "PNG"}
		pdf.RegisterImageOptionsReader("client_sig", opt, bytes.NewReader(data.SignaturePNG))
		pdf.ImageOptions("client_sig", 112, sigY+19, 76, 22, false, opt, 0, "")
	}

	pdf.SetY(sigY + 46)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(100, 116, 139)
	meta := fmt.Sprintf("Signé le %s", data.SignedAt.Format("02/01/2006 à 15:04:05"))
	if data.SignedIP != "" {
		meta += fmt.Sprintf(" — IP : %s", data.SignedIP)
	}
	pdf.MultiCell(0, 4, tr(meta), "", "L", false)
	pdf.Ln(2)
	pdf.MultiCell(0, 4, tr(
		"Conformément au règlement eIDAS (UE) n° 910/2014 et aux articles 1366 et 1367 du Code civil, "+
			"la signature électronique apposée ci-dessus a la même valeur juridique qu'une signature manuscrite."), "", "L", false)

	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 7)
	pdf.SetTextColor(160, 174, 192)
	pdf.CellFormat(0, 5,
		tr("Document généré le "+time.Now().Format("02/01/2006 à 15:04")+" — UpcycleConnect — "+data.Number),
		"", 0, "C", false, 0, "")

	filename := fmt.Sprintf("%s.pdf", data.Number)
	path := filepath.Join(s.baseDir, filename)
	if err := pdf.OutputFileAndClose(path); err != nil {
		return "", fmt.Errorf("generate contract pdf: %w", err)
	}
	return path, nil
}

func sectionTitle(pdf *fpdf.Fpdf, tr func(string) string, title string) {
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(0, 109, 53)
	pdf.Cell(0, 6, tr(title))
	pdf.Ln(7)
	pdf.SetTextColor(30, 41, 59)
}
