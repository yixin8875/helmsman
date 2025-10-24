package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"helmsman/internal/cache"
	"helmsman/internal/database"
	"helmsman/internal/model"
)

var _ TradeTagsDao = (*tradeTagsDao)(nil)

// TradeTagsDao defining the dao interface
type TradeTagsDao interface {
	Create(ctx context.Context, table *model.TradeTags) error
	DeleteByTradeID(ctx context.Context, tradeID int) error
	UpdateByTradeID(ctx context.Context, table *model.TradeTags) error
	GetByTradeID(ctx context.Context, tradeID int) (*model.TradeTags, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.TradeTags, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.TradeTags) (int, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, tradeID int) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.TradeTags) error
}

type tradeTagsDao struct {
	db    *gorm.DB
	cache cache.TradeTagsCache // if nil, the cache is not used.
	sfg   *singleflight.Group  // if cache is nil, the sfg is not used.
}

// NewTradeTagsDao creating the dao interface
func NewTradeTagsDao(db *gorm.DB, xCache cache.TradeTagsCache) TradeTagsDao {
	if xCache == nil {
		return &tradeTagsDao{db: db}
	}
	return &tradeTagsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *tradeTagsDao) deleteCache(ctx context.Context, tradeID int) error {
	if d.cache != nil {
		return d.cache.Del(ctx, tradeID)
	}
	return nil
}

// Create a new tradeTags, insert the record and the tradeID value is written back to the table
func (d *tradeTagsDao) Create(ctx context.Context, table *model.TradeTags) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByTradeID delete a tradeTags by tradeID
func (d *tradeTagsDao) DeleteByTradeID(ctx context.Context, tradeID int) error {
	err := d.db.WithContext(ctx).Where("trade_id = ?", tradeID).Delete(&model.TradeTags{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, tradeID)

	return nil
}

// UpdateByTradeID update a tradeTags by tradeID
func (d *tradeTagsDao) UpdateByTradeID(ctx context.Context, table *model.TradeTags) error {
	err := d.updateDataByTradeID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.TradeID)

	return err
}

func (d *tradeTagsDao) updateDataByTradeID(ctx context.Context, db *gorm.DB, table *model.TradeTags) error {
	if table.TradeID < 1 {
		return errors.New("tradeID cannot be 0")
	}

	update := map[string]interface{}{}

	if table.TradeID != 0 {
		update["trade_id"] = table.TradeID
	}
	if table.TagID != 0 {
		update["tag_id"] = table.TagID
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByTradeID get a tradeTags by tradeID
func (d *tradeTagsDao) GetByTradeID(ctx context.Context, tradeID int) (*model.TradeTags, error) {
	// no cache
	if d.cache == nil {
		record := &model.TradeTags{}
		err := d.db.WithContext(ctx).Where("trade_id = ?", tradeID).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, tradeID)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same tradeID, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.IntToStr(tradeID), func() (interface{}, error) {

			table := &model.TradeTags{}
			err = d.db.WithContext(ctx).Where("trade_id = ?", tradeID).First(table).Error
			if err != nil {
				// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, database.ErrRecordNotFound) {
					if err = d.cache.SetPlaceholder(ctx, tradeID); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("tradeID", tradeID))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, tradeID, table, cache.TradeTagsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("tradeID", tradeID))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.TradeTags)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get a paginated list of tradeTags by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *tradeTagsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.TradeTags, int64, error) {
	if params.Sort == "" {
		params.Sort = "-trade_id"
	}
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.TradeTagsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.TradeTags{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.TradeTags{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *tradeTagsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.TradeTags) (int, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.TradeID, err
}

// DeleteByTx delete a record by tradeID in the database using the provided transaction
func (d *tradeTagsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, tradeID int) error {
	err := tx.WithContext(ctx).Where("trade_id = ?", tradeID).Delete(&model.TradeTags{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, tradeID)

	return nil
}

// UpdateByTx update a record by tradeID in the database using the provided transaction
func (d *tradeTagsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.TradeTags) error {
	err := d.updateDataByTradeID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.TradeID)

	return err
}
