package view

import (
	"fmt"
	"net/url"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/data/validation"
	"github.com/KevinJGard/MusicDB/src/controller"
)

func Run_View() {
	controller := controller.NewController()
	var (
		progress *widget.ProgressBar
		progressContainer *fyne.Container
	)
	myApp := app.NewWithID("com.kevingard.musicdatabase")
	myWindow := myApp.NewWindow("Music Data Base")
	myWindow.SetIcon(theme.MediaMusicIcon())
	myWindow.Resize(fyne.NewSize(1000, 600))

	searchContainer := createSearchContainer(myWindow)
	cont, updateList := createListContainer(controller)
	progress = widget.NewProgressBar()
	loading := widget.NewLabel("Getting metadata...")
	loading.TextStyle = fyne.TextStyle{Monospace: true}
	progressContainer = container.NewVBox(loading, progress,)
	progressContainer.Hide()
	mineMetadata := func() {
		progress.SetValue(0) 
		progressContainer.Show() 

		go func() {
			err := controller.MineMetadata(
				func(pro int) {
					progress.SetValue(float64(pro) / 100.0)
					myWindow.Content().Refresh()
				}, 
				func() {
					dialog.ShowInformation("Completed", "Data was mined.", myWindow)
					progressContainer.Hide()
					updateList()
				},
			)

			if err != nil {
				dialog.ShowError(err, myWindow)
				progressContainer.Hide()
			}
		}()
	}

	menu := createMainMenu(myApp, myWindow, mineMetadata, controller)
	myWindow.SetMainMenu(menu)
	contSouth := createMusicControlContainer()

	content := container.New(layout.NewBorderLayout(searchContainer, contSouth, nil, nil),
		searchContainer,
		contSouth,
		cont,
		container.NewCenter(progressContainer),
	)
	
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func createSearchContainer(myWindow fyne.Window) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search...")

	searchButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		query := searchEntry.Text
		if query != "" {
			fmt.Println("Search:", query)
			searchEntry.SetText("")
		} else {
			err := errors.New("Enter your search.")
			dialog.ShowError(err, myWindow)
		}
	})

	return container.NewGridWithColumns(2, searchEntry, searchButton)
}

func openSettingsWindow(myApp fyne.App) {
	var editButton *widget.Button
	settingsWindow := myApp.NewWindow("Settings")
	settingsWindow.SetIcon(theme.SettingsIcon())
	settingsWindow.Resize(fyne.NewSize(600, 500))

	//form
	name := widget.NewEntry()
	name.SetPlaceHolder("Your username")
	name.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	email := widget.NewEntry()
	email.SetPlaceHolder("example@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")
	password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "password can only contain letters, numbers, '_', and '-'")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "UserName", Widget: name, HintText: "Your username"},
			{Text: "Email", Widget: email, HintText: "Your email address"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
	}
	editButton = widget.NewButtonWithIcon("Login", theme.AccountIcon() ,func() {
		form.Show()
		editButton.Hide()
	})

	form.OnSubmit = func() {
			fmt.Println("Form submitted")
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Music Data Base",
				Content: "Hello " + name.Text,
			})
			editButton.Show()
			form.Hide()
		}

	form.Append("Password", password)

	themes := createThemeButtons(myApp)
	quit := widget.NewButton("Close", func() {settingsWindow.Close()})

	labelSettings := widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	iconSettings := widget.NewIcon(theme.SettingsIcon())
	north := container.NewHBox(labelSettings, iconSettings)
	center := container.NewCenter(north)

	settingsContent := container.New(layout.NewBorderLayout(center, themes, quit, nil),
		center, themes, quit, form, editButton)

	settingsWindow.SetContent(settingsContent)
	settingsWindow.Show()
}

func createThemeButtons(myApp fyne.App) *fyne.Container {
	darkButton := widget.NewButton("Dark", func() {
		myApp.Settings().SetTheme(theme.DarkTheme())
	})
	lightButton := widget.NewButton("Light", func() {
		myApp.Settings().SetTheme(theme.LightTheme())
	})

	return container.NewGridWithColumns(2, darkButton, lightButton)
}

func createMainMenu(myApp fyne.App, myWindow fyne.Window, mineMetadata func(), controller *controller.Controller) *fyne.MainMenu {
	menuItemFull := fyne.NewMenuItem("Full screen", func() {
		myWindow.SetFullScreen(!myWindow.FullScreen())
	})
	menuItemFull.Icon = theme.ViewFullScreenIcon()

	menu := fyne.NewMenu("Screen", menuItemFull)

	menuItemSettings := fyne.NewMenuItem("Settings", func() {
		openSettingsWindow(myApp)
	})
	menuItemSettings.Icon = theme.SettingsIcon()
	menuItemHelp := fyne.NewMenuItem("Help", func() {
		url, _ := url.Parse("https://github.com/KevinJGard/MusicDB")
		_ = myApp.OpenURL(url)
	})
	menuItemHelp.Icon = theme.HelpIcon()

	newMenu2 := fyne.NewMenu("Options", menuItemSettings, menuItemHelp)

	menuItemSetPath := fyne.NewMenuItem("Set path", func() {
		setPath(myWindow, controller)
	})
	menuItemSetPath.Icon = theme.FolderIcon()
	menuItemMineMetadata := fyne.NewMenuItem("Mine metadata", mineMetadata)
	menuItemMineMetadata.Icon = theme.UploadIcon()

	newMenu3 := fyne.NewMenu("Miner", menuItemSetPath, menuItemMineMetadata)
	return fyne.NewMainMenu(menu, newMenu2, newMenu3)
}

func setPath(myWindow fyne.Window, controller *controller.Controller) {
	dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err == nil && uri != nil {
			fmt.Println("Selected path:", uri.Path())
			if err := controller.SetMusicDirectory(uri.Path()); err != nil {
				dialog.ShowError(err, myWindow)
			} else {
				dialog.ShowInformation("Path", "Selected path that you can mine.", myWindow)
			}
		} else {
			fmt.Println("Error selecting path:", err)
		}
	}, myWindow).Show()
}

func createListContainer(controller *controller.Controller) (*container.Split, func()) {
	data := make([]string, 0)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.MediaMusicIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)
	
	updateList := func() {
		songs, err := controller.GetSongs()
		if err == nil {
			data = data[:0]
			for _, song := range songs {
				data = append(data, song.Title)
			}
		} else {
			data = append(data, "Error loading songs")
		}
		list.Refresh()
	}

	icon := widget.NewIcon(theme.FileAudioIcon())
	label := widget.NewLabel("Select An Item From The List")
	label.TextStyle = fyne.TextStyle{Bold: true, Italic: true}
	hbox := container.NewHBox(icon, label)

	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id])
		icon.SetResource(theme.MediaMusicIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	updateList()

	return container.NewHSplit(list, container.NewCenter(hbox)), updateList
}

func createMusicControlContainer() *container.Split {
	music := widget.NewLabel("Your Music.")
	iconPlay := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { fmt.Println("Play.") })
	iconNext := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() { fmt.Println("Next.") })
	iconPrevious := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() { fmt.Println("Previous.") })
	iconStop := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() { fmt.Println("Stop.") })
	iconPlay.Disable()
	iconNext.Disable()
	iconPrevious.Disable()
	iconStop.Disable()

	contentIcons := container.NewHBox(iconPrevious, iconPlay, iconNext)
	contentIcons2 := container.NewVBox(contentIcons, iconStop)

	return container.NewHSplit(music, container.NewCenter(contentIcons2))
}