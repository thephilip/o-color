package printer

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/hidetatz/kubecolor/color"
)

// DescribePrinter is a specific printer to print kubectl describe format.
type DescribePrinter struct {
	DarkBackground bool
	TablePrinter   *TablePrinter
}

// Define route-specific keywords and colors at package level or within the struct if preferred
var (
	routeDetectionKeywords = []string{"Requested Host:", "TLS Termination:", "Ingress:"} // Keywords highly specific to routes for detection
	routeSpecificKeys      = map[string]bool{
		"Name:":            true,
		"Namespace:":       true,
		"Created:":         true,
		"Labels:":          true,
		"Annotations:":     true,
		"Requested Host:":  true,
		"Path:":            true,
		"TLS Termination:": true,
		"Service:":         true,
		"Weight:":          true,
		"Endpoints:":       true,
		"Ingress:":         true,
	}
	routeKeyColor            = color.Yellow
	routeResourceNameColor   = color.Green
	routeEndpointColor       = color.Cyan
	routeTLSColorEdge        = color.Blue
	routeTLSColorPassthrough = color.Yellow
	routeTLSColorReencrypt   = color.Yellow
	routeCommaColor          = color.White

	// aliases for backward compatibility with older variable names
	ocRouteKeyColor            = routeKeyColor
	ocRouteResourceNameColor   = routeResourceNameColor
	ocRouteEndpointColor       = routeEndpointColor
	ocRouteTLSColorEdge        = routeTLSColorEdge
	ocRouteTLSColorPassthrough = routeTLSColorPassthrough
	ocRouteTLSColorReencrypt   = routeTLSColorReencrypt
	ocRouteCommaColor          = routeCommaColor
)

func (dp *DescribePrinter) Print(r io.Reader, w io.Writer) {
	basicIndentWidth := 2 // according to kubectl describe format
	scanner := bufio.NewScanner(r)
	isRoute := false // Flag to indicate if current resource is likely a route

	for scanner.Scan() {
		line := scanner.Text()

		// Attempt to detect if this is a route description
		if !isRoute {
			for _, keyword := range routeDetectionKeywords {
				if strings.Contains(line, keyword) {
					isRoute = true
					break
				}
			}
		}

		if line == "" {
			fmt.Fprintln(w)
			continue
		}

		// Split a line by spaces to colorize and render them
		// For example:
		// e.g. 1-----------------
		// Status:         Running
		// -----------------------
		// spacesIndices: [[7, 15]] // <- where spaces locate
		// columns: ["Status:", "Running"]
		//
		// e.g. 2--------------------------------------------
		//     Ports:          10001/TCP, 5000/TCP, 18000/TCP
		// --------------------------------------------------
		// spacesIndices: [[0, 3], [10, 19]] // <- where spaces locate
		// columns: ["Ports:", "10001/TCP, 5000/TCP, 18000/TCP"]
		//
		// So now, we know where to render which column.
		spacesIndices := spaces.FindAllStringIndex(line, -1)
		columns := spaces.Split(line, -1)
		// when the line has indent (spaces on left), the first item will be
		// just a "" and we don't need it so remove
		if len(columns) > 0 {
			if columns[0] == "" {
				columns = columns[1:]
			}
		}

		// First, identify if there is an indent
		indentCnt := findIndent(line)
		indent := toSpaces(indentCnt)
		if indentCnt > 0 {
			// TODO: Remove this condition for workaround
			// Basically, kubectl describe output has its indentation level
			// with **2** spaces, but "Resource Quota" section in
			// `kubectl describe ns` output has only 1 space at the head.
			// Because of it, indentCnt is still 1, but the indent space is not in `spacesIndices` (see regex definition of `spaces`)
			// So it must be checked here
			// https://github.com/hidetatz/kubecolor/issues/36
			// When https://github.com/kubernetes/kubectl/issues/1005#issuecomment-758385759 is fixed
			// this is not needed anymore.
			if indentCnt > 1 {
				// when an indent exists, removes it because it's already captured by "indent" var
				spacesIndices = spacesIndices[1:]
			}
		}

		spacesCnt := 0
		if len(spacesIndices) > 0 {
			spacesCnt = spacesIndices[0][1] - spacesIndices[0][0]
		}

		// when there are multiple columns, treat is as table format
		if len(columns) > 2 {
			dp.TablePrinter.printLineAsTableFormat(w, line, getColorsByBackground(dp.DarkBackground))
			continue
		}

		// Determine effective colors to use, initially based on generic logic
		effectiveKeyColor := getColorByKeyIndent(indentCnt, basicIndentWidth, dp.DarkBackground)
		effectiveValColor := getColorByValueType(columns[0], dp.DarkBackground) // Default for single-column lines
		if len(columns) > 1 {
			effectiveValColor = getColorByValueType(columns[1], dp.DarkBackground) // Default for value part of key-value
		}

		if isRoute && len(columns) > 0 {
			keyPart := columns[0]
			// Trim space for map lookups and switch, but preserve original keyPart for printing if needed.
			trimmedSpaceKeyPart := strings.TrimSpace(keyPart)
			// For map lookups, ensure the key has a colon if it's a primary field.
			// This logic assumes routeSpecificKeys stores keys exactly as they appear, e.g. "Name:", "  Host:"
			keyForMapLookup := keyPart

			// Override Key Color for Routes
			if _, ok := routeSpecificKeys[keyForMapLookup]; ok {
				effectiveKeyColor = ocRouteKeyColor
			} else if indentCnt > basicIndentWidth {
				// For indented keys within a route section (e.g. "Host:" under "Ingress:")
				// If the space-trimmed version (e.g. "Host:") is in our specific map, use route key color.
				// This allows routeSpecificKeys to define "Host:" for sub-sections without leading spaces.
				if _, okSubKey := routeSpecificKeys[trimmedSpaceKeyPart+":"]; okSubKey {
					effectiveKeyColor = ocRouteKeyColor
				} else {
					// Otherwise, use generic indentation logic for unknown sub-keys.
					effectiveKeyColor = getColorByKeyIndent(indentCnt, basicIndentWidth+1, dp.DarkBackground)
				}
			}
			// else, default effectiveKeyColor from generic logic applies.

			// Override Value Color for Routes (if a value column exists)
			if len(columns) > 1 {
				valuePart := columns[1]
				// Use space-trimmed key (without colon) for switch cases for cleaner matching.
				switch strings.TrimSuffix(trimmedSpaceKeyPart, ":") {
				case "Name", "Requested Host":
					effectiveValColor = ocRouteResourceNameColor
				case "Service":
					effectiveValColor = ocRouteResourceNameColor
				case "Endpoints":
					effectiveValColor = ocRouteEndpointColor
				case "TLS Termination":
					if strings.HasPrefix(strings.ToLower(valuePart), "edge") {
						effectiveValColor = ocRouteTLSColorEdge
					} else if strings.HasPrefix(strings.ToLower(valuePart), "passthrough") {
						effectiveValColor = ocRouteTLSColorPassthrough
					} else if strings.HasPrefix(strings.ToLower(valuePart), "reencrypt") {
						effectiveValColor = ocRouteTLSColorReencrypt
					}
					// else, effectiveValColor from generic typing is used if no prefix matches.
				}
				// For other keys not listed, effectiveValColor (from generic typing) is used.
			} else if len(columns) == 1 && !strings.HasSuffix(keyPart, ":") {
				// Single column, not a key. If isRoute, it's likely descriptive text.
				effectiveValColor = getColorByValueType(keyPart, dp.DarkBackground)
				effectiveKeyColor = effectiveValColor // Print this line with the value's color
			}
		}

		// TODO: Remove this if statement for workaround (Kubectl 1.19.3 bug)
		// This workaround is for lines like " Resource Used Hard" where "Resource" has one leading space.
		// The current column splitting logic might already handle this if columns[0] becomes empty and is removed.
		// If columns[0] is the first word (e.g. "Resource") but indentCnt is 1 due to the bug,
		// this check might be needed.
		keyToPrint := columns[0]
		if indentCnt == 1 && strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "  ") && columns[0] != "" {
			// Heuristic for the single-space indent bug if columns[0] isn't empty.
			// If columns[0] is already trimmed by the initial split logic, this specific check might not be needed.
			// The original code did `columns[0] = strings.TrimLeft(columns[0], " ")` which is more aggressive.
			// Let's ensure keyToPrint is the actual content intended for the key column.
		} else if strings.HasPrefix(keyToPrint, " ") && indentCnt > 0 && len(spacesIndices) > 0 && spacesIndices[0][0] == 0 {
			// This was the original logic to handle the workaround, let's refine it.
			// If the first column still has leading spaces AND it's an indented line, trim.
			// This should ideally be handled by the columns splitting logic or findIndent.
			// For now, if columns[0] (keyToPrint) has leading spaces and it's not a non-indented line,
			// it might be an artifact of the split.
		}

		// Apply coloring to the key part that will be printed
		// Use the original columns[0] for TrimRight because keyToPrint might have been trimmed.
		// However, the content to color should be from keyToPrint if it was modified.
		finalKeyString := strings.TrimRight(keyToPrint, ":")
		coloredKeyOutput := color.Apply(finalKeyString, effectiveKeyColor)
		if strings.HasSuffix(keyToPrint, ":") {
			coloredKeyOutput += ":"
		}

		if len(columns) == 1 {
			fmt.Fprintf(w, "%s%s\n", indent, coloredKeyOutput) // Print single column line with its determined color
			continue
		}

		// len(columns) > 1, so there is a value part
		valueOutput := columns[1]

		// Specific rendering for multi-part values in routes
		if isRoute {
			// Use TrimSpace for switch to handle keys like "  Service:" correctly
			trimmedKeyForSwitch := strings.TrimSuffix(strings.TrimSpace(columns[0]), ":")
			switch trimmedKeyForSwitch {
			case "Service":
				if strings.Contains(valueOutput, "(") && strings.Contains(valueOutput, "%") {
					parts := strings.SplitN(valueOutput, " ", 2)
					serviceNameColored := color.Apply(parts[0], ocRouteResourceNameColor)
					weightPartColored := ""
					if len(parts) > 1 {
						weightPartColored = " " + color.Apply(parts[1], getColorByValueType(parts[1], dp.DarkBackground))
					}
					fmt.Fprintf(w, "%s%s%s%s\n", indent, coloredKeyOutput, toSpaces(spacesCnt), serviceNameColored+weightPartColored)
					continue
				}
			case "Endpoints":
				endpointParts := strings.Split(valueOutput, ",")
				var coloredEpStrings []string
				for i, ep := range endpointParts {
					trimmedEp := strings.TrimSpace(ep)
					coloredEpStrings = append(coloredEpStrings, color.Apply(trimmedEp, ocRouteEndpointColor))
					if i < len(endpointParts)-1 {
						coloredEpStrings = append(coloredEpStrings, color.Apply(",", ocRouteCommaColor)+" ")
					}
				}
				fmt.Fprintf(w, "%s%s%s%s\n", indent, coloredKeyOutput, toSpaces(spacesCnt), strings.Join(coloredEpStrings, ""))
				continue
			case "TLS Termination":
				var finalTLSOutput string
				spaceIdx := strings.Index(valueOutput, " ")
				firstWord := valueOutput
				remainingText := ""
				if spaceIdx != -1 {
					firstWord = valueOutput[:spaceIdx]
					remainingText = valueOutput[spaceIdx:]
				}
				lowerFirstWord := strings.ToLower(firstWord)
				coloredFirstWord := ""

				if lowerFirstWord == "edge" {
					coloredFirstWord = color.Apply(firstWord, ocRouteTLSColorEdge)
				} else if lowerFirstWord == "passthrough" {
					coloredFirstWord = color.Apply(firstWord, ocRouteTLSColorPassthrough)
				} else if lowerFirstWord == "reencrypt" {
					coloredFirstWord = color.Apply(firstWord, ocRouteTLSColorReencrypt)
				} else {
					// If no specific keyword, color the first word with the general value color determined earlier
					coloredFirstWord = color.Apply(firstWord, effectiveValColor)
				}
				finalTLSOutput = coloredFirstWord
				if remainingText != "" {
					finalTLSOutput += color.Apply(remainingText, getColorByValueType(remainingText, dp.DarkBackground))
				}
				fmt.Fprintf(w, "%s%s%s%s\n", indent, coloredKeyOutput, toSpaces(spacesCnt), finalTLSOutput)
				continue
			}
		}

		// Default key-value printing if no special route handling took over or if not a route
		fmt.Fprintf(w, "%s%s%s%s\n", indent, coloredKeyOutput, toSpaces(spacesCnt), color.Apply(valueOutput, effectiveValColor))
	}
}
