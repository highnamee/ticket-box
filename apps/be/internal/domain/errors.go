package domain

import "errors"

var (
	ErrNotFound          = errors.New("record not found")
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrTicketSoldOut     = errors.New("ticket is sold out")
	ErrInsufficientStock = errors.New("insufficient ticket stock")
	ErrInvalidQuantity   = errors.New("invalid quantity requested")
)
