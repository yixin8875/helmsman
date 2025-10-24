package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"helmsman/internal/cache"
	"helmsman/internal/dao"
	"helmsman/internal/database"
	"helmsman/internal/ecode"
	"helmsman/internal/model"
	"helmsman/internal/types"
)

var _ StrategiesHandler = (*strategiesHandler)(nil)

// StrategiesHandler defining the handler interface
type StrategiesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type strategiesHandler struct {
	iDao dao.StrategiesDao
}

// NewStrategiesHandler creating the handler interface
func NewStrategiesHandler() StrategiesHandler {
	return &strategiesHandler{
		iDao: dao.NewStrategiesDao(
			database.GetDB(), // db driver is sqlite
			cache.NewStrategiesCache(database.GetCacheType()),
		),
	}
}

// Create a new strategies
// @Summary Create a new strategies
// @Description Creates a new strategies entity using the provided data in the request body.
// @Tags strategies
// @Accept json
// @Produce json
// @Param data body types.CreateStrategiesRequest true "strategies information"
// @Success 200 {object} types.CreateStrategiesReply{}
// @Router /api/v1/strategies [post]
// @Security BearerAuth
func (h *strategiesHandler) Create(c *gin.Context) {
	form := &types.CreateStrategiesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	strategies := &model.Strategies{}
	err = copier.Copy(strategies, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateStrategies)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, strategies)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": strategies.ID})
}

// DeleteByID delete a strategies by id
// @Summary Delete a strategies by id
// @Description Deletes a existing strategies identified by the given id in the path.
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteStrategiesByIDReply{}
// @Router /api/v1/strategies/{id} [delete]
// @Security BearerAuth
func (h *strategiesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getStrategiesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a strategies by id
// @Summary Update a strategies by id
// @Description Updates the specified strategies by given id in the path, support partial update.
// @Tags strategies
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateStrategiesByIDRequest true "strategies information"
// @Success 200 {object} types.UpdateStrategiesByIDReply{}
// @Router /api/v1/strategies/{id} [put]
// @Security BearerAuth
func (h *strategiesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getStrategiesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateStrategiesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	strategies := &model.Strategies{}
	err = copier.Copy(strategies, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDStrategies)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, strategies)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a strategies by id
// @Summary Get a strategies by id
// @Description Gets detailed information of a strategies specified by the given id in the path.
// @Tags strategies
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetStrategiesByIDReply{}
// @Router /api/v1/strategies/{id} [get]
// @Security BearerAuth
func (h *strategiesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getStrategiesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	strategies, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.StrategiesObjDetail{}
	err = copier.Copy(data, strategies)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDStrategies)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"strategies": data})
}

// List get a paginated list of strategiess by custom conditions
// @Summary Get a paginated list of strategiess by custom conditions
// @Description Returns a paginated list of strategies based on query filters, including page number and size.
// @Tags strategies
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListStrategiessReply{}
// @Router /api/v1/strategies/list [post]
// @Security BearerAuth
func (h *strategiesHandler) List(c *gin.Context) {
	form := &types.ListStrategiessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	strategiess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertStrategiess(strategiess)
	if err != nil {
		response.Error(c, ecode.ErrListStrategies)
		return
	}

	response.Success(c, gin.H{
		"strategiess": data,
		"total":       total,
	})
}

func getStrategiesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertStrategies(strategies *model.Strategies) (*types.StrategiesObjDetail, error) {
	data := &types.StrategiesObjDetail{}
	err := copier.Copy(data, strategies)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertStrategiess(fromValues []*model.Strategies) ([]*types.StrategiesObjDetail, error) {
	toValues := []*types.StrategiesObjDetail{}
	for _, v := range fromValues {
		data, err := convertStrategies(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
