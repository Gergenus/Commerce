package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gergenus/commerce/product-service/internal/models"
	dbpkg "github.com/Gergenus/commerce/product-service/pkg/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrProductAlreadyExists  = errors.New("product already exists")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryIDNotFound    = errors.New("category id not found")
	ErrProductIDNotFound     = errors.New("product id not found")
	ErrNoSuchCategoryExists  = errors.New("no such category exists")
	ErrNoSuchProductExists   = errors.New("no such product exists")
)

type PostgresRepository struct {
	db dbpkg.PostgresDB
}

func NewPostgresRepository(DB dbpkg.PostgresDB) PostgresRepository {
	return PostgresRepository{
		db: DB,
	}
}

func (p *PostgresRepository) AddCategory(ctx context.Context, category string) (int, error) {
	const op = "repository.AddCategory"
	var id int
	err := p.db.DB.QueryRow(ctx, "INSERT INTO categories (category) VALUES($1) RETURNING id", category).Scan(&id)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" {
				return -1, fmt.Errorf("%s: %w", op, ErrCategoryAlreadyExists)
			}
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *PostgresRepository) GetCategoryID(ctx context.Context, category string) (int, error) {
	const op = "repository.GetCategoryID"
	var id int
	err := p.db.DB.QueryRow(ctx, "SELECT id FROM categories WHERE category = $1", category).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, fmt.Errorf("%s: %w", op, ErrNoSuchCategoryExists)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (p *PostgresRepository) DeleteCategoryByID(ctx context.Context, id int) error {
	const op = "repository.DeleteCategoryByID"
	res, err := p.db.DB.Exec(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrCategoryIDNotFound)
	}
	return nil
}

func (p *PostgresRepository) GetStockByID(ctx context.Context, id int) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) AddStockByID(ctx context.Context, id int, number int) (int, error) {
	panic("implement")
}

func (p *PostgresRepository) ReduceStock(ctx context.Context, id int, number int) (int, error) {
	panic("implement")
}

// One seller cannot hold any items with the same names
func (p *PostgresRepository) CreateProduct(ctx context.Context, product models.Product) (int, error) {
	const op = "repository.CreateProduct"
	var id int
	err := p.db.DB.QueryRow(ctx, "INSERT INTO product_list (product_name, price, seller_id, category_id) VALUES($1, $2, $3, $4) RETURNING id",
		product.ProductName, product.Price, product.SellerID, product.CategoryID).Scan(&id)
	if err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) {
			if pgxErr.Code == "23505" {
				return -1, fmt.Errorf("%s: %w", op, ErrProductAlreadyExists)
			}
			if pgxErr.Code == "23503" {
				return -1, fmt.Errorf("%s: %w", op, ErrNoSuchCategoryExists)
			}
			return -1, fmt.Errorf("%s: %w", op, err)
		}

	}
	return id, nil
}

func (p *PostgresRepository) DeleteProduct(ctx context.Context, id int) error {
	const op = "repository.DeleteProduct"
	res, err := p.db.DB.Exec(ctx, "DELETE FROM product_list WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrProductIDNotFound)
	}
	return nil
}

// Seller cannot be changed
func (p *PostgresRepository) UpdateProduct(ctx context.Context, product models.Product, product_id int) error {
	const op = "repository.UpdateProduct"
	res, err := p.db.DB.Exec(ctx, "UPDATE product_list SET product_name = $1, price = $2, category_id = $3 WHERE id = $4",
		product.ProductName, product.Price, product.CategoryID, product_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrNoSuchProductExists)
	}
	return nil
}

func (p *PostgresRepository) GetProductByID(ctx context.Context, id int) (models.Product, error) {
	var product models.Product
	const op = "repository.GetProductByID"
	err := p.db.DB.QueryRow(ctx, "SELECT * FROM product_list WHERE id = $1", id).Scan(&product.ID, &product.ProductName, &product.Price,
		&product.SellerID, &product.CategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, fmt.Errorf("%s: %w", op, ErrNoSuchProductExists)
		}
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}
	return product, nil
}

func (p *PostgresRepository) GetProductsByCategory(ctx context.Context, category string) ([]models.Product, error) {
	const op = "repository.GetProductsByCategory"

	id, err := p.GetCategoryID(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var products []models.Product
	rows, err := p.db.DB.Query(ctx, "SELECT * FROM product_list WHERE category_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.ProductName, &product.Price, &product.SellerID, &product.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return products, nil
}
