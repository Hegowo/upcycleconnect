package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"upcycleconnect/backend/config"

	"github.com/gin-gonic/gin"
)

type SiretHandler struct {
	Cfg *config.Config
}

func (h *SiretHandler) Lookup(c *gin.Context) {
	siret := c.Param("siret")

	matched, _ := regexp.MatchString(`^\d{14}$`, siret)
	if !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SIRET invalide (14 chiffres requis)"})
		return
	}

	url := fmt.Sprintf("https://api.societe.com/api/v1/entreprise/%s/infoslegales", siret)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne"})
		return
	}
	req.Header.Set("X-Authorization", "socapi "+h.Cfg.SocieteApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Impossible de joindre l'API SIRET"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lecture réponse"})
		return
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Clé API SIRET non autorisée"})
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entreprise non trouvée"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("Erreur API SIRET: %d", resp.StatusCode)})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Réponse API invalide"})
		return
	}

	c.JSON(http.StatusOK, data)
}
