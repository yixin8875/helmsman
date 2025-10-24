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

func newTradesCache() *gotest.Cache {
	record1 := &model.Trades{}
	record1.ID = 1
	record2 := &model.Trades{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewTradesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_tradesCache_Set(t *testing.T) {
	c := newTradesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Trades)
	err := c.ICache.(TradesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(TradesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_tradesCache_Get(t *testing.T) {
	c := newTradesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Trades)
	err := c.ICache.(TradesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TradesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(TradesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_tradesCache_MultiGet(t *testing.T) {
	c := newTradesCache()
	defer c.Close()

	var testData []*model.Trades
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Trades))
	}

	err := c.ICache.(TradesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TradesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Trades))
	}
}

func Test_tradesCache_MultiSet(t *testing.T) {
	c := newTradesCache()
	defer c.Close()

	var testData []*model.Trades
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Trades))
	}

	err := c.ICache.(TradesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_tradesCache_Del(t *testing.T) {
	c := newTradesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Trades)
	err := c.ICache.(TradesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_tradesCache_SetCacheWithNotFound(t *testing.T) {
	c := newTradesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Trades)
	err := c.ICache.(TradesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(TradesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewTradesCache(t *testing.T) {
	c := NewTradesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewTradesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewTradesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
