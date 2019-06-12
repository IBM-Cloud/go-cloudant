package cloudant

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	request "github.com/parnurzeal/gorequest"
	couchdb "github.com/timjacobi/go-couchdb"
)

// Client ...
type Client struct {
	Client   *couchdb.Client
	username string
	password string
}

// DB ...
type DB struct {
	*couchdb.DB
	username string
	password string
	path     string
}

// DB returns the DB object without verifying its existence.
func (c *Client) DB(name string) *DB {
	dbPath := c.Client.URL() + "/" + name
	return &DB{c.Client.DB(name), c.username, c.password, dbPath}
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
func NewClient(username string, password string) (*Client, error) {
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
func (c *Client) CreateDB(dbName string) (*DB, error) {
	var db *couchdb.DB
	var err error
	if db, err = c.Client.CreateDB(dbName); err != nil {
		return nil, err
	}
	dbPath := c.Client.URL() + "/" + dbName
	return &DB{db, c.username, c.password, dbPath}, nil
}

// EnsureDB ensures that a database with the given name exists.
func (c *Client) EnsureDB(name string) (*DB, error) {
	var db *couchdb.DB
	var err error
	if db, err = c.Client.EnsureDB(name); err != nil {
		return nil, err
	}
	dbPath := c.Client.URL() + "/" + name
	return &DB{db, c.username, c.password, dbPath}, nil
}

// DeleteDB ...
func (c *Client) DeleteDB(dbName string) error {
	return c.Client.DeleteDB(dbName)
}

// CreateDocument ...
func (db *DB) CreateDocument(doc interface{}) (string, string, error) {
	return db.Post(doc)
}

// DeleteDocument ...
func (db *DB) DeleteDocument(id string, rev string) (string, error) {
	return db.Delete(id, rev)
}

// UpdateDocument ...
func (db *DB) UpdateDocument(id string, rev string, doc interface{}) (string, error) {
	return db.Put(id, doc, rev)
}

// GetDocument ...
func (db *DB) GetDocument(id string, doc interface{}, opts Options) error {
	return db.Get(id, doc, couchdb.Options(opts))
}

// GetDocumentRev gets the current document revision.
func (db *DB) GetDocumentRev(id string) (string, error) {
	return db.Rev(id)
}

// GetAllDocument ...
func (db *DB) GetAllDocument(result interface{}, opts Options) error {
	return db.AllDocs(result, couchdb.Options(opts))
}

// SearchDocument ...
func (db *DB) SearchDocument(query Query) (result []interface{}, err error) {
	req := request.New()
	path := "/_find"

	var data struct {
		Docs     []interface{}
		Bookmark string `json:"bookmark"`
	}
	_, _, errs := req.SetBasicAuth(db.username, db.password).Post(db.path + path).Send(query).EndStruct(&data)

	if errs != nil {
		return nil, errs[0]
	}
	return data.Docs, nil
}

// SetIndex ...
func (db *DB) SetIndex(index Index) error {
	req := request.New()
	path := "/_index"

	resp, _, errs := req.SetBasicAuth(db.username, db.password).Post(db.path + path).Send(index).End()
	if errs != nil {
		return errs[0]
	}
	if resp.StatusCode >= 400 {
		return errors.New("Error in setting index: " + strconv.Itoa(resp.StatusCode))
	}
	return nil
}

// CreateDesignDoc ...
func (db *DB) CreateDesignDoc(name string, designJSON string) error {
	var data struct {
		Ok  bool   `json:"ok"`
		ID  string `json:"id"`
		Rev string `json:"rev"`
	}
	req := request.New()
	path := "/_design/" + name
	_, _, errs := req.SetBasicAuth(db.username, db.password).Put(db.path + path).SendString(designJSON).EndStruct(&data)
	if errs != nil {
		return errs[0]
	}
	if data.Ok != true {
		return errors.New("Error in creating design doc")
	}
	return nil
}

// DesignDocument ...
type DesignDocument struct {
	ID      string                 `json:"_id"`
	Indexes map[string]interface{} `json:"indexes,omitempty"`
	Views   map[string]interface{} `json:"views,omitempty"`
}

// NewDesignDocument ...
func NewDesignDocument(name string) *DesignDocument {
	return &DesignDocument{ID: "_design/" + name}
}

// Get design document.
func (ddoc *DesignDocument) Get(db *DB) error {
	return db.GetDocument(ddoc.ID, ddoc, Options{})
}

type searchRow struct {
	ID     string      `json:"id"`
	Order  []float64   `json:"order"`
	Fields interface{} `json:"fields"`
}

// SearchResp ...
type SearchResp struct {
	Num      int         `json:"total_rows"`
	Bookmark string      `json:"bookmark"`
	Rows     []searchRow `json:"rows"`
}

// Search indexes, defined in design documents.
// Cloudant doc: https://docs.cloudant.com/search.html
func (ddoc *DesignDocument) Search(db *DB, index, query, bookmark string, limit int) (*SearchResp, error) {
	path := "/" + ddoc.ID + "/_search/" + index
	body := &SearchResp{}
	req := request.New().
		SetBasicAuth(db.username, db.password).
		Get(db.path + path).
		Query("query=" + query).
		Query("limit=" + strconv.Itoa(limit))
	if bookmark != "" {
		req = req.Query("bookmark=" + bookmark)
	}
	if _, _, errs := req.EndStruct(body); errs != nil {
		return nil, errs[len(errs)-1]
	}
	return body, nil
}

// ViewResp ...
type ViewResp struct {
	Num    int           `json:"total_rows"`
	Offset int           `json:"offset"`
	Rows   []interface{} `json:"rows"`
}

// View ...
// Cloudant doc: https://docs.cloudant.com/creating_views.html
func (ddoc *DesignDocument) View(db *DB, view string) (*ViewResp, error) {
	path := "/" + ddoc.ID + "/_view/" + view
	body := &ViewResp{}
	req := request.New().
		SetBasicAuth(db.username, db.password).
		Get(db.path + path)
	if _, _, errs := req.EndStruct(body); errs != nil {
		return nil, errs[len(errs)-1]
	}
	return body, nil
}

// UpdateDocumentByUpdateHandler ...
// Cloudant doc: https://cloud.ibm.com/docs/services/Cloudant?topic=cloudant-design-documents#update-handlers
func (db *DB) UpdateDocumentByUpdateHandler(id, designID, updateHandlerNamer string, document interface{}) error {
	req := request.New()
	path := "/_design/" + designID + "/_update/" + updateHandlerNamer + "/" + id
	resp, _, errs := req.SetBasicAuth(db.username, db.password).Put(db.path + path).Send(document).End()
	if errs != nil {
		return errs[0]
	}
	if resp.StatusCode >= 400 {
		return errors.New("Error in updating document: " + strconv.Itoa(resp.StatusCode))
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// It is expeted that update handler should return 'ok' on success else return error.
	if string(responseData) != "ok" {
		return errors.New(string(responseData))
	}
	return nil
}
