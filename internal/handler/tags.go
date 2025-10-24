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

var _ TagsHandler = (*tagsHandler)(nil)

// TagsHandler defining the handler interface
type TagsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type tagsHandler struct {
	iDao dao.TagsDao
}

// NewTagsHandler creating the handler interface
func NewTagsHandler() TagsHandler {
	return &tagsHandler{
		iDao: dao.NewTagsDao(
			database.GetDB(), // db driver is sqlite
			cache.NewTagsCache(database.GetCacheType()),
		),
	}
}

// Create a new tags
// @Summary Create a new tags
// @Description Creates a new tags entity using the provided data in the request body.
// @Tags tags
// @Accept json
// @Produce json
// @Param data body types.CreateTagsRequest true "tags information"
// @Success 200 {object} types.CreateTagsReply{}
// @Router /api/v1/tags [post]
// @Security BearerAuth
func (h *tagsHandler) Create(c *gin.Context) {
	form := &types.CreateTagsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	tags := &model.Tags{}
	err = copier.Copy(tags, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateTags)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, tags)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": tags.ID})
}

// DeleteByID delete a tags by id
// @Summary Delete a tags by id
// @Description Deletes a existing tags identified by the given id in the path.
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteTagsByIDReply{}
// @Router /api/v1/tags/{id} [delete]
// @Security BearerAuth
func (h *tagsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getTagsIDFromPath(c)
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

// UpdateByID update a tags by id
// @Summary Update a tags by id
// @Description Updates the specified tags by given id in the path, support partial update.
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateTagsByIDRequest true "tags information"
// @Success 200 {object} types.UpdateTagsByIDReply{}
// @Router /api/v1/tags/{id} [put]
// @Security BearerAuth
func (h *tagsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getTagsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateTagsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	tags := &model.Tags{}
	err = copier.Copy(tags, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDTags)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, tags)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a tags by id
// @Summary Get a tags by id
// @Description Gets detailed information of a tags specified by the given id in the path.
// @Tags tags
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetTagsByIDReply{}
// @Router /api/v1/tags/{id} [get]
// @Security BearerAuth
func (h *tagsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getTagsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	tags, err := h.iDao.GetByID(ctx, id)
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

	data := &types.TagsObjDetail{}
	err = copier.Copy(data, tags)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDTags)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"tags": data})
}

// List get a paginated list of tagss by custom conditions
// @Summary Get a paginated list of tagss by custom conditions
// @Description Returns a paginated list of tags based on query filters, including page number and size.
// @Tags tags
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListTagssReply{}
// @Router /api/v1/tags/list [post]
// @Security BearerAuth
func (h *tagsHandler) List(c *gin.Context) {
	form := &types.ListTagssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	tagss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertTagss(tagss)
	if err != nil {
		response.Error(c, ecode.ErrListTags)
		return
	}

	response.Success(c, gin.H{
		"tagss": data,
		"total": total,
	})
}

func getTagsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertTags(tags *model.Tags) (*types.TagsObjDetail, error) {
	data := &types.TagsObjDetail{}
	err := copier.Copy(data, tags)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertTagss(fromValues []*model.Tags) ([]*types.TagsObjDetail, error) {
	toValues := []*types.TagsObjDetail{}
	for _, v := range fromValues {
		data, err := convertTags(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
