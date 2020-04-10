package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"kkagitala/go-rest-api/service"
)

var (
	ErrRepository = errors.New("unable to handle request")
)

type repository struct {
	db     *sql.DB
	logger log.Logger
}

// New returns a concrete repository backed by CockroachDB
func New(db *sql.DB, logger log.Logger) (service.Repository, error) {
	// return  repository
	return &repository{
		db:     db,
		logger: log.With(logger, "rep", "cockroachdb"),
	}, nil
}

// CreateOrder inserts a new order and its order items into db
func (repo *repository) CreateBill(ctx context.Context, bill service.Bill) error {
	// Run a transaction to sync the query model.
	err := crdb.ExecuteTx(ctx, repo.db, nil, func(tx *sql.Tx) error {
		return createBill(tx, bill)
	})
	if err != nil {
		return err
	}
	return nil
}

func createBill(tx *sql.Tx, bill service.Bill) error {

	// Insert order into the "orders" table.
	sql := `
			INSERT INTO bills (bill_id, amount_cents, campaign_id, pledge_id, user_id)
			VALUES ($1,$2,$3,$4,$5)`
	_, err := tx.Exec(sql, bill.BillID, bill.AmountCents, bill.CampaignID, bill.PledgeID, bill.UserID)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderByID query the order by given id
func (repo *repository) GetBillByID(ctx context.Context, id string) (service.Bill, error) {
	var billRow = service.Bill{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT bill_id, amount_cents, campaign_id, pledge_id, user_id FROM bills WHERE bill_id = $1",
		id).
		Scan(
			&billRow.BillID, &billRow.AmountCents, &billRow.CampaignID, &billRow.PledgeID, &billRow.UserID,
		); err != nil {
		level.Error(repo.logger).Log("err", err.Error())
		return billRow, err
	}
	// ToDo: Query order items from orderitems table
	return billRow, nil
}

// Close implements DB.Close
func (repo *repository) Close() error {
	return repo.db.Close()
}
