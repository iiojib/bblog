package main

import (
	"strings"
	"testing"
)

func TestTextToStyledSegments_ParsesAnsiColor(t *testing.T) {
	segments := textToStyledSegments("\x1b[31mERR\x1b[0m")

	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}

	if segments[0].text != "ERR" {
		t.Fatalf("unexpected segment text: %q", segments[0].text)
	}

	if !strings.Contains(segments[0].style, "color:#cc0000") {
		t.Fatalf("expected red color style, got %q", segments[0].style)
	}
}

func TestTextToStyledSegments_StripsNonSgrAnsi(t *testing.T) {
	segments := textToStyledSegments("\x1b[31mERR\x1b[0m\x1b[2Kdone")

	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}

	if segments[0].text != "ERR" {
		t.Fatalf("unexpected first segment text: %q", segments[0].text)
	}

	if segments[1].text != "done" {
		t.Fatalf("unexpected second segment text: %q", segments[1].text)
	}
}

func TestTextToStyledSegments_MergesAdjacentSameStyle(t *testing.T) {
	segments := textToStyledSegments("abc")

	if len(segments) != 1 {
		t.Fatalf("expected 1 merged segment, got %d", len(segments))
	}

	if segments[0].text != "abc" {
		t.Fatalf("unexpected merged text: %q", segments[0].text)
	}
}

func TestTextToStyledSegments_ParsesEightBitAnsiColor(t *testing.T) {
	segments := textToStyledSegments("\x1b[38;5;196mred\x1b[0m")

	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}

	if segments[0].text != "red" {
		t.Fatalf("unexpected segment text: %q", segments[0].text)
	}

	if !strings.Contains(segments[0].style, "color:#ff0000") {
		t.Fatalf("expected 256-color style, got %q", segments[0].style)
	}
}

func TestTextToStyledSegments_ParsesBrightEightBitAnsiColor(t *testing.T) {
	segments := textToStyledSegments("\x1b[38;5;11myellow\x1b[0m")

	if len(segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(segments))
	}

	if !strings.Contains(segments[0].style, "color:#ffcc33") {
		t.Fatalf("expected bright yellow 256-color style, got %q", segments[0].style)
	}
}

func TestTextToStyledSegments_ParsesExtendedAnsiColors(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		wantStyle string
	}{
		{
			name:      "RGB foreground",
			text:      "\x1b[38;2;255;128;0morange\x1b[0m",
			wantStyle: "color:#ff8000",
		},
		{
			name:      "RGB background",
			text:      "\x1b[48;2;12;34;56mbackground\x1b[0m",
			wantStyle: "background-color:#0c2238",
		},
		{
			name:      "wrapped RGB values",
			text:      "\x1b[38;2;1140;290;512mwrapped\x1b[0m",
			wantStyle: "color:#742200",
		},
		{
			name:      "wrapped palette index",
			text:      "\x1b[38;5;452mwrapped\x1b[0m",
			wantStyle: "color:#ff0000",
		},
		{
			name:      "all empty RGB values",
			text:      "\x1b[38;2;;;mblack\x1b[0m",
			wantStyle: "color:#000000",
		},
		{
			name:      "individual empty RGB values",
			text:      "\x1b[38;2;;34;mgreen\x1b[0m",
			wantStyle: "color:#002200",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			segments := textToStyledSegments(test.text)

			if len(segments) != 1 {
				t.Fatalf("expected 1 segment, got %d", len(segments))
			}

			if !strings.Contains(segments[0].style, test.wantStyle) {
				t.Fatalf("expected style %q, got %q", test.wantStyle, segments[0].style)
			}
		})
	}
}

func TestSegmentsToPayload_ProducesFormatAndStyles(t *testing.T) {
	payload := segmentsToPayload([]styledSegment{
		{text: "A", style: "color:#ff0000"},
		{text: "B", style: ""},
	})

	parts := strings.Split(payload, "\n")

	if len(parts) < 2 {
		t.Fatalf("expected at least format + style lines, got %q", payload)
	}

	if parts[0] != "%cA%cB" {
		t.Fatalf("unexpected format line: %q", parts[0])
	}

	if parts[1] != "color:#ff0000" {
		t.Fatalf("unexpected first style line: %q", parts[1])
	}
}
