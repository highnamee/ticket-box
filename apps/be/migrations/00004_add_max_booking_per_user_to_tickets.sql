-- +goose Up
-- +goose StatementBegin
ALTER TABLE tickets ADD COLUMN IF NOT EXISTS max_booking_per_user INT;
ALTER TABLE tickets ADD CONSTRAINT chk_tickets_max_booking_per_user_positive CHECK (max_booking_per_user IS NULL OR max_booking_per_user > 0);
CREATE INDEX IF NOT EXISTS idx_bookings_user_ticket_status ON bookings (user_id, ticket_id, status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_bookings_user_ticket_status;
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS chk_tickets_max_booking_per_user_positive;
ALTER TABLE tickets DROP COLUMN IF EXISTS max_booking_per_user;
-- +goose StatementEnd
