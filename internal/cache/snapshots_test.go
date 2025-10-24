package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"helmsman/internal/database"
	"helmsman/internal/model"
)

func newSnapshotsCache() *gotest.Cache {
	record1 := &model.Snapshots{}
	record1.ID = 1
	record2 := &model.Snapshots{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewSnapshotsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_snapshotsCache_Set(t *testing.T) {
	c := newSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Snapshots)
	err := c.ICache.(SnapshotsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(SnapshotsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_snapshotsCache_Get(t *testing.T) {
	c := newSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Snapshots)
	err := c.ICache.(SnapshotsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(SnapshotsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(SnapshotsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_snapshotsCache_MultiGet(t *testing.T) {
	c := newSnapshotsCache()
	defer c.Close()

	var testData []*model.Snapshots
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Snapshots))
	}

	err := c.ICache.(SnapshotsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(SnapshotsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Snapshots))
	}
}

func Test_snapshotsCache_MultiSet(t *testing.T) {
	c := newSnapshotsCache()
	defer c.Close()

	var testData []*model.Snapshots
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Snapshots))
	}

	err := c.ICache.(SnapshotsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_snapshotsCache_Del(t *testing.T) {
	c := newSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Snapshots)
	err := c.ICache.(SnapshotsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_snapshotsCache_SetCacheWithNotFound(t *testing.T) {
	c := newSnapshotsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Snapshots)
	err := c.ICache.(SnapshotsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(SnapshotsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewSnapshotsCache(t *testing.T) {
	c := NewSnapshotsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewSnapshotsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewSnapshotsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
