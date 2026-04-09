# Granary

Exports meeting notes and transcripts from [Granola](https://www.granola.so)'s local cache to markdown files.

Granary exports AI-generated notes and full transcripts to markdown. It auto-detects the latest Granola cache version, only writes changed files, and preserves transcripts even after Granola purges them from its cache. A built-in macOS LaunchAgent can run exports automatically every 2 hours.

## 🛠️ Installation

### Homebrew

```bash
brew install wassimk/tap/granary
```

### From source

```bash
go install github.com/wassimk/granary@latest
```

## 💻 Usage

### Export meeting notes

```bash
granary run
```

By default, Granary reads from `~/Library/Application Support/Granola/cache-v*.json` and exports markdown files to `~/.local/share/granola-transcripts/`. Each file is named `YYYY-MM-DD_Meeting_Title.md`.

#### Options

```
-o, --output-dir   Custom output directory (default: ~/.local/share/granola-transcripts)
```

### Background service (LaunchAgent)

Install a macOS LaunchAgent that automatically exports every 2 hours:

```bash
granary install
```

Check the service status:

```bash
granary status
```

Remove the background service:

```bash
granary uninstall
```

### Other commands

```bash
granary version    # Show version
granary help       # Show help
```

## ⚠️ Transcript availability

Granola does not keep all transcripts in its local cache. Transcripts are fetched from Granola's servers on demand when you open a meeting, and older ones are periodically purged. New meetings will have transcripts in cache after you view them in Granola, but previously viewed meetings may not.

Once Granary exports a transcript, it preserves it permanently. On future runs it merges the latest AI notes with any previously exported transcript, so you never lose data.

## 📄 Output format

```markdown
# Meeting Title
Date: 2025-01-24 14:30
Meeting ID: abc-123

---

## AI-Generated Notes

[Granola's AI-generated meeting notes and summaries]

---

## Transcript

**Me:** [Your words]

**Them:** [Other participant's words]
```

## 📝 Disclaimer

This project is not affiliated with, endorsed by, or connected to [Granola](https://www.granola.so) in any way. I love Granola and use it every day. This is just a personal utility to export my meeting data.
