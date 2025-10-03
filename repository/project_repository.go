package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := "SELECT id_product, name, price FROM product"
	rows, err := pr.connection.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsList []model.Product
	var productObj model.Product

	for rows.Next() {
		err := rows.Scan(
			&productObj.ID,
			&productObj.Name,
			&productObj.Price)

		if err != nil {
			return nil, err
		}

		productsList = append(productsList, productObj)
	}

	return productsList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	query, err := pr.connection.Prepare("INSERT INTO product (name, price) VALUES ($1, $2) RETURNING id_product")
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	defer query.Close()

	var id int
	err = query.QueryRow(product.Name, product.Price).Scan(&id)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	return id, nil
}

func (pr *ProductRepository) GetProductByID(id_product int) (*model.Product, error) {
	query, err := pr.connection.Prepare("SELECT id_product, name, price FROM product WHERE id_product = $1")
	if err != nil {
		return nil, err
	}

	var produto model.Product
	defer query.Close()

	err = query.QueryRow(id_product).Scan(&produto.ID, &produto.Name, &produto.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &produto, nil
}
