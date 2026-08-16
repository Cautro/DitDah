package check

import (
	user "ditdah/internal/features/user"
)

func CheckToAdmin(user *user.UserEntity) bool {
	if user.IsAdmin {
		return true
	}
	return false 
}