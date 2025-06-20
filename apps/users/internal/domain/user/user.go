package user

import (
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/pkg/utils"
)

// UserStatus defines the status of a user.
type UserStatus string

const (
	UserStatusActive              UserStatus = "ACTIVE"
	UserStatusInactive            UserStatus = "INACTIVE"
	UserStatusPendingVerification UserStatus = "PENDING_VERIFICATION"
	UserStatusSuspended           UserStatus = "SUSPENDED"
)

// User represents the User Aggregate Root.
// It encapsulates the core identity of a user and acts as a consistency boundary.
type User struct {
	ID           ids.UserID // Value Object for User ID
	LineUserID   string     // LINE User ID from LIFF
	Email        string     // Email, can be empty if only LINE login is used
	PasswordHash string     // Hashed password, if traditional login is supported
	Status       UserStatus // Current status of the user
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time // Pointer to allow nil (no login yet)

	Profile               Profile                // Value Object for user's profile details
	Roles                 []role.Role            // Value Object for user's Roles
	identityVerifications []IdentityVerification // Entities related to identity verification attempts
}

// NewUser creates a new User aggregate.
// This is the factory method for creating User objects, ensuring initial invariants.
func NewUser(userID ids.UserID, lineUserID string, email string, password string) (*User, error) {
	passwordHash := password
	now := time.Now()
	user := &User{
		ID:                    userID,
		LineUserID:            lineUserID,
		Email:                 email,
		PasswordHash:          passwordHash,
		Status:                UserStatusPendingVerification, // Default initial status
		CreatedAt:             now,
		UpdatedAt:             now,
		Roles:                 []role.Role{},
		identityVerifications: []IdentityVerification{},
	}
	user.Profile = *NewProfile(userID, lineUserID)

	// Domain events can be published here, but typically handled by the application layer after persistence
	// e.g., user.RecordEvent(UserRegistered{UserID: userID})

	return user, nil
}

func (u *User) IsVerified() bool {
	return utils.ArraySome(u.identityVerifications, func(idv IdentityVerification) bool {
		return idv.VerifiedAt != nil
	})
}

func (u *User) SubmitVerify(idv IdentityVerification) {
	u.identityVerifications = append(u.identityVerifications, idv)
}

// UpdateProfile updates the user's profile details.
func (u *User) UpdateProfile(displayName, firstName, lastName, bio, avatarURL, phoneNumber, address *string, preferences map[string]any) {
	if displayName != nil {
		u.Profile.DisplayName = *displayName
	}
	u.Profile.FirstName = firstName
	u.Profile.LastName = lastName
	u.Profile.Bio = bio
	u.Profile.AvatarURL = avatarURL
	u.Profile.PhoneNumber = phoneNumber
	u.Profile.Address = address
	u.Profile.Preferences = preferences
	u.UpdatedAt = time.Now()
	// u.RecordEvent(ProfileUpdated{UserID: u.ID}) // Mark event for publication
}

// SetStatus changes the status of the user.
func (u *User) SetStatus(newStatus UserStatus) error {
	// Add domain rules for status transitions here
	// e.g., if u.Status == UserStatusSuspended && newStatus == UserStatusActive { return ErrCannotActivateSuspendedUser }
	u.Status = newStatus
	u.UpdatedAt = time.Now()
	// u.RecordEvent(UserStatusChanged{UserID: u.ID, NewStatus: newStatus})
	return nil
}

// HasRole checks if the user has a specific role.
func (u *User) HasRole(roleName role.RoleName) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

// AddRole assigns a role to the user.
func (u *User) AddRole(role role.Role) error {
	if u.HasRole(role.Name) {
		return ErrRoleAlreadyAssigned
	}
	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now()
	return nil
}

// RecordLogin updates the last login timestamp.
func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}
