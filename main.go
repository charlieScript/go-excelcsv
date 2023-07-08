package main

import (
	"fmt"
	"image/color"
	"strings"

	// "strings"

	// "strings"

	// "io/fs"
	// "io/ioutil"
	// "os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/charliescript/go-excel/converters"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Gbedu")
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.Title()

	var selectedOption string
	var hasHeaders bool = false
	var errorMessage string

	// Create a dropdown select widget
	selectItems := []string{"excel", "csv"}
	dropdown := widget.NewSelect(selectItems, func(selected string) {
		selectedOption = selected
	})

	checkbox := widget.NewCheck("", func(checked bool) {
		hasHeaders = checked
	})

	textArea := widget.NewMultiLineEntry()
	textArea.MultiLine = true
	textArea.Resize(fyne.NewSize(600, 900))

	showError := func() {
		alertContent := container.NewVBox()
		dialog.ShowCustom(errorMessage, "OK", alertContent, myWindow)
	}

	// Create a button to display the selected option
	showButton := widget.NewButton("Open file", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				var fileExt string
				filters := []string{".csv", ".xlxs"}
				fileExt = reader.URI().Extension()

				fileExt = strings.Split(fileExt, ".")[1]

				fmt.Println(fileExt)
				fmt.Println("Comparsion", Contains(filters, fileExt))

				if !Contains(filters, fileExt) {
					errorMessage = "Invalid file Extension"
					showError()
					return
				}

				if selectedOption == "csv" && fileExt == ".csv" {
					fileExt = "csv"
				} else {
					fileExt = "excel"
				}

				fmt.Println(selectedOption, fileExt)

				if selectedOption != fileExt {
					errorMessage = "Invalid file format"
					showError()
					return
				}

				// if !strings.Contains(selectedOption, fileExt) {
				// 	showError()
				// 	return
				// }

				var data []byte

				switch selectedOption {
				case "csv":
					data = converters.ConvertCSVToJSON(reader.URI().Path(), hasHeaders)
					break
				case "excel":
					data = converters.ConvertExcelToJSON(reader.URI().Path(), hasHeaders)
					break
				}
				textArea.SetText(string(data))

				defer reader.Close()

			} else {
				fmt.Println("File selection canceled or error occurred:", err)
			}
		}, myWindow)
	})

	text1 := canvas.NewText("File Format", color.White)
	text2 := canvas.NewText("Header (Set if first row it heder)", color.White)
	grid := container.New(layout.NewGridLayout(2), text1, dropdown, text2, checkbox)

	content := container.NewVBox(
		grid,
		showButton,
		textArea,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
