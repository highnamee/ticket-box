-- +goose Up
-- +goose StatementBegin
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS max_booking_per_user INT;
CREATE INDEX IF NOT EXISTS idx_bookings_user_ticket_status ON bookings (user_id, ticket_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_bookings_user_ticket_status;
ALTER TABLE tickets DROP COLUMN IF EXISTS max_booking_per_user;
-- +goose StatementEnd
