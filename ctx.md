Startup flow (with Bubble Tea + fzf):

---

## On start (bootstrap)

1. Detect system package manager:

* try `pacman`, `apt`, `yay`, `rpm` in order
* first found becomes active provider

---

2. Immediately open fzf menu (no TUI list yet):

Options fed into fzf:

```
Explore installed
Install new
Update system
Select provider
Quit
```

---

## fzf selection result routing

### Explore installed

* enter Bubble Tea main list view
* load `ListInstalled()`

### Install new

* run fzf again:

  * input: `ListAllPackages()`
* selected package → `Install(pkg)`
* refresh installed list

### Update system

* run `Update()` on active manager
* refresh installed list

### Select provider

* fzf list:

```
apt
pacman
yay
rpm
```

* switch manager
* reload installed list

### Quit

* exit program

---

## Bubble Tea structure after startup

* state 1: bootstrap (fzf menu only, no UI yet)
* state 2: main list (installed packages)
* state 3: modal actions (install/update/remove)

---

## Key rule

* fzf is only used for:

  * first action menu
  * install selection
  * provider selection
* Bubble Tea handles everything else (state + navigation)

---

## Minimal flow

```
start
  ↓
detect PM
  ↓
fzf main menu
  ↓
┌──────────────┐
│ installed    │ → Bubble Tea list
│ install      │ → fzf package list
│ update       │ → run update
│ provider     │ → fzf provider list
└──────────────┘
```
