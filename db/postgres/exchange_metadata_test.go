package postgres_test

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"testing"

// 	"bitbucket.org/hornetdefiant/core-engine/factory"
// 	_ "github.com/lib/pq"
// 	"github.com/ory/dockertest/v3"
// 	"github.com/stretchr/testify/assert"

// 	"bitbucket.org/hornetdefiant/core-engine/db/postgres"
// )

// var db *sql.DB
// var resource *dockertest.Resource

// func TestMain(m *testing.M) {
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		log.Fatalf("Could not connect to Docker: %s", err)
// 	}

// 	resource, err = pool.Run("postgres", "latest", []string{
// 		"POSTGRES_USER=user",
// 		"POSTGRES_PASSWORD=pass",
// 		"POSTGRES_DB=database",
// 	})
// 	if err != nil {
// 		log.Fatalf("Could not start PostgreSQL container: %s", err)
// 	}

// 	if err := pool.Retry(func() error {
// 		var err error
// 		connStr := fmt.Sprintf("user=myuser password=mypassword host=localhost port=%s dbname=mydb sslmode=disable", resource.GetPort("5432/tcp"))
// 		db, err = sql.Open("postgres", connStr)
// 		if err != nil {
// 			return err
// 		}
// 		return db.Ping()
// 	}); err != nil {
// 		log.Fatalf("Could not connect to PostgreSQL: %s", err)
// 	}

// 	code := m.Run()

// 	if err := pool.Purge(resource); err != nil {
// 		log.Fatalf("Could not purge PostgreSQL container: %s", err)
// 	}

// 	os.Exit(code)
// }

// func TestAddExchangeMetaData(t *testing.T) {
// 	p := postgres.NewPostgres(db)
// 	data := &factory.ExchangeMetaDataInsert{
// 		Name:          "TestExchange",
// 		Country:       "TestCountry",
// 		ContactEmail:  "test@example.com",
// 		ContactNumber: "1234567890",
// 		WalletAddress: "Wallet1",
// 	}

// 	err := p.AddExchangeMetaData(data)
// 	assert.NoError(t, err)
// }

// func TestGetExchangeMetaData(t *testing.T) {
// 	p := postgres.NewPostgres(db)

// 	data := &factory.ExchangeMetaDataInsert{
// 		Name:          "TestExchange",
// 		Country:       "TestCountry",
// 		ContactEmail:  "test@example.com",
// 		ContactNumber: "1234567890",
// 		WalletAddress: "Wallet2",
// 	}
// 	err := p.AddExchangeMetaData(data)
// 	assert.NoError(t, err)

// 	exchanges, err := p.GetExchangeMetaData()
// 	assert.NoError(t, err)
// 	assert.Len(t, exchanges, 1)

// 	expectedExchange := factory.Exchange{
// 		Name:          "TestExchange",
// 		Country:       "TestCountry",
// 		ContactEmail:  "test@example.com",
// 		ContactNumber: "1234567890",
// 		WalletAddress: "Wallet3",
// 	}

// 	assert.Equal(t, expectedExchange, exchanges[0])
// }
