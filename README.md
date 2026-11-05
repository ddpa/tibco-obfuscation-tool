# TiBCO Obfuscation Tool (macOS)

A small Fyne GUI to obfuscate / deobfuscate TIBCO fixed-key (`#!`) secrets.

## Build the app

In **macOS Terminal** (not inside Claude):

```bash
cd ~/dev/tibco-gui
./build-macos.sh
```

This produces `TiBCO Obfuscation Tool.app` in this folder. Then:

```bash
open "TiBCO Obfuscation Tool.app"          # run it
cp -R "TiBCO Obfuscation Tool.app" /Applications/   # install to Launchpad
```

### Requirements (one-time)
- **Go** — https://go.dev/dl or `brew install go`
- **Xcode command-line tools** — `xcode-select --install` (Fyne uses CGO)

The build script checks for both and installs the `fyne` packager automatically.

## Alternatives
- `make run` — launch the GUI without packaging (quick test)
- `make universal` — build a single Intel + Apple-Silicon universal `.app`

## Note on the crypto
This is TIBCO's **fixed-key** obfuscation: the 3DES key ships inside TIBCO's own
libraries, so `#!` values are reversible by anyone and are effectively
plaintext. Use only on secrets from systems you administer. Secrets made with a
custom key or machine key cannot be recovered with this tool.

## Unsigned app
The `.app` is unsigned. On first launch, Gatekeeper may block it:
right-click the app → **Open** → **Open**, or clear quarantine:
`xattr -dr com.apple.quarantine "TiBCO Obfuscation Tool.app"`.
