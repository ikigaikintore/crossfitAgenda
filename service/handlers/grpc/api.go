package grpc

import (
	"context"
	"errors"
	"github.com/ervitis/crossfitAgenda/ports"
	"github.com/ervitis/crossfitAgenda/service/domain"
	"github.com/ervitis/crossfitAgenda/service/handlers/models"
	"github.com/ervitis/crossfitAgenda/service/usecases"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CrossfitAgendaSvc DefaultServiceServer

type handler struct {
	svc usecases.AgendaCrossfit
}

func (h handler) CredentialsGoogle(_ context.Context, _ *emptypb.Empty) (*CredentialsGoogle200Response, error) {
	return &CredentialsGoogle200Response{Link: h.svc.GetCredentialsUrl()}, nil
}

func (h handler) SetTokenGoogle(ctx context.Context, request *SetTokenGoogleRequest) (*emptypb.Empty, error) {
	if err := h.svc.SaveToken(ctx, request.GetToken()); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h handler) StartCrossfitAgenda(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := h.svc.Book(ctx); err != nil {
		switch {
		case errors.Is(err, domain.ErrTokenCredentialsExpired):
			return nil, models.NewErrNotAuth()
		default:
			return nil, models.NewErrInternal()
		}
	}
	return nil, nil
}

func (h handler) Status(_ context.Context, _ *emptypb.Empty) (*Status200Response, error) {
	status := h.svc.Status()
	return &Status200Response{
		Id:       Into(*status.Status),
		Detail:   status.Detail,
		Date:     uint64(status.Date.Unix()),
		Complete: status.IsComplete(),
	}, nil
}

func New(rm ports.ResourceManager, mgr ports.IManager) CrossfitAgendaSvc {
	return &handler{usecases.New(rm, mgr)}
}
