// +build ignore
package cloudant

import (
	"os"
	"testing"
)

var username = os.Getenv("CLOUDANT_USER_NAME")
var password = os.Getenv("CLOUDANT_PASSWORD")

const dbname = "test_db"

func TestConnectingCloudant(t *testing.T) {
	t.Log("Testing Connecting Cloudant")
	client, _  := NewCloudantClient(username, password, dbname)
	err := client.CheckAlive()
	if err != nil {
		t.Error("Error connecting cloudant!")
		t.Error(err)
	}
}

func TestCloudantClient_DeleteDB(t *testing.T) {
	t.Log("Deleting DB")
	client, _ := NewCloudantClient(username, password,dbname)
	err := client.DeleteDB(dbname)
	if err != nil {
		t.Error("Error Deleting DB")
		t.Error(err)
	}

}

func TestCloudantClient_CreateDB(t *testing.T) {
	t.Log("Testing Creating DB")
	client, _ := NewCloudantClient(username, password , dbname)
	db, err := client.CreateDB(dbname)
	if err != nil || db == nil {
		t.Error("Create Failed")
		t.Error(err)
	}
}

func TestCloudantClient_CreateDB_Dup(t *testing.T) {
	t.Log("Testing Creating DB Duplicated/ Should fail!")
	client, _ := NewCloudantClient(username, password , dbname)
	_, err := client.CreateDB(dbname)
	if err == nil {
		t.Error("Create Successful with Duplicate!")
		t.Error(err)
	}
}

func TestCloudantClient_CRUD_Map(t *testing.T) {
	//Step 1. Create document with map
	t.Log("Testing creating doc with map")
	testData := make(map[string]string)
	testData["name"] = "test"
	testData["id"] = "123"
	client, _ := NewCloudantClient(username, password,  dbname)
	id, rev, err := client.CreateDocument(dbname, testData)
	if err != nil {
		t.Error("Error creating Document with map")
		t.Error(err)
	}
	//Step 2. Fetch Document with id
	t.Log("Testing get doc with map")
	resultData := make(map[string]string)
	err = client.GetDocument(dbname, id, &resultData, Options{})
	if resultData["name"] != "test" {
		t.Error("Error fetching Document with map")
		t.Error(err)
	}
	//Step 3. Update Document with id
	t.Log("Testing update doc with map")
	testData["id"] = "updated123"
	newRev, err := client.UpdateDocument(dbname, id, rev, testData)
	resultData = make(map[string]string)
	err = client.GetDocument(dbname, id, &resultData, Options{})
	if resultData["id"] != "updated123" {
		t.Error("Error fetching Document with updated map")
		t.Error(err)
	}
	//Step 4. Delete Document with id
	t.Log("Testing delete doc with map")
	_, err = client.DeleteDocument(dbname, id, newRev)
	if err != nil {
		t.Error("Error deleting Document")
		t.Error(err)
	}

}

func TestCloudantClient_CRUD_Struct(t *testing.T) {
	//Step 1. Create document with struct
	t.Log("Testing creating doc with struct")
	type Data struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	testData := &Data{
		Id:   "1",
		Name: "test2",
	}
	client, _ := NewCloudantClient(username, password , dbname)
	id, rev, err := client.CreateDocument(dbname, testData)
	if err != nil {
		t.Error("Error creating Document with struct")
		t.Error(err)
	}
	//Step 2. Fetch Document with id
	t.Log("Testing get doc with struct")
	resultData := Data{}
	err = client.GetDocument(dbname, id, &resultData, Options{} )
	if resultData.Name != "test2" {
		t.Error("Error fetching Document with struct")
		t.Error(err)
	}
	//Step 3. Update Document with id
	t.Log("Testing update doc with struct")
	testData.Id = "updated123"
	newRev, err := client.UpdateDocument(dbname, id, rev, testData )
	resultData = Data{}
	err = client.GetDocument(dbname, id, &resultData, Options{}  )
	if resultData.Id != "updated123" {
		t.Error("Error fetching Document with updated struct")
		t.Error(err)
	}
	//Step 4. Delete Document with id
	t.Log("Testing delete doc with map")
	_, err = client.DeleteDocument(dbname, id, newRev )
	if err != nil {
		t.Error("Error deleting Document")
		t.Error(err)
	}
}

func TestCloudantClient_SetIndex(t *testing.T) {
	t.Log("Testing setting index for DB")
	index := IndexStruct{}
	index.Index.Fields = []string{"id"}
	client,_ := NewCloudantClient(username, password, dbname)
	b, err := client.SetIndex(dbname, index)
	if !b {
		t.Error("Error setting index")
		t.Error(err)
	}
}

func TestCloudantClient_SearchDocument(t *testing.T) {
	t.Log("Testing search documents")
	//Step 1. Create document with struct
	t.Log("Testing creating doc with struct")
	type Data struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}
	testData1 := &Data{
		Id:   "1",
		Name: "test3-1",
	}
	testData2 := &Data{
		Id:   "11",
		Name: "test3-2",
	}
	testData3 := &Data{
		Id:   "111",
		Name: "test3-3",
	}
	client, _ := NewCloudantClient(username, password, dbname)

	_, _, err1 := client.CreateDocument(dbname, testData1)
	_, _, err2 := client.CreateDocument(dbname, testData2)
	_, _, err3 := client.CreateDocument(dbname, testData3)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Error creating documents for search")
	}
	query := Query{}
	query.Selector = make(map[string]interface{})
	query.Selector["id"] = "11"

	result, err := client.SearchDocument(dbname, query)
	if err != nil {
		t.Error("Error searching documents")
		t.Error(err)
	}
	for _, element := range result {
		r := element.(map[string]interface{})
		if r["id"] != "11" {
			t.Error("Error searching documents")
		}
	}
}
