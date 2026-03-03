# ufWall

> Because `sudo ufw status` is for people who hate themselves.

A terminal UI for UFW that makes firewall management feel less like defusing a bomb blindfolded.

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=flat)
![Platform](https://img.shields.io/badge/Platform-Linux-FCC624?style=flat&logo=linux)

## Why?

If you genuinely enjoy typing `sudo ufw status numbered` followed by `sudo ufw delete 7` while praying you counted correctly, this isn't for you. Go in peace.

## Install

```bash
git clone https://github.com/The-Robin-Hood/ufWall.git
cd ufWall
make build
sudo ./ufwall
```

Or if you trust the internet:
```bash
go install github.com/The-Robin-Hood/ufWall/cmd/ufWall@latest
sudo ufWall
```

## Project Structure

```
ufWall/
├── cmd/ufWall/         # Entry point
├── internal/
│   ├── app/            # Bubble Tea model, view, update loop
│   ├── sections/       # Stats, Policy, Rules : each with their own problems
│   ├── ufw/            # UFW command wrappers (the real hero)
│   ├── ui/             # Styles, menus, boxes — the pretty stuff
│   └── keys/           # Keybindings
└── go.mod              # Dependencies (Bubble Tea, Lip Gloss, regret)
```

## Requirements

- Linux with UFW installed (obviously)
- Go 1.20+ (to build)
- sudo privileges (we're managing a firewall, not a todo list)
- A terminal made after 1995

## FAQ

**Q: Will this break my firewall?**  
A: It runs the same `ufw` commands you'd run manually. If you break it, you were going to break it anyway. We just made it prettier.

**Q: Why not just use `ufw`?**  
A: Why drive a car when you can mass a horse? You do you.

**Q: Does this work on Mac/Windows?**  
A: UFW is a Linux thing. So no. Get a real OS. (Or use something else, I'm not your dad.)

**Q: I found a bug.**  
A: That's not a question, but open an issue anyway.

## License

MIT — Do whatever you want. Just don't blame me when you lock yourself out.

---

*Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) because terminals deserve to look good while you panic.*
