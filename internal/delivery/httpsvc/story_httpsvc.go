package httpsvc

import (
	"net/http"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/labstack/echo/v4"
)

type StoryHandler struct {
	storyUsecase model.StoryUsecase
}

// NewStoryHandler is used to init story handler
func NewStoryHandler(storyUsecase model.StoryUsecase) *StoryHandler {
	return &StoryHandler{
		storyUsecase: storyUsecase,
	}
}

// RegisterRoutes is used to register routes of story handler
func (h *StoryHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/stories")

	g.GET("", h.findAll)
	g.GET("/:id", h.findByID)
	g.POST("", h.create)
	g.PUT("/:id", h.update)
	g.DELETE("/:id", h.delete)
}

func (h *StoryHandler) findAll(c echo.Context) error {
	results, _, err := h.storyUsecase.FindAll(c.Request().Context(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, results)
}

func (h *StoryHandler) findByID(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, "find by id")
}

func (h *StoryHandler) create(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, "create")
}

func (h *StoryHandler) update(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, "update by id")
}

func (h *StoryHandler) delete(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, "delete by id")
}
