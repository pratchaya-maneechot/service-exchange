package user

import "github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"

// Profile represents a user's profile details as a Value Object.
// It has no independent identity and belongs solely to a User.
type Profile struct {
	UserID      ids.UserID // Owner of this profile
	DisplayName string
	FirstName   *string // Pointers for optional fields
	LastName    *string
	Bio         *string
	AvatarURL   *string
	PhoneNumber *string
	Address     *string
	Preferences map[string]any // Stored as JSONB in DB, map in Go
}

// NewProfile creates a default Profile for a new User.
func NewProfile(userID ids.UserID, defaultDisplayName string) *Profile {
	return &Profile{
		UserID:      userID,
		DisplayName: defaultDisplayName, // Often initialized from LINE display name
		Preferences: make(map[string]any),
	}
}

// WithFirstName sets the first name for the profile. Example of a fluent setter.
func (p *Profile) WithFirstName(name string) *Profile {
	p.FirstName = &name
	return p
}

// TODO: Add more methods for updating specific profile fields or preferences.
