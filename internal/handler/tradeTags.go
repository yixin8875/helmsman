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

var _ TradeTagsHandler = (*tradeTagsHandler)(nil)

// TradeTagsHandler defining the handler interface
type TradeTagsHandler interface {
	Create(c *gin.Context)
	DeleteByTradeID(c *gin.Context)
	UpdateByTradeID(c *gin.Context)
	GetByTradeID(c *gin.Context)
	List(c *gin.Context)
}

type tradeTagsHandler struct {
	iDao dao.TradeTagsDao
}

// NewTradeTagsHandler creating the handler interface
func NewTradeTagsHandler() TradeTagsHandler {
	return &tradeTagsHandler{
		iDao: dao.NewTradeTagsDao(
			database.GetDB(), // db driver is sqlite
			cache.NewTradeTagsCache(database.GetCacheType()),
		),
	}
}

// Create a new tradeTags
// @Summary Create a new tradeTags
// @Description Creates a new tradeTags entity using the provided data in the request body.
// @Tags tradeTags
// @Accept json
// @Produce json
// @Param data body types.CreateTradeTagsRequest true "tradeTags information"
// @Success 200 {object} types.CreateTradeTagsReply{}
// @Router /api/v1/tradeTags [post]
// @Security BearerAuth
func (h *tradeTagsHandler) Create(c *gin.Context) {
	form := &types.CreateTradeTagsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	tradeTags := &model.TradeTags{}
	err = copier.Copy(tradeTags, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateTradeTags)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, tradeTags)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"tradeID": tradeTags.TradeID})
}

// DeleteByTradeID delete a tradeTags by tradeID
// @Summary Delete a tradeTags by tradeID
// @Description Deletes a existing tradeTags identified by the given tradeID in the path.
// @Tags tradeTags
// @Accept json
// @Produce json
// @Param tradeID path string true "tradeID"
// @Success 200 {object} types.DeleteTradeTagsByTradeIDReply{}
// @Router /api/v1/tradeTags/{tradeID} [delete]
// @Security BearerAuth
func (h *tradeTagsHandler) DeleteByTradeID(c *gin.Context) {
	tradeID, isAbort := getTradeTagsTradeIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByTradeID(ctx, tradeID)
	if err != nil {
		logger.Error("DeleteByTradeID error", logger.Err(err), logger.Any("tradeID", tradeID), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByTradeID update a tradeTags by tradeID
// @Summary Update a tradeTags by tradeID
// @Description Updates the specified tradeTags by given tradeID in the path, support partial update.
// @Tags tradeTags
// @Accept json
// @Produce json
// @Param tradeID path string true "tradeID"
// @Param data body types.UpdateTradeTagsByTradeIDRequest true "tradeTags information"
// @Success 200 {object} types.UpdateTradeTagsByTradeIDReply{}
// @Router /api/v1/tradeTags/{tradeID} [put]
// @Security BearerAuth
func (h *tradeTagsHandler) UpdateByTradeID(c *gin.Context) {
	tradeID, isAbort := getTradeTagsTradeIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateTradeTagsByTradeIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.TradeID = tradeID

	tradeTags := &model.TradeTags{}
	err = copier.Copy(tradeTags, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByTradeIDTradeTags)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByTradeID(ctx, tradeTags)
	if err != nil {
		logger.Error("UpdateByTradeID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByTradeID get a tradeTags by tradeID
// @Summary Get a tradeTags by tradeID
// @Description Gets detailed information of a tradeTags specified by the given tradeID in the path.
// @Tags tradeTags
// @Param tradeID path string true "tradeID"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetTradeTagsByTradeIDReply{}
// @Router /api/v1/tradeTags/{tradeID} [get]
// @Security BearerAuth
func (h *tradeTagsHandler) GetByTradeID(c *gin.Context) {
	tradeID, isAbort := getTradeTagsTradeIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	tradeTags, err := h.iDao.GetByTradeID(ctx, tradeID)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByTradeID not found", logger.Err(err), logger.Any("tradeID", tradeID), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByTradeID error", logger.Err(err), logger.Any("tradeID", tradeID), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.TradeTagsObjDetail{}
	err = copier.Copy(data, tradeTags)
	if err != nil {
		response.Error(c, ecode.ErrGetByTradeIDTradeTags)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"tradeTags": data})
}

// List get a paginated list of tradeTags by custom conditions
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
// @Summary Get a paginated list of tradeTags by custom conditions
// @Description Returns a paginated list of tradeTags based on query filters, including page number and size.
// @Tags tradeTags
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListTradeTagsReply{}
// @Router /api/v1/tradeTags/list [post]
// @Security BearerAuth
func (h *tradeTagsHandler) List(c *gin.Context) {
	form := &types.ListTradeTagsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	tradeTags, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertTradeTags(tradeTags)
	if err != nil {
		response.Error(c, ecode.ErrListTradeTags)
		return
	}

	response.Success(c, gin.H{
		"tradeTags": data,
		"total":     total,
	})
}

func getTradeTagsTradeIDFromPath(c *gin.Context) (int, bool) {
	tradeIDStr := c.Param("tradeID")

	tradeID, err := utils.StrToIntE(tradeIDStr)
	if err != nil || tradeIDStr == "" {
		logger.Warn("StrToIntE error: ", logger.String("tradeIDStr", tradeIDStr), middleware.GCtxRequestIDField(c))
		return 0, true
	}
	return tradeID, false

}

func convertTradeTag(tradeTags *model.TradeTags) (*types.TradeTagsObjDetail, error) {
	data := &types.TradeTagsObjDetail{}
	err := copier.Copy(data, tradeTags)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertTradeTags(fromValues []*model.TradeTags) ([]*types.TradeTagsObjDetail, error) {
	toValues := []*types.TradeTagsObjDetail{}
	for _, v := range fromValues {
		data, err := convertTradeTag(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
