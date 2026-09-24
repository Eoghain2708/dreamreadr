package nav

type Navigator interface {
	ShowDreams()
	ShowInsights()
	ShowDream(id string)
	ShowSettings()

	CreateDream()
}
