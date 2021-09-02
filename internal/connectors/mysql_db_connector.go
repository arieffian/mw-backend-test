package connectors

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arieffian/mw-backend-test/internal/config"
)

var (
	mysqlLog        = log.WithField("file", "mysql_db_connector.go")
	mySQLDbInstance *MySQLDB
)

// GetMySQLDBInstance initializes the MySQL.DB instance
func GetMySQLDBInstance() *MySQLDB {
	if mySQLDbInstance == nil {
		host := config.Get("db.host")
		port := config.GetInt("db.port")
		user := config.Get("db.user")
		password := config.Get("db.password")
		database := config.Get("db.database")
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", user, password, host, port, database))
		if err != nil {
			mysqlLog.WithField("func", "GetMySQLDBInstance").Fatalf("sql.Open got %s", err.Error())
		}

		mySQLDbInstance = &MySQLDB{
			instance: db,
		}
	}
	return mySQLDbInstance
}

// MySQLDB db instance
type MySQLDB struct {
	instance *sql.DB
}

// GetBrandByID retrieves an BrandRecord from database where the brand id is specified.
func (db *MySQLDB) GetBrandByID(ctx context.Context, brandID int) (*BrandRecord, error) {
	fLog := mysqlLog.WithField("func", "GetBrandByID")
	brand := &BrandRecord{}

	row := db.instance.QueryRowContext(ctx, "SELECT id, name FROM brands WHERE id = ?", brandID)
	err := row.Scan(&brand.ID, &brand.Name)
	if err != nil {
		fLog.Errorf("row.Scan got %s", err.Error())
		return nil, err
	}

	return brand, nil
}

// CreateBrand insert an entity record of brand into database.
func (db *MySQLDB) CreateBrand(ctx context.Context, rec *BrandRecord) (string, error) {
	fLog := mysqlLog.WithField("func", "CreateBrand")

	_, err := db.instance.ExecContext(ctx, "INSERT INTO brands(name) VALUES(?)", rec.Name)
	if err != nil {
		fLog.Errorf("db.instance.ExecContext got %s", err.Error())
		return "", err
	}

	return "brand created successfully", nil
}

// CreateProduct insert an entity record of product into database.
func (db *MySQLDB) CreateProduct(ctx context.Context, rec *ProductRecord) (string, error) {
	fLog := mysqlLog.WithField("func", "CreateProduct")

	_, err := db.instance.ExecContext(ctx, "INSERT INTO products(brand_id, name, qty, price) VALUES(?,?,?,?)", rec.BrandID, rec.Name, rec.Qty, rec.Price)
	if err != nil {
		fLog.Errorf("db.instance.ExecContext got %s", err.Error())
		return "", err
	}

	return "product created successfully", nil
}

// GetProductByID retrieves an ProductRecord from database where the product id is specified.
func (db *MySQLDB) GetProductByID(ctx context.Context, productID int) (*ProductRecord, error) {
	fLog := mysqlLog.WithField("func", "GetProductByID")
	product := &ProductRecord{}

	row := db.instance.QueryRowContext(ctx, "SELECT id, brand_id, name, price, qty FROM products WHERE id = ?", productID)
	err := row.Scan(&product.ID, &product.BrandID, &product.Name, &product.Price, &product.Qty)
	if err != nil {
		fLog.Errorf("row.Scan got %s", err.Error())
		return nil, err
	}

	return product, nil
}

// GetProductByBrandID retrieves an array of ProductRecord from database where the brand id is specified.
func (db *MySQLDB) GetProductByBrandID(ctx context.Context, brandID int) ([]*ProductRecord, error) {
	fLog := mysqlLog.WithField("func", "GetProductByBrandID")

	q := fmt.Sprintf("SELECT id, brand_id, name, price, qty FROM products WHERE brand_id = %v", brandID)
	rows, err := db.instance.QueryContext(ctx, q)
	if err != nil {
		fLog.Errorf("db.instance.QueryContext got %s", err.Error())
		return nil, err
	}
	productList := make([]*ProductRecord, 0)
	for rows.Next() {
		product := &ProductRecord{}
		err := rows.Scan(&product.ID, &product.BrandID, &product.Name, &product.Price, &product.Qty)
		if err != nil {
			fLog.Errorf("rows.Scan got %s", err.Error())
		} else {
			productList = append(productList, product)
		}
	}
	return productList, nil
}

// GetTransactionByTransactionID retrieves the detail of a transaction from database where the transaction id is specified.
func (db *MySQLDB) GetTransactionByTransactionID(ctx context.Context, transactionID int) (*TransactionRecord, error) {
	fLog := mysqlLog.WithField("func", "GetTransactionByTransactionID")
	transaction := &TransactionRecord{}

	row := db.instance.QueryRowContext(ctx, "SELECT id, user_id, date, grand_total FROM transactions WHERE id = ?", transactionID)
	err := row.Scan(&transaction.ID, &transaction.UserID, &transaction.Date, &transaction.GrandTotal)
	if err != nil {
		fLog.Errorf("row.Scan got %s", err.Error())
		return nil, err
	}

	q := fmt.Sprintf("SELECT transaction_id, product_id, qty, sub_total FROM transaction_detail WHERE transaction_id = %v", transactionID)
	rows, err := db.instance.QueryContext(ctx, q)
	if err != nil {
		fLog.Errorf("db.instance.QueryContext got %s", err.Error())
		return nil, err
	}

	tDetail := make([]*TransactionDetailRecord, 0)
	for rows.Next() {
		tD := &TransactionDetailRecord{}
		err := rows.Scan(&tD.TransactionID, &tD.ProductID, &tD.Qty, &tD.SubTotal)
		if err != nil {
			fLog.Errorf("rows.Scan got %s", err.Error())
		} else {
			tDetail = append(tDetail, tD)
		}
	}

	transaction.TransactionDetail = tDetail

	return transaction, nil
}

// CreateTransaction insert an entity record of transaction into database.
func (db *MySQLDB) CreateTransaction(ctx context.Context, rec *TransactionRecord) (string, error) {
	return "", nil
}

// GetUserByID retrieves an UserRecord from database where the user id is specified.
func (db *MySQLDB) GetUserByID(ctx context.Context, userID int) (*UserRecord, error) {
	fLog := mysqlLog.WithField("func", "GetUserByID")
	user := &UserRecord{}

	row := db.instance.QueryRowContext(ctx, "SELECT id, name, email, address FROM users WHERE id = ?", userID)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Address)
	if err != nil {
		fLog.Errorf("row.Scan got %s", err.Error())
		return nil, err
	}

	return user, nil
}
