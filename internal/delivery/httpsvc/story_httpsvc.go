package httpsvc

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	g.Use(checkErrorMiddleware)

	g.GET("", h.findAll)
	g.GET("/:id", h.findByID)
	g.POST("", h.create, bearerAuthMiddleware)
	g.PUT("/:id", h.update, bearerAuthMiddleware)
	g.DELETE("/:id", h.delete, bearerAuthMiddleware)
}

func (h *StoryHandler) findAll(c echo.Context) error {
	opt := new(model.StoryOptions)
	if err := c.Bind(opt); err != nil {
		return model.NewErrorBadRequest(err.Error())
	}

	results, totalItems, err := h.storyUsecase.FindAll(c.Request().Context(), opt)
	if err != nil {
		return err
	}

	resp := response{
		Status: "success",
		Metadata: map[string]any{
			"total_items": totalItems,
		},
		Data: results,
	}

	if len(results) > 0 {
		lastCursor := results[len(results)-1].CreatedAt.Format(time.RFC3339)
		encodedCursor := base64.StdEncoding.EncodeToString([]byte(lastCursor))
		c.Response().Header().Set("X-Cursor", encodedCursor)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *StoryHandler) findByID(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, "find by id")
}

func (h *StoryHandler) create(c echo.Context) error {
	var bodyReq model.Story
	if err := c.Bind(&bodyReq); err != nil {
		return model.NewErrorBadRequest(err.Error())
	}

	claims, ok := c.Request().Context().Value(model.JWTKey).(*model.CustomClaims)
	if !ok {
		return model.NewErrorUnAuthorized("invalid token")
	}
	bodyReq.Author.ID = claims.UserID

	// validate struct
	err := validate.Struct(bodyReq)
	if err != nil {
		return model.NewErrorBadRequest(fmt.Sprintf("error validation: %s", err.Error()))
	}

	insertedData, err := h.storyUsecase.Create(c.Request().Context(), bodyReq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, insertedData)
}

func (h *StoryHandler) update(c echo.Context) error {
	var bodyReq model.Story
	if err := c.Bind(&bodyReq); err != nil {
		return model.NewErrorBadRequest(err.Error())
	}

	claims, ok := c.Request().Context().Value(model.JWTKey).(*model.CustomClaims)
	if !ok {
		return model.NewErrorUnAuthorized("invalid token")
	}
	bodyReq.Author.ID = claims.UserID

	storyID, _ := strconv.Atoi(c.Param("id"))
	bodyReq.ID = int64(storyID)

	updatedData, err := h.storyUsecase.Update(c.Request().Context(), bodyReq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedData)
}

func (h *StoryHandler) delete(c echo.Context) error {
	// TODO
	return c.JSON(http.StatusOK, "delete by id")
}
