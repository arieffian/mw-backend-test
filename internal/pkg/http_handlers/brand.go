package http_handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/arieffian/mw-backend-test/internal/connectors"
	"github.com/arieffian/mw-backend-test/internal/constants/response"
	"github.com/arieffian/mw-backend-test/pkg/helpers"
	"github.com/go-playground/validator"
)

type BrandHandler struct{}

type brandRequest struct {
	Name string `json:"name" validate:"required"`
}

var (
	BrandRepo connectors.BrandRepository
	validate  *validator.Validate
)

func (b *BrandHandler) BrandHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost:
		b.CreateBrand(w, r)
	default:
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusMethodNotAllowed, "Method not Allowed", nil, nil, nil)
	}
}

func (b *BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	brand := &brandRequest{}

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

	err = json.Unmarshal(body, &brand)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Error processing request", nil, nil, nil)
		return
	}

	//validate json input
	validate = validator.New()
	err = validate.Struct(brand)
	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Invalid json structure", nil, nil, nil)
		return
	}

	bRecord := &connectors.BrandRecord{
		Name: brand.Name,
	}

	// insert to database
	result, err := BrandRepo.CreateBrand(r.Context(), bRecord)

	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Internal server error", nil, nil, nil)
		return
	}

	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, result, nil, nil, nil)
}
