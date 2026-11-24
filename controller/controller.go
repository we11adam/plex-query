package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"plex-query/database"

	"github.com/labstack/echo/v4"
)

const (
	IMDB = "imdb"
	TMDB = "tmdb"
	TVDB = "tvdb"
)

type Controller struct {
	queries database.Querier
}

func New(q database.Querier) *Controller {
	return &Controller{
		queries: q,
	}
}

func (ctrl *Controller) GetMediaByTag(c echo.Context) error {
	typ := c.QueryParam("type")
	id := c.QueryParam("id")

	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "id is required")
	}

	if typ != "" && typ != IMDB && typ != TMDB && typ != TVDB {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid type: %s", typ))
	}

	if typ == "" {
		typ = IMDB
	}

	if typ == IMDB && !strings.HasPrefix(id, "tt") {
		id = "tt" + id
	}

	tag := fmt.Sprintf("%s://%s", typ, id)

	files, err := ctrl.queries.GetMediaByTag(c.Request().Context(), sql.NullString{String: tag, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to get media by tag: %w", err)
	}

	ret := make([]string, 0, len(files))
	for _, file := range files {
		if file.Valid {
			ret = append(ret, file.String)
		}
	}

	return c.JSON(http.StatusOK, ret)
}
