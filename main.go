// A small Fyne GUI for obfuscating/deobfuscating TIBCO "#!" fixed-key secrets.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const appTitle = "TiBCO Obfuscation Tool"

func main() {
	a := app.New()
	w := a.NewWindow(appTitle)
	w.SetContent(makeUI(w))
	w.Resize(fyne.NewSize(640, 560))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func makeUI(w fyne.Window) fyne.CanvasObject {
	title := canvas.NewText(appTitle, theme.PrimaryColor())
	title.TextSize = 30
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	link := widget.NewHyperlink(
		"TIBCO fixed-key (#!) obfuscation — reversible by design",
		nil,
	)
	link.Alignment = fyne.TextAlignCenter

	input := widget.NewMultiLineEntry()
	input.Wrapping = fyne.TextWrapBreak
	input.SetPlaceHolder("Plaintext password, or a #! obfuscated secret")

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapBreak
	output.SetPlaceHolder("Result")

	setResult := func(s string, err error) {
		if err != nil {
			output.SetText("error: " + err.Error())
			return
		}
		output.SetText(s)
	}

	obfBtn := widget.NewButtonWithIcon("Obfuscate", theme.MediaSkipNextIcon(), func() {
		setResult(Obfuscate(input.Text))
	})
	obfBtn.Importance = widget.HighImportance

	clearBtn := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
		input.SetText("")
		output.SetText("")
	})

	deobfBtn := widget.NewButtonWithIcon("Deobfuscate", theme.MediaSkipPreviousIcon(), func() {
		setResult(Deobfuscate(input.Text))
	})
	deobfBtn.Importance = widget.HighImportance

	copyBtn := widget.NewButtonWithIcon("Copy Result", theme.ContentCopyIcon(), func() {
		if output.Text != "" {
			w.Clipboard().SetContent(output.Text)
		}
	})
	copyBtn.Importance = widget.WarningImportance

	buttonRow := container.NewGridWithColumns(3, obfBtn, clearBtn, deobfBtn)
	controls := container.NewVBox(buttonRow, copyBtn)

	topPanel := container.NewBorder(nil, controls, nil, nil, input)
	center := container.NewGridWithRows(2, topPanel, output)

	return container.NewBorder(
		container.NewPadded(title),
		container.NewPadded(link),
		nil, nil,
		container.NewPadded(center),
	)
}
