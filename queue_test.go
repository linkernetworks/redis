package redis

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestTask struct {
	V1 int     `json:"v1"`
	V2 string  `json:"v2"`
	V3 *string `json:"v3"`
}
type QueueTestSuite struct {
	suite.Suite
	service *Service
	conn    *Connection
}

func (suite *QueueTestSuite) SetupTest() {
	suite.service = NewTestRedisService()
	assert.NotNil(suite.T(), suite.service)
	suite.conn = suite.service.GetConnection()
	assert.NotNil(suite.T(), suite.conn)
}

func (suite *QueueTestSuite) TearDownTest() {
	defer suite.conn.Close()
}

func (suite *QueueTestSuite) TestSingleQueue() {
	q := NewQueue(suite.conn)
	assert.NotNil(suite.T(), q)

	_, err := q.RemoveAll("test1")
	assert.NoError(suite.T(), err)

	n, err := q.EnqueueString("test1", "value1")
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), n)

	n2, err2 := q.EnqueueString("test1", "value2")
	assert.NoError(suite.T(), err2)
	assert.NotZero(suite.T(), n2)

	v1, err3 := q.DequeueString("test1")
	assert.NoError(suite.T(), err3)
	assert.Equal(suite.T(), "value1", v1)

	v2, err4 := q.DequeueString("test1")
	assert.NoError(suite.T(), err4)
	assert.Equal(suite.T(), "value2", v2)
}

func (suite *QueueTestSuite) TestSingleQueueInterface() {
	q := NewQueue(suite.conn)
	assert.NotNil(suite.T(), q)

	v3 := "v3"
	task := &TestTask{
		V1: 5,
		V2: "v2",
		V3: &v3,
	}

	_, err := q.RemoveAll("test3")
	assert.NoError(suite.T(), err)

	n, err := q.EnqueueJSON("test3", task)
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), n)

	byt, err3 := q.DequeueJSON("test3")
	v := &TestTask{}
	err = json.Unmarshal(byt, v)
	assert.NoError(suite.T(), err3)
	assert.Equal(suite.T(), "v2", v.V2)
	assert.Equal(suite.T(), 5, v.V1)
}

func (suite *QueueTestSuite) TestMultipleQueue() {
	q := NewQueue(suite.conn)
	assert.NotNil(suite.T(), q)

	_, err := q.RemoveAll("test1")
	assert.NoError(suite.T(), err)

	_, err = q.RemoveAll("test2")
	assert.NoError(suite.T(), err)

	n, err := q.EnqueueString("test1", "value1")
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), n)

	n, err = q.EnqueueString("test2", "value3")
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), n)

	n, err = q.EnqueueString("test1", "value2")
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), n)

	n, err = q.EnqueueString("test2", "value4")
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), n)

	v1, err := q.DequeueString("test1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "value1", v1)

	v2, err := q.DequeueString("test1")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "value2", v2)

	v3, err := q.DequeueString("test2")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "value3", v3)

	v4, err := q.DequeueString("test2")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "value4", v4)
}

func TestQueueTestSuite(t *testing.T) {
	suite.Run(t, new(QueueTestSuite))
}
