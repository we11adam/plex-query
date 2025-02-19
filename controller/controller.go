package controller

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"plex-query/database"
	"strings"
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
		return fmt.Errorf("id is required")
	}

	if typ != "" && typ != IMDB && typ != TMDB && typ != TVDB {
		return fmt.Errorf("invalid type: %s", typ)
	}

	if typ == "" {
		typ = IMDB
	}

	if typ == IMDB && !strings.HasPrefix(id, "tt") {
		id = "tt" + id
	}

	tag := fmt.Sprintf("%s://%s", typ, id)

	files, err := ctrl.queries.GetMediaByTag(context.TODO(), sql.NullString{String: tag, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to get media by tag: %v", err)
	}

	var ret []string
	for _, file := range files {
		if file.Valid {
			ret = append(ret, file.String)
		}
	}

	return c.JSON(200, ret)
}
