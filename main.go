package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	cacheDirName   = "dayfact"
	cacheFileName  = "last_shown"
	userAgent      = "dailyTerminal"
	requestTimeout = 2 * time.Second
)

var mesesEs = map[time.Month]string{
	time.January:   "enero",
	time.February:  "febrero",
	time.March:     "marzo",
	time.April:     "abril",
	time.May:       "mayo",
	time.June:      "junio",
	time.July:      "julio",
	time.August:    "agosto",
	time.September: "septiembre",
	time.October:   "octubre",
	time.November:  "noviembre",
	time.December:  "diciembre",
}

type holidayEntry struct {
	Text string `json:"text"`
}

type holidaysResponse struct {
	Holidays []holidayEntry `json:"holidays"`
}

// cachePath devuelve la ruta al archivo donde guardamos la última fecha
// en la que ya se mostró la efeméride (normalmente ~/.cache/dayfact/last_shown).
func cachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	full := filepath.Join(dir, cacheDirName)
	if err := os.MkdirAll(full, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(full, cacheFileName), nil
}

func alreadyShownToday(path, today string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(data)) == today
}

func markShown(path, today string) error {
	return os.WriteFile(path, []byte(today), 0o644)
}

// fetchHolidays consulta el endpoint "onthisday/holidays" de Wikipedia,
// que es público y no requiere API key, solo un User-Agent identificable.
func fetchHolidays(month, day string) ([]string, error) {
	url := fmt.Sprintf("https://en.wikipedia.org/api/rest_v1/feed/onthisday/holidays/%s/%s", month, day)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia respondió con status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseHolidays(body)
}

// parseHolidays toma el JSON crudo de la API y devuelve los nombres de las
// festividades, omitiendo las festividades cristianas. Está separada de
// fetchHolidays a propósito: no depende de red, así que se puede probar
// con datos falsos (ver main_test.go).
func parseHolidays(body []byte) ([]string, error) {
	var parsed holidaysResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	const feastDayPrefix = "Christian feast day"

	names := make([]string, 0, len(parsed.Holidays))
	for _, h := range parsed.Holidays {
		if strings.HasPrefix(h.Text, feastDayPrefix) {
			continue
		}
		names = append(names, h.Text)
	}
	return names, nil
}

func main() {
	now := time.Now()
	today := now.Format("2006-01-02")

	path, err := cachePath()
	if err != nil {
		// Sin acceso al cache: no arriesgamos, dejamos que el bashrc use la quote.
		os.Exit(1)
	}

	if alreadyShownToday(path, today) {
		os.Exit(1)
	}

	month := now.Format("01")
	day := now.Format("02")

	holidays, err := fetchHolidays(month, day)
	if err != nil || len(holidays) == 0 {
		// Sin internet o sin datos para hoy: fallamos en silencio,
		// el bashrc caerá a la quote de siempre.
		os.Exit(1)
	}

	fmt.Printf("Hoy es %d de %s, se celebra:\n", now.Day(), mesesEs[now.Month()])
	for _, h := range holidays {
		fmt.Printf("  • %s\n", h)
	}

	// Si falla el guardado no es grave, simplemente se volvería a mostrar.
	_ = markShown(path, today)
	os.Exit(0)
}
