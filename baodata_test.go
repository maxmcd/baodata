package baodata

import (
	"fmt"
	"net/http"
	"testing"
)

type locTest struct {
	loc      string
	id       int
	resource string
}

var locTests []locTest

func init() {
	locTests = []locTest{
		{"/emails/1", 1, "emails"},
		{"emails/1", 1, "emails"},
		{"/emails/1/", 1, "emails"},
		{"/emails", 0, "emails"},
		{"emails", 0, "emails"},
		{"emails/", 0, "emails"},
	}
}

func TestMockBaoData(t *testing.T) {
	basicFunctinality(t)

	// clean up
	mockData = make(map[string][]Data)
}

func TestBoaData(t *testing.T) {
	Connect()
	secret = "not blank"
	baoId = "testId"
	mock = false

	http.HandleFunc("/", Handler)
	go http.ListenAndServe(":8337", nil)
	endpoint = "http://localhost:8337/"

	basicFunctinality(t)

	// clean up
	bd.C(baoId).DropCollection()
}

func TestLocationParsing(t *testing.T) {
	for _, lt := range locTests {
		data, err := parseLocation(lt.loc)
		if err != nil {
			t.Error(err)
		}
		if data.Id != lt.id {
			t.Errorf("id is incorrect")
		}
		if data.ResourceName != lt.resource {
			t.Errorf("resourceName is incorrect")
		}
	}
}

func basicFunctinality(t *testing.T) {
	putData, err := Put("/users", Data{"email": "max@max.com"})
	if err != nil {
		t.Error(err)
	}
	if putData["email"] != "max@max.com" {
		t.Errorf("incorrect email for initial put")
	}
	if putData["id"] != "1" {
		t.Errorf(
			"incorrect initial id should be 1, is: %s",
			putData["id"],
		)
	}

	putData2, err := Put("/users", Data{"email": "steve@max.com"})
	if err != nil {
		t.Error(err)
	}
	if putData2["email"] != "steve@max.com" {
		t.Errorf("incorrect email for second put")
	}
	if putData2["id"] != "2" {
		t.Errorf(
			"incorrect second id should be 2, is: %s",
			putData2["id"],
		)
	}

	putUpdateData, err := Put(
		"/users/"+putData2["id"],
		Data{"email": "notsteve@max.com"},
	)
	if err != nil {
		t.Error(err)
	}
	if putUpdateData["email"] != "notsteve@max.com" {
		t.Errorf("incorrect email for second put")
	}
	if putUpdateData["id"] != "2" {
		t.Errorf(
			"incorrect second id upon update should be 2, is: %s",
			putUpdateData["id"],
		)
	}

	getAllData, err := Get(fmt.Sprintf("/users"))
	if err != nil {
		t.Error(err)
	}
	if len(getAllData) != 2 {
		t.Errorf("not enough entires in the get all data call")
	}
	// TODO further tests?

	getbyId, err := Get(fmt.Sprintf("/users/" + putData["id"]))
	if err != nil {
		t.Error(err)
	}
	if len(getbyId) > 0 {
		if getbyId[0]["email"] != "max@max.com" {
			t.Errorf("incorrect email for initial put")
		}
		if getbyId[0]["id"] != "1" {
			t.Errorf(
				"incorrect initial id should be 1, is: %s",
				getbyId[0]["id"],
			)
		}
	}

	err = Delete(fmt.Sprintf("/users/" + putData["id"]))
	if err != nil {
		t.Error(err)
	}

	_, err = Get(fmt.Sprintf("/users/" + putData["id"]))
	if err == nil || err.Error() != "not found" {
		t.Errorf("should raise error when no items are returned")
	}
}
