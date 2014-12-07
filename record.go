package dnsmadeeasy

import (
	"fmt"
	//	"github.com/imdario/mergo"
	"strconv"
)

// DataResponse is the response from a GET ie all records for
// a domainID
// TODO rename Record to Data, as suitable
type DataResponse struct {
	Data []Record `json:"data"`
}

// Record is used to represent a retrieved Record. All properties
// are set as strings.
// TODO ID should be called RecordID here
type Record struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	ID           int64  `json:"id"`
	Type         string `json:"type"`
	Source       int64  `json:"source"`
	SourceID     int64  `json:"sourceId"`
	DynamicDNS   bool   `json:"dynamicDns"`
	Password     string `json:"password"`
	TTL          int64  `json:"ttl"`
	Monitor      bool   `json:"monitor"`
	Failover     bool   `json:"failover"`
	Failed       bool   `json:"failed"`
	GtdLocation  string `json:"gtdLocation"`
	Description  string `json:"description"`
	Keywords     string `json:"keywords"`
	Title        string `json:"title"`
	RedirectType string `json:redirectType:`
	Hardlink     bool   `json:"hardLink"`
	MXLevel      int64  `json:"mxLevel"`
	Weight       int64  `json:"weight"`
	Priority     int64  `json:"priority"`
	Port         int64  `json:"port"`
}

// StringID returns the id as a string
func (r *Record) StringID() string {
	return strconv.FormatInt(r.ID, 10)
}

// ChangeRecord contains the request parameters to create or update a
// record.
type ChangeRecord struct {
	Name         string
	Value        string
	Type         string
	Source       int64
	SourceID     int64
	DynamicDNS   bool
	Password     string
	TTL          int64
	Monitor      bool
	Failover     bool
	Failed       bool
	GtdLocation  string
	Description  string
	Keywords     string
	Title        string
	RedirectType string
	Hardlink     bool
	MXLevel      int64
	Weight       int64
	Priority     int64
	Port         int64
}

type requestType int

const (
	create requestType = iota
	retrieve
	update
	destroy
)

func (rt requestType) endpoint(domainID string, recordID string) (result string) {
	switch rt {
	case create, retrieve:
		result = fmt.Sprintf("/dns/managed/%s/records/", domainID)
	case update, destroy:
		result = fmt.Sprintf("/dns/managed/%s/records/%s/", domainID, recordID)
	}
	return result
}

// CreateRecord creates a DNS record on DNSMadeEasy
func (c *Client) CreateRecord(domainID string, opts *ChangeRecord) (string, error) {
	// Make the request parameters
	params := make(map[string]interface{})

	params["name"] = opts.Name
	params["type"] = opts.Type
	params["value"] = opts.Value
	params["ttl"] = strconv.FormatInt(opts.TTL, 10)

	ep := create.endpoint(domainID, "")

	req, err := c.NewRequest(params, "POST", ep, "")
	if err != nil {
		return "", fmt.Errorf("Error from NewRequest: %s", err)
	}

	resp, err := checkResp(c.HTTP.Do(req))
	if err != nil {
		return "", fmt.Errorf("Error creating record: %s", err)
	}

	record := new(Record)

	err = decodeBody(resp, &record)
	if err != nil {
		return "", fmt.Errorf("Error parsing record response: %s", err)
	}

	// The request was successful
	return record.StringID(), nil
}

// UpdateRecord updated a record from the parameters specified and
// returns an error if it fails.
func (c *Client) UpdateRecord(domainID string, recordID string, opts *ChangeRecord) (string, error) {
	// Make the request parameters
	params := make(map[string]interface{})

	if opts.Name != "" {
		params["name"] = opts.Name
	}
	if opts.Type != "" {
		params["type"] = opts.Type
	}
	if opts.Value != "" {
		params["value"] = opts.Value
	}
	params["ttl"] = strconv.FormatInt(opts.TTL, 10)

	ep := update.endpoint(domainID, recordID)

	req, err := c.NewRequest(params, "PUT", ep, "")
	if err != nil {
		return "", err
	}

	_, err = checkResp(c.HTTP.Do(req))
	if err != nil {
		return "", fmt.Errorf("Error updating record: %s", err)
	}

	// The request was successful
	return recordID, nil
}

// DestroyRecord destroys a record by the ID specified and
// returns an error if it fails. If no error is returned,
// the Record was succesfully destroyed.
func (c *Client) DestroyRecord(domainID string, recordID string) error {
	var body map[string]interface{}
	ep := destroy.endpoint(domainID, recordID)
	req, err := c.NewRequest(body, "DELETE", ep, "")
	if err != nil {
		return err
	}

	_, err = checkResp(c.HTTP.Do(req))
	if err != nil {
		return fmt.Errorf("Error destroying record: %s", err)
	}

	// The request was successful
	return nil
}

// RetrieveRecord gets a record by the ID specified and returns a Record and an
// error. An error will be returned for failed requests with a nil Record.
func (c *Client) RetrieveRecord(domainID string, recordID string) (*Record, error) {
	var body map[string]interface{}
	ep := retrieve.endpoint(domainID, recordID)
	req, err := c.NewRequest(body, "GET", ep, "")
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
	id, _ := strconv.ParseInt(recordID, 10, 64) // TODO
	for _, record := range dataResp.Data {
		if record.ID == id {
			result = record // not pointer, so data copied
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("Unable to find record %s", recordID)
	}
	return &result, nil
}
