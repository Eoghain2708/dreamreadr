package views

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/Eoghain2708/dreamreadr/internal/service"
	nav "github.com/Eoghain2708/dreamreadr/internal/ui/fyneui/navigation"
)

type CreateDreamView struct {
	service   *service.DreamService
	container fyne.Container
	dream     *dream.Dream
}

func NewCreateDreamView(ds *service.DreamService, n nav.Navigator) (*CreateDreamView, error) {
	v := &CreateDreamView{service: ds}

	titleEntry := widget.NewEntry()

	titleContainer := container.NewVBox(
		widget.NewButton("Back", func() {
			n.ShowDreams()
		}),
		widget.NewLabel("Title"),
		titleEntry,
	)

	contentInput := widget.NewMultiLineEntry()
	contentInput.Wrapping = fyne.TextWrapWord

	var err error
	submitButton := container.NewHBox(
		widget.NewButton("Submit", func() {
			v.dream, err = ds.CreateDream(titleEntry.Text, contentInput.Text)
			n.ShowDream(v.dream.ID)
		}),
	)

	if err != nil {
		return nil, err
	}

	v.container = *container.NewBorder(
		titleContainer,
		submitButton,
		nil,
		nil,
		contentInput,
	)

	return v, nil
}

func (v *CreateDreamView) CanvasObject() fyne.CanvasObject {
	return &v.container
}
