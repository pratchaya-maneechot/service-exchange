package entities

type LineProfile struct {
	ID            string // line user id eg. U4af4980629...
	UserID        string // user model
	DisplayName   string
	PictureURL    string
	StatusMessage *string
}
