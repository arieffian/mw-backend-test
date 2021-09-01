package http_handlers

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/arieffian/mw-backend-test/internal/connectors"
	"github.com/arieffian/mw-backend-test/pkg/helpers"
)

var (
	ProductRepo connectors.ProductRepository

	productRegExp      = regexp.MustCompile(`^\/product[\/]*$`)
	productBrandRegExp = regexp.MustCompile(`^\/product\/brand[\/]*$`)
)

type productRequest struct {
	BrandID int    `json:"brand_id" validate:"required,numeric,gt=0"`
	Name    string `json:"name" validate:"required"`
	Qty     int    `json:"qty" validate:"required,numeric,gte=0"`
	Price   int    `json:"price" validate:"required,numeric,gte=0"`
}
type ProductHandler struct{}

func (p *ProductHandler) ProductHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost && productRegExp.MatchString(r.URL.Path):
		p.CreateProduct(w, r)
	case r.Method == http.MethodGet && productRegExp.MatchString(r.URL.Path):
		p.GetProductByID(w, r)
	case r.Method == http.MethodGet && productBrandRegExp.MatchString(r.URL.Path):
		p.GetProductByBrandID(w, r)
	default:
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusNotFound, "404 page not found", nil, nil, nil)
	}
}

func (b *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Create Product", nil, nil, nil)
}

func (b *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
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

	product, err := ProductRepo.GetProductByID(r.Context(), id)

	if err != nil {
		helpers.WriteHTTPResponse(r.Context(), w, http.StatusInternalServerError, "Error fetching the product", nil, nil, nil)
		return
	}

	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Success", nil, product, nil)
}

func (b *ProductHandler) GetProductByBrandID(w http.ResponseWriter, r *http.Request) {
	helpers.WriteHTTPResponse(r.Context(), w, http.StatusOK, "Get Product By Brand ID", nil, nil, nil)
}
