package exporter

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestExporter(t *testing.T) {
	t.Run("skips documents with only transcript and no notes", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Doc with transcript only", CreatedAt: "2026-01-21T10:00:00Z"},
				"doc2": {ID: "doc2", Title: "Doc without anything", CreatedAt: "2026-01-21T10:00:00Z"},
			},
			Transcripts: map[string][]TranscriptEntry{
				"doc1": {{Text: "Hello", Source: "microphone"}},
			},
		}

		result, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Written != 0 {
			t.Errorf("Expected 0 written (transcript-only docs have no notes to export), got %d", result.Written)
		}
	})

	t.Run("filters documents with notes > 10 chars", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Doc with notes", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "This is a long enough note to export"},
				"doc2": {ID: "doc2", Title: "Doc with short notes", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "Short"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		result, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Written != 1 {
			t.Errorf("Expected 1 written, got %d", result.Written)
		}
	})

	t.Run("exports shared documents", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "My Meeting", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "# My notes here"},
			},
			SharedDocuments: map[string]Document{
				"doc2": {ID: "doc2", Title: "Shared Meeting", CreatedAt: "2026-01-21T11:00:00Z", NotesMarkdown: "# Shared notes here"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		result, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Written != 2 {
			t.Errorf("Expected 2 written (1 owned + 1 shared), got %d", result.Written)
		}

		// Verify shared document file was created
		sharedPath := filepath.Join(tmpDir, "2026-01-21_Shared Meeting.md")
		content, err := os.ReadFile(sharedPath)
		if err != nil {
			t.Fatalf("Failed to read shared document file: %v", err)
		}
		if !strings.Contains(string(content), "Shared notes here") {
			t.Error("Expected shared document notes in output")
		}
	})

	t.Run("skips documents with neither notes nor transcript", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Empty doc", CreatedAt: "2026-01-21T10:00:00Z"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		result, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Written != 0 {
			t.Errorf("Expected 0 written, got %d", result.Written)
		}
	})

	t.Run("skips writing unchanged files", func(t *testing.T) {
		tmpDir := t.TempDir()
		exp := NewExporter(tmpDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Test", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "Some notes here"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		// First export
		result1, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result1.Written != 1 {
			t.Errorf("First export: expected 1 written, got %d", result1.Written)
		}

		// Second export with same content
		result2, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result2.Skipped != 1 {
			t.Errorf("Second export: expected 1 skipped, got %d", result2.Skipped)
		}
		if result2.Written != 0 {
			t.Errorf("Second export: expected 0 written, got %d", result2.Written)
		}
	})

	t.Run("creates output directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputDir := filepath.Join(tmpDir, "nested", "output", "dir")
		exp := NewExporter(outputDir)

		state := &CacheState{
			Documents: map[string]Document{
				"doc1": {ID: "doc1", Title: "Test", CreatedAt: "2026-01-21T10:00:00Z", NotesMarkdown: "Some notes here"},
			},
			Transcripts: map[string][]TranscriptEntry{},
		}

		_, err := exp.Export(state, false)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Verify directory was created
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			t.Error("Output directory was not created")
		}
	})
}

func TestDefaultOutputDir(t *testing.T) {
	dir := DefaultOutputDir()
	if runtime.GOOS == "windows" {
		if !strings.Contains(dir, "00-inbox") {
			t.Errorf("Unexpected default output dir on Windows: %s", dir)
		}
	} else {
		if !strings.Contains(dir, ".local") || !strings.Contains(dir, "granola-transcripts") {
			t.Errorf("Unexpected default output dir: %s", dir)
		}
	}
}
