package repository

import (
	"errors"
	"suger-clickup/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func FindUserByCredentials(email, password string) (*models.User, error) {
	query := `SELECT * FROM accounts WHERE email=$1`
	user := new(models.User)
	err := models.GetDB().Get(user, query, email)
	log.Debugf("email:%s passwd:%s user:%+v", email, password, user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func CreateUserByCredentials(email, password string) error {
	sql := `INSERT INTO accounts (email, password, created_at) VALUES ($1, $2, $3)`
	db := models.GetDB()

	res := db.MustExec(sql, email, password, time.Now())
	if count, err := res.RowsAffected(); err == nil && count > 0 {
		return nil
	}
	return errors.New("fail to insert user")
}
