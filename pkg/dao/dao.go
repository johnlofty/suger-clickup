package dao

import (
	"errors"
	"suger-clickup/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type DBDao interface {
	CreateUser(email, password string) error
	GetUser(email, password string) (*models.User, error)
}
type postgresDao struct {
	db *sqlx.DB
}

func NewDBDao(db *sqlx.DB) DBDao {
	return &postgresDao{
		db: db,
	}
}
func (d *postgresDao) GetUser(email, password string) (*models.User, error) {
	query := `SELECT * FROM accounts WHERE email=$1 and password=$2`
	user := models.User{}
	err := d.db.Get(&user, query, email, password)
	log.Debugf("email:%v passwd:%v user:%+v err:%#+v", email, password, user, err)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *postgresDao) CreateUser(email, password string) error {
	sql := `INSERT INTO accounts (email, password, created_at) VALUES ($1, $2, $3)`

	res := d.db.MustExec(sql, email, password, time.Now())
	if count, err := res.RowsAffected(); err == nil && count > 0 {
		return nil
	}
	return errors.New("fail to insert user")
}
