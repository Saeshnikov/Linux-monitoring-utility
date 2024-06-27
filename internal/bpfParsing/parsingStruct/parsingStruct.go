package parsingstruct

type Interaction interface {
	String() string
}

type ParsingData struct {
	PathsOfExecutableFiles [2]string
	WayOfInteraction       Interaction
}
