package user

import "github.com/Benyam-S/onemembership/entity"

// IService is an interface that defines all the service methods of a user struct
type IService interface {
	AddUser(newUser *entity.User) error
	ValidateUserProfile(user *entity.User) entity.ErrMap
	FindUser(identifier string) (*entity.User, error)
	AllUsers() []*entity.User
	AllUsersWithPagination(pageNum int64) ([]*entity.User, int64)
	SearchUsers(key string, pageNum int64, extra ...string) ([]*entity.User, int64)
	TotalUsers() int64
	UpdateUser(user *entity.User) error
	UpdateUserSingleValue(userID, columnName string, columnValue interface{}) error
	DeleteUser(userID string) (*entity.User, error)

	AddUserPassword(newUserPassword *entity.UserPassword) error
	VerifyUserPassword(userPassword *entity.UserPassword, verifyPassword string) error
	FindUserPassword(userID string) (*entity.UserPassword, error)
	UpdateUserPassword(userPassword *entity.UserPassword) error
	DeleteUserPassword(userID string) (*entity.UserPassword, error)
}
