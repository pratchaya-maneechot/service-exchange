package query

import (
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
)

// GetUserProfileQuery represents a query to retrieve a user's profile.
type GetUserProfileQuery struct {
	UserID ids.UserID `validate:"required"`
}

// UserProfileDTO represents the data transfer object for a user's profile.
// This is the "Read Model" or "DTO" that is returned to the client.
type UserProfileDTO struct {
	UserID      string          `json:"userId"`
	LineUserID  string          `json:"lineUserId"`
	Email       string          `json:"email"`
	DisplayName string          `json:"displayName"`
	FirstName   *string         `json:"firstName,omitempty"`
	LastName    *string         `json:"lastName,omitempty"`
	Bio         *string         `json:"bio,omitempty"`
	AvatarURL   *string         `json:"avatarUrl,omitempty"`
	PhoneNumber *string         `json:"phoneNumber,omitempty"`
	Address     *string         `json:"address,omitempty"`
	Preferences map[string]any  `json:"preferences,omitempty"`
	Status      user.UserStatus `json:"status"`
	Roles       []string        `json:"roles"` // List of role names
	IsVerified  bool            `json:"isVerified"`
	LastLoginAt *time.Time      `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
}

// GetUserIdentityVerificationQuery represents a query to retrieve a user's identity verification status.
type GetUserIdentityVerificationQuery struct {
	UserID ids.UserID `validate:"required"`
}

// IdentityVerificationDTO represents the data transfer object for identity verification status.
type IdentityVerificationDTO struct {
	ID              string                  `json:"id"`
	UserID          string                  `json:"userId"`
	DocumentType    user.DocumentType       `json:"documentType"`
	DocumentURLs    []string                `json:"documentUrls"`
	Status          user.VerificationStatus `json:"status"`
	SubmittedAt     time.Time               `json:"submittedAt"`
	VerifiedAt      *time.Time              `json:"verifiedAt,omitempty"`
	RejectionReason *string                 `json:"rejectionReason,omitempty"`
}
