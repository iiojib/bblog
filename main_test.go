package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestSanitizePrefix_StripsNewlinesAndAnsiWhenEnabled(t *testing.T) {
	got := sanitizePrefix("api\n\x1b[31mred\x1b[0m", true)
	want := "api red"

	if got != want {
		t.Fatalf("sanitizePrefix() = %q, want %q", got, want)
	}
}

func TestSanitizePrefix_KeepsAnsiWhenStripDisabled(t *testing.T) {
	got := sanitizePrefix("api\x1b[31mred\x1b[0m", false)

	if !strings.Contains(got, "\x1b[31m") {
		t.Fatalf("sanitizePrefix() should keep ANSI escapes when strip disabled, got %q", got)
	}
}

func TestTranslateLine_StripEscapeNoTimestamp(t *testing.T) {
	got := translateLine("\x1b[31mERR\x1b[0m", "api", true, true)
	want := "[api] ERR\n"

	if got != want {
		t.Fatalf("translateLine() = %q, want %q", got, want)
	}
}

func TestTranslateLine_StripEscapeWithTimestamp(t *testing.T) {
	got := translateLine("ok", "", false, true)
	matched, err := regexp.MatchString(`^\[[0-9]{2}:[0-9]{2}:[0-9]{2}\] ok\n$`, got)

	if err != nil {
		t.Fatalf("regex failed: %v", err)
	}

	if !matched {
		t.Fatalf("translateLine() format mismatch, got %q", got)
	}
}

func TestTranslateLine_NoStripReturnsStyledPayload(t *testing.T) {
	got := translateLine("plain", "", true, false)

	if !strings.HasPrefix(got, "%cplain") {
		t.Fatalf("expected styled payload prefix, got %q", got)
	}

	if !strings.Contains(got, "\n") {
		t.Fatalf("expected payload to contain format/style separator newline, got %q", got)
	}
}
