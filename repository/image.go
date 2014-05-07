package repository

type Image struct {
	OriginalName string
	ID           string
	localRoot    string
}

func NewImage(name, workRoot string) *Image {
	return &Image{
		OriginalName: name,
		localRoot:    workRoot,
	}
}
