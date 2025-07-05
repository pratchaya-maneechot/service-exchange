package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/shared/ids"
)

type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "PENDING"
	VerificationStatusApproved VerificationStatus = "APPROVED"
	VerificationStatusRejected VerificationStatus = "REJECTED"
)

type DocumentType string

const (
	DocumentTypeNationalID    DocumentType = "NATIONAL_ID"
	DocumentTypePassport      DocumentType = "PASSPORT"
	DocumentTypeDriverLicense DocumentType = "DRIVER_LICENSE"
)

type IdentityVerification struct {
	ID              uuid.UUID
	UserID          ids.UserID
	DocumentType    DocumentType
	DocumentNumber  string
	DocumentURLs    []string
	Status          VerificationStatus
	SubmittedAt     time.Time
	VerifiedAt      *time.Time
	ReviewerID      *ids.UserID
	RejectionReason string
}

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
		Status:          VerificationStatusPending,
		SubmittedAt:     time.Now(),
		RejectionReason: "",
	}, nil
}

func NewIdentityVerificationFromRepository(
	id string,
	userID string,
	documentType string,
	documentNumber string,
	documentURLs []string,
	status string,
	submittedAt time.Time,
	verifiedAt *time.Time,
	reviewerID *string,
	rejectionReason string,
) (*IdentityVerification, error) {
	var rwID *ids.UserID
	if reviewerID != nil {
		uid := ids.UserID(*reviewerID)
		rwID = &uid
	}

	return &IdentityVerification{
		ID:              uuid.MustParse(id),
		UserID:          ids.UserID(userID),
		DocumentType:    DocumentType(documentType),
		DocumentNumber:  documentNumber,
		DocumentURLs:    documentURLs,
		Status:          VerificationStatus(status),
		SubmittedAt:     submittedAt,
		VerifiedAt:      verifiedAt,
		ReviewerID:      rwID,
		RejectionReason: rejectionReason,
	}, nil
}

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
