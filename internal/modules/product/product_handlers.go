package product

import (
	"net/http"

	"github.com/ulshv/go-service/internal/utils/httputils"
)

type productHandlers struct {
	svc *productSvc
}

func newProductHandlers(svc *productSvc) *productHandlers {
	return &productHandlers{
		svc: svc,
	}
}

func (h *productHandlers) RegisterHandlers(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("GET /api/v1/products/:id", h.getProductByIdHandler)
	mux.HandleFunc("POST /api/v1/products", h.createProductHandler)
	return mux
}

func (h *productHandlers) createProductHandler(w http.ResponseWriter, r *http.Request) {
	// body, _ := io.ReadAll(r.Body)
	// bodyStr := string(body)
	// fmt.Println("bodyStr:", bodyStr)
	// httputils.WriteErrorJson(w, "not implemented", http.StatusInternalServerError)
	// return
	var dto createProductDto
	err := httputils.DecodeBody(w, r, &dto)
	if err != nil {
		return
	}
	p := newProduct(dto.Name, dto.Desc, dto.Price)
	created, err := h.svc.createProduct(p)
	if err != nil {
		httputils.WriteErrorJson(w, err.Error(), http.StatusBadRequest)
		return
	}
	httputils.WriteJson(w, created)
}

func (h *productHandlers) getProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented yet."))
}
