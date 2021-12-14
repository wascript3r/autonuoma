package http

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/autonuoma/pkg/cars"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
)

type HTTPHandler struct {
	carsUcase cars.Usecase
}

func NewHTTPHandler(r *httprouter.Router, fu cars.Usecase) {
	handler := &HTTPHandler{
		carsUcase: fu,
	}

	r.GET("/api/cars/list", handler.AllCars)
	r.POST("/api/cars/single", handler.SingleCar)
	r.POST("/api/cars/remove", handler.RemoveCar)
	r.POST("/api/cars/add", handler.AddCar)
	r.POST("/api/cars/update", handler.UpdateCar)
	r.POST("/api/cars/trips", handler.CarTrips)
	r.GET("/api/cars/statistics", handler.Statistics)
}

func serveError(w http.ResponseWriter, err error) {
	if err == cars.InvalidInputError {
		httpjson.BadRequestCustom(w, cars.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, cars.UnknownError)
	if code == cars.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) AllCars(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.carsUcase.GetAll(r.Context())
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) SingleCar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &cars.SingleCarReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.carsUcase.GetSingle(r.Context(), req.CarID)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) RemoveCar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &cars.SingleCarReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.carsUcase.RemoveCar(r.Context(), req.CarID)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) AddCar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &cars.AddCarReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.carsUcase.AddCar(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) UpdateCar(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &cars.UpdateCarReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.carsUcase.UpdateCar(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) CarTrips(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &cars.SingleCarReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.carsUcase.CarTrips(r.Context(), req.CarID)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) Statistics(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.carsUcase.Statistics(r.Context())
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
