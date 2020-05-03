package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/prometheus/common/log"

	"github.com/go-kit/kit/endpoint"
	gokithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/hellerox/AcCatalog/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Endpoints includes all endpoints for the service.
type Endpoints struct {
	GetAllMaterialsEndpoint endpoint.Endpoint
	GetFullCostume          endpoint.Endpoint
	GetAllCostumes          endpoint.Endpoint
	CreateMaterial          endpoint.Endpoint
	CreateCostume           endpoint.Endpoint
}

type restResponse struct {
	statusCode int
	Errors     []string    `json:"errors,omitempty"`
	Data       interface{} `json:"data"`
}

// StatusCode func to implement StatusCoder interface
func (rr restResponse) StatusCode() int {
	return rr.statusCode
}

type Service interface {
	GetAllMaterials() (ms []model.Material, err error)
	GetFullCostume(cID int) (c *model.Costume, err error)
	GetAllCostumes() (cs []model.Costume, err error)
	CreateMaterial(material model.Material) (mID int64, err error)
	CreateCostume(c model.Costume) (cID int, err error)
}

// MakeHTTPHandlers makes handlers
func MakeHTTPHandlers(s Service) http.Handler {
	r := mux.NewRouter()

	endpoints := makeServerEndpoints(s)

	GetAllMaterials := gokithttp.NewServer(
		endpoints.GetAllMaterialsEndpoint,
		decodeGetAllMaterialsRequest,
		encodeResponse,
		gokithttp.ServerErrorEncoder(encodeError),
		gokithttp.ServerBefore(gokithttp.PopulateRequestContext),
	)

	GetFullCostume := gokithttp.NewServer(
		endpoints.GetFullCostume,
		decodeGetFullCostumeRequest,
		encodeResponse,
		gokithttp.ServerErrorEncoder(encodeError),
		gokithttp.ServerBefore(gokithttp.PopulateRequestContext),
	)

	GetAllCostumes := gokithttp.NewServer(
		endpoints.GetAllCostumes,
		decodeGetAllCostumesRequest,
		encodeResponse,
		gokithttp.ServerErrorEncoder(encodeError),
		gokithttp.ServerBefore(gokithttp.PopulateRequestContext),
	)

	CreateMaterial := gokithttp.NewServer(
		endpoints.CreateMaterial,
		decodeCreateMaterialRequest,
		encodeResponse,
		gokithttp.ServerErrorEncoder(encodeError),
		gokithttp.ServerBefore(gokithttp.PopulateRequestContext),
	)

	r.Methods("GET").Path("/v1/materials").Handler(GetAllMaterials)
	r.Methods("GET").Path("/v1/costumes/{costume_id:[0-9]+}").Handler(GetFullCostume)
	r.Methods("GET").Path("/v1/costumes").Handler(GetAllCostumes)
	r.Methods("POST").Path("/v1/materials").Handler(CreateMaterial)

	// prometheus endpoint
	r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

	healthcheckHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"status":"OK"}`))
			if err != nil {
				log.Errorf("error writing response: %+v", err)
			}
		},
	)

	r.Methods("GET").Path("/healthcheck").Handler(healthcheckHandler)

	return r
}

// makeServerEndpoints creates the endpoints needed for the service.
func makeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetAllMaterialsEndpoint: makeGetAllMaterialsEndpoint(s),
		GetFullCostume:          makeGetFullCostume(s),
		GetAllCostumes:          makeGetAllCostumes(s),
		CreateMaterial:          makeCreateMaterial(s),
		CreateCostume:           makeCreateCostume(s),
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	resp := restResponse{
		Data:       nil,
		Errors:     []string{err.Error()},
		statusCode: http.StatusInternalServerError, // default status code
	}

	if statusCoder, ok := err.(gokithttp.StatusCoder); ok {
		resp.statusCode = statusCoder.StatusCode()
	}

	if encodeErr := gokithttp.EncodeJSONResponse(ctx, w, resp); encodeErr != nil {
		logrus.WithFields(logrus.Fields{
			"function": "errorEncoder",
			"step":     "EncodeJSONResponse",
			"error":    encodeErr,
		}).Error()
		gokithttp.DefaultErrorEncoder(ctx, encodeErr, w)
	}
}
