package dao

import (
	"errors"
	"suger-clickup/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
)

type DBDao interface {
	// user
	CreateUser(email, password string, orgID int32) error
	UpdateUserOrg(userID, orgID int32) error
	GetUser(email, password string) (*models.User, error)
	GetUserByID(userID int32) (*models.User, error)

	// org
	CreateOrg(orgName string) error
	GetOrg(orgId int32) (*models.Org, error)
	GetOrgNotification(orgId int32) (*models.OrgNotiRequest, error)
	ModOrgNotification(orgID int32, req *models.OrgNotiRequest) error

	// ticket
	CreateTicket(ticket *models.Ticket) error
	GetTicketsByOrgID(orgID, page, pageSize int32) ([]models.Ticket, error)
	GetTicketsCount(orgID int32) (int32, error)
	GetTicket(ticketID string) (models.Ticket, error)
	DeleteTicket(ticketID string) error
	SetTicketWatcher(ticketID string, watcher []int32) error
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

func (d *postgresDao) GetUserByID(userId int32) (*models.User, error) {
	query := `SELECT * FROM accounts WHERE user_id=$1`
	user := models.User{}
	err := d.db.Get(&user, query, userId)
	log.Debugf("userid:%v user:%+v err:%#+v", userId, user, err)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *postgresDao) CreateUser(email, password string, orgID int32) error {
	sql := `INSERT INTO accounts (email, password, created_at, org_id) VALUES ($1, $2, $3, $4)`

	res := d.db.MustExec(sql, email, password, time.Now(), orgID)
	if count, err := res.RowsAffected(); err == nil && count > 0 {
		return nil
	}
	return errors.New("fail to insert user")
}

func (d *postgresDao) UpdateUserOrg(userID, orgID int32) error {
	sql := `UPDATE accounts SET org_id=$1 WHERE user_id=$2`
	_, err := d.db.Exec(sql, orgID, userID)
	return err
}

func (d *postgresDao) CreateOrg(orgName string) error {
	sql := `INSERT INTO orgs (org_name) VALUES ($1)`
	res, err := d.db.Exec(sql, orgName)
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); count < 1 || err != nil {
		return errors.New("fail to insert")
	}
	return nil
}

func (d *postgresDao) ModOrgNotification(orgID int32,
	req *models.OrgNotiRequest) error {
	sql := `INSERT INTO notifications (org_id, status_change, content_change)
	 VALUES ($1, $2, $3)
	 ON CONFLICT (org_id) DO UPDATE
		SET status_change=EXCLUDED.status_change,
		content_change=EXCLUDED.content_change;	
	 `
	res, err := d.db.Exec(sql, orgID, req.StatusChange, req.ContentChange)
	if err != nil {
		log.Debugf("sql:%s err:%v", sql, err)
		return err
	}
	if count, err := res.RowsAffected(); count < 1 || err != nil {
		return errors.New("fail to insert")
	}
	return nil
}

func (d *postgresDao) GetOrgNotification(orgId int32) (*models.OrgNotiRequest, error) {
	sql := `SELECT * FROM notifications WHERE org_id=$1`
	org := models.OrgNotiRequest{}
	err := d.db.Get(&org, sql, orgId)
	return &org, err
}

func (d *postgresDao) GetOrg(orgId int32) (*models.Org, error) {
	sql := `SELECT * FROM orgs WHERE org_id=$1`
	org := models.Org{}
	err := d.db.Get(&org, sql, orgId)
	return &org, err
}

func (d *postgresDao) GetTicketsCount(orgID int32) (int32, error) {
	var count int32
	if err := d.db.Get(&count, "SELECT COUNT(*) FROM tickets WHERE org_id=$1",
		orgID); err != nil {
		return 0, err
	}
	return count, nil
}

func (d *postgresDao) GetTicketsByOrgID(orgID, page, pageSize int32) ([]models.Ticket, error) {
	offset := (page - 1) * pageSize
	sql := `SELECT * FROM tickets WHERE org_id=$1 ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	var result []models.Ticket
	err := d.db.Select(&result, sql, orgID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *postgresDao) CreateTicket(ticket *models.Ticket) error {
	sql := `INSERT INTO tickets (ticket_id, user_id, org_id, created_at, watcher)
	 VALUES ($1, $2, $3, $4, $5)`

	res := d.db.MustExec(sql, ticket.TicketID, ticket.UserID, ticket.OrgID,
		time.Now(), ticket.Watcher)
	if count, err := res.RowsAffected(); err == nil && count > 0 {
		return nil
	}
	return errors.New("fail to insert user")
}

func (d *postgresDao) GetTicket(ticketID string) (models.Ticket, error) {
	sql := `SELECT * FROM tickets WHERE ticket_id=$1`
	ticket := models.Ticket{}
	if err := d.db.Get(&ticket, sql, ticketID); err != nil {
		return ticket, err
	}
	return ticket, nil
}

func (d *postgresDao) SetTicketWatcher(ticketID string, watcher []int32) error {
	sql := `UPDATE tickets SET watcher=$2 WHERE ticket_id=$1`
	if _, err := d.db.Exec(sql, ticketID, watcher); err != nil {
		log.Debugf("ticket:%s sql:%s array:%v", ticketID, sql, watcher)
		return err
	}
	return nil
}
func (d *postgresDao) DeleteTicket(ticketID string) error {
	sql := `DELETE FROM tickets WHERE ticket_id=$1`
	if _, err := d.db.Exec(sql, ticketID); err != nil {
		return err
	}
	return nil
}
