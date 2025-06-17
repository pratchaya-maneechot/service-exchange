package user

import (
	"time"

	"github.com/google/uuid" // For internal ID
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
)

// VerificationStatus defines the status of an identity verification.
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "PENDING"
	VerificationStatusApproved VerificationStatus = "APPROVED"
	VerificationStatusRejected VerificationStatus = "REJECTED"
)

// DocumentType defines the type of document submitted for verification.
type DocumentType string

const (
	DocumentTypeNationalID    DocumentType = "NATIONAL_ID"
	DocumentTypePassport      DocumentType = "PASSPORT"
	DocumentTypeDriverLicense DocumentType = "DRIVER_LICENSE"
)

// IdentityVerification represents an identity verification attempt by a user.
// It's an Entity within the User Aggregate.
type IdentityVerification struct {
	ID              uuid.UUID  // Internal ID for this verification record
	UserID          ids.UserID // Reference to the owning User Aggregate
	DocumentType    DocumentType
	DocumentNumber  string   // Encrypted or hashed for PII
	DocumentURLs    []string // URLs to uploaded document images
	Status          VerificationStatus
	SubmittedAt     time.Time
	VerifiedAt      *time.Time  // Nullable
	ReviewerID      *ids.UserID // ID of the admin who reviewed it (can be internal UserID)
	RejectionReason string      // Reason if rejected
}

// NewIdentityVerification creates a new IdentityVerification entity.
func NewIdentityVerification(
	userID ids.UserID,
	docType DocumentType,
	docNumber string,
	docURLs []string,
) (*IdentityVerification, error) {
	if len(docURLs) == 0 {
		return nil, ErrMissingDocumentURLs
	}
	if docType == "" {
		return nil, ErrMissingDocumentType
	}

	return &IdentityVerification{
		ID:              uuid.New(),
		UserID:          userID,
		DocumentType:    docType,
		DocumentNumber:  docNumber,
		DocumentURLs:    docURLs,
		Status:          VerificationStatusPending, // Initial status
		SubmittedAt:     time.Now(),
		RejectionReason: "",
	}, nil
}

// Approve marks the identity verification as approved.
func (iv *IdentityVerification) Approve(reviewerID ids.UserID) error {
	if iv.Status != VerificationStatusPending {
		return ErrInvalidVerificationStatusTransition
	}
	iv.Status = VerificationStatusApproved
	now := time.Now()
	iv.VerifiedAt = &now
	iv.ReviewerID = &reviewerID
	return nil
}

// Reject marks the identity verification as rejected.
func (iv *IdentityVerification) Reject(reviewerID ids.UserID, reason string) error {
	if iv.Status != VerificationStatusPending {
		return ErrInvalidVerificationStatusTransition
	}
	iv.Status = VerificationStatusRejected
	now := time.Now()
	iv.VerifiedAt = &now
	iv.ReviewerID = &reviewerID
	iv.RejectionReason = reason
	return nil
}
