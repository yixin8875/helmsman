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

func newTagsCache() *gotest.Cache {
	record1 := &model.Tags{}
	record1.ID = 1
	record2 := &model.Tags{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewTagsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_tagsCache_Set(t *testing.T) {
	c := newTagsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Tags)
	err := c.ICache.(TagsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(TagsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_tagsCache_Get(t *testing.T) {
	c := newTagsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Tags)
	err := c.ICache.(TagsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TagsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(TagsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_tagsCache_MultiGet(t *testing.T) {
	c := newTagsCache()
	defer c.Close()

	var testData []*model.Tags
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Tags))
	}

	err := c.ICache.(TagsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(TagsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Tags))
	}
}

func Test_tagsCache_MultiSet(t *testing.T) {
	c := newTagsCache()
	defer c.Close()

	var testData []*model.Tags
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Tags))
	}

	err := c.ICache.(TagsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_tagsCache_Del(t *testing.T) {
	c := newTagsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Tags)
	err := c.ICache.(TagsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_tagsCache_SetCacheWithNotFound(t *testing.T) {
	c := newTagsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Tags)
	err := c.ICache.(TagsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(TagsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewTagsCache(t *testing.T) {
	c := NewTagsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewTagsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewTagsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
