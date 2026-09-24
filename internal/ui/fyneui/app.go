package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"github.com/Eoghain2708/dreamreadr/internal/service"
	"github.com/Eoghain2708/dreamreadr/internal/ui/fyneui/views"
)

type App struct {
	w       fyne.Window
	service *service.DreamService
}

func NewApp(service *service.DreamService) *App {
	return &App{
		service: service,
	}
}

func (a *App) Run() {
	fyneApp := app.New()
	window := fyneApp.NewWindow("Dreamreadr")
	a.w = window

	view, err := views.NewDreamListView(a.service, a)
	if err != nil {
		panic(err)
	}

	a.w.SetContent(view.CanvasObject())

	a.w.Resize(fyne.NewSize(400, 700))
	a.w.ShowAndRun()
}

func (a *App) ShowDreams() {
	v, err := views.NewDreamListView(a.service, a)
	if err != nil {
		a.showError(err)
		return
	}

	a.w.SetContent(v.CanvasObject())
}

func (a *App) ShowDream(id string) {
	v, err := views.NewDreamView(a.service, a, id)
	if err != nil {
		a.showError(err)
		return
	}

	a.w.SetContent(v.CanvasObject())
}

func (a *App) ShowInsights() {

}

func (a *App) ShowSettings() {

}

func (a *App) CreateDream() {
	v, err := views.NewCreateDreamView(a.service, a)
	if err != nil {
		a.showError(err)
	}

	a.w.SetContent(v.CanvasObject())
}

func (a *App) showError(err error) {
	dialog.ShowError(err, a.w)
}
