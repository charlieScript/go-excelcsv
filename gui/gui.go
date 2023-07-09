package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"

	"github.com/charliescript/go-excel/converters"
	"github.com/charliescript/go-excel/utils"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func RenderUI() {
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
	var fileExt string
	var jsonData []byte

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

	saveFileDialog := func() {
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				fmt.Println("Error while saving file:", err)
				return
			}

			if writer == nil {
				fmt.Println("Save file dialog was canceled")
				return
			}

			writer.Write(jsonData)
		}, myWindow)
		timestamp := time.Now().Format("20060102-150405")

		fileName := "data-" + timestamp + ".json"
		saveDialog.SetFileName(fileName)
		saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".xlxs", ".csv"}))
		saveDialog.Show()
		saveDialog.SetOnClosed(func() {
			errorMessage = "File save successfully"
			showError()
		})
	}

	downloadBtn := widget.NewButton("Download JSON", func() {
		saveFileDialog()

	})

	downloadBtn.Disable()

	showButton := widget.NewButton("Open file", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				filters := []string{"csv", "xlsx"}
				fileExt = reader.URI().Extension()

				fileExt = strings.Split(fileExt, ".")[1]

				fmt.Println(fileExt)
				fmt.Println("Comparsion", utils.Contains(filters, fileExt))

				if !utils.Contains(filters, fileExt) {
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

				switch fileExt {
				case "csv":
					jsonData = converters.ConvertCSVToJSON(reader.URI().Path(), hasHeaders)
				case "excel":
					jsonData = converters.ConvertExcelToJSON(reader.URI().Path(), hasHeaders)
				}

				var formattedJSON bytes.Buffer
				err := json.Indent(&formattedJSON, []byte(jsonData), "", "  ")
				if err != nil {
					// Handle error
					errorMessage = "File selection canceled or error occurred:"
					showError()
					return
				}

				textArea.SetText(formattedJSON.String())
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
