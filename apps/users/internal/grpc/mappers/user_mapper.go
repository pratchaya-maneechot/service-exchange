package mappers

import (
	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

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

func MapUserProfileDTOToProto(internalDTO *query.UserProfileDTO) *pb.UserProfileDTO {
	if internalDTO == nil {
		return nil
	}
	var lastLoginAt *timestamppb.Timestamp
	if internalDTO.LastLoginAt != nil {
		lastLoginAt = timestamppb.New(*internalDTO.LastLoginAt)
	}
	protoDTO := &pb.UserProfileDTO{
		UserId:      internalDTO.UserID,
		LineUserId:  internalDTO.LineUserID,
		Email:       internalDTO.Email,
		DisplayName: internalDTO.DisplayName,
		Status:      MapDomainUserStatusToProto(internalDTO.Status),
		IsVerified:  internalDTO.IsVerified,
		CreatedAt:   timestamppb.New(internalDTO.CreatedAt),
		LastLoginAt: lastLoginAt,
		FirstName:   utils.GetStringValue(internalDTO.FirstName),
		LastName:    utils.GetStringValue(internalDTO.LastName),
		Bio:         utils.GetStringValue(internalDTO.Bio),
		AvatarUrl:   utils.GetStringValue(internalDTO.AvatarURL),
		PhoneNumber: utils.GetStringValue(internalDTO.PhoneNumber),
		Address:     utils.GetStringValue(internalDTO.Address),
		Preferences: utils.GetInterfaceString(internalDTO.Preferences),
		Roles:       internalDTO.Roles,
	}
	return protoDTO
}

func MapIdentityVerificationDTOToProto(internalDTO *query.IdentityVerificationDTO) *pb.IdentityVerificationDTO {
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
