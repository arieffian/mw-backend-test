package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/arieffian/mw-backend-test/internal/connectors"
	"github.com/arieffian/mw-backend-test/pkg/helpers"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransactionByID(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	urlEndPoint := "/order"
	method := "GET"
	Router = http.NewServeMux()
	InitializeRouter()

	t.Run("error-query-param-not-present", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		createRequest := httptest.NewRequest(method, urlEndPoint, nil)
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		httpCode := http.StatusInternalServerError
		dataExpect := &helpers.ResponseJSON{
			Status:  httpCode,
			Message: "Parameter ID not found",
			Data:    nil,
			Error: &helpers.ErrorJSON{
				Message:      "Internal Server Error",
				Reason:       "",
				ErrTittleMsg: "Internal Server Error",
				ErrBodyMsg:   "Internal Server Error",
			},
		}

		assert.Equal(t, httpCode, recorder.Code)
		assert.Equal(t, dataExpect, resBody)
	})

	t.Run("error-query-param-not-numeric", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		createRequest := httptest.NewRequest(method, urlEndPoint, nil)
		q := createRequest.URL.Query()
		q.Add("id", "a")
		createRequest.URL.RawQuery = q.Encode()
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		httpCode := http.StatusInternalServerError
		dataExpect := &helpers.ResponseJSON{
			Status:  httpCode,
			Message: "Parameter ID is not numeric",
			Data:    nil,
			Error: &helpers.ErrorJSON{
				Message:      "Internal Server Error",
				Reason:       "",
				ErrTittleMsg: "Internal Server Error",
				ErrBodyMsg:   "Internal Server Error",
			},
		}

		assert.Equal(t, httpCode, recorder.Code)
		assert.Equal(t, dataExpect, resBody)
	})

	t.Run("success", func(t *testing.T) {
		TransactionRepoMock := new(connectors.MockDBType)
		TransactionRepoMock.On("GetTransactionByTransactionID", mock.Anything, mock.Anything).Return(&connectors.TransactionRecord{}, nil).Once()
		TransactionRepo = TransactionRepoMock

		recorder := httptest.NewRecorder()
		createRequest := httptest.NewRequest(method, urlEndPoint, nil)
		q := createRequest.URL.Query()
		q.Add("id", "1")
		createRequest.URL.RawQuery = q.Encode()
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		if recorder.Code != http.StatusOK {
			t.Errorf("expecting code 200 but got %d. Body %s", recorder.Code, recorder.Body.String())
			t.FailNow()
		}
	})

}

func TestCreateTransaction(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	urlEndPoint := "/order"
	method := "POST"
	Router = http.NewServeMux()
	InitializeRouter()

	t.Run("error-unmarshal", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		createRequest := httptest.NewRequest(method, urlEndPoint, iotest.DataErrReader(bytes.NewReader(nil)))
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		httpCode := http.StatusInternalServerError
		dataExpect := &helpers.ResponseJSON{
			Status:  httpCode,
			Message: "Error processing request",
			Data:    nil,
			Error: &helpers.ErrorJSON{
				Message:      "Internal Server Error",
				Reason:       "",
				ErrTittleMsg: "Internal Server Error",
				ErrBodyMsg:   "Internal Server Error",
			},
		}

		assert.Equal(t, httpCode, recorder.Code)
		assert.Equal(t, dataExpect, resBody)
	})

	t.Run("error-user-not-found", func(t *testing.T) {
		UserRepoMock := new(connectors.MockDBType)
		UserRepoMock.On("GetUserByID", mock.Anything, mock.Anything).Return(&connectors.UserRecord{}, fmt.Errorf("product not found")).Once()
		UserRepo = UserRepoMock

		recorder := httptest.NewRecorder()
		s := `{"user_id": 1,"detail": [{"product_id": 1,"qty": 1},{"product_id": 2,"qty": 1},{"product_id": 3,"qty": 1}]}`
		raw := json.RawMessage(s)
		req, _ := raw.MarshalJSON()
		createRequest := httptest.NewRequest(method, urlEndPoint, bytes.NewReader(req))
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		httpCode := http.StatusInternalServerError
		dataExpect := &helpers.ResponseJSON{
			Status:  httpCode,
			Message: "User ID not found",
			Data:    nil,
			Error: &helpers.ErrorJSON{
				Message:      "Internal Server Error",
				Reason:       "",
				ErrTittleMsg: "Internal Server Error",
				ErrBodyMsg:   "Internal Server Error",
			},
		}

		assert.Equal(t, httpCode, recorder.Code)
		assert.Equal(t, dataExpect, resBody)
	})

	t.Run("error-product-not-found", func(t *testing.T) {
		UserRepoMock := new(connectors.MockDBType)
		UserRepoMock.On("GetUserByID", mock.Anything, mock.Anything).Return(&connectors.UserRecord{}, nil).Once()
		UserRepo = UserRepoMock

		ProductRepoMock := new(connectors.MockDBType)
		ProductRepoMock.On("GetProductByID", mock.Anything, mock.Anything).Return(&connectors.ProductRecord{}, fmt.Errorf("product not found")).Once()
		ProductRepo = ProductRepoMock

		recorder := httptest.NewRecorder()
		s := `{"user_id": 1,"detail": [{"product_id": 1,"qty": 1},{"product_id": 2,"qty": 1},{"product_id": 3,"qty": 1}]}`
		raw := json.RawMessage(s)
		req, _ := raw.MarshalJSON()
		createRequest := httptest.NewRequest(method, urlEndPoint, bytes.NewReader(req))
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		httpCode := http.StatusInternalServerError
		dataExpect := &helpers.ResponseJSON{
			Status:  httpCode,
			Message: "Product ID not found",
			Data:    nil,
			Error: &helpers.ErrorJSON{
				Message:      "Internal Server Error",
				Reason:       "",
				ErrTittleMsg: "Internal Server Error",
				ErrBodyMsg:   "Internal Server Error",
			},
		}

		assert.Equal(t, httpCode, recorder.Code)
		assert.Equal(t, dataExpect, resBody)
	})

	t.Run("error-qty-not-enough", func(t *testing.T) {
		UserRepoMock := new(connectors.MockDBType)
		UserRepoMock.On("GetUserByID", mock.Anything, mock.Anything).Return(&connectors.UserRecord{}, nil).Once()
		UserRepo = UserRepoMock

		ProductRepoMock := new(connectors.MockDBType)
		ProductRepoMock.On("GetProductByID", mock.Anything, mock.Anything).Return(&connectors.ProductRecord{}, nil).Times(3)
		ProductRepo = ProductRepoMock

		TransactionRepoMock := new(connectors.MockDBType)
		TransactionRepoMock.On("CreateTransaction", mock.Anything, mock.Anything).Return("", fmt.Errorf("product qty is not enoug")).Once()
		TransactionRepo = TransactionRepoMock

		recorder := httptest.NewRecorder()
		s := `{"user_id": 1,"detail": [{"product_id": 1,"qty": 10000},{"product_id": 2,"qty": 1},{"product_id": 3,"qty": 1}]}`
		raw := json.RawMessage(s)
		req, _ := raw.MarshalJSON()
		createRequest := httptest.NewRequest(method, urlEndPoint, bytes.NewReader(req))
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		httpCode := http.StatusInternalServerError
		dataExpect := &helpers.ResponseJSON{
			Status:  httpCode,
			Message: "Internal Server Error",
			Data:    nil,
			Error: &helpers.ErrorJSON{
				Message:      "Internal Server Error",
				Reason:       "",
				ErrTittleMsg: "Internal Server Error",
				ErrBodyMsg:   "Internal Server Error",
			},
		}

		assert.Equal(t, httpCode, recorder.Code)
		assert.Equal(t, dataExpect, resBody)
	})

	t.Run("success", func(t *testing.T) {
		UserRepoMock := new(connectors.MockDBType)
		UserRepoMock.On("GetUserByID", mock.Anything, mock.Anything).Return(&connectors.UserRecord{}, nil).Once()
		UserRepo = UserRepoMock

		ProductRepoMock := new(connectors.MockDBType)
		ProductRepoMock.On("GetProductByID", mock.Anything, mock.Anything).Return(&connectors.ProductRecord{}, nil).Times(3)
		ProductRepo = ProductRepoMock

		TransactionRepoMock := new(connectors.MockDBType)
		TransactionRepoMock.On("CreateTransaction", mock.Anything, mock.Anything).Return("", nil).Once()
		TransactionRepo = TransactionRepoMock

		recorder := httptest.NewRecorder()
		s := `{"user_id": 1,"detail": [{"product_id": 1,"qty": 1},{"product_id": 2,"qty": 1},{"product_id": 3,"qty": 1}]}`
		raw := json.RawMessage(s)
		req, _ := raw.MarshalJSON()
		createRequest := httptest.NewRequest(method, urlEndPoint, bytes.NewReader(req))
		createRequest.Header.Add("Content-Type", "application/json")
		Router.ServeHTTP(recorder, createRequest)

		rawBody, _ := ioutil.ReadAll(recorder.Body)
		resBody := &helpers.ResponseJSON{}
		json.Unmarshal(rawBody, resBody)

		if recorder.Code != http.StatusOK {
			t.Errorf("expecting code 200 but got %d. Body %s", recorder.Code, recorder.Body.String())
			t.FailNow()
		}

	})

}
