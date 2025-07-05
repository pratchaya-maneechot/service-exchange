package user

import (
	"slices"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/role"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
)

type UserStatus string

const (
	UserStatusActive              UserStatus = "ACTIVE"
	UserStatusInactive            UserStatus = "INACTIVE"
	UserStatusPendingVerification UserStatus = "PENDING_VERIFICATION"
	UserStatusSuspended           UserStatus = "SUSPENDED"
)

type User struct {
	ID           ids.UserID
	LineUserID   string
	Email        *string
	PasswordHash *string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time

	Profile               Profile
	Roles                 []role.Role
	identityVerifications []IdentityVerification
}

func NewUser(userID ids.UserID, lineUserID string, email, password *string) (*User, error) {
	passwordHash := password
	now := time.Now()
	user := &User{
		ID:                    userID,
		LineUserID:            lineUserID,
		Email:                 email,
		PasswordHash:          passwordHash,
		Status:                UserStatusPendingVerification,
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

func NewUserFromRepository(
	id string,
	lineUserID string,
	email *string,
	passwordHash *string,
	status string,
	createdAt time.Time,
	updatedAt time.Time,
	lastLoginAt *time.Time,
	profile Profile,
	roles []role.Role,
) (*User, error) {
	return &User{
		ID:           ids.UserID(id),
		LineUserID:   lineUserID,
		Email:        email,
		PasswordHash: passwordHash,
		Status:       UserStatus(status),
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		LastLoginAt:  lastLoginAt,
		Profile:      profile,
		Roles:        roles,
	}, nil
}

func (u *User) IsVerified() bool {
	return slices.ContainsFunc(u.identityVerifications, func(idv IdentityVerification) bool {
		return idv.VerifiedAt != nil
	})
}

func (u *User) SubmitVerify(idv IdentityVerification) {
	u.identityVerifications = append(u.identityVerifications, idv)
}

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

func (u *User) SetStatus(newStatus UserStatus) error {
	u.Status = newStatus
	u.UpdatedAt = time.Now()
	// u.RecordEvent(UserStatusChanged{UserID: u.ID, NewStatus: newStatus})
	return nil
}

func (u *User) HasRole(roleName role.RoleName) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

func (u *User) AddRole(role role.Role) error {
	if u.HasRole(role.Name) {
		return ErrRoleAlreadyAssigned
	}
	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) RecordLogin() {
	now := time.Now()
	u.LastLoginAt = &now
	u.UpdatedAt = now
}
