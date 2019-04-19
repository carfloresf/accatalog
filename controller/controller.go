package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/buaazp/fasthttprouter"
	"github.com/hellerox/AcCatalog/model"
	"github.com/hellerox/AcCatalog/service"
	"github.com/hellerox/AcCatalog/storage"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// Controller controller
type Controller struct {
	Router  *fasthttprouter.Router
	Storage storage.Storage
	Service service.Service
}

// InitializeRoutes route initialize
func (c *Controller) InitializeRoutes() {

	c.Router.HandleMethodNotAllowed = true
	c.Router.NotFound = c.notFound
	c.Router.MethodNotAllowed = c.methodNotAllowed
	c.Router.PanicHandler = c.panic

	c.Router.GET("/healthcheck", c.healthcheck)
	c.Router.POST("/materials", c.basicAuth(c.createMaterial))
	c.Router.POST("/costumes", c.basicAuth(c.createCostume))
	c.Router.GET("/costumes/:costume_id", c.basicAuth(c.getCostume))
	c.Router.GET("/costumes", c.basicAuth(c.getAllCostumes))
	c.Router.POST("/materialTypes", c.basicAuth(c.createMaterialType))
	c.Router.POST("/costumeMaterial", c.basicAuth(c.createCostumeMaterialRelation))

	log.Info("starting routes")
}

func (c *Controller) basicAuth(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		apikey := string(ctx.Request.Header.Peek("apikey"))
		log.Infof("apikey received: %s", apikey)
		permission, user, _ := c.Storage.GetPermission(apikey)
		if permission == "a" {
			// Delegate request to the given handle
			log.Debugf("login user %s", user)
			h(ctx)
			return
		}

		c.unauthorized(ctx)
	})
}

func (c *Controller) createMaterial(ctx *fasthttp.RequestCtx) {
	var m model.Material
	if err := json.Unmarshal(ctx.Request.Body(), &m); err != nil {
		respond(ctx, fasthttp.StatusBadRequest, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	materialID, err := c.Storage.InsertMaterial(m)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respond(ctx, http.StatusCreated, fmt.Sprintf(`{"materialID": %d}`, materialID))
}

func (c *Controller) createMaterialType(ctx *fasthttp.RequestCtx) {
	var m model.MaterialType
	if err := json.Unmarshal(ctx.Request.Body(), &m); err != nil {
		respond(ctx, fasthttp.StatusBadRequest, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	materialTypeID, err := c.Storage.InsertMaterialType(m)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respond(ctx, http.StatusCreated, fmt.Sprintf(`{"materialTypeID": %d}`, materialTypeID))
}

func (c *Controller) createCostume(ctx *fasthttp.RequestCtx) {
	var co model.Costume
	if err := json.Unmarshal(ctx.Request.Body(), &co); err != nil {
		respond(ctx, fasthttp.StatusBadRequest, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	costumeID, err := c.Storage.InsertCostume(co)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respond(ctx, http.StatusCreated, fmt.Sprintf(`{"costumeID": %d}`, costumeID))
}

func (c *Controller) createCostumeMaterialRelation(ctx *fasthttp.RequestCtx) {
	var cm model.CostumeMaterialRelation
	if err := json.Unmarshal(ctx.Request.Body(), &cm); err != nil {
		respond(ctx, fasthttp.StatusBadRequest, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	err := c.Storage.InsertCostumeMaterialRelation(cm)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respond(ctx, http.StatusCreated, fmt.Sprintf(`{"status": "OK"}`))
}

func (c *Controller) getCostume(ctx *fasthttp.RequestCtx) {
	costumeIDParam := ctx.UserValue("costume_id").(string)
	costumeID, err := strconv.Atoi(costumeIDParam)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}
	cm, err := c.Service.GetFullCostume(costumeID)
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respondInterface(ctx, http.StatusCreated, cm)
}

func (c *Controller) getAllCostumes(ctx *fasthttp.RequestCtx) {

	cs, err := c.Service.GetAllCostumes()
	if err != nil {
		respond(ctx, fasthttp.StatusInternalServerError, fmt.Sprintf(`{"error":"%s"}`, err))
		return
	}

	respondInterface(ctx, http.StatusCreated, cs)
}
