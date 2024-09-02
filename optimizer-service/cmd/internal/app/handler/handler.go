package handler

import (
	"net/http"
	"optimizer-service/cmd/internal/types"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Container *types.AppContainer
}

func NewHandler(c *types.AppContainer) *Handler {
	return &Handler{Container: c}
}

func (h *Handler) Index(c echo.Context) error {
	response := map[string]interface{}{
		"message": "Index Welcome",
	}

	return h.Container.Utils.WriteSuccessResponse(c, http.StatusOK, "success", response)
}
