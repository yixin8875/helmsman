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

var _ TradesHandler = (*tradesHandler)(nil)

// TradesHandler defining the handler interface
type TradesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type tradesHandler struct {
	iDao dao.TradesDao
}

// NewTradesHandler creating the handler interface
func NewTradesHandler() TradesHandler {
	return &tradesHandler{
		iDao: dao.NewTradesDao(
			database.GetDB(), // db driver is sqlite
			cache.NewTradesCache(database.GetCacheType()),
		),
	}
}

// Create a new trades
// @Summary Create a new trades
// @Description Creates a new trades entity using the provided data in the request body.
// @Tags trades
// @Accept json
// @Produce json
// @Param data body types.CreateTradesRequest true "trades information"
// @Success 200 {object} types.CreateTradesReply{}
// @Router /api/v1/trades [post]
// @Security BearerAuth
func (h *tradesHandler) Create(c *gin.Context) {
	form := &types.CreateTradesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	trades := &model.Trades{}
	err = copier.Copy(trades, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateTrades)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, trades)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": trades.ID})
}

// DeleteByID delete a trades by id
// @Summary Delete a trades by id
// @Description Deletes a existing trades identified by the given id in the path.
// @Tags trades
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteTradesByIDReply{}
// @Router /api/v1/trades/{id} [delete]
// @Security BearerAuth
func (h *tradesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getTradesIDFromPath(c)
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

// UpdateByID update a trades by id
// @Summary Update a trades by id
// @Description Updates the specified trades by given id in the path, support partial update.
// @Tags trades
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateTradesByIDRequest true "trades information"
// @Success 200 {object} types.UpdateTradesByIDReply{}
// @Router /api/v1/trades/{id} [put]
// @Security BearerAuth
func (h *tradesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getTradesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateTradesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	trades := &model.Trades{}
	err = copier.Copy(trades, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDTrades)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, trades)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a trades by id
// @Summary Get a trades by id
// @Description Gets detailed information of a trades specified by the given id in the path.
// @Tags trades
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetTradesByIDReply{}
// @Router /api/v1/trades/{id} [get]
// @Security BearerAuth
func (h *tradesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getTradesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	trades, err := h.iDao.GetByID(ctx, id)
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

	data := &types.TradesObjDetail{}
	err = copier.Copy(data, trades)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDTrades)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"trades": data})
}

// List get a paginated list of tradess by custom conditions
// @Summary Get a paginated list of tradess by custom conditions
// @Description Returns a paginated list of trades based on query filters, including page number and size.
// @Tags trades
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListTradessReply{}
// @Router /api/v1/trades/list [post]
// @Security BearerAuth
func (h *tradesHandler) List(c *gin.Context) {
	form := &types.ListTradessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	tradess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertTradess(tradess)
	if err != nil {
		response.Error(c, ecode.ErrListTrades)
		return
	}

	response.Success(c, gin.H{
		"tradess": data,
		"total":   total,
	})
}

func getTradesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertTrades(trades *model.Trades) (*types.TradesObjDetail, error) {
	data := &types.TradesObjDetail{}
	err := copier.Copy(data, trades)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertTradess(fromValues []*model.Trades) ([]*types.TradesObjDetail, error) {
	toValues := []*types.TradesObjDetail{}
	for _, v := range fromValues {
		data, err := convertTrades(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
