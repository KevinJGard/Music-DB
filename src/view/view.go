package view

import (
	"fmt"
	"net/url"
	"errors"
	"strconv"
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

// Run_View initializes and starts the main application window.
func Run_View() {
	controller := controller.NewController()
	var (
		progress *widget.ProgressBar
		progressContainer *fyne.Container
	)
	myApp := app.NewWithID("com.kevingard.musicdatabase")
	myWindow := myApp.NewWindow("Music DB")
	myWindow.SetIcon(theme.MediaMusicIcon())
	myWindow.Resize(fyne.NewSize(1150, 600))

	searchContainer := createSearchContainer(myWindow, controller, myApp)
	cont, contSouth, updateList := createListContainer(controller, myWindow, myApp)
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

	content := container.New(layout.NewBorderLayout(searchContainer, contSouth, nil, nil),
		searchContainer,
		contSouth,
		cont,
		container.NewCenter(progressContainer),
	)
	
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

// createSearchContainer creates a container for the search functionality.
func createSearchContainer(myWindow fyne.Window, controller *controller.Controller, myApp fyne.App) *fyne.Container {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search...")

	searchButton := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		query := searchEntry.Text
		if query != "" {
			fmt.Println("Search:", query)
			openSongsFound(controller, myApp, searchEntry.Text)
			searchEntry.SetText("")
		} else {
			err := errors.New("Enter your search.")
			dialog.ShowError(err, myWindow)
		}
	})

	return container.NewGridWithColumns(2, searchEntry, searchButton)
}

// openSongsFound creates a new window displaying the found songs based on the search query.
func openSongsFound(controller *controller.Controller, myApp fyne.App, search string) {
	songsFound := myApp.NewWindow("Songs Found")
	songsFound.SetIcon(theme.SearchIcon())
	songsFound.Resize(fyne.NewSize(1000, 600))

	songsLabel := widget.NewLabelWithStyle("Songs Found", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	songsIcon := widget.NewIcon(theme.SearchIcon())
	north := container.NewHBox(songsLabel, songsIcon)
	center := container.NewCenter(north)
	cont, contSouth := createListContainerBySearch(controller, songsFound, myApp, search)
	editContent := container.New(layout.NewBorderLayout(center, contSouth, nil, nil),
		center, cont, contSouth)

	songsFound.SetContent(editContent)
	songsFound.CenterOnScreen()
	songsFound.Show()
}

// openSettingsWindow creates a settings window for the application.
func openSettingsWindow(myApp fyne.App) {
	settingsWindow := myApp.NewWindow("Settings")
	settingsWindow.SetIcon(theme.SettingsIcon())
	settingsWindow.Resize(fyne.NewSize(600, 500))

	themes := createThemeButtons(myApp)
	quit := widget.NewButton("Close", func() {settingsWindow.Close()})

	labelSettings := widget.NewLabelWithStyle("Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	iconSettings := widget.NewIcon(theme.SettingsIcon())
	north := container.NewHBox(labelSettings, iconSettings)
	center := container.NewCenter(north)

	settingsContent := container.New(layout.NewBorderLayout(center, quit, nil, nil),
		center, themes, quit)

	settingsWindow.SetContent(settingsContent)
	settingsWindow.Show()
}

// createThemeButtons creates buttons to switch between light and dark theme.
func createThemeButtons(myApp fyne.App) *fyne.Container {
	darkButton := widget.NewButton("Dark", func() {
		myApp.Settings().SetTheme(theme.DarkTheme())
	})
	lightButton := widget.NewButton("Light", func() {
		myApp.Settings().SetTheme(theme.LightTheme())
	})

	return container.NewGridWithColumns(2, darkButton, lightButton)
}

// createMainMenu sets up the main menu of the application.
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

// setPath allows the user to select a directory for music files.
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

// createListContainer creates a container to display the list of songs.
func createListContainer(controller *controller.Controller, myWindow fyne.Window, myApp fyne.App) (*container.Split, *container.Split, func()) {
	var (
		songEdit *widget.Button
		albumEdit *widget.Button
		performerEdit *widget.Button
	)
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
	songEdit = widget.NewButtonWithIcon("Edit Song", theme.DocumentCreateIcon(), nil)
	performerLabel := widget.NewLabel("Artist:  ")
	performerEdit = widget.NewButtonWithIcon("Edit P.", theme.DocumentCreateIcon(), nil)
	performerCont := container.NewGridWithColumns(2, performerLabel, performerEdit)
	albumLabel := widget.NewLabel("Album: ")
	albumEdit = widget.NewButtonWithIcon("Edit A.", theme.DocumentCreateIcon(), nil)
	albumCont := container.NewGridWithColumns(2, albumLabel, albumEdit)
	trackLabel := widget.NewLabel("Track: ")
	yearLabel := widget.NewLabel("Year: ")
	genreLabel := widget.NewLabel("Genre: ")
	detailsCont := container.NewVBox(widget.NewSeparator(), performerCont, widget.NewSeparator(), albumCont, widget.NewSeparator(), 
				trackLabel, widget.NewSeparator(), yearLabel, widget.NewSeparator(), genreLabel, widget.NewSeparator(), songEdit)
	detailsCont.Hide()
	detailsContainer := container.NewVBox(hbox, detailsCont)

	music := widget.NewLabel("Your Music.")
	musicIcon := widget.NewIcon(theme.FileAudioIcon())
	iconPlay := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { fmt.Println("Play.") })
	iconNext := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() { fmt.Println("Next.") })
	iconPrevious := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() { fmt.Println("Previous.") })
	iconStop := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() { fmt.Println("Stop.") })
	iconPlay.Disable()
	iconNext.Disable()
	iconPrevious.Disable()
	iconStop.Disable()

	yourMusic := container.NewHBox(musicIcon, music)
	contentIcons := container.NewHBox(iconPrevious, iconPlay, iconNext)
	contentIcons2 := container.NewVBox(contentIcons, iconStop)

	list.OnSelected = func(id widget.ListItemID) {
		songs, err := controller.GetSongs()
		if err != nil {
			dialog.ShowError(err, myWindow) 
			return
		}
		song := songs[id]
		label.SetText(song.Title)
		music.SetText(song.Title)
		icon.SetResource(theme.MediaMusicIcon())
		musicIcon.SetResource(theme.MediaMusicIcon())
		performerLabel.SetText("Artist: " + song.PerformerName)
		albumLabel.SetText("Album: " + song.AlbumName)
		trackLabel.SetText("Track: " + fmt.Sprintf("%d", song.Track))
		yearLabel.SetText("Year: " + fmt.Sprintf("%d", song.Year))
		genreLabel.SetText("Genre: " + song.Genre)
		detailsCont.Show()
		songEdit.OnTapped = func() {
			openEditSongWindow(myApp, controller, song.ID, updateList)
		}
		albumEdit.OnTapped = func() {
			openEditAlbumWindow(myApp, controller, song.AlbumID)
		}
		performerEdit.OnTapped = func() {
			openEditPerformerWindow(myApp, controller, song.PerformerID)
		}
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	updateList()

	return container.NewHSplit(list, container.NewCenter(detailsContainer)), container.NewHSplit(container.NewCenter(yourMusic), container.NewCenter(contentIcons2)),updateList
}

// createListContainerBySearch creates a container to display songs based on the search query.
func createListContainerBySearch(controller *controller.Controller, myWindow fyne.Window, myApp fyne.App, search string) (*container.Split, *container.Split) {
	var (
		songEdit *widget.Button
		albumEdit *widget.Button
		performerEdit *widget.Button
	)
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
		songs, err := controller.GetSearchSongs(search)
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
	songEdit = widget.NewButtonWithIcon("Edit Song", theme.DocumentCreateIcon(), nil)
	performerLabel := widget.NewLabel("Artist:  ")
	performerEdit = widget.NewButtonWithIcon("Edit P.", theme.DocumentCreateIcon(), nil)
	performerCont := container.NewGridWithColumns(2, performerLabel, performerEdit)
	albumLabel := widget.NewLabel("Album: ")
	albumEdit = widget.NewButtonWithIcon("Edit A.", theme.DocumentCreateIcon(), nil)
	albumCont := container.NewGridWithColumns(2, albumLabel, albumEdit)
	trackLabel := widget.NewLabel("Track: ")
	yearLabel := widget.NewLabel("Year: ")
	genreLabel := widget.NewLabel("Genre: ")
	detailsCont := container.NewVBox(widget.NewSeparator(), performerCont, widget.NewSeparator(), albumCont, widget.NewSeparator(), 
				trackLabel, widget.NewSeparator(), yearLabel, widget.NewSeparator(), genreLabel, widget.NewSeparator(), songEdit)
	detailsCont.Hide()
	detailsContainer := container.NewVBox(hbox, detailsCont)

	music := widget.NewLabel("Your Music.")
	musicIcon := widget.NewIcon(theme.FileAudioIcon())
	iconPlay := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() { fmt.Println("Play.") })
	iconNext := widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), func() { fmt.Println("Next.") })
	iconPrevious := widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), func() { fmt.Println("Previous.") })
	iconStop := widget.NewButtonWithIcon("", theme.MediaStopIcon(), func() { fmt.Println("Stop.") })
	iconPlay.Disable()
	iconNext.Disable()
	iconPrevious.Disable()
	iconStop.Disable()

	yourMusic := container.NewHBox(musicIcon, music)
	contentIcons := container.NewHBox(iconPrevious, iconPlay, iconNext)
	contentIcons2 := container.NewVBox(contentIcons, iconStop)

	list.OnSelected = func(id widget.ListItemID) {
		songs, err := controller.GetSearchSongs(search)
		if err != nil {
			dialog.ShowError(err, myWindow) 
			return
		}
		song := songs[id]
		label.SetText(song.Title)
		music.SetText(song.Title)
		icon.SetResource(theme.MediaMusicIcon())
		musicIcon.SetResource(theme.MediaMusicIcon())
		performerLabel.SetText("Artist: " + song.PerformerName)
		albumLabel.SetText("Album: " + song.AlbumName)
		trackLabel.SetText("Track: " + fmt.Sprintf("%d", song.Track))
		yearLabel.SetText("Year: " + fmt.Sprintf("%d", song.Year))
		genreLabel.SetText("Genre: " + song.Genre)
		detailsCont.Show()
		songEdit.OnTapped = func() {
			openEditSongWindow(myApp, controller, song.ID, updateList)
		}
		albumEdit.OnTapped = func() {
			openEditAlbumWindow(myApp, controller, song.AlbumID)
		}
		performerEdit.OnTapped = func() {
			openEditPerformerWindow(myApp, controller, song.PerformerID)
		}
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	updateList()

	return container.NewHSplit(list, container.NewCenter(detailsContainer)), container.NewHSplit(container.NewCenter(yourMusic), container.NewCenter(contentIcons2))
}

// openEditSongWindow opens a window to edit song information.
func openEditSongWindow(myApp fyne.App, controller *controller.Controller, id int64, updateList func()) {
	editS := myApp.NewWindow("Edit")
	editS.SetIcon(theme.DocumentCreateIcon())
	editS.Resize(fyne.NewSize(600, 500))

	title := widget.NewEntry()
	title.SetPlaceHolder("New Title")
	title.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "Title can only contain letters, numbers, '_' and '-'")
	track := widget.NewEntry()
	track.SetPlaceHolder("Track number")
	track.Validator = validation.NewRegexp(`^[0-9]+$`, "Track can only contain numbers.")
	year := widget.NewEntry()
	year.SetPlaceHolder("Year number")
	year.Validator = validation.NewRegexp(`^[0-9]+$`, "Year can only contain numbers.")
	genre := widget.NewEntry()
	genre.SetPlaceHolder("New Genre")
	genre.Validator = validation.NewRegexp(`^[A-Za-z]+$`, "Genre can only contain letters.")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Title", Widget: title, HintText: "Title change."},
			{Text: "Track", Widget: track, HintText: "Track number change."},
			{Text: "Year", Widget: year, HintText: "Year change."},
			{Text: "Genre", Widget: genre, HintText: "Genre change."},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
			editS.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			trackNum, _ := strconv.Atoi(track.Text)
			yearNum, _ := strconv.Atoi(year.Text)
			if err := controller.EditSong(id, title.Text, genre.Text, trackNum, yearNum); err != nil {
				dialog.ShowError(err, editS)
			} else {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Music DB",
					Content: "Modified Song: " + title.Text + ".\n Track: " + track.Text + ".\n Year:" + year.Text + ".\n Genre: " + genre.Text,
				})
				updateList()
			}
			editS.Close()
		},
	}

	editLabel := widget.NewLabelWithStyle("Edit Song", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	editIcon := widget.NewIcon(theme.DocumentCreateIcon())
	north := container.NewHBox(editLabel, editIcon)
	center := container.NewCenter(north)
	editContent := container.New(layout.NewBorderLayout(center, nil, nil, nil),
		center, form)

	editS.SetContent(editContent)
	editS.Show()
}

// openEditAlbumWindow opens a window to edit album information.
func openEditAlbumWindow(myApp fyne.App, controller *controller.Controller, id int64) {
	editA := myApp.NewWindow("Edit")
	editA.SetIcon(theme.DocumentCreateIcon())
	editA.Resize(fyne.NewSize(600, 500))

	name := widget.NewEntry()
	name.SetPlaceHolder("New Name")
	name.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "Name can only contain letters, numbers, '_' and '-'")
	year := widget.NewEntry()
	year.SetPlaceHolder("Year number")
	year.Validator = validation.NewRegexp(`^[0-9]+$`, "Year can only contain numbers.")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name, HintText: "Name change."},
			{Text: "Year", Widget: year, HintText: "Year change."},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
			editA.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			yearNum, _ := strconv.Atoi(year.Text)
			if err := controller.EditAlbum(id, name.Text, yearNum); err != nil {
				dialog.ShowError(err, editA)
			} else {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Music DB",
					Content: "Modified Album: " + name.Text + ".\n Year: " + year.Text,
				})
			}
			editA.Close()
		},
	}

	editLabel := widget.NewLabelWithStyle("Edit Album", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	editIcon := widget.NewIcon(theme.DocumentCreateIcon())
	north := container.NewHBox(editLabel, editIcon)
	center := container.NewCenter(north)
	editContent := container.New(layout.NewBorderLayout(center, nil, nil, nil),
		center, form)

	editA.SetContent(editContent)
	editA.Show()
}

// openEditPerformerWindow opens a window to edit performer information, also allows you to set it as a person 
// or as a group and if you set it as a person it allows you to add it to a group.
func openEditPerformerWindow(myApp fyne.App, controller *controller.Controller, id int64) {
	var (
		person *widget.Check
		group *widget.Check
		inGroup *widget.Check
		noDef *widget.Check
	)
	editP := myApp.NewWindow("Edit Performer")
	editP.SetIcon(theme.DocumentCreateIcon())
	editP.Resize(fyne.NewSize(700, 600))

	name := widget.NewEntry()
	name.SetPlaceHolder("Stage Name")
	name.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "Name can only contain letters, numbers, '_' and '-'")
	name.Disable()
	realName := widget.NewEntry()
	realName.SetPlaceHolder("Stage Name")
	realName.Validator = validation.NewRegexp(`^[A-Za-z]+$`, "Name can only contain letters.")
	realName.Disable()
	birth := widget.NewEntry()
	birth.SetPlaceHolder("Birth date")
	birth.Validator = validation.NewRegexp(`^[0-9]+$`, "Date can only contain numbers.")
	birth.Disable()
	death := widget.NewEntry()
	death.SetPlaceHolder("Death date")
	death.Validator = validation.NewRegexp(`^[0-9]+$`, "Date can only contain numbers.")
	death.Disable()
	person = widget.NewCheck("Define as a person", func(b bool) {
		if b {
			name.Enable()
			realName.Enable()
			birth.Enable()
			death.Enable()
			group.Disable()
			inGroup.Enable()
			noDef.Disable()
		} else {
			name.Disable()
			realName.Disable()
			birth.Disable()
			death.Disable()
			group.Enable()
			inGroup.Disable()
			noDef.Enable()
		}
	})
	nameInG := widget.NewEntry()
	nameInG.SetPlaceHolder("Group name")
	nameInG.Disable()
	inGroup = widget.NewCheck("Put it in a group", func(b bool) {
		if b {
			nameInG.Enable()
		} else {
			nameInG.Disable()
		}
	})
	inGroup.Disable()
	nameG := widget.NewEntry()
	nameG.SetPlaceHolder("Name")
	nameG.Validator = validation.NewRegexp(`^[A-Za-z]+$`, "Name can only contain letters.")
	nameG.Disable()
	start := widget.NewEntry()
	start.SetPlaceHolder("Start date")
	start.Validator = validation.NewRegexp(`^[0-9]+$`, "Date can only contain numbers.")
	start.Disable()
	end := widget.NewEntry()
	end.SetPlaceHolder("End date")
	end.Validator = validation.NewRegexp(`^[0-9]+$`, "Date can only contain numbers.")
	end.Disable()
	group = widget.NewCheck("Define as a group", func(b bool) {
		if b {
			nameG.Enable()
			start.Enable()
			end.Enable()
			person.Disable()
			noDef.Disable()
		} else {
			nameG.Disable()
			start.Disable()
			end.Disable()
			person.Enable()
			noDef.Enable()
		}
	})
	newName := widget.NewEntry()
	newName.SetPlaceHolder("Name")
	newName.Validator = validation.NewRegexp(`^[A-Za-z]+$`, "Name can only contain letters.")
	newName.Disable()
	noDef = widget.NewCheck("No def", func(b bool) {
		if b{
			newName.Enable()
			group.Disable()
			person.Disable()
		} else {
			newName.Disable()
			group.Enable()
			person.Enable()
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Stage name", Widget: name, HintText: "Put a name."},
			{Text: "Real name", Widget: realName, HintText: "Put a name."},
			{Text: "Birth date", Widget: birth, HintText: "Put a birth date."},
			{Text: "death date", Widget: death, HintText: "Put a death date. 0 if alive."},
			{Text: "Group", Widget: nameInG, HintText: "Enter a name of an existing group."},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
			editP.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			if nameInG.Text != "" {
				err := controller.DefPerson(id, name.Text, realName.Text, birth.Text, death.Text)
				err = controller.AddPersonToGroup(name.Text, realName.Text, birth.Text, death.Text, nameInG.Text)
				if err != nil {
					dialog.ShowError(err, editP)
				} else {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   "Music DB",
						Content: "Modified Performer: " + name.Text + ".\n Real name: " + realName.Text + ".\n Birth date: " + birth.Text + ".\n Death date: " + death.Text + ".\n Add to group: " + nameInG.Text,
					})
				}
			} else {
				if err := controller.DefPerson(id, name.Text, realName.Text, birth.Text, death.Text); err != nil {
					dialog.ShowError(err, editP)
				} else {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   "Music DB",
						Content: "Modified Performer: " + name.Text + ".\n Real name: " + realName.Text + ".\n Birth date: " + birth.Text + ".\n Death date: " + death.Text,
					})
				}
			}
			editP.Close()
		},
	}
	form.Append("Person", person)
	form.Append("In a Group", inGroup)

	form2 := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Group name", Widget: nameG, HintText: "Put a name."},
			{Text: "Start date", Widget: start, HintText: "Put a start date."},
			{Text: "End date", Widget: end, HintText: "Put a end date. 0 if alive"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
			editP.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			if err := controller.DefGroup(id, nameG.Text, start.Text, end.Text); err != nil {
				dialog.ShowError(err, editP)
			} else {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Music DB",
					Content: "Modified Performer: " + nameG.Text + ".\n Start date: " + start.Text + ".\n End date: " + end.Text,
				})
			}
			editP.Close()
		},
	}
	form2.Append("Group", group)

	form3 := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "New Name", Widget: newName, HintText: "Put a new name"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
			editP.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			if err := controller.EditPerf(id, newName.Text); err != nil {
				dialog.ShowError(err, editP)
			} else {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Music DB",
					Content: "Modified Performer: " + newName.Text,
				})
			}
			editP.Close()
		},
	}
	form3.Append("Undefined", noDef)

	forms := container.NewVBox(form2, widget.NewSeparator(), form3)

	editLabel := widget.NewLabelWithStyle("Edit Performer", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	editIcon := widget.NewIcon(theme.DocumentCreateIcon())
	north := container.NewHBox(editLabel, editIcon)
	center := container.NewCenter(north)
	editContent := container.New(layout.NewBorderLayout(center, nil, nil, nil),
		center, container.NewHSplit(form, forms))

	editP.SetContent(editContent)
	editP.Show()
}