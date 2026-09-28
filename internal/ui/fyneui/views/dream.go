package views

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Eoghain2708/dreamreadr/internal/dream"
	"github.com/Eoghain2708/dreamreadr/internal/service"
	"github.com/Eoghain2708/dreamreadr/internal/ui/fyneui/helpers"
	nav "github.com/Eoghain2708/dreamreadr/internal/ui/fyneui/navigation"
)

type DreamView struct {
	service   *service.DreamService
	container *fyne.Container
	dream     dream.Dream
	analysis  dream.DreamAnalysis
}

func NewDreamView(ds *service.DreamService, n nav.Navigator, dreamID string) (*DreamView, error) {
	d, err := ds.GetDream(dreamID)
	if err != nil {
		return nil, err
	}

	v := &DreamView{
		service: ds,
		dream:   d,
	}

	analysis, err := ds.GetDreamAnalysis(v.dream.ID)
	if err != nil {
		return nil, err
	}

	var button widget.Button
	label := widget.NewLabel(v.dream.RawText)
	label.Wrapping = fyne.TextWrapWord

	var analysisText string

	if analysis == nil {
		analysisText = "This dream has not been analysed yet"
	} else {
		analysisText = analysis.Summary
	}

	analysisLabel := widget.NewLabel(analysisText)
	analysisLabel.Wrapping = fyne.TextWrapWord

	if analysis == nil {
		button = *widget.NewButton("Analyse", func() {
			analysis, err = ds.AnalyseDream(context.Background(), v.dream.ID)
			if err != nil {
				return
			}
			analysisLabel.SetText(analysis.Summary)
			analysisLabel.Refresh()
		})
	} else {
		button = *widget.NewButton("Re-analyse", func() {
			analysis, err = ds.AnalyseDream(context.Background(), v.dream.ID)
			if err != nil {
				return
			}
			analysisLabel.SetText(analysis.Summary)
			analysisLabel.Refresh()
		})
	}

	v.container = container.NewVBox(
		container.NewHBox(
			widget.NewButton("Back", n.ShowDreams),
			widget.NewButton("Delete", func() {
				_ = ds.DeleteDream(v.dream.ID)
				n.ShowDreams()
			})),
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("%s - %s", v.dream.Title, helpers.FormatCreatedAt(v.dream.CreatedAt))),
			label,
		),
		container.NewVBox(
			widget.NewLabel("Analysis"),
			analysisLabel,
			&button,
		),
	)

	return v, nil
}

func (d *DreamView) CanvasObject() fyne.CanvasObject {
	return d.container
}
