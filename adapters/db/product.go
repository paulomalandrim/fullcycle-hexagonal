package db

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
	"github.com/paulomalandrim/go-hexagonal/application"
)

type ProductDb struct {
	db *sql.DB
}

func NewProductDb(db *sql.DB) *ProductDb {
	return &ProductDb{db: db}
}

func (p *ProductDb) Get(id string) (application.ProductInterface, error) {

	if p.db == nil {
		return nil, errors.New("database is not initialized")
	}

	var product application.Product
	stmt, err := p.db.Prepare("select id, name, price, status from products where id = ?")
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price, &product.Status)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *ProductDb) Save(product application.ProductInterface) (application.ProductInterface, error) {
	if p.db == nil {
		return nil, errors.New("database is not initialized")
	}

	var rows int
	p.db.QueryRow("SELECT COUNT(*) FROM products WHERE id =?", product.GetID()).Scan(&rows)

	if rows > 0 {
		return p.update(product)
	} else {
		return p.create(product)
	}
}

func (p *ProductDb) create(product application.ProductInterface) (application.ProductInterface, error) {
	if p.db == nil {
		return nil, errors.New("database is not initialized")
	}

	stmt, err := p.db.Prepare("INSERT INTO products (id, name, price, status) VALUES (?,?,?,?)")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(product.GetID(), product.GetName(), product.GetPrice(), product.GetStatus())
	if err != nil {
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductDb) update(product application.ProductInterface) (application.ProductInterface, error) {
	if p.db == nil {
		return nil, errors.New("database is not initialized")
	}

	stmt, err := p.db.Prepare("UPDATE products SET name =?, price =?, status =? WHERE id =?")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(product.GetName(), product.GetPrice(), product.GetStatus(), product.GetID())
	if err != nil {
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		return nil, err
	}

	return product, nil
}
