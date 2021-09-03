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
	UserID int                        `json:"user_id" validate:"required,numeric,gt=0"`
	Detail []trasanctionDetailRequest `json:"detail" validate:"required"`
}

type trasanctionDetailRequest struct {
	ProductID int `json:"product_id" validate:"required,numeric,gt=0"`
	Qty       int `json:"qty" validate:"required,numeric,gt=0"`
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

	// validate user id exists
	_, err = UserRepo.GetUserByID(r.Context(), transaction.UserID)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "User ID not found", nil, nil, nil)
		return
	}

	trans := &connectors.TransactionRecord{
		UserID: transaction.UserID,
		Date:   time.Now(),
	}

	detail := []*connectors.TransactionDetailRecord{}
	for i := 0; i < len(transaction.Detail); i++ {
		// validate product id exists
		_, err = ProductRepo.GetProductByID(r.Context(), transaction.Detail[i].ProductID)
		if err != nil {
			helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Product ID not found", nil, nil, nil)
			return
		}
		detail = append(detail, &connectors.TransactionDetailRecord{
			ProductID: transaction.Detail[i].ProductID,
			Qty:       transaction.Detail[i].Qty,
		})
	}

	trans.TransactionDetail = detail

	_, err = TransactionRepo.CreateTransaction(r.Context(), trans)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Internal Server Error", nil, nil, nil)
		return
	}

	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Success", nil, nil, nil)
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
