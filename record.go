package dnsmadeeasy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
)

// DataResponse is the response from a GET ie all records for
// a domainID
type DataResponse struct {
	Data []Record `json:"data"`
}

// Record is used to represent a retrieved Record.
type Record struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	RecordID    int64  `json:"id"`
	Type        string `json:"type"`
	Source      int64  `json:"source"`
	SourceID    int64  `json:"sourceId"`
	DynamicDNS  bool   `json:"dynamicDns"`
	Password    string `json:"password"`
	TTL         int64  `json:"ttl"`
	Monitor     bool   `json:"monitor"`
	Failover    bool   `json:"failover"`
	Failed      bool   `json:"failed"`
	GtdLocation string `json:"gtdLocation"`
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Title       string `json:"title"`
	Hardlink    bool   `json:"hardLink"`
	MXLevel     int64  `json:"mxLevel"`
	Weight      int64  `json:"weight"`
	Priority    int64  `json:"priority"`
	Port        int64  `json:"port"`
}

type requestType int

const (
	create requestType = iota
	retrieve
	update
	destroy
)

func (rt requestType) endpoint(domainID int64, recordID int64) (result string) {
	switch rt {
	case create, retrieve:
		result = fmt.Sprintf("/dns/managed/%d/records/", domainID)
	case update, destroy:
		result = fmt.Sprintf("/dns/managed/%d/records/%d/", domainID, recordID)
	}
	return result
}

// CRUD - Create, Read, Update, Delete

// CreateRecord creates a DNS record on DNSMadeEasy
func (c *Client) CreateRecord(domainID int64, cr map[string]interface{}) (int64, error) {

	path := create.endpoint(domainID, 0)
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(cr); err != nil {
		return 0, err
	}

	req, err := c.NewRequest("POST", path, buf, "")
	if err != nil {
		return 0, fmt.Errorf("Error from NewRequest: %s", err)
	}

	resp, err := checkResp(c.HTTP.Do(req))
	if err != nil {
		return 0, fmt.Errorf("Error creating record: %s", err)
	}

	record := new(Record)

	err = decodeBody(resp, &record)
	if err != nil {
		return 0, fmt.Errorf("Error parsing record response: %s", err)
	}

	// The request was successful
	return record.RecordID, nil
}

// ReadRecord gets a record by the ID specified and returns a Record and an
// error.
func (c *Client) ReadRecord(domainID int64, recordID int64) (*Record, error) {
	body := bytes.NewBuffer(nil)
	path := retrieve.endpoint(domainID, recordID)
	req, err := c.NewRequest("GET", path, body, "")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.HTTP.Do(req))
	if err != nil {
		return nil, fmt.Errorf("Error retrieving record: %s", err)
	}

	dataResp := DataResponse{}
	err = decodeBody(resp, &dataResp)
	if err != nil {
		return nil, fmt.Errorf("Error decoding data response: %s", err)
	}
	var result Record
	var found bool
	for _, record := range dataResp.Data {
		if record.RecordID == recordID {
			result = record // not pointer, so data copied
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("Unable to find record %d", recordID)
	}
	return &result, nil
}

// UpdateRecord updated a record from the parameters specified and
// returns an error if it fails.
func (c *Client) UpdateRecord(domainID int64, recordID int64, cr map[string]interface{}) (int64, error) {

	current, err := c.ReadRecord(domainID, recordID)
	if err != nil {
		return 0, err
	}

	err = mergo.Map(current, cr)
	if err != nil {
		return 0, err
	}

	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(current); err != nil {
		return 0, err
	}

	path := update.endpoint(domainID, recordID)
	req, err := c.NewRequest("PUT", path, buf, "")
	if err != nil {
		return 0, err
	}

	_, err = checkResp(c.HTTP.Do(req))
	if err != nil {
		return 0, fmt.Errorf("Error updating record: %s", err)
	}

	// The request was successful
	return recordID, nil
}

// DeleteRecord destroys a record by the ID specified and
// returns an error if it fails. If no error is returned,
// the Record was succesfully destroyed.
func (c *Client) DeleteRecord(domainID int64, recordID int64) error {
	body := bytes.NewBuffer(nil)
	path := destroy.endpoint(domainID, recordID)
	req, err := c.NewRequest("DELETE", path, body, "")
	if err != nil {
		return err
	}

	_, err = checkResp(c.HTTP.Do(req))
	if err != nil {
		return fmt.Errorf("Unable to find record %d", recordID)
	}

	// The request was successful
	return nil
}
