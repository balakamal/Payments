package service

import "context"

type Bill struct {
	BillID      int32  `json:"bill_id,omitempty"`
	UserID      int32  `json:"user_id"`
	CampaignID  int32  `json:"campaign_id"`
	PledgeID    int32  `json:"pledge_id"`
	VatChargeID int32  `json:"vat_charge_id"`
	AmountCents int32  `json:"amount_cents"`
	Status      string `json:"status"`
}

type Repository interface {
	CreateBill(ctx context.Context, bill Bill) error
	GetBillByID(ctx context.Context, id string) (Bill, error)
}
