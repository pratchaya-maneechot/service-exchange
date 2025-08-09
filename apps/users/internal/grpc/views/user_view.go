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
		return pb.UserStatus_USER_STATUS_ACTIVE
	case user.UserStatusInactive:
		return pb.UserStatus_USER_STATUS_INACTIVE
	case user.UserStatusSuspended:
		return pb.UserStatus_USER_STATUS_SUSPENDED
	case user.UserStatusPendingVerification:
		return pb.UserStatus_USER_STATUS_PENDING_VERIFICATION
	default:
		return pb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

func UserProfile(payload *query.UserProfileDTO) *pb.UserProfile {
	if payload == nil {
		return nil
	}
	var lastLoginAt *timestamppb.Timestamp
	if payload.LastLoginAt != nil {
		lastLoginAt = timestamppb.New(*payload.LastLoginAt)
	}
	protoDTO := &pb.UserProfile{
		UserId:      payload.UserID,
		LineUserId:  payload.LineUserID,
		Email:       lg.PtrToStringValue(payload.Email),
		DisplayName: payload.DisplayName,
		Status:      domainUserStatusToProto(payload.Status),
		IsVerified:  payload.IsVerified,
		CreatedAt:   timestamppb.New(payload.CreatedAt),
		LastLoginAt: lastLoginAt,
		FirstName:   lg.PtrToStringValue(payload.FirstName),
		LastName:    lg.PtrToStringValue(payload.LastName),
		Bio:         lg.PtrToStringValue(payload.Bio),
		AvatarUrl:   lg.PtrToStringValue(payload.AvatarURL),
		PhoneNumber: lg.PtrToStringValue(payload.PhoneNumber),
		Address:     lg.PtrToStringValue(payload.Address),
		Preferences: lg.AnyMapToStringMap(payload.Preferences),
		Roles:       payload.Roles,
	}
	return protoDTO
}
