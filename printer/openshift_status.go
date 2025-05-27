package printer

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/hidetatz/kubecolor/color"
)

// OpenShiftStatusPrinter prints the output of 'oc status' with colors.
type OpenShiftStatusPrinter struct {
	DarkBackground bool
}

// Print reads r then write it to w with colors.
func (p *OpenShiftStatusPrinter) Print(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintln(w, p.colorizeLine(line))
	}
}

func (p *OpenShiftStatusPrinter) colorizeLine(line string) string {
	// Project context
	if strings.HasPrefix(line, "In project ") {
		return color.Apply(line, color.Cyan)
	}

	// Service names (e.g., svc/service-name)
	reService := regexp.MustCompile(`(svc/\S+)`)
	line = reService.ReplaceAllStringFunc(line, func(match string) string {
		return color.Apply(match, color.Green)
	})

	// Deployment config names (e.g., dc/deployment-config-name)
	reDc := regexp.MustCompile(`(dc/\S+)`)
	line = reDc.ReplaceAllStringFunc(line, func(match string) string {
		return color.Apply(match, color.Blue)
	})

	// URLs / routes
	// A simple regex for URLs, might need refinement
	reURL := regexp.MustCompile(`(https?://[^\s]+)`)
	line = reURL.ReplaceAllStringFunc(line, func(match string) string {
		return color.Apply(match, color.Magenta)
	})

	// Keywords
	// Using regex with word boundaries to avoid partial matches (e.g. "deployment" contains "deploy")
	keywords := map[string]color.Color{
		`\brunning\b`:  color.Green,
		`\bdeployed\b`: color.Green,
		`\bfailed\b`:   color.Red,
		// Add more keywords as needed
	}

	for keyword, c := range keywords {
		reKeyword := regexp.MustCompile(keyword)
		line = reKeyword.ReplaceAllStringFunc(line, func(match string) string {
			return color.Apply(match, c)
		})
	}

	return line
}
