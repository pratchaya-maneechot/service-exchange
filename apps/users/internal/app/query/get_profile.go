package query

type GetUserProfile struct {
	Id string
}

const GetUserProfileQuery = "GetUserProfile"

func (c GetUserProfile) QueryName() string {
	return GetUserProfileQuery
}
