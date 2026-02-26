# ufWall ‑ Your UFW Firewall, but a little prettier

> *“Because why would you *just* use `sudo ufw status` when you can stare at a very pretty little terminal program that tells you the same thing in a more dramatic way?”*


## What the Heck Is This? (Short Answer)

`ufWall` is a **Tiny, Go‑powered Bubbles** / **Lipgloss** based terminal UI to *manage* the Ubuntu [UFW](https://help.ubuntu.com/community/UFW) firewall.

> The only difference between this and `ufw` is that *you* feel the *warm glow of control*.


## Why You Might Need It

- You hate doing `sudo ufw status` and then `sudo ufw allow 22` *inline*.
- You like having a visual list of rules that you can scroll through, delete with a single key.
- You want a little more color and structure before you risk breaking the internet.

If you’re content with the vanilla UFW CLI, this is *unnecessary* – but if you want to add a small flourish to your terminal routine, give it a spin.

## 📦 Prerequisites

- ufw installed
- Go 1.20+ (for building from source)
- A terminal that supports ANSI colors (most modern terminals do)
- Sudo privileges to manage the firewall rules

## 🚀 Quick Start

```bash
git clone https://github.com/The-Robin-Hood/ufWall.git
cd ufWall
go build ./cmd/ufWall
sudo ./ufWall
```

You can also install it with `go install`:
```bash
go install ./cmd/ufWall@latest
sudo ufWall    
```


## 📁 Project Structure

```
ufWall/
├─ internal/ufw            # wrappers around `sudo ufw` commands and helpers
├─ internal/sections       # logic for each UI section (stats, policy, rules)
├─ internal/ui             # styling & helpers for lipgloss
├─ internal/app            # Bubbletea model, update loop, view rendering
├─ cmd/ufWall/main.go      # tiny entry point
└─ go.mod                  # Go dependencies
```

> Nothing fancy: pure Go, a handful of third‑party Bubbles & Lipgloss.

## 📜 TL;DR

> **`ufWall`** is literally a **fancy wrapper** around `ufw`. It adds colors, tables, and a *little* interactive feel so that you *feel* better while you *break* your own firewall. If you’re on a machine that already has UFW and you’re comfortable with the CLI, you’re probably fine. 
Otherwise, enjoy the new UI and *don’t forget* to `sudo`.
---
