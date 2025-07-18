package views

import (
	pb "github.com/pratchaya-maneechot/service-exchange/apps/users/api/proto/user"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/app/query"
	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/domain/user"
	lg "github.com/pratchaya-maneechot/service-exchange/libs/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func domainUserStatusToProto(domainStatus user.UserStatus) pb.UserStatus {
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

func UserProfile(internalDTO *query.UserProfileDTO) *pb.UserProfileDTO {
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
		Email:       lg.PtrToStringValue(internalDTO.Email),
		DisplayName: internalDTO.DisplayName,
		Status:      domainUserStatusToProto(internalDTO.Status),
		IsVerified:  internalDTO.IsVerified,
		CreatedAt:   timestamppb.New(internalDTO.CreatedAt),
		LastLoginAt: lastLoginAt,
		FirstName:   lg.PtrToStringValue(internalDTO.FirstName),
		LastName:    lg.PtrToStringValue(internalDTO.LastName),
		Bio:         lg.PtrToStringValue(internalDTO.Bio),
		AvatarUrl:   lg.PtrToStringValue(internalDTO.AvatarURL),
		PhoneNumber: lg.PtrToStringValue(internalDTO.PhoneNumber),
		Address:     lg.PtrToStringValue(internalDTO.Address),
		Preferences: lg.AnyMapToStringMap(internalDTO.Preferences),
		Roles:       internalDTO.Roles,
	}
	return protoDTO
}
