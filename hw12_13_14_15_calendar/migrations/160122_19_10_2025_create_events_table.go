package main

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up, Down)
}

func Up(tx *sql.Tx) error {
	// Создаем таблицу events (для PostgreSQL 13+)
	_, err := tx.Exec(`
        CREATE TABLE events (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            title VARCHAR(255) NOT NULL,
            start_time TIMESTAMPTZ NOT NULL,
            end_time TIMESTAMPTZ NOT NULL,
            description TEXT,
            user_id UUID NOT NULL,
            notify_before INTERVAL,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW()
        );
    `)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        CREATE INDEX idx_events_user_id ON events(user_id);
        CREATE INDEX idx_events_start_time ON events(start_time);
        CREATE INDEX idx_events_end_time ON events(end_time);
        CREATE INDEX idx_events_time_range ON events(start_time, end_time);
    `)
	if err != nil {
		return err
	}

	// Триггер для обновления updated_at
	_, err = tx.Exec(`
        CREATE OR REPLACE FUNCTION trigger_set_updated_at()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = NOW();
            RETURN NEW;
        END;
        $$ LANGUAGE plpgsql;

        CREATE TRIGGER set_updated_at
            BEFORE UPDATE ON events
            FOR EACH ROW
            EXECUTE FUNCTION trigger_set_updated_at();
    `)

	return err
}

func Down(tx *sql.Tx) error {
	_, err := tx.Exec(`
        DROP TRIGGER IF EXISTS set_updated_at ON events;
        DROP FUNCTION IF EXISTS trigger_set_updated_at;
        DROP TABLE IF EXISTS events;
    `)

	return err
}
