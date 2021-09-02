package connectors

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
)

func TestCreateBrand(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	t.Run("error-create-brand", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectExec("INSERT INTO brands").WillReturnError(fmt.Errorf("Error DB"))

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		brand := &BrandRecord{
			Name: "test brand",
		}

		_, err = mySQL.CreateBrand(context.Background(), brand)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectExec("INSERT INTO brands").WillReturnResult(sqlmock.NewResult(12, 1))

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		brand := &BrandRecord{
			Name: "test brand",
		}

		_, err = mySQL.CreateBrand(context.Background(), brand)
		if err != nil {
			t.Error("error shouldnt be occurs")
			t.FailNow()
		}
	})
}

func TestGetBrandByID(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	t.Run("error-exec-query-row-context", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM brands").WillReturnError(fmt.Errorf("Error DB"))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetBrandByID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("error-not-found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM brands").WillReturnError(sql.ErrNoRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetBrandByID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test")

		mock.ExpectQuery("SELECT (.+) FROM brands").WillReturnRows(rows)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetBrandByID(context.Background(), 1)
		if err != nil {
			t.Error("error shouldnt be occurs")
			t.FailNow()
		}
	})
}

func TestCreateProduct(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	t.Run("error-create-product", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectExec("INSERT INTO products").WillReturnError(fmt.Errorf("Error DB"))

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		product := &ProductRecord{
			BrandID: 1,
			Name:    "test product",
			Qty:     1,
			Price:   1000,
		}

		_, err = mySQL.CreateProduct(context.Background(), product)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectExec("INSERT INTO products").WillReturnResult(sqlmock.NewResult(12, 1))

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		product := &ProductRecord{
			BrandID: 1,
			Name:    "test product",
			Qty:     1,
			Price:   1000,
		}

		_, err = mySQL.CreateProduct(context.Background(), product)
		if err != nil {
			t.Error("error shouldnt be occurs")
			t.FailNow()
		}
	})
}

func TestGetProductByID(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	t.Run("error-exec-query-row-context", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnError(fmt.Errorf("Error DB"))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetProductByID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("error-not-found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnError(sql.ErrNoRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetProductByID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		rows := sqlmock.NewRows([]string{"id", "product_id", "name", "qty", "price"}).AddRow(1, 1, "name", 1, 1000)

		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetProductByID(context.Background(), 1)
		if err != nil {
			t.Error("error shouldnt be occurs")
			t.FailNow()
		}
	})
}

func TestGetProductByBrandID(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	t.Run("error-exec-query-row-context", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnError(fmt.Errorf("Error DB"))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetProductByBrandID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("error-not-found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnError(sql.ErrNoRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetProductByBrandID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		rows := sqlmock.NewRows([]string{"id", "product_id", "name", "qty", "price"}).
			AddRow(1, 1, "name 1", 1, 1000).
			AddRow(2, 1, "name 2", 2, 1100).
			AddRow(3, 1, "name 3", 3, 1200)

		mock.ExpectQuery("SELECT (.+) FROM products").WillReturnRows(rows)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetProductByBrandID(context.Background(), 1)
		if err != nil {
			t.Error("error shouldnt be occurs")
			t.FailNow()
		}
	})
}

func TestGetTransactionByTransactionID(t *testing.T) {

}

func TestCreateTransaction(t *testing.T) {

}

func TestGetUserByID(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	t.Run("error-exec-query-row-context", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(fmt.Errorf("Error DB"))
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetUserByID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("error-not-found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		mock.ExpectQuery("SELECT (.+) FROM users").WillReturnError(sql.ErrNoRows)
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetUserByID(context.Background(), 1)
		if err == nil {
			t.Error("error should be occurs")
			t.FailNow()
		}
	})

	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		rows := sqlmock.NewRows([]string{"id", "email", "name", "address"}).AddRow(1, "donny@arieffian.com", "donny", "surabaya")

		mock.ExpectQuery("SELECT (.+) FROM users").WillReturnRows(rows)

		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()
		// inject sqlmock.DB into MySQLDB
		mySQL := MySQLDB{
			instance: db,
		}

		_, err = mySQL.GetUserByID(context.Background(), 1)
		if err != nil {
			t.Error("error shouldnt be occurs")
			t.FailNow()
		}
	})
}
