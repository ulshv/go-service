package product

import (
	"log/slog"
	"net/http"

	"github.com/ulshv/go-service/internal/core/httperrs"
	"github.com/ulshv/go-service/pkg/logs"
	"github.com/ulshv/go-service/pkg/mw"
	"github.com/ulshv/go-service/pkg/utils/httputils"
)

type productHandlers struct {
	svc    *productSvc
	logger *slog.Logger
}

func newProductHandlers(svc *productSvc) *productHandlers {
	return &productHandlers{
		svc:    svc,
		logger: logs.NewLogger("product/handlers"),
	}
}

func (h *productHandlers) RegisterHandlers(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("GET /api/v1/products/:id", h.getProductByIDHandler)
	mux.HandleFunc("POST /api/v1/products", mw.Chain(mw.Authenticate, mw.AuthRequired, h.createProductHandler))
	return mux
}

func (h *productHandlers) createProductHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("createProduct - start")
	userID, ok := mw.GetUserId(r)
	h.logger.Debug("create product for", "user_id", userID)
	if !ok {
		h.logger.Debug("createProduct - unauthrized")
		httputils.WriteErrorJSON(w, httperrs.ErrUnauthorized, httperrs.ErrCodeUnautorized, http.StatusUnauthorized)
		return
	}
	var dto createProductDto
	h.logger.Debug("createProduct - start decoding body")
	err := httputils.DecodeBody(w, r, &dto)
	if err != nil {
		h.logger.Debug("createProduct - err while decoding body", "error", err)
		return
	}
	h.logger.Debug("createProduct - creating new product instance", "dto", dto)
	p := newProduct(userID, dto.Name, dto.Desc, dto.Price)
	h.logger.Debug("createProduct - created new instance, before svc.createProduct", "product", p)
	created, err := h.svc.createProduct(p)
	if err != nil {
		h.logger.Debug("createProduct - error in svc.createProduct", "error", err)
		httputils.WriteErrorJSON(w, err, httperrs.ErrCodeUnknown, http.StatusBadRequest)
		return
	}
	h.logger.Debug("successfuly created a product", "product", created)
	httputils.WriteJSON(w, created)
}

func (h *productHandlers) getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not implemented yet."))
}

func (h *productHandlers) updateProduct(w http.ResponseWriter, r *http.Request) {
}
