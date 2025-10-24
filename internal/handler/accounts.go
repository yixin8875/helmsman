package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

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

var _ AccountsHandler = (*accountsHandler)(nil)

// AccountsHandler defining the handler interface
type AccountsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type accountsHandler struct {
	iDao dao.AccountsDao
}

// NewAccountsHandler creating the handler interface
func NewAccountsHandler() AccountsHandler {
	return &accountsHandler{
		iDao: dao.NewAccountsDao(
			database.GetDB(), // db driver is sqlite
			cache.NewAccountsCache(database.GetCacheType()),
		),
	}
}

// Create a new accounts
// @Summary Create a new accounts
// @Description Creates a new accounts entity using the provided data in the request body.
// @Tags accounts
// @Accept json
// @Produce json
// @Param data body types.CreateAccountsRequest true "accounts information"
// @Success 200 {object} types.CreateAccountsReply{}
// @Router /api/v1/accounts [post]
// @Security BearerAuth
func (h *accountsHandler) Create(c *gin.Context) {
	form := &types.CreateAccountsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	accounts := &model.Accounts{}
	err = copier.Copy(accounts, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateAccounts)
		return
	}
	claim, ok := middleware.GetClaims(c)
	if !ok {
		response.Error(c, ecode.ErrCreateStrategies)
		return
	}
	accounts.UserID = cast.ToInt(claim.UID)
	accounts.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	accounts.UpdatedAt = accounts.CreatedAt
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, accounts)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": accounts.ID})
}

// DeleteByID delete a accounts by id
// @Summary Delete a accounts by id
// @Description Deletes a existing accounts identified by the given id in the path.
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteAccountsByIDReply{}
// @Router /api/v1/accounts/{id} [delete]
// @Security BearerAuth
func (h *accountsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getAccountsIDFromPath(c)
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

// UpdateByID update a accounts by id
// @Summary Update a accounts by id
// @Description Updates the specified accounts by given id in the path, support partial update.
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateAccountsByIDRequest true "accounts information"
// @Success 200 {object} types.UpdateAccountsByIDReply{}
// @Router /api/v1/accounts/{id} [put]
// @Security BearerAuth
func (h *accountsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getAccountsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateAccountsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	accounts := &model.Accounts{}
	err = copier.Copy(accounts, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDAccounts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, accounts)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a accounts by id
// @Summary Get a accounts by id
// @Description Gets detailed information of a accounts specified by the given id in the path.
// @Tags accounts
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetAccountsByIDReply{}
// @Router /api/v1/accounts/{id} [get]
// @Security BearerAuth
func (h *accountsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getAccountsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	accounts, err := h.iDao.GetByID(ctx, id)
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

	data := &types.AccountsObjDetail{}
	err = copier.Copy(data, accounts)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDAccounts)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"accounts": data})
}

// List get a paginated list of accountss by custom conditions
// @Summary Get a paginated list of accountss by custom conditions
// @Description Returns a paginated list of accounts based on query filters, including page number and size.
// @Tags accounts
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListAccountssReply{}
// @Router /api/v1/accounts/list [post]
// @Security BearerAuth
func (h *accountsHandler) List(c *gin.Context) {
	form := &types.ListAccountssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	accountss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertAccountss(accountss)
	if err != nil {
		response.Error(c, ecode.ErrListAccounts)
		return
	}

	response.Success(c, gin.H{
		"accountss": data,
		"total":     total,
	})
}

func getAccountsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertAccounts(accounts *model.Accounts) (*types.AccountsObjDetail, error) {
	data := &types.AccountsObjDetail{}
	err := copier.Copy(data, accounts)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertAccountss(fromValues []*model.Accounts) ([]*types.AccountsObjDetail, error) {
	toValues := []*types.AccountsObjDetail{}
	for _, v := range fromValues {
		data, err := convertAccounts(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
