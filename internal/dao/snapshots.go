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

var _ SnapshotsDao = (*snapshotsDao)(nil)

// SnapshotsDao defining the dao interface
type SnapshotsDao interface {
	Create(ctx context.Context, table *model.Snapshots) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Snapshots) error
	GetByID(ctx context.Context, id uint64) (*model.Snapshots, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Snapshots, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Snapshots) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Snapshots) error
}

type snapshotsDao struct {
	db    *gorm.DB
	cache cache.SnapshotsCache // if nil, the cache is not used.
	sfg   *singleflight.Group  // if cache is nil, the sfg is not used.
}

// NewSnapshotsDao creating the dao interface
func NewSnapshotsDao(db *gorm.DB, xCache cache.SnapshotsCache) SnapshotsDao {
	if xCache == nil {
		return &snapshotsDao{db: db}
	}
	return &snapshotsDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *snapshotsDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new snapshots, insert the record and the id value is written back to the table
func (d *snapshotsDao) Create(ctx context.Context, table *model.Snapshots) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a snapshots by id
func (d *snapshotsDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Snapshots{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a snapshots by id, support partial update
func (d *snapshotsDao) UpdateByID(ctx context.Context, table *model.Snapshots) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *snapshotsDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Snapshots) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.TradeID != 0 {
		update["trade_id"] = table.TradeID
	}
	if table.Type != "" {
		update["type"] = table.Type
	}
	if table.ImageURL != "" {
		update["image_url"] = table.ImageURL
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a snapshots by id
func (d *snapshotsDao) GetByID(ctx context.Context, id uint64) (*model.Snapshots, error) {
	// no cache
	if d.cache == nil {
		record := &model.Snapshots{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.Snapshots{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				if errors.Is(err, database.ErrRecordNotFound) {
					// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.SnapshotsExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Snapshots)
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

// GetByColumns get a paginated list of snapshotss by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *snapshotsDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Snapshots, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.SnapshotsColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Snapshots{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Snapshots{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *snapshotsDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Snapshots) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *snapshotsDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Snapshots{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *snapshotsDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Snapshots) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
