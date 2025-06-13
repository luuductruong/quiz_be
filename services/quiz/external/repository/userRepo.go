package repository

import (
	"github.com/quiz_be/services/core/context"
	quizDomain "github.com/quiz_be/services/core/domain/quiz"
	"github.com/quiz_be/services/core/helper/sql/query"
)

type userGorm struct {
	UserID      string `gorm:"primary_key"`
	UserName    string
	PhoneNumber string
}

func mapUserToDomain(source *userGorm) *quizDomain.User {
	if source == nil {
		return nil
	}
	return &quizDomain.User{
		UserID:      source.UserID,
		UserName:    source.UserName,
		PhoneNumber: source.PhoneNumber,
	}
}

func (u userGorm) TableName() string {
	return "users"
}

func NewUserRepo() quizDomain.UserRepo {
	return &userRepo{}
}

type userRepo struct {
}

type userQuery struct {
	query.BaseQuery
}

func (u *userRepo) Query(ctx context.Context) quizDomain.UserQuery {
	return &userQuery{query.NewBQ(ctx.GetDbTx().Model(&userGorm{}))}
}

func (u *userQuery) ByUserID(userID string) quizDomain.UserQuery {
	return query.Where(u, "user_id = ?", userID)
}

func (u *userQuery) ByPhoneNumber(phoneNumber string) quizDomain.UserQuery {
	return query.Where(u, "phone_number = ?", phoneNumber)
}

func (u *userQuery) Result() (*quizDomain.User, error) {
	return query.Result(u, mapUserToDomain)
}

func (u *userQuery) ResultList() ([]*quizDomain.User, error) {
	return query.ResultList(u, mapUserToDomain)
}
