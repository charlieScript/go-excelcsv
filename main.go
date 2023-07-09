package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
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
	myWindow.Resize(fyne.NewSize(800, 500))

	file, _ := os.ReadFile("icon.jpeg")

	pic := fyne.NewStaticResource("name", file)

	myApp.SetIcon(pic)
	myWindow.SetIcon(pic)

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

	textArea := widget.NewEntry()
	textArea.MultiLine = true

	scroll := container.NewVScroll(textArea)
	scroll.SetMinSize(fyne.NewSize(200, 400))

	showError := func() {
		alertContent := container.NewVBox()
		dialog.ShowCustom(errorMessage, "OK", alertContent, myWindow)
	}

	downloadBtn := widget.NewButton("Download JSON", func() {

	})

	downloadBtn.Disable()

	// Create a button to display the selected option
	showButton := widget.NewButton("Open file", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				var fileExt string
				filters := []string{"csv", "xlsx"}
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

				var data []byte

				switch fileExt {
				case "csv":
					data = converters.ConvertCSVToJSON(reader.URI().Path(), hasHeaders)
				case "excel":
					data = converters.ConvertExcelToJSON(reader.URI().Path(), hasHeaders)
				}

				var formattedJSON bytes.Buffer
				err := json.Indent(&formattedJSON, []byte(data), "", "  ")
				if err != nil {
					// Handle error
					errorMessage = "File selection canceled or error occurred:"
					showError()
					return
				}

				textArea.SetText(formattedJSON.String())
				fmt.Print(data)
				downloadBtn.Enable()
				defer reader.Close()

			} else {
				errorMessage = "File selection canceled or error occurred:"
				showError()
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
		scroll,
		downloadBtn,
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
