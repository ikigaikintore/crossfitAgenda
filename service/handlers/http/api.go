package http

import (
	"github.com/ervitis/crossfitAgenda/service/domain"
	"github.com/ervitis/crossfitAgenda/service/usecases"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

type CrossfitHandler ServerInterface

type crossfitHandler struct {
	agenda usecases.Crossfit
}

func (c crossfitHandler) StartCrossfitAgenda(ctx echo.Context) error {
	if err := c.agenda.Book(ctx.Request().Context()); err != nil {
		log.Printf("booking error: %+v\n", err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c crossfitHandler) Status(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, toStatus(c.agenda.Status(), ctx.Request().Header.Get(echo.HeaderXRequestID)))
}

func toStatus(st domain.Status, reqID string) Status {
	apiSt := Status{
		Complete: st.IsComplete(),
		Date:     time.Now().Unix(),
		ID:       reqID,
	}

	switch st {
	case domain.Finished:
		apiSt.Status = Finished
		apiSt.Detail = Finished.ToString()
		break
	case domain.Working:
		apiSt.Status = Working
		apiSt.Detail = Working.ToString()
		break
	case domain.Failed:
		apiSt.Status = Failed
		apiSt.Detail = Failed.ToString()
		break
	}
	return apiSt
}

func NewHandler(agenda usecases.Crossfit) CrossfitHandler {
	return &crossfitHandler{agenda}
}
