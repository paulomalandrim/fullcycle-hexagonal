package db_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/paulomalandrim/go-hexagonal/adapters/db"
	"github.com/paulomalandrim/go-hexagonal/application"
	"github.com/stretchr/testify/require"
)

var Db *sql.DB

func setUp() {
	Db, _ := sql.Open("sqlite3", ":memory:")

	createTable(Db)
	createProduct(Db)
}

func createTable(db *sql.DB) {
	createTableQuery := `CREATE TABLE products (
        "id" string,
        "name" string,
		"price" float,
        "status" string
        
    );`

	stmt, err := db.Prepare(createTableQuery)

	if err != nil {
		log.Fatal(err.Error())
	}
	stmt.Exec()
}

func createProduct(db *sql.DB) {
	insertQuery := `INSERT INTO products (id, name, price, status) VALUES ("abc","Product 1",1,"Disabled")`
	stmt, err := db.Prepare(insertQuery)

	if err != nil {
		panic(err.Error())
	}
	stmt.Exec()
}

func TestProductDb_Get(t *testing.T) {

	setUp()

	defer Db.Close()

	productDb := db.NewProductDb(Db)

	product, err := productDb.Get("abc")

	require.Nil(t, err)
	require.Equal(t, "Product 1", product.GetName())
	require.Equal(t, application.DISABLED, product.GetStatus())
	require.Equal(t, float64(1), product.GetPrice())

}
