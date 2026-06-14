package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"upcycleconnect/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SearchHandler struct {
	DB *gorm.DB
}

type SearchResult struct {
	Type         string `json:"type"`
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle,omitempty"`
	Meta         string `json:"meta,omitempty"`
	URL          string `json:"url"`
	Slug         string `json:"slug,omitempty"`
	CategorySlug string `json:"category_slug,omitempty"`
}

func (h *SearchHandler) Search(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	empty := gin.H{
		"events":      []SearchResult{},
		"prestations": []SearchResult{},
		"threads":     []SearchResult{},
		"providers":   []SearchResult{},
	}
	if len([]rune(q)) < 2 {
		c.JSON(http.StatusOK, empty)
		return
	}

	like := "%" + q + "%"
	limit := 5

	var events []models.Event
	h.DB.Where(
		"(title LIKE ? OR description LIKE ? OR location LIKE ?) AND status = 'published'",
		like, like, like,
	).Limit(limit).Find(&events)

	eventResults := make([]SearchResult, len(events))
	for i, e := range events {
		sub := ""
		if e.Location != nil {
			sub = *e.Location
		}
		meta := ""
		if !e.StartDate.IsZero() {
			meta = e.StartDate.Format("02 Jan 2006")
		}
		eventResults[i] = SearchResult{
			Type:     "event",
			ID:       e.ID,
			Title:    e.Title,
			Subtitle: sub,
			Meta:     meta,
			URL:      fmt.Sprintf("/evenements/%d", e.ID),
		}
	}

	var prestations []models.Prestation
	h.DB.Where(
		"(title LIKE ? OR description LIKE ?) AND status = 'published'",
		like, like,
	).Limit(limit).Find(&prestations)

	prestationResults := make([]SearchResult, len(prestations))
	for i, p := range prestations {
		meta := ""
		if p.Price != nil {
			meta = *p.Price + " €"
		}
		prestationResults[i] = SearchResult{
			Type:  "prestation",
			ID:    p.ID,
			Title: p.Title,
			Meta:  meta,
			URL:   fmt.Sprintf("/prestations/%d", p.ID),
		}
	}

	var threads []models.ForumThread
	h.DB.Where("(title LIKE ? OR content LIKE ?) AND deleted_at IS NULL", like, like).
		Preload("Category").Limit(limit).Find(&threads)

	threadResults := make([]SearchResult, len(threads))
	for i, t := range threads {
		catName := ""
		catSlug := "_"
		if t.Category != nil {
			catName = t.Category.Name
			catSlug = t.Category.Slug
		}
		meta := ""
		if !t.CreatedAt.IsZero() {
			if time.Since(t.CreatedAt) < 24*time.Hour {
				meta = "Aujourd'hui"
			} else if time.Since(t.CreatedAt) < 7*24*time.Hour {
				meta = fmt.Sprintf("Il y a %d j", int(time.Since(t.CreatedAt).Hours()/24))
			} else {
				meta = t.CreatedAt.Format("02 Jan")
			}
		}
		threadResults[i] = SearchResult{
			Type:         "thread",
			ID:           t.ID,
			Title:        t.Title,
			Subtitle:     catName,
			Meta:         meta,
			URL:          fmt.Sprintf("/communaute/forum/%s/%d", catSlug, t.ID),
			CategorySlug: catSlug,
		}
	}

	type provRow struct {
		ID          uint
		CompanyName string
		Description *string
		UserID      uint
		FirstName   string
		LastName    string
	}
	var provRows []provRow
	h.DB.Table("provider_profiles").
		Select("provider_profiles.id, provider_profiles.company_name, provider_profiles.description, provider_profiles.user_id, users.first_name, users.last_name").
		Joins("JOIN users ON users.id = provider_profiles.user_id").
		Where("(provider_profiles.company_name LIKE ? OR provider_profiles.description LIKE ?) AND provider_profiles.status = 'approved'", like, like).
		Limit(limit).
		Scan(&provRows)

	providerResults := make([]SearchResult, len(provRows))
	for i, p := range provRows {
		sub := strings.TrimSpace(p.FirstName + " " + p.LastName)
		providerResults[i] = SearchResult{
			Type:     "provider",
			ID:       p.ID,
			Title:    p.CompanyName,
			Subtitle: sub,
			URL:      fmt.Sprintf("/prestations?provider=%d", p.UserID),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"events":      eventResults,
		"prestations": prestationResults,
		"threads":     threadResults,
		"providers":   providerResults,
	})
}
