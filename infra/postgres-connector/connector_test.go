package pgconnector

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

var validConnString string // example postgres://test:test@localhost:5432/adv_test_db

var config *pgxpool.Config

var (
	attempts     = 5
	attemptDelay = time.Second * 2
	ctxTimeout   = time.Second
)

func TestMain(m *testing.M) {
	validConnString = os.Getenv("PGCONN")
	if validConnString == "" {
		panic("env PGCONN variable (postgres connection string) not set")
	}

	var err error
	config, err = pgxpool.ParseConfig(validConnString)
	if err != nil {
		panic("env PGCONN variable should be valid postgres connection string")
	}
	os.Exit(m.Run())
}

func TestConnectValid(t *testing.T) {
	assert := assert.New(t)
	conn, err := Connect(validConnString, attempts, attemptDelay, ctxTimeout)
	if assert.Nilf(err, "Connect(...) with valid PGCONN string should return nil, got: %w", err) {
		assert.NotNilf(conn, "If Connect(...) returns nil error, connection shouldn't be nil")
	}
	if conn != nil {
		conn.Close()
	}
}

func TestConnectConfigValid(t *testing.T) {
	assert := assert.New(t)
	conn, err := ConnectConfig(config, attempts, attemptDelay, ctxTimeout)
	if assert.Nilf(err, "Connect(...) with valid PGCONN string should return nil, got: %w", err) {
		assert.NotNilf(conn, "If Connect(...) returns nil error, connection shouldn't be nil")
	}
	if conn != nil {
		conn.Close()
	}
}

func TestInvalid(t *testing.T) {
	assert := assert.New(t)
	// Empty connection string
	emptyConnString := ""
	conn, err := Connect(emptyConnString, attempts, attemptDelay, ctxTimeout)
	assert.NotNilf(err, "Connect(...) with empty PGCONN string should return not-nil error")
	if conn != nil {
		conn.Close()
	}

	// Invalid connection string
	invalidConnString := strings.Replace(validConnString, "postgres", "mongodb", 1)
	conn, err = Connect(invalidConnString, attempts, attemptDelay, ctxTimeout)
	assert.NotNilf(err, "Connect(...) with invalid PGCONN string should return not-nil error")
	if conn != nil {
		conn.Close()
	}
}
