package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"upcycleconnect/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"
	"gorm.io/gorm"
)

type ExportHandler struct {
	DB    *gorm.DB
	Audit *services.AuditService
}

type dataset struct {
	Label   string
	Headers []string
	Query   func(db *gorm.DB, from, to *time.Time) [][]string
}

func eur(cents int64) string { return fmt.Sprintf("%.2f", float64(cents)/100.0) }

func applyRange(q *gorm.DB, col string, from, to *time.Time) *gorm.DB {
	if from != nil {
		q = q.Where(col+" >= ?", *from)
	}
	if to != nil {
		q = q.Where(col+" <= ?", *to)
	}
	return q
}

func (h *ExportHandler) registry() map[string]dataset {
	return map[string]dataset{
		"users": {"Utilisateurs", []string{"ID", "Email", "Prénom", "Nom", "Téléphone", "Statut", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					ID                                 uint
					Email, FirstName, LastName, Status string
					Phone                              *string
					CreatedAt                          time.Time
				}
				var rs []R
				applyRange(db.Table("users").Where("deleted_at IS NULL"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					ph := ""
					if r.Phone != nil {
						ph = *r.Phone
					}
					out = append(out, []string{strconv.Itoa(int(r.ID)), r.Email, r.FirstName, r.LastName, ph, r.Status, r.CreatedAt.Format("2006-01-02 15:04")})
				}
				return out
			}},
		"providers": {"Prestataires", []string{"Société", "SIRET", "Statut", "Email", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					CompanyName, Status, Email string
					Siret                      *string
					CreatedAt                  time.Time
				}
				var rs []R
				applyRange(db.Table("provider_profiles").
					Joins("LEFT JOIN users ON users.id = provider_profiles.user_id"), "provider_profiles.created_at", from, to).
					Select("provider_profiles.company_name, provider_profiles.status, provider_profiles.siret, users.email, provider_profiles.created_at").
					Order("provider_profiles.created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					si := ""
					if r.Siret != nil {
						si = *r.Siret
					}
					out = append(out, []string{r.CompanyName, si, r.Status, r.Email, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"prestations": {"Prestations / Annonces", []string{"ID", "Titre", "Prix", "Type", "Statut", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					ID                       uint
					Title, PriceType, Status string
					Price                    *string
					CreatedAt                time.Time
				}
				var rs []R
				applyRange(db.Table("prestations").Where("deleted_at IS NULL"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					pr := ""
					if r.Price != nil {
						pr = *r.Price
					}
					out = append(out, []string{strconv.Itoa(int(r.ID)), r.Title, pr, r.PriceType, r.Status, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"reservations": {"Réservations", []string{"ID", "Montant", "Commission", "Net pro", "Statut", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					ID                                     uint
					AmountCents, CommissionCents, NetCents int64
					Status                                 string
					CreatedAt                              time.Time
				}
				var rs []R
				applyRange(db.Table("reservations").Where("deleted_at IS NULL"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					out = append(out, []string{strconv.Itoa(int(r.ID)), eur(r.AmountCents), eur(r.CommissionCents), eur(r.NetCents), r.Status, r.CreatedAt.Format("2006-01-02 15:04")})
				}
				return out
			}},
		"invoices": {"Factures & Devis", []string{"Numéro", "Type", "Intitulé", "Montant", "TVA %", "Statut", "Émise le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					Number, Type, PrestationTitle, Status string
					AmountCents                           int64
					TVAPercent                            float64
					CreatedAt                             time.Time
				}
				var rs []R
				applyRange(db.Table("invoices"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					out = append(out, []string{r.Number, r.Type, r.PrestationTitle, eur(r.AmountCents), fmt.Sprintf("%.0f", r.TVAPercent), r.Status, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"contracts": {"Contrats", []string{"Numéro", "Client", "Prestation", "Montant", "Statut", "Signé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					Number, CustomerName, PrestationTitle, Status string
					AmountCents                                   int64
					SignedAt                                      time.Time
				}
				var rs []R
				applyRange(db.Table("contracts").Where("deleted_at IS NULL"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					out = append(out, []string{r.Number, r.CustomerName, r.PrestationTitle, eur(r.AmountCents), r.Status, r.SignedAt.Format("2006-01-02")})
				}
				return out
			}},
		"deposits": {"Dépôts en conteneur", []string{"ID", "Titre", "Statut", "Code-barres", "Collecté le", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					ID            uint
					Title, Status string
					QRCode        *string
					CollectedAt   *time.Time
					CreatedAt     time.Time
				}
				var rs []R
				applyRange(db.Table("deposit_requests"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					qr := ""
					if r.QRCode != nil {
						qr = *r.QRCode
					}
					col := ""
					if r.CollectedAt != nil {
						col = r.CollectedAt.Format("2006-01-02")
					}
					out = append(out, []string{strconv.Itoa(int(r.ID)), r.Title, r.Status, qr, col, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"events": {"Événements & Formations", []string{"ID", "Titre", "Début", "Statut", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					ID                   uint
					Title, Status        string
					StartDate, CreatedAt time.Time
				}
				var rs []R
				applyRange(db.Table("platform_events").Where("deleted_at IS NULL"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					out = append(out, []string{strconv.Itoa(int(r.ID)), r.Title, r.StartDate.Format("2006-01-02"), r.Status, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"campaigns": {"Campagnes publicitaires", []string{"ID", "Titre", "Budget", "Statut", "Payée le", "Créée le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					ID            uint
					Title, Status string
					BudgetCents   int64
					PaidAt        *time.Time
					CreatedAt     time.Time
				}
				var rs []R
				applyRange(db.Table("campaigns").Where("deleted_at IS NULL"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					pd := ""
					if r.PaidAt != nil {
						pd = r.PaidAt.Format("2006-01-02")
					}
					out = append(out, []string{strconv.Itoa(int(r.ID)), r.Title, eur(r.BudgetCents), r.Status, pd, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"subscriptions": {"Abonnements", []string{"User ID", "Plan", "Statut", "Fin période", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					UserID           uint
					Plan, Status     string
					CurrentPeriodEnd *time.Time
					CreatedAt        time.Time
				}
				var rs []R
				applyRange(db.Table("subscriptions"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					pe := ""
					if r.CurrentPeriodEnd != nil {
						pe = r.CurrentPeriodEnd.Format("2006-01-02")
					}
					out = append(out, []string{strconv.Itoa(int(r.UserID)), r.Plan, r.Status, pe, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"collection_points": {"Points de collecte", []string{"Nom", "Adresse", "Ville", "CP", "Actif", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					Name, Address, City, PostalCode string
					IsActive                        bool
					CreatedAt                       time.Time
				}
				var rs []R
				applyRange(db.Table("collection_points"), "created_at", from, to).Order("city ASC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					act := "Oui"
					if !r.IsActive {
						act = "Non"
					}
					out = append(out, []string{r.Name, r.Address, r.City, r.PostalCode, act, r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
		"notifications": {"Notifications envoyées", []string{"Type", "Titre", "Destinataire", "Lu", "Envoyée le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					Type, Title, Email string
					ReadAt             *time.Time
					CreatedAt          time.Time
				}
				var rs []R
				applyRange(db.Table("notifications").
					Joins("LEFT JOIN users ON users.id = notifications.user_id").Where("notifications.deleted_at IS NULL"), "notifications.created_at", from, to).
					Select("notifications.type, notifications.title, users.email, notifications.read_at, notifications.created_at").
					Order("notifications.created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					lu := "Non"
					if r.ReadAt != nil {
						lu = "Oui"
					}
					out = append(out, []string{r.Type, r.Title, r.Email, lu, r.CreatedAt.Format("2006-01-02 15:04")})
				}
				return out
			}},
		"audit_logs": {"Journal d'audit", []string{"Action", "Utilisateur", "Ressource", "Réf.", "IP", "Date"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					Action, ResourceType string
					UserID               *uint
					ResourceID           *uint
					IPAddress            *string
					CreatedAt            time.Time
				}
				var rs []R
				applyRange(db.Table("audit_logs"), "created_at", from, to).Order("created_at DESC").Limit(5000).Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					adm := ""
					if r.UserID != nil {
						adm = strconv.Itoa(int(*r.UserID))
					}
					rid := ""
					if r.ResourceID != nil {
						rid = strconv.Itoa(int(*r.ResourceID))
					}
					ip := ""
					if r.IPAddress != nil {
						ip = *r.IPAddress
					}
					out = append(out, []string{r.Action, adm, r.ResourceType, rid, ip, r.CreatedAt.Format("2006-01-02 15:04")})
				}
				return out
			}},
		"tips": {"Conseils", []string{"Titre", "Catégorie", "Statut", "Vues", "Créé le"},
			func(db *gorm.DB, from, to *time.Time) [][]string {
				type R struct {
					Title, Category, Status string
					ViewCount               uint
					CreatedAt               time.Time
				}
				var rs []R
				applyRange(db.Table("tips").Where("deleted_at IS NULL"), "created_at", from, to).Order("created_at DESC").Scan(&rs)
				out := [][]string{}
				for _, r := range rs {
					out = append(out, []string{r.Title, r.Category, r.Status, strconv.Itoa(int(r.ViewCount)), r.CreatedAt.Format("2006-01-02")})
				}
				return out
			}},
	}
}

func (h *ExportHandler) Datasets(c *gin.Context) {
	reg := h.registry()
	out := []gin.H{}
	for key, ds := range reg {
		out = append(out, gin.H{"key": key, "label": ds.Label, "columns": len(ds.Headers)})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func parseDate(s string, endOfDay bool) *time.Time {
	if s == "" {
		return nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		if endOfDay {
			t = t.Add(24*time.Hour - time.Second)
		}
		return &t
	}
	return nil
}

func (h *ExportHandler) Export(c *gin.Context) {
	key := c.Query("dataset")
	format := c.DefaultQuery("format", "csv")
	from := parseDate(c.Query("from"), false)
	to := parseDate(c.Query("to"), true)

	ds, ok := h.registry()[key]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Jeu de données inconnu"})
		return
	}

	rows := ds.Query(h.DB, from, to)

	if h.Audit != nil {
		h.Audit.Log(c, "data.exported", "Export", nil, nil, map[string]interface{}{
			"dataset": key, "format": format, "rows": len(rows),
		})
	}

	stamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("upcycleconnect-%s-%s", key, stamp)

	if format == "pdf" {
		h.writePDF(c, ds, rows, from, to, filename)
		return
	}
	h.writeCSV(c, ds, rows, filename)
}

func (h *ExportHandler) writeCSV(c *gin.Context, ds dataset, rows [][]string, filename string) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.csv"`, filename))

	c.Writer.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(c.Writer)
	w.Write(ds.Headers)
	for _, r := range rows {
		w.Write(r)
	}
	w.Flush()
}

func (h *ExportHandler) writePDF(c *gin.Context, ds dataset, rows [][]string, from, to *time.Time, filename string) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 12, 10)
	pdf.SetAutoPageBreak(true, 12)
	pdf.AddPage()
	tr := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 109, 53)
	pdf.Cell(0, 8, tr("UpcycleConnect — "+ds.Label))
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(100, 116, 139)
	period := "Toutes périodes"
	if from != nil || to != nil {
		f := "…"
		t := "…"
		if from != nil {
			f = from.Format("02/01/2006")
		}
		if to != nil {
			t = to.Format("02/01/2006")
		}
		period = "Période : " + f + " → " + t
	}
	pdf.Cell(0, 5, tr(fmt.Sprintf("%s · %d lignes · généré le %s", period, len(rows), time.Now().Format("02/01/2006 15:04"))))
	pdf.Ln(8)

	pageW, _ := pdf.GetPageSize()
	usable := pageW - 20
	n := len(ds.Headers)
	colW := usable / float64(n)

	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(0, 109, 53)
	pdf.SetTextColor(255, 255, 255)
	for _, hd := range ds.Headers {
		pdf.CellFormat(colW, 7, tr(hd), "1", 0, "L", true, 0, "")
	}
	pdf.Ln(7)

	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(30, 41, 59)
	fill := false
	for _, r := range rows {
		if fill {
			pdf.SetFillColor(244, 247, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		for i := 0; i < n; i++ {
			val := ""
			if i < len(r) {
				val = r[i]
			}

			if len(val) > 60 {
				val = val[:57] + "..."
			}
			pdf.CellFormat(colW, 6, tr(val), "1", 0, "L", true, 0, "")
		}
		pdf.Ln(6)
		fill = !fill
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, filename))
	_ = pdf.Output(c.Writer)
}
