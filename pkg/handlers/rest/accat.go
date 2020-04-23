package rest

import (
	"context"
	"net/http"

	"github.com/spf13/cast"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"

	"github.com/hellerox/AcCatalog/model"
)

type errorResponse struct {
	error error
}

func decodeGetAllMaterialsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getAllMaterialsRequest
	return req, nil
}

type getAllMaterialsRequest struct{}

type getAllMaterialsResponse struct {
	Materials []model.Material `json:"materials"`
}

func makeGetAllMaterialsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ms, err := s.GetAllMaterials()
		if err != nil {
			return errorResponse{error: err}, err
		}

		return getAllMaterialsResponse{Materials: ms}, nil
	}
}

type getFullCostumeRequest struct {
	CostumeID int `json:"costumeID"`
}

type getFullCostumeResponse struct {
	Costume *model.Costume `json:"costume"`
}

func decodeGetFullCostumeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	req := getFullCostumeRequest{CostumeID: cast.ToInt(params["costume_id"])}

	return req, nil
}

func makeGetFullCostume(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		c, err := s.GetFullCostume(request.(getFullCostumeRequest).CostumeID)
		if err != nil {
			return errorResponse{error: err}, err
		}

		return getFullCostumeResponse{Costume: c}, nil
	}
}

type getAllCostumesRequest struct{}

type getAllCostumesResponse struct {
	Costumes []model.Costume `json:"costumes"`
}

func decodeGetAllCostumesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getAllCostumesRequest
	return req, nil
}

func makeGetAllCostumes(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cs, err := s.GetAllCostumes()
		if err != nil {
			return errorResponse{error: err}, err
		}

		return getAllCostumesResponse{Costumes: cs}, nil
	}
}

type createMaterialRequest struct {
	Material model.Material `json:"material"`
}

type createMaterialResponse struct {
	MaterialID int64 `json:"materialID"`
}

func decodeCreateMaterialRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := createMaterialRequest{Material: model.Material{}}

	return req, nil
}

func makeCreateMaterial(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		mID, err := s.CreateMaterial(request.(createMaterialRequest).Material)
		if err != nil {
			return errorResponse{error: err}, err
		}

		return createMaterialResponse{MaterialID: mID}, nil
	}
}
