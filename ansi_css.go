package main

import (
	"regexp"
	"strconv"
	"strings"
)

var ansiEscapeRE = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)
var ansiSgrRE = regexp.MustCompile(`\x1b\[([0-9;]*)m`)

type sgrState struct {
	fg        string
	bg        string
	bold      bool
	dim       bool
	italic    bool
	underline bool
	strike    bool
}

type styledSegment struct {
	text  string
	style string
}

func stripAnsi(text string) string {
	return ansiEscapeRE.ReplaceAllString(text, "")
}

func defaultSgrState() sgrState {
	return sgrState{}
}

func cssFromSgrState(s sgrState) string {
	parts := make([]string, 0, 6)

	if s.fg != "" {
		parts = append(parts, "color:"+s.fg)
	}
	if s.bg != "" {
		parts = append(parts, "background-color:"+s.bg)
	}
	if s.bold {
		parts = append(parts, "font-weight:700")
	}
	if s.italic {
		parts = append(parts, "font-style:italic")
	}
	if s.dim {
		parts = append(parts, "opacity:0.7")
	}

	deco := make([]string, 0, 2)

	if s.underline {
		deco = append(deco, "underline")
	}
	if s.strike {
		deco = append(deco, "line-through")
	}
	if len(deco) > 0 {
		parts = append(parts, "text-decoration:"+strings.Join(deco, " "))
	}

	return strings.Join(parts, ";")
}

func ansiFgColor(code int) string {
	switch code {
	case 30:
		return "#000000"
	case 31:
		return "#cc0000"
	case 32:
		return "#00aa00"
	case 33:
		return "#aa7700"
	case 34:
		return "#1e5eff"
	case 35:
		return "#b000b0"
	case 36:
		return "#009999"
	case 37:
		return "#cccccc"
	case 90:
		return "#666666"
	case 91:
		return "#ff4d4d"
	case 92:
		return "#33cc66"
	case 93:
		return "#ffcc33"
	case 94:
		return "#66a3ff"
	case 95:
		return "#e066ff"
	case 96:
		return "#33cccc"
	case 97:
		return "#ffffff"
	}
	return ""
}

func ansiBgColor(code int) string {
	switch code {
	case 40:
		return "#000000"
	case 41:
		return "#cc0000"
	case 42:
		return "#00aa00"
	case 43:
		return "#aa7700"
	case 44:
		return "#1e5eff"
	case 45:
		return "#b000b0"
	case 46:
		return "#009999"
	case 47:
		return "#cccccc"
	case 100:
		return "#666666"
	case 101:
		return "#ff4d4d"
	case 102:
		return "#33cc66"
	case 103:
		return "#ffcc33"
	case 104:
		return "#66a3ff"
	case 105:
		return "#e066ff"
	case 106:
		return "#33cccc"
	case 107:
		return "#ffffff"
	}
	return ""
}

func applySgrCode(state *sgrState, code int) {
	switch code {
	case 0:
		*state = defaultSgrState()
	case 1:
		state.bold = true
	case 2:
		state.dim = true
	case 3:
		state.italic = true
	case 4:
		state.underline = true
	case 9:
		state.strike = true
	case 22:
		state.bold = false
		state.dim = false
	case 23:
		state.italic = false
	case 24:
		state.underline = false
	case 29:
		state.strike = false
	case 39:
		state.fg = ""
	case 49:
		state.bg = ""
	default:
		if color := ansiFgColor(code); color != "" {
			state.fg = color
			return
		}

		if color := ansiBgColor(code); color != "" {
			state.bg = color
		}
	}
}

func appendSegment(segments []styledSegment, text, style string) []styledSegment {
	if text == "" {
		return segments
	}

	if len(segments) > 0 && segments[len(segments)-1].style == style {
		segments[len(segments)-1].text += text
		return segments
	}

	return append(segments, styledSegment{text: text, style: style})
}

func textToStyledSegments(text string) []styledSegment {
	segments := make([]styledSegment, 0, 4)
	state := defaultSgrState()
	last := 0
	matches := ansiSgrRE.FindAllStringSubmatchIndex(text, -1)

	for _, m := range matches {
		if m[0] > last {
			segments = appendSegment(segments, text[last:m[0]], cssFromSgrState(state))
		}

		codesRaw := ""
		if m[2] >= 0 && m[3] >= 0 {
			codesRaw = text[m[2]:m[3]]
		}

		if codesRaw == "" {
			applySgrCode(&state, 0)
		} else {
			for _, part := range strings.Split(codesRaw, ";") {
				code := 0

				if part != "" {
					parsed, err := strconv.Atoi(part)
					if err == nil {
						code = parsed
					}
				}

				applySgrCode(&state, code)
			}
		}

		last = m[1]
	}

	if last < len(text) {
		segments = appendSegment(segments, text[last:], cssFromSgrState(state))
	}

	return segments
}

func segmentsToPayload(segments []styledSegment) string {
	if len(segments) == 0 {
		return ""
	}

	formatParts := make([]string, 0, len(segments))
	styles := make([]string, 0, len(segments))

	for _, seg := range segments {
		if seg.text == "" {
			continue
		}

		escapedText := strings.ReplaceAll(seg.text, "%", "%%")
		formatParts = append(formatParts, "%c"+escapedText)
		styles = append(styles, seg.style)
	}

	if len(styles) == 0 {
		return ""
	}

	return strings.Join(formatParts, "") + "\n" + strings.Join(styles, "\n")
}
