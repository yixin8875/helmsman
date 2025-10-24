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

func newStrategiesCache() *gotest.Cache {
	record1 := &model.Strategies{}
	record1.ID = 1
	record2 := &model.Strategies{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewStrategiesCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_strategiesCache_Set(t *testing.T) {
	c := newStrategiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Strategies)
	err := c.ICache.(StrategiesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(StrategiesCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_strategiesCache_Get(t *testing.T) {
	c := newStrategiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Strategies)
	err := c.ICache.(StrategiesCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(StrategiesCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(StrategiesCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_strategiesCache_MultiGet(t *testing.T) {
	c := newStrategiesCache()
	defer c.Close()

	var testData []*model.Strategies
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Strategies))
	}

	err := c.ICache.(StrategiesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(StrategiesCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Strategies))
	}
}

func Test_strategiesCache_MultiSet(t *testing.T) {
	c := newStrategiesCache()
	defer c.Close()

	var testData []*model.Strategies
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Strategies))
	}

	err := c.ICache.(StrategiesCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_strategiesCache_Del(t *testing.T) {
	c := newStrategiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Strategies)
	err := c.ICache.(StrategiesCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_strategiesCache_SetCacheWithNotFound(t *testing.T) {
	c := newStrategiesCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Strategies)
	err := c.ICache.(StrategiesCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(StrategiesCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewStrategiesCache(t *testing.T) {
	c := NewStrategiesCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewStrategiesCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewStrategiesCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
