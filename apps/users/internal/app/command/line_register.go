package command

type LineRegister struct {
	LineRefID     string
	DisplayName   string
	PictureURL    string
	StatusMessage *string
}

const LineRegisterCommand = "LineRegister"

func (c LineRegister) CommandName() string {
	return LineRegisterCommand
}
