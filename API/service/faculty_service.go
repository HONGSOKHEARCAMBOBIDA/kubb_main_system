package service

import (
	"context"
	"mysql/request"
	"mysql/response"
)

type FacultyService interface {
	GetFaculty(ctx context.Context, pf request.Pagination, filter map[string]string) ([]response.FacultyResponse, error)
}
