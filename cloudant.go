package cloudant

import (
	"fmt"
	"strconv"

	request "github.com/parnurzeal/gorequest"
	couchdb "github.com/timjacobi/go-couchdb"

	"errors"
)

// Client ...
type Client struct {
	Client   *couchdb.Client
	username string
	password string
}

// Options ...
type Options couchdb.Options

// Query ...
type Query struct {
	Selector map[string]interface{} `json:"selector"`
	Fields   []string               `json:"fields,omitempty"`
	Sort     []interface{}          `json:"sort,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	Skip     int                    `json:"skip,omitempty"`
}

// Index query struct
type Index struct {
	Index struct {
		Fields interface{} `json:"fields"`
	} `json:"index"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Ddoc string `json:"ddoc,omitempty"`
}

// NewClient ...
func NewClient(username string, password string, dbName string) (*Client, error) {
	auth := couchdb.BasicAuth(username, password)
	url := fmt.Sprintf("https://%s.cloudant.com", username)
	couchClient, err := couchdb.NewClient(url, nil)
	couchClient.SetAuth(auth)
	return &Client{Client: couchClient, username: username, password: password}, err
}

// IsAlive check whether a server is alive.
func (c *Client) IsAlive() error {
	return c.Client.Ping()
}

// CreateDB ensures that a database with the given name exists.
func (c *Client) CreateDB(dbName string) (*couchdb.DB, error) {
	return c.Client.CreateDB(dbName)
}

// DeleteDB ...
func (c *Client) DeleteDB(dbName string) error {
	return c.Client.DeleteDB(dbName)
}

// CreateDocument ...
func (c *Client) CreateDocument(dbName string, doc interface{}) (string, string, error) {
	db := c.Client.DB(dbName)
	return db.Post(doc)
}

// DeleteDocument ...
func (c *Client) DeleteDocument(dbName string, id string, rev string) (string, error) {
	db := c.Client.DB(dbName)
	return db.Delete(id, rev)
}

// UpdateDocument ...
func (c *Client) UpdateDocument(dbName string, id string, rev string, doc interface{}) (string, error) {
	db := c.Client.DB(dbName)
	return db.Put(id, doc, rev)
}

// GetDocument ...
func (c *Client) GetDocument(dbName string, id string, doc interface{}, opts Options) error {
	db := c.Client.DB(dbName)
	return db.Get(id, doc, couchdb.Options(opts))
}

// GetRawDocument ...
func (c *Client) GetRawDocument(dbName string, id string) (string, error) {
	db := c.Client.DB(dbName)
	return db.Rev(id)
}

// GetAllDocument ...
func (c *Client) GetAllDocument(dbName string, result interface{}, opts Options) error {
	db := c.Client.DB(dbName)
	return db.AllDocs(result, couchdb.Options(opts))
}

// SearchDocument ...
func (c *Client) SearchDocument(dbName string, query Query) (result []interface{}, err error) {
	db := c.Client.DB(dbName)
	req := request.New()
	path := "/" + db.Name() + "/_find"

	var data struct {
		Docs     []interface{}
		Bookmark string `json:"bookmark"`
	}
	_, _, errs := req.SetBasicAuth(c.username, c.password).Post(c.Client.URL() + path).Send(query).EndStruct(&data)

	if errs != nil {
		return nil, errs[0]
	}
	return data.Docs, nil
}

// SetIndex ...
func (c *Client) SetIndex(dbName string, index Index) error {
	db := c.Client.DB(dbName)
	req := request.New()
	path := "/" + db.Name() + "/_index"

	resp, _, err := req.SetBasicAuth(c.username, c.password).Post(c.Client.URL() + path).Send(index).End()
	if err != nil {
		return err[0]
	}
	if resp.StatusCode >= 400 {
		return errors.New("Error in SetIndex: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}
