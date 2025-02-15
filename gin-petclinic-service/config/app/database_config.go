package app

import validation "github.com/go-ozzo/ozzo-validation/v4"

type DatabaseConfig struct {
	Postgres     PostgresConfig
	Sqlite       SqliteConfig
	MaxIdleConns int
	MaxOpenConns int
	MaxIdleTime  int
}

func (db DatabaseConfig) Validate() error {
	return validation.ValidateStruct(&db,
		validation.Field(&db.Postgres, validation.Required),
		validation.Field(&db.MaxIdleConns, validation.Min(0)),
		validation.Field(&db.MaxOpenConns, validation.Required),
		validation.Field(&db.MaxIdleTime, validation.Min(0)),
	)
}

type PostgresConfig struct {
	Driver string
	Dsn    string
}

func (pc PostgresConfig) Validate() error {
	return validation.ValidateStruct(&pc,
		validation.Field(&pc.Driver, validation.Required),
		validation.Field(&pc.Dsn, validation.Required),
	)
}

type SqliteConfig struct {
	Driver string
	Dsn    string
}

func (lite SqliteConfig) Validate() error {
	return validation.ValidateStruct(&lite,
		validation.Field(&lite.Driver, validation.Required),
		validation.Field(&lite.Dsn, validation.Required),
	)
}
