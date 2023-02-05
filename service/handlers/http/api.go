package http

import (
	"errors"
	"github.com/ervitis/crossfitAgenda/service/domain"
	"github.com/ervitis/crossfitAgenda/service/handlers/models"
	"github.com/ervitis/crossfitAgenda/service/usecases"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

type CrossfitHandler ServerInterface

type crossfitHandler struct {
	agenda usecases.AgendaCrossfit
}

func (c crossfitHandler) CredentialsGoogle(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, GoogleCredentials{Link: c.agenda.GetCredentialsUrl()})
}

func (c crossfitHandler) SetTokenGoogle(ctx echo.Context) error {
	var body GoogleToken
	if err := ctx.Bind(body); err != nil {
		return ctx.JSON(http.StatusBadRequest, Error{
			Date:    time.Now().Format(time.RFC3339),
			Message: "invalid body",
			Status:  http.StatusBadRequest,
		})
	}
	if body.Token == "" {
		return ctx.JSON(http.StatusBadRequest, Error{
			Date:    time.Now().Format(time.RFC3339),
			Message: "empty token body",
			Status:  http.StatusBadRequest,
		})
	}
	err := c.agenda.SaveToken(ctx.Request().Context(), body.Token)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Error{
			Date:    time.Now().Format(time.RFC3339),
			Message: "error saving token",
			Status:  http.StatusInternalServerError,
		})
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (c crossfitHandler) StartCrossfitAgenda(ctx echo.Context) error {
	if err := c.agenda.Book(ctx.Request().Context()); err != nil {
		switch {
		case errors.Is(err, domain.ErrTokenCredentialsExpired):
			return ctx.JSON(http.StatusUnauthorized, toError(models.NewErrNotAuth()))
		default:
			log.Printf("booking error: %+v\n", err)
			return ctx.JSON(http.StatusInternalServerError, toError(models.NewErrInternal()))
		}
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (c crossfitHandler) Status(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, toStatus(c.agenda.Status(), ctx.Request().Header.Get(echo.HeaderXRequestID)))
}

func toStatus(st domain.ProcessStatus, reqID string) ProcessStatus {
	apiSt := ProcessStatus{
		Complete: st.IsComplete(),
		Date:     uint64(time.Now().Unix()),
	}

	switch *st.Status {
	case domain.Finished:
		apiSt.Id = Finished
		apiSt.Detail = Finished.ToString()
		break
	case domain.Working:
		apiSt.Id = Working
		apiSt.Detail = Working.ToString()
		break
	case domain.Failed:
		apiSt.Id = Failed
		apiSt.Detail = Failed.ToString()
		break
	}
	return apiSt
}

func toError(err error) Error {
	cErr, ok := err.(models.AgendaError)
	if !ok {
		return Error{
			Date:    time.Now().Format(time.RFC3339),
			Message: "internal error",
			Status:  http.StatusInternalServerError,
		}
	}
	return Error{
		Date:    time.Now().Format(time.RFC3339),
		Message: cErr.Error(),
		Status:  uint32(cErr.Code()),
	}
}

func NewHandler(agenda usecases.AgendaCrossfit) CrossfitHandler {
	return &crossfitHandler{agenda}
}
