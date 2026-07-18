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

func TestTextToStyledSegments_MergesAdjacentSameStyle(t *testing.T) {
	segments := textToStyledSegments("abc")

	if len(segments) != 1 {
		t.Fatalf("expected 1 merged segment, got %d", len(segments))
	}

	if segments[0].text != "abc" {
		t.Fatalf("unexpected merged text: %q", segments[0].text)
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
