package views

import (
	"context"
	"fmt"
	"log"

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
	container *container.Scroll
	dream     *dream.Dream
	analysis  *dream.DreamAnalysis
}

func NewDreamView(ds *service.DreamService, n nav.Navigator, dreamID string) (*DreamView, error) {
	d, err := ds.GetDream(dreamID)
	if err != nil {
		return nil, err
	}

	v := &DreamView{
		service: ds,
		dream:   &d,
	}

	analysis, err := ds.GetDreamAnalysis(dreamID)
	if err != nil {
		return nil, err
	}

	v.analysis = analysis

	dreamLabel := widget.NewLabel(v.dream.RawText)
	dreamLabel.Wrapping = fyne.TextWrapWord

	analysisSummary := widget.NewLabel("")
	analysisSummary.Wrapping = fyne.TextWrapWord
	analysisContent := widget.NewCard("Summary", "", analysisSummary)

	emotionsCard, emotionsContent := makeFeatureSection("Emotions")
	themesCard, themesContent := makeFeatureSection("Themes")
	locationsCard, locationsContent := makeFeatureSection("Locations")
	peopleCard, peopleContent := makeFeatureSection("People")
	symbolsCard, symbolsContent := makeFeatureSection("Symbols")

	// This is the part of the page that we'll replace/update.
	analysisContainer := container.NewVBox(
		widget.NewLabel("This dream has not been analysed yet"),
	)

	// If we already have an analysis, populate everything now.
	if analysis != nil {
		analysisSummary.SetText(analysis.Summary)

		populateFeatureSection(emotionsContent, analysis.Emotions)
		populateFeatureSection(themesContent, analysis.Themes)
		populateFeatureSection(locationsContent, analysis.Locations)
		populateFeatureSection(peopleContent, analysis.People)
		populateFeatureSection(symbolsContent, analysis.Symbols)

		analysisContainer.Objects = []fyne.CanvasObject{
			analysisContent,
			emotionsCard,
			themesCard,
			locationsCard,
			peopleCard,
			symbolsCard,
		}

		analysisContainer.Refresh()
	}

	var analyseButton *widget.Button

	analyseButton = widget.NewButton("Analyse", func() {
		analysis, err := ds.AnalyseDream(
			context.Background(),
			dreamID,
		)
		if err != nil {
			log.Println(err)
			return
		}

		v.analysis = analysis

		// Update the summary.
		analysisSummary.SetText(analysis.Summary)

		// Update the feature sections.
		populateFeatureSection(emotionsContent, analysis.Emotions)
		populateFeatureSection(themesContent, analysis.Themes)
		populateFeatureSection(locationsContent, analysis.Locations)
		populateFeatureSection(peopleContent, analysis.People)
		populateFeatureSection(symbolsContent, analysis.Symbols)

		// Replace "not analysed" with the actual analysis UI.
		analysisContainer.Objects = []fyne.CanvasObject{
			analysisContent,
			emotionsCard,
			themesCard,
			locationsCard,
			peopleCard,
			symbolsCard,
		}

		analysisContainer.Refresh()

		analyseButton.SetText("Re-analyse")
	})

	if analysis != nil {
		analyseButton.SetText("Re-analyse")
	}

	content := container.NewVBox(
		container.NewHBox(
			widget.NewButton("Back", n.ShowDreams),
			widget.NewButton("Delete", func() {
				_ = ds.DeleteDream(v.dream.ID)
				n.ShowDreams()
			}),
			analyseButton,
		),

		container.NewVBox(
			widget.NewLabel(
				fmt.Sprintf(
					"%s - %s",
					v.dream.Title,
					helpers.FormatCreatedAt(v.dream.CreatedAt),
				),
			),
			dreamLabel,
		),

		analysisContainer,
	)

	v.container = container.NewVScroll(content)

	return v, nil
}

func (d *DreamView) CanvasObject() fyne.CanvasObject {
	return d.container
}

func makeFeatureSection(title string) (*widget.Card, *fyne.Container) {
	content := container.NewVBox()
	card := widget.NewCard(title, "", content)

	return card, content
}

func populateFeatureSection(content *fyne.Container, values []string) {
	content.Objects = nil

	for _, value := range values {
		content.Add(widget.NewLabel("• " + value))
	}

	content.Refresh()
}
