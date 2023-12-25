package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

func CreateNewEntry() (*widget.Entry, binding.String) {
	str := binding.NewString()
	entry := widget.NewEntryWithData(str)

	return entry, str
}

type Row struct {
	Date time.Time `json:"date"`
	Name string    `json:"name"`
}

var debugLog = widget.NewMultiLineEntry()

func debug(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	debugLog.Append(fmt.Sprintf("%s: %s\n", time.Now().Format(time.DateTime), message))
	log.Printf(message)
}

var URL = "http://example.com"

func main() {
	a := app.New()
	w := a.NewWindow("Report Builder")

	debugLog.Disable()
	debug("Starting application")

	entry1, startDate := CreateNewEntry()
	entry2, endDate := CreateNewEntry()

	details := binding.NewSprintf("Report will be created from %s to %s.", startDate, endDate)
	detailsLabel := widget.NewLabelWithData(details)

	now := time.Now().Format(time.DateOnly)
	entry1.SetText(now)
	entry2.SetText(now)
	var selected_company string
	company_select := widget.NewSelect([]string{"lizing", "ravis"}, func(s string) {
		selected_company = s
	})

	form := widget.NewForm()
	form.Append("Start Date: ", entry1)
	form.Append("End Date: ", entry2)
	form.Append("Company: ", company_select)

	var data []Row
	var modal *widget.PopUp

	table := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			name := data[i].Name
			date := data[i].Date.Format(time.DateTime)
			o.(*widget.Label).SetText(fmt.Sprintf("%d\t\t%s\t\t%s", i, date, name))
		},
	)

	exportBtn := widget.NewButton("Export to Excel", func() {
		form := widget.NewForm()
		fileName := widget.NewEntry()
		form.Append("File: ", fileName)
		content := container.NewVBox(
			layout.NewSpacer(),
			form,
			layout.NewSpacer(),
			widget.NewButton("Export", func() {
				log.Printf("exporting...")
				debug("Exporting to excel...")
				f := excelize.NewFile()
				defer func() {
					if err := f.Close(); err != nil {
						debug("Cannot export to file")
					}
				}()

				index, err := f.NewSheet("Report")
				if err != nil {
					debug("Error creating new excel sheet")
					return
				}

				if len(data) == 0 {
					debug("Nothing to export")
					return
				}
				rows, _ := f.GetRows("Report")
				debug("ROWS: %#v", rows)
				for i, row := range data {
					dateCell := fmt.Sprintf("A%d", i)
					nameCell := fmt.Sprintf("B%d", i)
					debug("Setting to %s %s", dateCell, nameCell)
					f.SetCellStr("Report", dateCell, row.Date.Format(time.DateTime))
					f.SetCellStr("Report", nameCell, row.Name)
				}

				f.SetActiveSheet(index)

				if err := f.SaveAs(fileName.Text); err != nil {
					debug("Error while saving excel file")
				}

				modal.Hide()
			}),
			widget.NewButton("Cancel", func() {
				modal.Hide()
			}),
		)
		modal = widget.NewModalPopUp(content, w.Canvas())
		modal.Resize(fyne.NewSize(400, 200))

		modal.Show()
	})

	btn := widget.NewButton("Get Report", func() {
		if err := form.Validate(); err != nil {
			return
		}
		start := entry1.Text
		end := entry2.Text

		client := &http.Client{
			Timeout: time.Second * 3,
		}

		body, _ := json.Marshal(map[string]string{
			"start": start,
			"end":   end,
		})

		debug("selected: company = %s", selected_company)

		req, err := http.NewRequest(http.MethodGet, URL + selected_company, bytes.NewReader(body))
		if err != nil {
			debug("request: %s", err.Error())
		}

		req.Header.Set("passphrase", "passphrase")

		res, err := client.Do(req)
		if err != nil {
			debug("client: %s", err.Error())
		}

		defer res.Body.Close()

		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			debug("json: %s", err.Error())
		}

		if len(data) > 0 {
			exportBtn.Enable()
		}

		str, _ := json.MarshalIndent(data, "", "  ")
		debug("response: %s", string(str))
	})

	validator := func(value string) error {
		_, err := time.Parse(time.DateOnly, value)
		if err != nil {
			btn.Disable()
			return fmt.Errorf("Correct format is 'YYYY-MM-DD'")
		}
		return nil
	}
	entry1.Validator = validator
	entry2.Validator = validator

	form.SetOnValidationChanged(func(err error) {
		if err != nil {
			btn.Disable()
		} else {
			btn.Enable()
		}
	})

	btn.Disable()
	exportBtn.Disable()

	top := container.NewVBox()
	top.Add(form)
	top.Add(detailsLabel)
	top.Add(btn)
	top.Add(exportBtn)
	bottom := container.NewStack()
	bottom.Add(table)

	main := container.NewVSplit(top, bottom)
	debug := container.NewStack(debugLog)

	tabs := container.NewAppTabs(
		container.NewTabItem("Main", main),
		container.NewTabItem("Debug", debug),
	)

	w.SetContent(tabs)

	w.ShowAndRun()

}
