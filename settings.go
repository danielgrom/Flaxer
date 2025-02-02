package main

import (
	"flaxer/repository"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func (app *Config) CreateSettingsTab(settings repository.ProjectSettings) *fyne.Container {
	title := widget.NewRichText()
	title.ParseMarkdown("# Settings")
	separator := widget.NewSeparator()
	projLocationEntry := widget.NewEntry()
	savedSettings := GetSavedFlaxerSettings(app)
	projLocationEntry.SetText(savedSettings.ProjectsDirectory)
	projLocationEntry.PlaceHolder = "Location of Flax Projects"
	projLocation := widget.NewFormItem("Flax Projects Directory...", projLocationEntry)
	projdialog := dialog.NewFolderOpen(func(file fyne.ListableURI, err error) {
		if file != nil {
			savedSettings = GetSavedFlaxerSettings(app)
			projLocationEntry.SetText(file.Path())
			savedSettings.ProjectsDirectory = file.Path()
			SaveFlaxerSettings(app, savedSettings)
		}
	}, app.MainWindow)
	projButton := widget.NewButton("Browse Projects Directory...", func() { projdialog.Show() })
	projFormButton := widget.NewFormItem("", projButton)

	flaxLocationEntry := widget.NewEntry()
	flaxLocationEntry.Text = savedSettings.FlaxLocation
	flaxLocation := widget.NewFormItem("Flax Executable Location", flaxLocationEntry)
	flaxdialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if file != nil {
			savedSettings = GetSavedFlaxerSettings(app)
			flaxLocationEntry.SetText(file.URI().Path())
			savedSettings.FlaxLocation = file.URI().Path()
			SaveFlaxerSettings(app, savedSettings)
		}
	}, app.MainWindow)
	extensions := []string{".exe", ".bin", ""}
	flaxdialog.SetFilter(storage.NewExtensionFileFilter(extensions))
	flaxButton := widget.NewButton("Browse Flax Executable", func() { flaxdialog.Show() })
	flaxForButton := widget.NewFormItem("", flaxButton)

	form := widget.NewForm(projLocation, projFormButton, flaxLocation, flaxForButton)

	ret := container.NewVBox(title, separator, form)
	return ret
}

func GetSavedFlaxerSettings(app *Config) *repository.FlaxerSettings {
	settings, err := app.DB.GetFlaxerSettings()
	if err != nil {
		dialog.ShowError(err, app.MainWindow)
	}

	if settings == nil {
		settings := repository.FlaxerSettings{ID: 0, ProjectsDirectory: "", FlaxLocation: ""}
		app.DB.InsertFlaxerSettings(settings)
	}
	return settings
}

func SaveFlaxerSettings(app *Config, fs *repository.FlaxerSettings) {
	app.DB.UpdateFlaxerSettings(fs.ID, *fs)
}
