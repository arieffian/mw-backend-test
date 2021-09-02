package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/arieffian/mw-backend-test/internal/connectors"
	"github.com/arieffian/mw-backend-test/internal/constants/response"
	"github.com/arieffian/mw-backend-test/pkg/helpers"
	"github.com/go-playground/validator"
)

type TransactionHandler struct{}

var (
	TransactionRepo connectors.TransactionRepository
	UserRepo        connectors.UserRepository
)

type transactionRequest struct {
	UserID     int                        `json:"user_id" validate:"required,numeric,gt=0"`
	Date       time.Time                  `json:"date" validate:"required,datetime"`
	Detail     []trasanctionDetailRequest `json:"detail" validate:"required"`
	GrandTotal int                        `json:"grandtotal" validate:"required,numeric,gte=0"`
}

type trasanctionDetailRequest struct {
	ProductID int `json:"product_id" validate:"required,numeric,gt=0"`
	Price     int `json:"price" validate:"required,numeric,gt=0"`
	Qty       int `json:"qty" validate:"required,numeric,gt=0"`
	SubTotal  int `json:"subtotal" validate:"required,numeric,gte=0"`
}

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
	transaction := &transactionRequest{}

	//Unmarshal json
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errJSON := &helpers.ErrorJSON{
			Message:      "Error when parse Body request",
			Reason:       "internal_error",
			ErrTittleMsg: "Error parsing request",
			ErrBodyMsg:   response.Get("general", http.StatusInternalServerError, ""),
		}
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "", nil, nil, errJSON)
		return
	}

	err = json.Unmarshal(body, &transaction)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Error processing request", nil, nil, nil)
		return
	}

	//validate json input
	validate = validator.New()
	err = validate.Struct(transaction)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Invalid json structure", nil, nil, nil)
		return
	}

	//validate brand id exists
	// _, err = BrandRepo.GetBrandByID(r.Context(), product.BrandID)
	// if err != nil {
	// 	helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Brand ID not found", nil, nil, nil)
	// 	return
	// }
	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Create Product", nil, nil, nil)
}

func (t *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	//check if id present and greater than 0
	sID := query.Get("id")
	if sID == "" {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Parameter ID not found", nil, nil, nil)
		return
	}

	//check if id is number or not
	id, err := strconv.Atoi(sID)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Parameter ID is not numeric", nil, nil, nil)
		return
	}

	transaction, err := TransactionRepo.GetTransactionByTransactionID(r.Context(), id)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Error fetching the transaction", nil, nil, nil)
		return
	}
	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Success", nil, transaction, nil)
}
