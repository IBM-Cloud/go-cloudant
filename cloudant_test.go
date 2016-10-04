package cloudant

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var username = os.Getenv("CLOUDANT_USER_NAME")
var password = os.Getenv("CLOUDANT_PASSWORD")

const dbname = "test_db"

func TestConnection(t *testing.T) {
	t.Log("Testing Cloudant connection")
	client, _ := NewClient(username, password, dbname)
	err := client.IsAlive()
	assert.NoError(t, err, "Error connecting to cloudant")
}

func TestDeleteDB(t *testing.T) {
	t.Log("Testing DB delete")
	client, _ := NewClient(username, password, dbname)
	err := client.DeleteDB(dbname)
	assert.NoError(t, err, "Error deleting DB")
}

func TestCreateDB(t *testing.T) {
	t.Log("Testing DB create")
	client, _ := NewClient(username, password, dbname)
	_, err := client.CreateDB(dbname)
	assert.NoError(t, err, "Error creating DB")
}

func TestCreateExistingDB(t *testing.T) {
	t.Log("Testing existing DB create")
	client, _ := NewClient(username, password, dbname)
	_, err := client.CreateDB(dbname)
	assert.Error(t, err, "Unexpected DB create success with existing name")
}

func TestDocumentCRUDMap(t *testing.T) {
	// Step 1. Create document with map
	t.Log("Testing doc create with map")
	testData := make(map[string]string)
	testData["name"] = "test"
	testData["id"] = "123"
	client, _ := NewClient(username, password, dbname)
	id, rev, err := client.CreateDocument(dbname, testData)
	assert.NoError(t, err, "Error creating document with map")

	// Step 2. Fetch Document with id
	t.Log("Testing doc get with map")
	resultData := make(map[string]string)
	err = client.GetDocument(dbname, id, &resultData, Options{})
	assert.Equal(t, "test", resultData["name"])

	// Step 3. Update Document with id
	t.Log("Testing doc update with map")
	testData["id"] = "updated123"
	newRev, err := client.UpdateDocument(dbname, id, rev, testData)
	resultData = make(map[string]string)
	err = client.GetDocument(dbname, id, &resultData, Options{})
	assert.Equal(t, "updated123", resultData["id"])

	//Step 4. Delete Document with id
	t.Log("Testing doc delete with map")
	_, err = client.DeleteDocument(dbname, id, newRev)
	assert.NoError(t, err, "Error deleting document with map")
}

func TestDocumentCRUDStruct(t *testing.T) {
	// Step 1. Create document with struct
	t.Log("Testing doc create with struct")
	type data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	testData := &data{
		ID:   "1",
		Name: "test2",
	}
	client, _ := NewClient(username, password, dbname)
	id, rev, err := client.CreateDocument(dbname, testData)
	assert.NoError(t, err, "Error creating document with struct")

	// Step 2. Fetch Document with id
	t.Log("Testing doc get with struct")
	resultData := data{}
	err = client.GetDocument(dbname, id, &resultData, Options{})
	assert.Equal(t, "test2", resultData.Name)

	// Step 3. Update Document with id
	t.Log("Testing doc update with struct")
	testData.ID = "updated123"
	newRev, err := client.UpdateDocument(dbname, id, rev, testData)
	resultData = data{}
	err = client.GetDocument(dbname, id, &resultData, Options{})
	assert.Equal(t, "updated123", resultData.ID)

	// Step 4. Delete Document with id
	t.Log("Testing doc delete with struct")
	_, err = client.DeleteDocument(dbname, id, newRev)
	assert.NoError(t, err, "Error deleting document with struct")
}

func TestSetIndex(t *testing.T) {
	t.Log("Testing setting index for DB")
	index := Index{}
	index.Index.Fields = []string{"id"}
	client, _ := NewClient(username, password, dbname)
	err := client.SetIndex(dbname, index)
	assert.NoError(t, err, "Error setting index")
}

func TestSearchDocument(t *testing.T) {
	t.Log("Testing search documents")
	//Step 1. Create document with struct
	t.Log("Testing creating doc with struct")
	type data struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	testData1 := &data{
		ID:   "1",
		Name: "test3-1",
	}
	testData2 := &data{
		ID:   "11",
		Name: "test3-2",
	}
	testData3 := &data{
		ID:   "111",
		Name: "test3-3",
	}
	client, _ := NewClient(username, password, dbname)

	_, _, err1 := client.CreateDocument(dbname, testData1)
	assert.NoError(t, err1)
	_, _, err2 := client.CreateDocument(dbname, testData2)
	assert.NoError(t, err2)
	_, _, err3 := client.CreateDocument(dbname, testData3)
	assert.NoError(t, err3)

	query := Query{}
	query.Selector = make(map[string]interface{})
	query.Selector["id"] = "11"

	result, err := client.SearchDocument(dbname, query)
	assert.NoError(t, err, "Error searching documents")

	for _, element := range result {
		r := element.(map[string]interface{})
		assert.Equal(t, "11", r["id"])
	}
}
