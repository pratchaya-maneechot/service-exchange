package command

import (
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

// RegisterUserCommand represents a command to register a new user.
type RegisterUserCommand struct {
	LineUserID  string `json:"lineUserId" validate:"required"`             // Required for LINE LIFF integration
	Email       string `json:"email,omitempty" validate:"omitempty,email"` // Optional, but if present, must be valid email
	Password    string `json:"password,omitempty"`                         // Optional, if only LINE login is supported
	DisplayName string `json:"displayName" validate:"required"`            // From LINE profile or user input
	AvatarURL   string `json:"avatarUrl,omitempty"`                        // From LINE profile or user input
}

// UserLoginCommand represents a command to log in a user.
type UserLoginCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LineLoginCommand represents a command for user login/registration via LINE.
type LineLoginCommand struct {
	LineUserID string `json:"lineUserId" validate:"required"`
	IDToken    string `json:"idToken" validate:"required"` // LINE ID Token for verification
}

// UpdateUserProfileCommand represents a command to update an existing user's profile.
type UpdateUserProfileCommand struct {
	UserID      ids.UserID     `json:"-"` // Passed from JWT, not from body
	DisplayName *string        `json:"displayName,omitempty"`
	FirstName   *string        `json:"firstName,omitempty"`
	LastName    *string        `json:"lastName,omitempty"`
	Bio         *string        `json:"bio,omitempty"`
	AvatarURL   *string        `json:"avatarUrl,omitempty"`
	PhoneNumber *string        `json:"phoneNumber,omitempty" validate:"omitempty,e164"` // E.g., +66812345678
	Address     *string        `json:"address,omitempty"`
	Preferences map[string]any `json:"preferences,omitempty"`
}

// SubmitIdentityVerificationCommand represents a command to submit identity verification documents.
type SubmitIdentityVerificationCommand struct {
	UserID         ids.UserID        `json:"-"`
	DocumentType   user.DocumentType `json:"documentType" validate:"required,oneof=NATIONAL_ID PASSPORT DRIVER_LICENSE"`
	DocumentURLs   []string          `json:"documentUrls" validate:"required,min=1,dive,url"` // Array of URLs to uploaded docs
	DocumentNumber string            `json:"documentNumber,omitempty"`                        // Optional, but good to include for national ID
}
