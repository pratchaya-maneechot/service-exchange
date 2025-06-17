package handler

import (
	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// --- Helper Functions for Mapping ---

// GetValuePointer converts google.protobuf.StringValue to *string
func GetValuePointer(sv *wrapperspb.StringValue) *string {
	if sv == nil {
		return nil
	}
	val := sv.GetValue()
	return &val
}

// ToStringInterfaceMap converts map<string, string> from proto to map[string]any
func ToStringInterfaceMap(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any)
	for k, v := range m {
		res[k] = v
	}
	return res
}

// MapProtoDocumentTypeToDomain maps Protobuf DocumentType enum to Domain DocumentType string.
func MapProtoDocumentTypeToDomain(protoType pb.DocumentType) user.DocumentType {
	switch protoType {
	case pb.DocumentType_NATIONAL_ID:
		return user.DocumentTypeNationalID
	case pb.DocumentType_PASSPORT:
		return user.DocumentTypePassport
	case pb.DocumentType_DRIVER_LICENSE:
		return user.DocumentTypeDriverLicense
	default:
		return "" // Or handle error
	}
}

// MapDomainDocumentTypeToProto maps Domain DocumentType string to Protobuf DocumentType enum.
func MapDomainDocumentTypeToProto(domainType user.DocumentType) pb.DocumentType {
	switch domainType {
	case user.DocumentTypeNationalID:
		return pb.DocumentType_NATIONAL_ID
	case user.DocumentTypePassport:
		return pb.DocumentType_PASSPORT
	case user.DocumentTypeDriverLicense:
		return pb.DocumentType_DRIVER_LICENSE
	default:
		return pb.DocumentType_DOCUMENT_TYPE_UNSPECIFIED // Default or error
	}
}

// MapDomainUserStatusToProto maps Domain UserStatus string to Protobuf UserStatus enum.
func MapDomainUserStatusToProto(domainStatus user.UserStatus) pb.UserStatus {
	switch domainStatus {
	case user.UserStatusActive:
		return pb.UserStatus_ACTIVE
	case user.UserStatusInactive:
		return pb.UserStatus_INACTIVE
	case user.UserStatusSuspended:
		return pb.UserStatus_SUSPENDED
	case user.UserStatusPendingVerification:
		return pb.UserStatus_PENDING_VERIFICATION
	default:
		return pb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

// MapDomainVerificationStatusToProto maps Domain VerificationStatus string to Protobuf VerificationStatus enum.
func MapDomainVerificationStatusToProto(domainStatus user.VerificationStatus) pb.VerificationStatus {
	switch domainStatus {
	case user.VerificationStatusPending:
		return pb.VerificationStatus_PENDING
	case user.VerificationStatusApproved:
		return pb.VerificationStatus_VERIFIED
	case user.VerificationStatusRejected:
		return pb.VerificationStatus_REJECTED
	default:
		return pb.VerificationStatus_VERIFICATION_STATUS_UNSPECIFIED
	}
}

// MapUserProfileDTOToProto maps an internal UserProfileDTO to a Protobuf UserProfileDTO.
func (h *UserGRPCHandler) MapUserProfileDTOToProto(internalDTO *query.UserProfileDTO) *pb.UserProfileDTO {
	if internalDTO == nil {
		return nil
	}

	protoDTO := &pb.UserProfileDTO{
		UserId:      internalDTO.UserID,
		LineUserId:  internalDTO.LineUserID,
		Email:       internalDTO.Email,
		DisplayName: internalDTO.DisplayName,
		Status:      MapDomainUserStatusToProto(internalDTO.Status),
		IsVerified:  internalDTO.IsVerified,
		CreatedAt:   timestamppb.New(internalDTO.CreatedAt),
	}

	if internalDTO.FirstName != nil {
		protoDTO.FirstName = wrapperspb.String(*internalDTO.FirstName)
	}
	if internalDTO.LastName != nil {
		protoDTO.LastName = wrapperspb.String(*internalDTO.LastName)
	}
	if internalDTO.Bio != nil {
		protoDTO.Bio = wrapperspb.String(*internalDTO.Bio)
	}
	if internalDTO.AvatarURL != nil {
		protoDTO.AvatarUrl = wrapperspb.String(*internalDTO.AvatarURL)
	}
	if internalDTO.PhoneNumber != nil {
		protoDTO.PhoneNumber = wrapperspb.String(*internalDTO.PhoneNumber)
	}
	if internalDTO.Address != nil {
		protoDTO.Address = wrapperspb.String(*internalDTO.Address)
	}
	if internalDTO.LastLoginAt != nil {
		protoDTO.LastLoginAt = timestamppb.New(*internalDTO.LastLoginAt)
	}

	if internalDTO.Preferences != nil {
		protoDTO.Preferences = make(map[string]string)
		for k, v := range internalDTO.Preferences {
			if strVal, ok := v.(string); ok { // Assuming preferences values are string for proto mapping
				protoDTO.Preferences[k] = strVal
			}
		}
	}

	protoDTO.Roles = internalDTO.Roles // Roles are already []string

	return protoDTO
}

// MapIdentityVerificationDTOToProto maps an internal IdentityVerificationDTO to a Protobuf IdentityVerificationDTO.
func (h *UserGRPCHandler) MapIdentityVerificationDTOToProto(internalDTO *query.IdentityVerificationDTO) *pb.IdentityVerificationDTO {
	if internalDTO == nil {
		return nil
	}

	protoDTO := &pb.IdentityVerificationDTO{
		Id:           internalDTO.ID,
		UserId:       internalDTO.UserID,
		DocumentType: MapDomainDocumentTypeToProto(internalDTO.DocumentType),
		DocumentUrls: internalDTO.DocumentURLs,
		Status:       MapDomainVerificationStatusToProto(internalDTO.Status),
		SubmittedAt:  timestamppb.New(internalDTO.SubmittedAt),
	}

	if internalDTO.VerifiedAt != nil {
		protoDTO.VerifiedAt = timestamppb.New(*internalDTO.VerifiedAt)
	}
	if internalDTO.RejectionReason != nil {
		protoDTO.RejectionReason = wrapperspb.String(*internalDTO.RejectionReason)
	}

	return protoDTO
}
