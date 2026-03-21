DROP INDEX IF EXISTS unique_active_booking_per_slot;

DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS slots;
DROP TABLE IF EXISTS schedule_days;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS booking_status;
DROP TYPE IF EXISTS user_role;
