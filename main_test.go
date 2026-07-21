package main

import (
	"strings"
	"testing"
)

func TestTranslateLine_StripEscape(t *testing.T) {
	got := translateLine("\x1b[31mERR\x1b[0m", true)
	want := "ERR\n"

	if got != want {
		t.Fatalf("translateLine() = %q, want %q", got, want)
	}
}

func TestTranslateLine_NoStripReturnsStyledPayload(t *testing.T) {
	got := translateLine("plain", false)

	if !strings.HasPrefix(got, "%cplain") {
		t.Fatalf("expected styled payload prefix, got %q", got)
	}

	if !strings.Contains(got, "\n") {
		t.Fatalf("expected payload to contain format/style separator newline, got %q", got)
	}
}
