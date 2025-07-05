package user

import "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"

type Profile struct {
	UserID      ids.UserID
	DisplayName string
	FirstName   *string
	LastName    *string
	Bio         *string
	AvatarURL   *string
	PhoneNumber *string
	Address     *string
	Preferences map[string]any
}

func NewProfile(userID ids.UserID, defaultDisplayName string) *Profile {
	return &Profile{
		UserID:      userID,
		DisplayName: defaultDisplayName,
		Preferences: make(map[string]any),
	}
}

func NewProfileFromRepository(userID string, defaultDisplayName string, preferences map[string]any) Profile {
	return Profile{
		UserID:      ids.UserID(userID),
		DisplayName: defaultDisplayName,
		Preferences: preferences,
	}
}

func (p *Profile) WithFirstName(name string) *Profile {
	p.FirstName = &name
	return p
}
