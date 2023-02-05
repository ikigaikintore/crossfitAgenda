package grpc

import "github.com/ervitis/crossfitAgenda/service/domain"

func Into(st domain.Status) Status200Response_IdEnum {
	switch st {
	case domain.Finished:
		return Status200Response_IdEnum_FINISHED
	case domain.Working:
		return Status200Response_IdEnum_WORKING
	default:
		return Status200Response_IdEnum_FAILED
	}
}
