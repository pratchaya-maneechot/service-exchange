package query

type GetUserProfile struct {
	Id string
}

var GetUserProfileQuery = "GetUserProfile"

func (c GetUserProfile) QueryName() string {
	return GetUserProfileQuery
}
