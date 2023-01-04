package entity

// Image entity Image
type Image struct {
	File []byte
	Name string
}

// NewImage create a new entity Image
func NewImage() Image {
	return Image{
		File: nil,
		Name: "",
	}
}

// FileName
func (i Image) FileName() string {
	return i.Name
}

func (i Image) IsEmpty() bool {
	return i.File == nil
}
