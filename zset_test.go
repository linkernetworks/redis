package redis

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"bitbucket.org/linkernetworks/aurora/src/config"
	"bitbucket.org/linkernetworks/aurora/src/jobcontroller/types"
	"github.com/stretchr/testify/assert"
)

const testingConfigPath = "../../../config/testing.json"
const defaultTestQueue = "_test_queue_"

// AnyStruct is used for Redis testing
type AnyStruct struct {
	ID        string
	CreatedAt time.Time
	Priority  float64
}

func TestNewRedisZSet(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)
	assert.NotNil(t, rzset)
	assert.Equal(t, defaultTestQueue, rzset.key)
}

func TestZADD(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)

	// call and check
	any := AnyStruct{"abc123", time.Now(), types.ScoreHigh}
	data, err := json.Marshal(any)
	assert.Nil(t, err)
	n, err := rzset.Add(any.Priority, data)
	assert.Nil(t, err)
	assert.Equal(t, 1, n)

	// clean up
	_, err = rzset.RemoveAll()
	assert.Nil(t, err)
	assert.Equal(t, 0, rzset.Len())
}

func TestZRANGEBYSCORE(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)

	// add 3 test data
	arr := []AnyStruct{
		AnyStruct{"a1", time.Now(), types.ScoreHigh},
		AnyStruct{"a2", time.Now(), types.ScoreMedium},
		AnyStruct{"a3", time.Now(), types.ScoreLow},
	}
	for _, a := range arr {
		data, err := json.Marshal(a)
		assert.Nil(t, err)
		_, err = rzset.Add(a.Priority, data)
		assert.Nil(t, err)
	}
	// call and check
	min, max, offset, limit := types.ScoreHigh, types.ScoreLow, 0, 2
	members, err := rzset.RangeByScore(min, max, offset, limit)
	assert.Nil(t, err)
	assert.Equal(t, limit, len(members))
	// decode member to AnyStruct
	for i, m := range members {
		var a = AnyStruct{}
		err = json.Unmarshal(m.([]byte), &a)
		assert.Nil(t, err)
		// fmt.Printf("%+v\n", a)
		// it's a1, a2 (in order)
		assert.Equal(t, arr[i].ID, a.ID)
	}

	// clean up
	_, err = rzset.RemoveAll()
	assert.Nil(t, err)
	assert.Equal(t, 0, rzset.Len())
}

func TestLen(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)

	assert.Equal(t, 0, rzset.Len())

	// add some data
	a1 := AnyStruct{"a1", time.Now(), types.ScoreHigh}
	data, err := json.Marshal(a1)
	assert.Nil(t, err)
	_, err = rzset.Add(a1.Priority, data)
	assert.Nil(t, err)

	a2 := AnyStruct{"a2", time.Now(), types.ScoreMedium}
	data2, err := json.Marshal(a2)
	assert.Nil(t, err)
	_, err = rzset.Add(a2.Priority, data2)
	assert.Nil(t, err)

	// call and check
	assert.Equal(t, 2, rzset.Len())

	// clean up
	_, err = rzset.RemoveAll()
	assert.Nil(t, err)
	assert.Equal(t, 0, rzset.Len())
}

func TestZREM(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)

	assert.Equal(t, 0, rzset.Len())

	// add some data
	a1 := AnyStruct{"a1", time.Now(), types.ScoreHigh}
	data, err := json.Marshal(a1)
	assert.Nil(t, err)
	_, err = rzset.Add(a1.Priority, data)
	assert.Nil(t, err)

	// call and check
	n, err := rzset.Remove(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, n)
}

func TestRemoveAll(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)

	assert.Equal(t, 0, rzset.Len())

	// add some data
	arr := []AnyStruct{
		AnyStruct{"a1", time.Now(), types.ScoreHigh},
		AnyStruct{"a2", time.Now(), types.ScoreHigh},
		AnyStruct{"a3", time.Now(), types.ScoreHigh},
	}
	for _, a := range arr {
		data, err := json.Marshal(a)
		assert.NoError(t, err)
		_, err = rzset.Add(a.Priority, data)
		assert.NoError(t, err)
	}
	// call and check
	n, err := rzset.RemoveAll()
	assert.NoError(t, err)
	assert.Equal(t, n, len(arr))

	// clean up
}

func TestZPOP(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)

	// add 3 test data
	arr := []AnyStruct{
		AnyStruct{"a2", time.Now(), types.ScoreMedium},
		AnyStruct{"a1", time.Now(), types.ScoreHigh}, // highest priority
		AnyStruct{"a3", time.Now(), types.ScoreLow},
	}
	for _, a := range arr {
		data, err := json.Marshal(a)
		assert.Nil(t, err)
		_, err = rzset.Add(a.Priority, data)
		assert.Nil(t, err)
	}

	assert.Equal(t, 3, rzset.Len())
	// call and check
	result, err := rzset.ZPOP()
	assert.Nil(t, err)

	var a = AnyStruct{}
	err = json.Unmarshal(result.([]byte), &a)
	assert.Nil(t, err)

	// assert.Equal(t, arr[1], a)
	assert.Equal(t, arr[1].ID, a.ID)
	assert.Equal(t, arr[1].Priority, a.Priority)
	assert.Equal(t, arr[1].CreatedAt.Unix(), a.CreatedAt.Unix())
	assert.Equal(t, 2, rzset.Len())

	// clean up
	_, err = rzset.RemoveAll()
	assert.Nil(t, err)
	assert.Equal(t, 0, rzset.Len())
}

func TestQueryAll(t *testing.T) {
	if _, defined := os.LookupEnv("TEST_REDIS"); !defined {
		t.SkipNow()
		return
	}
	cf := config.MustRead(testingConfigPath)
	rds := New(cf.Redis)
	conn := rds.GetConnection()
	defer conn.Close()
	rzset := NewZSet(conn, defaultTestQueue)

	// add 3 test data
	arr := []AnyStruct{
		AnyStruct{"a1", time.Now(), types.ScoreHigh},
		AnyStruct{"a2", time.Now(), types.ScoreMedium},
		AnyStruct{"a3", time.Now(), types.ScoreLow},
	}
	for _, a := range arr {
		data, err := json.Marshal(a)
		assert.Nil(t, err)
		_, err = rzset.Add(a.Priority, data)
		assert.Nil(t, err)
	}

	// call and check
	members, err := rzset.QueryAll()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(members))
	// decode member to AnyStruct
	for i, m := range members {
		var a = AnyStruct{}
		err = json.Unmarshal(m.([]byte), &a)
		assert.Nil(t, err)
		// fmt.Printf("%+v\n", a)
		// it's a2, a1 and a3 (in order)
		assert.Equal(t, arr[i].ID, a.ID)
	}

	// clean up
	_, err = rzset.RemoveAll()
	assert.Nil(t, err)
	assert.Equal(t, 0, rzset.Len())
}
