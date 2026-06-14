package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// dataURIExt maps a data-URI mime type to a file extension.
var dataURIExt = map[string]string{
	"image/jpeg": ".jpg",
	"image/jpg":  ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

// SaveDataURIImage decodes a base64 data-URI ("data:image/jpeg;base64,...") and
// writes it as a file under dir, returning its public URL ("<urlPrefix>/<file>").
//
// If s is empty or is NOT a data URI (e.g. already a "/uploads/..." URL or an
// http URL), it is returned unchanged with converted=false. This makes the call
// idempotent and safe to run over already-migrated values.
func SaveDataURIImage(s, dir, urlPrefix string) (out string, converted bool, err error) {
	if s == "" || !strings.HasPrefix(s, "data:") {
		return s, false, nil
	}

	// Split "data:<mime>;base64,<payload>".
	comma := strings.IndexByte(s, ',')
	if comma < 0 {
		return s, false, fmt.Errorf("data URI sans virgule")
	}
	header := s[5:comma] // strip leading "data:"
	payload := s[comma+1:]

	mime := header
	if i := strings.IndexByte(header, ';'); i >= 0 {
		mime = header[:i]
	}
	mime = strings.ToLower(strings.TrimSpace(mime))
	ext := dataURIExt[mime]
	if ext == "" {
		ext = ".jpg"
	}

	if !strings.Contains(header, "base64") {
		return s, false, fmt.Errorf("data URI non base64")
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(payload))
	if err != nil {
		return s, false, fmt.Errorf("base64 invalide: %w", err)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return s, false, err
	}
	buf := make([]byte, 8)
	_, _ = rand.Read(buf)
	filename := hex.EncodeToString(buf) + ext
	if err := os.WriteFile(filepath.Join(dir, filename), raw, 0644); err != nil {
		return s, false, err
	}
	return strings.TrimRight(urlPrefix, "/") + "/" + filename, true, nil
}
