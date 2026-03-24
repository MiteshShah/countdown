# ⏳ Countdown Menubar

A native macOS menubar app written in **Go + CGo + Cocoa**.  
No external Go packages. No Python. No Electron.

Shows a live `Xd XXh XXm XXs` countdown ticking every second right in your menu bar.

---

## Requirements

| Tool | How to get it |
|------|---------------|
| **Go ≥ 1.21** | https://go.dev/dl/ |
| **Xcode Command Line Tools** | `xcode-select --install` |
| **macOS 12+** | – |

---

## Build & Install

```bash
chmod +x build.sh
./build.sh
```

The script compiles the binary and optionally copies it to `~/.local/bin`.

---

## Usage

```bash
# Basic — just a date
./build/countdown 2025-12-31

# With a label (shown in the dropdown menu)
./build/countdown 2025-12-31 "New Year 2026"

# Run in background (survives terminal close)
nohup ./build/countdown 2025-12-31 "New Year" > /dev/null 2>&1 &

# Stop it
pkill countdown
```

---

## What it looks like

**Menu bar** (updates every second):
```
⏳ 42d 07h 23m 55s
```

**Click to open dropdown**:
```
⏱  42 days  07 hrs  23 min  54 sec
🎯  New Year 2026  ·  Dec 31, 2025
─────────────────────────────────
Quit Countdown          ⌘Q
```

---

## Auto-launch on login (optional)

Create a launchd plist so it starts automatically:

```bash
# Replace the date and label for your use case
cat > ~/Library/LaunchAgents/com.user.countdown.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
  "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.user.countdown</string>
  <key>ProgramArguments</key>
  <array>
    <string>/Users/YOUR_USERNAME/.local/bin/countdown</string>
    <string>2025-12-31</string>
    <string>New Year 2026</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
  <key>KeepAlive</key>
  <false/>
</dict>
</plist>
EOF

# Load it
launchctl load ~/Library/LaunchAgents/com.user.countdown.plist
```

> Replace `YOUR_USERNAME` with your macOS username.

---

## How it works

- `main.go` uses **CGo** to call Cocoa's `NSStatusBar` API directly
- A Go goroutine ticks every second and calls `dispatch_async` on the main queue to update the title
- Zero dependencies beyond the Go standard library and macOS system frameworks
