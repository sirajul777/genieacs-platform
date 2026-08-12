package pgxpool

import (
	"context"
	"database/sql"
	"errors"
)

type Pool struct{ db *sql.DB }

type Config struct{ ConnString string }

func ParseConfig(connString string) (*Config, error) { return &Config{ConnString: connString}, nil }

func NewWithConfig(ctx context.Context, config *Config) (*Pool, error) {
	return &Pool{}, nil
}

func New(ctx context.Context, connString string) (*Pool, error) { return &Pool{}, nil }

func (p *Pool) Close() {}

func (p *Pool) Ping(ctx context.Context) error { return nil }

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) Row        { return Row{} }
func (p *Pool) Query(ctx context.Context, sql string, args ...any) (Rows, error) { return Rows{}, nil }
func (p *Pool) Exec(ctx context.Context, sql string, args ...any) (CommandTag, error) {
	return CommandTag{}, nil
}

type Row struct{}

func (r Row) Scan(dest ...any) error { return errors.New("pgxpool stub row has no data") }

type CommandTag struct{}

type Rows struct{}

func (r Rows) Close()                 {}
func (r Rows) Next() bool             { return false }
func (r Rows) Scan(dest ...any) error { return errors.New("pgxpool stub row has no data") }
func (r Rows) Err() error             { return nil }
