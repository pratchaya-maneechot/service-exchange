package command

type LineRegister struct {
	LineRefID     string
	DisplayName   string
	PictureURL    string
	StatusMessage *string
}

var LineRegisterCommand = "LineRegister"

func (c LineRegister) CommandName() string {
	return LineRegisterCommand
}
