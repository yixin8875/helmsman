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

var _ SnapshotsHandler = (*snapshotsHandler)(nil)

// SnapshotsHandler defining the handler interface
type SnapshotsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type snapshotsHandler struct {
	iDao dao.SnapshotsDao
}

// NewSnapshotsHandler creating the handler interface
func NewSnapshotsHandler() SnapshotsHandler {
	return &snapshotsHandler{
		iDao: dao.NewSnapshotsDao(
			database.GetDB(), // db driver is sqlite
			cache.NewSnapshotsCache(database.GetCacheType()),
		),
	}
}

// Create a new snapshots
// @Summary Create a new snapshots
// @Description Creates a new snapshots entity using the provided data in the request body.
// @Tags snapshots
// @Accept json
// @Produce json
// @Param data body types.CreateSnapshotsRequest true "snapshots information"
// @Success 200 {object} types.CreateSnapshotsReply{}
// @Router /api/v1/snapshots [post]
// @Security BearerAuth
func (h *snapshotsHandler) Create(c *gin.Context) {
	form := &types.CreateSnapshotsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	snapshots := &model.Snapshots{}
	err = copier.Copy(snapshots, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateSnapshots)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, snapshots)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": snapshots.ID})
}

// DeleteByID delete a snapshots by id
// @Summary Delete a snapshots by id
// @Description Deletes a existing snapshots identified by the given id in the path.
// @Tags snapshots
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteSnapshotsByIDReply{}
// @Router /api/v1/snapshots/{id} [delete]
// @Security BearerAuth
func (h *snapshotsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getSnapshotsIDFromPath(c)
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

// UpdateByID update a snapshots by id
// @Summary Update a snapshots by id
// @Description Updates the specified snapshots by given id in the path, support partial update.
// @Tags snapshots
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateSnapshotsByIDRequest true "snapshots information"
// @Success 200 {object} types.UpdateSnapshotsByIDReply{}
// @Router /api/v1/snapshots/{id} [put]
// @Security BearerAuth
func (h *snapshotsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getSnapshotsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateSnapshotsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	snapshots := &model.Snapshots{}
	err = copier.Copy(snapshots, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDSnapshots)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, snapshots)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a snapshots by id
// @Summary Get a snapshots by id
// @Description Gets detailed information of a snapshots specified by the given id in the path.
// @Tags snapshots
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetSnapshotsByIDReply{}
// @Router /api/v1/snapshots/{id} [get]
// @Security BearerAuth
func (h *snapshotsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getSnapshotsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	snapshots, err := h.iDao.GetByID(ctx, id)
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

	data := &types.SnapshotsObjDetail{}
	err = copier.Copy(data, snapshots)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDSnapshots)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"snapshots": data})
}

// List get a paginated list of snapshotss by custom conditions
// @Summary Get a paginated list of snapshotss by custom conditions
// @Description Returns a paginated list of snapshots based on query filters, including page number and size.
// @Tags snapshots
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListSnapshotssReply{}
// @Router /api/v1/snapshots/list [post]
// @Security BearerAuth
func (h *snapshotsHandler) List(c *gin.Context) {
	form := &types.ListSnapshotssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	snapshotss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertSnapshotss(snapshotss)
	if err != nil {
		response.Error(c, ecode.ErrListSnapshots)
		return
	}

	response.Success(c, gin.H{
		"snapshotss": data,
		"total":      total,
	})
}

func getSnapshotsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertSnapshots(snapshots *model.Snapshots) (*types.SnapshotsObjDetail, error) {
	data := &types.SnapshotsObjDetail{}
	err := copier.Copy(data, snapshots)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertSnapshotss(fromValues []*model.Snapshots) ([]*types.SnapshotsObjDetail, error) {
	toValues := []*types.SnapshotsObjDetail{}
	for _, v := range fromValues {
		data, err := convertSnapshots(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
