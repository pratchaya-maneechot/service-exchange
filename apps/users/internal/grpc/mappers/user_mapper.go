package mappers

import (
	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/grpc/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
		Email:       utils.GetStringValue(internalDTO.Email),
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
