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

func TestCreateProduct(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	urlEndPoint := "/product"
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

	t.Run("error-brand-not-found", func(t *testing.T) {
		BrandRepoMock := new(connectors.MockDBType)
		BrandRepoMock.On("GetBrandByID", mock.Anything, mock.Anything).Return(&connectors.BrandRecord{}, fmt.Errorf("product not found")).Once()
		BrandRepo = BrandRepoMock

		recorder := httptest.NewRecorder()
		s := `{"brand_id": 100, "name": "predator", "qty": 3, "price": 1050}`
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
			Message: "Brand ID not found",
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
		BrandRepoMock := new(connectors.MockDBType)
		BrandRepoMock.On("GetBrandByID", mock.Anything, mock.Anything).Return(&connectors.BrandRecord{}, nil).Once()
		BrandRepo = BrandRepoMock

		ProductRepoMock := new(connectors.MockDBType)
		ProductRepoMock.On("CreateProduct", mock.Anything, mock.Anything).Return("sucess", nil).Once()
		ProductRepo = ProductRepoMock

		recorder := httptest.NewRecorder()
		s := `{"brand_id": 1, "name": "predator", "qty": 3, "price": 1050}`
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

func TestGetProductByID(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	urlEndPoint := "/product"
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
		ProductRepoMock := new(connectors.MockDBType)
		ProductRepoMock.On("GetProductByID", mock.Anything, mock.Anything).Return(&connectors.ProductRecord{}, nil).Once()
		ProductRepo = ProductRepoMock

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

func TestGetProductByBrandID(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	urlEndPoint := "/product/brand"
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
		BrandRepoMock := new(connectors.MockDBType)
		BrandRepoMock.On("GetBrandByID", mock.Anything, mock.Anything).Return(&connectors.BrandRecord{}, nil).Once()
		BrandRepo = BrandRepoMock

		ProductRepoMock := new(connectors.MockDBType)
		ProductRepoMock.On("GetProductByBrandID", mock.Anything, mock.Anything).Return([]*connectors.ProductRecord{}, nil).Once()
		ProductRepo = ProductRepoMock

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
