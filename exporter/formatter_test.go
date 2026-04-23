package exporter

import (
	"strings"
	"testing"
)

func TestFormatDocumentMarkdown(t *testing.T) {
	t.Run("document with notes", func(t *testing.T) {
		doc := &Document{
			ID:        "abc12345-1234-5678-9abc-def012345678",
			Title:     "Engineering Team Stand-Up",
			CreatedAt: "2026-01-21T20:30:01.410Z",
		}
		notes := "# Action Items\n\n- Follow up on project timeline"

		result := FormatDocumentMarkdown(doc, notes)

		if !strings.Contains(result, "# Engineering Team Stand-Up") {
			t.Error("Expected title to be in output")
		}
		if !strings.Contains(result, "Date: 2026-01-21 20:30") {
			t.Error("Expected formatted date in output")
		}
		if !strings.Contains(result, "Meeting ID: abc12345-1234-5678-9abc-def012345678") {
			t.Error("Expected meeting ID in output")
		}
		if !strings.Contains(result, "## AI-Generated Notes") {
			t.Error("Expected AI-Generated Notes section")
		}
		if !strings.Contains(result, "Follow up on project timeline") {
			t.Error("Expected notes content in output")
		}
		if strings.Contains(result, "## Transcript") {
			t.Error("Should not have transcript section")
		}
	})

	t.Run("document with no notes", func(t *testing.T) {
		doc := &Document{
			ID:        "test-id",
			Title:     "Test Meeting",
			CreatedAt: "2026-01-21T10:00:00Z",
		}

		result := FormatDocumentMarkdown(doc, "")

		if strings.Contains(result, "## AI-Generated Notes") {
			t.Error("Should not have notes section when empty")
		}
	})

	t.Run("handles missing created_at", func(t *testing.T) {
		doc := &Document{
			ID:    "test",
			Title: "Test",
		}

		result := FormatDocumentMarkdown(doc, "")

		if !strings.Contains(result, "Date: Unknown date") {
			t.Error("Expected 'Unknown date' for missing created_at")
		}
	})

	t.Run("handles empty title", func(t *testing.T) {
		doc := &Document{
			ID:        "test",
			CreatedAt: "2026-01-21T10:00:00Z",
		}

		result := FormatDocumentMarkdown(doc, "")

		if !strings.Contains(result, "# Untitled") {
			t.Error("Expected 'Untitled' for empty title")
		}
	})
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		expected  string
	}{
		{
			name:      "ISO8601 with Z suffix",
			timestamp: "2026-01-21T20:30:01.410Z",
			expected:  "2026-01-21 20:30",
		},
		{
			name:      "ISO8601 without milliseconds",
			timestamp: "2026-01-21T20:30:01Z",
			expected:  "2026-01-21 20:30",
		},
		{
			name:      "empty string",
			timestamp: "",
			expected:  "Unknown date",
		},
		{
			name:      "malformed date",
			timestamp: "not-a-date",
			expected:  "Unknown date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDate(tt.timestamp)
			if result != tt.expected {
				t.Errorf("FormatDate(%q) = %q, want %q", tt.timestamp, result, tt.expected)
			}
		})
	}
}

func TestFormatDateForFilename(t *testing.T) {
	tests := []struct {
		name      string
		timestamp string
		expected  string
	}{
		{
			name:      "valid timestamp",
			timestamp: "2026-01-21T20:30:01.410Z",
			expected:  "2026-01-21",
		},
		{
			name:      "empty string",
			timestamp: "",
			expected:  "unknown-date",
		},
		{
			name:      "malformed",
			timestamp: "invalid",
			expected:  "unknown-date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDateForFilename(tt.timestamp)
			if result != tt.expected {
				t.Errorf("FormatDateForFilename(%q) = %q, want %q", tt.timestamp, result, tt.expected)
			}
		})
	}
}

func TestSourceToSpeaker(t *testing.T) {
	tests := []struct {
		source   string
		expected string
	}{
		{"microphone", "Me"},
		{"system", "Them"},
		{"speaker1", "Speaker1"},
		{"", "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			result := SourceToSpeaker(tt.source)
			if result != tt.expected {
				t.Errorf("SourceToSpeaker(%q) = %q, want %q", tt.source, result, tt.expected)
			}
		})
	}
}

func TestNumberWithCommas(t *testing.T) {
	tests := []struct {
		n        int
		expected string
	}{
		{0, "0"},
		{100, "100"},
		{1000, "1,000"},
		{10000, "10,000"},
		{100000, "100,000"},
		{1000000, "1,000,000"},
		{1234567, "1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := NumberWithCommas(tt.n)
			if result != tt.expected {
				t.Errorf("NumberWithCommas(%d) = %q, want %q", tt.n, result, tt.expected)
			}
		})
	}
}
