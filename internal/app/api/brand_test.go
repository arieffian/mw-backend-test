package api

import (
	"bytes"
	"encoding/json"
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

func TestCreateBrand(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(ioutil.Discard)

	urlEndPoint := "/brand"
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

	t.Run("success", func(t *testing.T) {
		BrandRepoMock := new(connectors.MockDBType)
		BrandRepoMock.On("CreateBrand", mock.Anything, mock.Anything).Return("success", nil).Once()
		BrandRepo = BrandRepoMock

		recorder := httptest.NewRecorder()
		s := `{"name": "predator"}`
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
