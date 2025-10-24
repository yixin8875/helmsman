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

func newAccountsCache() *gotest.Cache {
	record1 := &model.Accounts{}
	record1.ID = 1
	record2 := &model.Accounts{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewAccountsCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_accountsCache_Set(t *testing.T) {
	c := newAccountsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Accounts)
	err := c.ICache.(AccountsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(AccountsCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_accountsCache_Get(t *testing.T) {
	c := newAccountsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Accounts)
	err := c.ICache.(AccountsCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(AccountsCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(AccountsCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_accountsCache_MultiGet(t *testing.T) {
	c := newAccountsCache()
	defer c.Close()

	var testData []*model.Accounts
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Accounts))
	}

	err := c.ICache.(AccountsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(AccountsCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Accounts))
	}
}

func Test_accountsCache_MultiSet(t *testing.T) {
	c := newAccountsCache()
	defer c.Close()

	var testData []*model.Accounts
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Accounts))
	}

	err := c.ICache.(AccountsCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_accountsCache_Del(t *testing.T) {
	c := newAccountsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Accounts)
	err := c.ICache.(AccountsCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_accountsCache_SetCacheWithNotFound(t *testing.T) {
	c := newAccountsCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Accounts)
	err := c.ICache.(AccountsCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(AccountsCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewAccountsCache(t *testing.T) {
	c := NewAccountsCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewAccountsCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewAccountsCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
