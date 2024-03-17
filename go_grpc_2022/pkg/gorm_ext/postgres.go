package gorm_ext

import (
	"challenge/pkg/env"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	testcontainers_postgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/* DSN */
type PostgresDSN struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDSN() PostgresDSN {
	return PostgresDSN{
		Host:     env.GetOrDefault("POSTGRES_HOST", "localhost"),
		Port:     env.GetOrDefault("POSTGRES_PORT", "5432"),
		User:     env.GetOrDefault("POSTGRES_USER", "postgres"),
		Password: env.GetOrDefault("POSTGRES_PASSWORD", "postgres"),
		DBName:   env.GetOrDefault("POSTGRES_DB", "challenge"),
		SSLMode:  env.GetOrDefault("POSTGRES_SSLMODE", "disable"),
	}
}

func NewTestPostgresDSN(t *testing.T) PostgresDSN {
	ctx := context.Background()
	const postgres = "postgres"

	ctn, err := testcontainers_postgres.RunContainer(
		ctx,
		testcontainers.WithImage("postgres:14"),
		testcontainers_postgres.WithDatabase(postgres),
		testcontainers_postgres.WithUsername(postgres),
		testcontainers_postgres.WithPassword(postgres),
	)
	require.NoError(t, err)

	host, err := ctn.Host(ctx)
	require.NoError(t, err)

	port, err := ctn.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	return PostgresDSN{
		Host:     host,
		Port:     port.Port(),
		User:     postgres,
		Password: postgres,
		DBName:   postgres,
		SSLMode:  "disable",
	}
}

/* DIALECTOR */
func NewPostgresDialector(dsn PostgresDSN) gorm.Dialector {
	return postgres.Open(
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dsn.Host, dsn.Port, dsn.User, dsn.Password, dsn.DBName, dsn.SSLMode,
		),
	)
}
