package http_handlers

import (
	"net/http"

	"github.com/arieffian/mw-backend-test/internal/connectors"
	"github.com/arieffian/mw-backend-test/pkg/helpers"
)

type TransactionHandler struct{}

var (
	TransactionRepo connectors.TransactionRepository
)

func (t *TransactionHandler) TransactionHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost:
		t.CreateTransaction(w, r)
	case r.Method == http.MethodGet:
		t.GetTransactionByID(w, r)
	default:
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusNotFound, "404 page not found", nil, nil, nil)
	}
}

func (t *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Create Product", nil, nil, nil)
}

func (t *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Get Product By ID", nil, nil, nil)
}
