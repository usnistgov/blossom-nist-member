package model

import "time"

type HistorySnapshot struct {
	TxId      string    `json:"txid"`
	Timestamp time.Time `json:"timestamp"`
	Value     Account   `json:"account"`
}
