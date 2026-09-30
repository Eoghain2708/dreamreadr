package views

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/Eoghain2708/dreamreadr/internal/service"
	"github.com/Eoghain2708/dreamreadr/internal/ui/fyneui/helpers"
	nav "github.com/Eoghain2708/dreamreadr/internal/ui/fyneui/navigation"
)

type DreamListView struct {
	service   *service.DreamService
	container *fyne.Container
	dreams    []dream.Dream
}

func NewDreamListView(ds *service.DreamService, n nav.Navigator) (*DreamListView, error) {
	dreams, err := ds.ListDreams()
	if err != nil {
		return nil, err
	}

	v := &DreamListView{
		service: ds,
		dreams:  dreams,
	}

	header := widget.NewLabel("My dreams")
	list := widget.NewList(
		func() int {
			return len(v.dreams)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("", nil)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			button := obj.(*widget.Button)
			button.Alignment = widget.ButtonAlignCenter
			dream := v.dreams[id]
			button.SetText(fmt.Sprintf("%s | %s", v.dreams[id].Title, helpers.FormatCreatedAt(v.dreams[id].CreatedAt)))
			button.OnTapped = func() {
				n.ShowDream(dream.ID)
			}

		},
	)

	scrollableList := container.NewVScroll(list)

	footer := container.NewHBox(
		widget.NewButton("Insights", n.ShowInsights),
		widget.NewButton("+ New Dream", n.CreateDream),
		widget.NewButton("Settings", n.ShowSettings),
	)

	v.container = container.NewBorder(
		header,
		footer,
		nil,
		nil,
		scrollableList,
	)

	return v, nil
}

func (v *DreamListView) CanvasObject() fyne.CanvasObject {
	return v.container
}

type DreamListItem struct {
	widget.BaseWidget
	button   *widget.Button
	onTapped func()
}
