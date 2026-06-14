package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
)

func (h *UserHandler) Export(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var user models.User
	if err := h.DB.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ressource introuvable"})
		return
	}

	var deposits []models.DepositRequest
	h.DB.Preload("Category").Where("user_id = ?", id).Order("created_at desc").Find(&deposits)

	var logs []models.AuditLog
	h.DB.Where("resource_type = 'User' AND resource_id = ?", id).Order("created_at desc").Limit(50).Find(&logs)

	doc := buildUserPDF(user, deposits, logs)

	var buf bytes.Buffer
	if err := doc.Output(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la generation du PDF."})
		return
	}

	filename := fmt.Sprintf("utilisateur_%d_%s.pdf", user.ID, time.Now().Format("20060102"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", strconv.Itoa(buf.Len()))
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}

func buildUserPDF(user models.User, deposits []models.DepositRequest, logs []models.AuditLog) *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)
	pdf.AddPage()

	tr := pdf.UnicodeTranslatorFromDescriptor("")

	if len(logoBytes) > 0 {
		pdf.RegisterImageOptionsReader("logo", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(logoBytes))
		pdf.ImageOptions("logo", 20, 12, 28, 0, false, fpdf.ImageOptions{}, 0, "")
	}

	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(0, 29, 50)
	pdf.SetXY(54, 16)
	pdf.CellFormat(130, 10, "UpcycleConnect", "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(64, 97, 127)
	pdf.SetXY(54, 26)
	pdf.CellFormat(130, 5, tr("Fiche utilisateur — exportée le "+time.Now().Format("02/01/2006")), "", 0, "L", false, 0, "")

	pdf.SetDrawColor(230, 234, 241)
	pdf.SetLineWidth(0.5)
	pdf.Line(20, 40, 190, 40)

	y := 48.0

	y = pdfSectionTitle(pdf, y, tr("Informations personnelles"))

	phone := "-"
	if user.Phone != nil && *user.Phone != "" {
		phone = *user.Phone
	}
	address := "-"
	if user.Address != nil && *user.Address != "" {
		address = tr(*user.Address)
	}
	verified := tr("Non vérifié")
	if user.EmailVerifiedAt != nil {
		verified = tr("Vérifié le ") + user.EmailVerifiedAt.Format("02/01/2006")
	}

	rows := [][2]string{
		{tr("Prénom"), tr(user.FirstName)},
		{"Nom", tr(user.LastName)},
		{"Email", user.Email},
		{tr("Téléphone"), phone},
		{"Adresse", address},
		{"Statut", tr(statusLabel(user.Status))},
		{tr("Email vérifié"), verified},
		{"Membre depuis", user.CreatedAt.Format("02/01/2006")},
	}
	for _, row := range rows {
		y = pdfRow(pdf, y, row[0], row[1])
	}

	y += 8
	y = pdfSectionTitle(pdf, y, tr(fmt.Sprintf("Dépôts d'objets (%d)", len(deposits))))

	if len(deposits) == 0 {
		y = pdfEmpty(pdf, y, tr("Aucun dépôt enregistré."))
	} else {
		pdf.SetFillColor(245, 248, 252)
		pdf.SetFont("Helvetica", "B", 8)
		pdf.SetTextColor(100, 116, 139)
		pdf.SetXY(20, y)
		pdf.CellFormat(10, 7, "#", "B", 0, "C", true, 0, "")
		pdf.CellFormat(70, 7, "Titre", "B", 0, "L", true, 0, "")
		pdf.CellFormat(25, 7, "Statut", "B", 0, "C", true, 0, "")
		pdf.CellFormat(30, 7, tr("Catégorie"), "B", 0, "L", true, 0, "")
		pdf.CellFormat(35, 7, "Date", "B", 0, "C", true, 0, "")
		y += 7

		pdf.SetFont("Helvetica", "", 8)
		for i, d := range deposits {
			if y > 265 {
				pdf.AddPage()
				y = 20
			}
			fill := i%2 == 0
			if fill {
				pdf.SetFillColor(250, 252, 255)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			pdf.SetTextColor(30, 41, 59)
			pdf.SetXY(20, y)

			cat := "-"
			if d.Category != nil {
				cat = tr(truncate(d.Category.Name, 15))
			}
			pdf.CellFormat(10, 6, fmt.Sprintf("%d", d.ID), "", 0, "C", fill, 0, "")
			pdf.CellFormat(70, 6, tr(truncate(d.Title, 42)), "", 0, "L", fill, 0, "")
			pdf.CellFormat(25, 6, tr(depositStatusLabel(d.Status)), "", 0, "C", fill, 0, "")
			pdf.CellFormat(30, 6, cat, "", 0, "L", fill, 0, "")
			pdf.CellFormat(35, 6, d.CreatedAt.Format("02/01/2006"), "", 0, "C", fill, 0, "")
			y += 6
		}
	}

	y += 8
	if y > 240 {
		pdf.AddPage()
		y = 20
	}
	y = pdfSectionTitle(pdf, y, tr(fmt.Sprintf("Journaux d'activité (%d)", len(logs))))

	if len(logs) == 0 {
		pdfEmpty(pdf, y, tr("Aucun log enregistré."))
	} else {
		pdf.SetFillColor(245, 248, 252)
		pdf.SetFont("Helvetica", "B", 8)
		pdf.SetTextColor(100, 116, 139)
		pdf.SetXY(20, y)
		pdf.CellFormat(40, 7, "Date", "B", 0, "L", true, 0, "")
		pdf.CellFormat(60, 7, "Action", "B", 0, "L", true, 0, "")
		pdf.CellFormat(70, 7, "Adresse IP", "B", 0, "L", true, 0, "")
		y += 7

		pdf.SetFont("Helvetica", "", 8)
		for i, l := range logs {
			if y > 265 {
				pdf.AddPage()
				y = 20
			}
			fill := i%2 == 0
			if fill {
				pdf.SetFillColor(250, 252, 255)
			} else {
				pdf.SetFillColor(255, 255, 255)
			}
			pdf.SetTextColor(30, 41, 59)
			pdf.SetXY(20, y)

			ip := "-"
			if l.IPAddress != nil {
				ip = *l.IPAddress
			}
			pdf.CellFormat(40, 6, l.CreatedAt.Format("02/01/2006 15:04"), "", 0, "L", fill, 0, "")
			pdf.CellFormat(60, 6, l.Action, "", 0, "L", fill, 0, "")
			pdf.CellFormat(70, 6, ip, "", 0, "L", fill, 0, "")
			y += 6
		}
	}

	pdf.SetFont("Helvetica", "I", 7)
	pdf.SetTextColor(160, 174, 192)
	pdf.SetXY(20, 285)
	pdf.CellFormat(170, 5, tr("Document confidentiel — UpcycleConnect — "+time.Now().Format("02/01/2006")), "", 0, "C", false, 0, "")

	return pdf
}

func pdfSectionTitle(pdf *fpdf.Fpdf, y float64, title string) float64 {
	pdf.SetFillColor(0, 109, 53)
	pdf.SetFont("Helvetica", "B", 9)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(20, y)
	pdf.CellFormat(170, 7, "  "+title, "", 1, "L", true, 0, "")
	return y + 10
}

func pdfRow(pdf *fpdf.Fpdf, y float64, label, value string) float64 {
	if y > 268 {
		pdf.AddPage()
		y = 20
	}
	pdf.SetFont("Helvetica", "B", 8)
	pdf.SetTextColor(100, 116, 139)
	pdf.SetXY(20, y)
	pdf.CellFormat(45, 6, label, "", 0, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(30, 41, 59)
	pdf.CellFormat(125, 6, value, "", 1, "L", false, 0, "")
	return y + 6
}

func pdfEmpty(pdf *fpdf.Fpdf, y float64, msg string) float64 {
	pdf.SetFont("Helvetica", "I", 8)
	pdf.SetTextColor(148, 163, 184)
	pdf.SetXY(20, y)
	pdf.CellFormat(170, 6, msg, "", 1, "L", false, 0, "")
	return y + 6
}

func statusLabel(s string) string {
	switch s {
	case "active":
		return "Actif"
	case "banned":
		return "Banni"
	case "inactive":
		return "Inactif"
	case "pending":
		return "En attente"
	default:
		return s
	}
}

func depositStatusLabel(s string) string {
	switch s {
	case "pending":
		return "En attente"
	case "approved":
		return "Approuvé"
	case "rejected":
		return "Refusé"
	case "collected":
		return "Collecté"
	default:
		return s
	}
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "..."
}
