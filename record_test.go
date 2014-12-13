package dnsmadeeasy

import (
	"github.com/soniah/dnsmadeeasy/testutil"
	"testing"

	. "github.com/motain/gocheck"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	client *Client
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	var err error
	s.client, err = NewClient("aaaaaa1a-11a1-1aa1-a101-11a1a11aa1aa",
		"11a0a11a-a1a1-111a-a11a-a11110a11111")
	s.client.URL = "http://localhost:4444"
	if err != nil {
		panic(err)
	}
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func (s *S) Test_endpoint(c *C) {
	c.Assert(create.endpoint("1", ""), Equals, "/dns/managed/1/records/")
	c.Assert(retrieve.endpoint("1", ""), Equals, "/dns/managed/1/records/")
	c.Assert(update.endpoint("1", "2"), Equals, "/dns/managed/1/records/2/")
	c.Assert(destroy.endpoint("1", "2"), Equals, "/dns/managed/1/records/2/")
}

func (s *S) Test_CreateRecord(c *C) {
	testServer.Response(201, nil, recordCreate)
	cr := map[string]interface{}{
		"Name":  "test",
		"Value": "1.1.1.1",
	}
	id, err := s.client.CreateRecord("870073", cr)
	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(id, Equals, "10022989")
}

func (s *S) Test_ReadRecord(c *C) {
	testServer.Response(200, nil, recordRead)
	record, err := s.client.ReadRecord("870073", "10039429")
	_ = testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(record.ID, Equals, int64(10039429))
}

func (s *S) Test_UpdateRecord(c *C) {
	testServer.Response(200, nil, recordRead)
	testServer.Response(200, nil, "")
	cr := map[string]interface{}{
		"Name": "test-update",
	}
	recordID, err := s.client.UpdateRecord("870073", "10039429", cr)
	_ = testServer.WaitRequest()
	c.Assert(err, IsNil)
	c.Assert(recordID, Equals, int64(10039429))
}

func (s *S) Test_DeleteRecord(c *C) {
	testServer.Response(200, nil, "")
	err := s.client.DeleteRecord("870073", "10039429")
	c.Assert(err, IsNil)
}

var recordCreate = `{
  "name":"test",
  "value":"1.1.1.1",
  "id":10022989,
  "type":"A",
  "source":1,
  "failover":false,
  "monitor":false,
  "sourceId":870073,
  "dynamicDns":false,
  "failed":false,
  "gtdLocation":"DEFAULT",
  "hardLink":false,
  "ttl":86400
}`

var recordRead = `{
  "data":[
    {
      "name":"test",
      "value":"1.1.1.1",
      "id":10039428,
      "type":"A",
      "source":1,
      "failover":false,
      "monitor":false,
      "sourceId":870073,
      "dynamicDns":false,
      "failed":false,
      "gtdLocation":"DEFAULT",
      "ha    rdLink":false,
      "ttl":86400
    },
    {
      "name":"test",
      "value":"1.1.1.2",
      "id":10039429,
      "type":"A",
      "source":1,
      "failover":false,
      "monitor":false,
      "sourceId":870073,
      "dynamicDns":false,
      "failed":false,
      "gtdLocation":"DEFAULT",
      "ha    rdLink":false,
      "ttl":86400
    }
  ],
  "page":0,
  "totalPages":1,
  "totalRecords":2
}`
