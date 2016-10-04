package cloudant

import (
	couchdb "github.com/timjacobi/go-couchdb"
	request "github.com/parnurzeal/gorequest"
	"fmt"

	"strings"
	"errors"
)

// CloudantClient ...
type CloudantClient struct {
	Client *couchdb.Client
	username string
	password string
}

type Options couchdb.Options

type Query struct {
	Selector map[string] interface{} `json:"selector"`
	Fields []string `json:"fields,omitempty"`
	Sort []interface{} `json:"sort,omitempty"`
	Limit int `json:"limit,omitempty"`
	Skip int `json:"skip,omitempty"`
}

//struct for index query
type IndexStruct struct {
	Index struct {
		Fields interface{} `json:"fields"`
	} `json:"index"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Ddoc string `json:"ddoc,omitempty"`
}

// NewCloudantClient ...
func NewCloudantClient(username string, password string, dbName string) (*CloudantClient, error) {
	auth := couchdb.BasicAuth(username, password)
	url := fmt.Sprintf("https://%s.cloudant.com", username)
	couchClient, err := couchdb.NewClient(url, nil)
	couchClient.SetAuth(auth)
	return &CloudantClient{Client:couchClient, username: username, password: password}, err
}

// CheckAlive ...
func (c *CloudantClient) CheckAlive() error {
	return c.Client.Ping()
}


// CreateDB ensures that a database with the given name exists.
func (c *CloudantClient) CreateDB(dbName string) (*couchdb.DB, error) {
	return c.Client.CreateDB(dbName)
}

// DeleteDB ...
func (c *CloudantClient) DeleteDB(dbName string) error {
	return c.Client.DeleteDB(dbName)
}

// CreateDocument ...
func (c *CloudantClient) CreateDocument(dbName string, doc interface{}) (string, string, error) {
	curDB := c.Client.DB(dbName)
	return curDB.Post(doc)
}

// DeleteDocument ...
func (c *CloudantClient) DeleteDocument(dbName string, id string, rev string ) (string, error) {
	curDB := c.Client.DB(dbName)
	return curDB.Delete(id, rev)
}

// UpdateDocument ...
func (c *CloudantClient) UpdateDocument(dbName string, id string, rev string, doc interface{} ) (string, error) {
	curDB := c.Client.DB(dbName)
	return curDB.Put(id, doc, rev)
}
// GetDocument ...
func (c *CloudantClient) GetDocument(dbName string, id string, doc interface{}, opts Options ) error {
	curDB := c.Client.DB(dbName)
	return curDB.Get(id, doc, couchdb.Options(opts))
}

// GetRawDocument ...
func (c *CloudantClient) GetRawDocument(dbName string, id string ) (string, error) {
	curDB := c.Client.DB(dbName)
	return curDB.Rev(id)
}

func (c *CloudantClient) GetAllDocument(dbName string, result interface{}, opts Options ) error {
	curDB := c.Client.DB(dbName)
	return curDB.AllDocs(result, couchdb.Options(opts))
}

func (c *CloudantClient) SearchDocument(dbName string, query Query) (result []interface{}, err error) {
	curDB := c.Client.DB(dbName)
	req := request.New()
	path := "/" + curDB.Name() + "/_find"

	var data struct {
		Docs []interface{}
		Bookmark string `json:"bookmark"`
	}
	_, _, errs := req.SetBasicAuth(c.username, c.password).Post(c.Client.URL() + path).Send(query).EndStruct(&data)

	if errs != nil {
		return nil, errs[0]
	}else{
		return data.Docs, nil
	}
}

func (c *CloudantClient) SetIndex(dbName string, index IndexStruct) (bool, error) {
	curDB := c.Client.DB(dbName)
	req := request.New()
	path := "/" + curDB.Name() + "/_index"

	resp , _, err := req.SetBasicAuth(c.username, c.password).Post(c.Client.URL() + path).Send(index).End()
	if err!= nil {
		return false, err[0]
	}
	if !strings.Contains(resp.Status, "200"){
		return false, errors.New("NOT AUTH")
	}
	return true, nil
}

