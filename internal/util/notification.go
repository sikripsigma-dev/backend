package util

func IsNotificationOwner(userID *string, userId string, companyId string) bool {
	if userID == nil {
		return false
	}
	return *userID == userId || *userID == companyId
}