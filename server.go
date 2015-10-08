package baodata

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var bd *mgo.Database

func Connect() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	bd = session.DB("baodata")
}

var sF secretHandler

type secretHandler func(string, string) (bool, error)

func SecretHandler(f func(baoId, secret string) (idValid bool, err error)) {
	sF = secretHandler(f)
}

func init() {
	SecretHandler(func(baoId, secret string) (isValid bool, err error) {
		isValid = false
		return
	})
}

func Handler(w http.ResponseWriter, req *http.Request) {
	bres, err := handler(w, req)
	if err != nil {
		bres.Error = err.Error()
	}
	jsonBytes, err := json.Marshal(bres)
	w.Write(jsonBytes)
}

func handler(w http.ResponseWriter, req *http.Request) (bres BaoResponse, err error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	var breq BaoRequest
	err = json.Unmarshal(body, &breq)
	if err != nil {
		return
	}

	if breq.Method == "get" {
		if breq.Path.Id == 0 {
			bres.Data, err = dbGetAll(
				breq.BaoId,
				breq.Path.ResourceName,
			)
			return
		} else {
			bres.Data, err = dbGet(
				breq.BaoId,
				breq.Path.ResourceName,
				breq.Path.Id,
			)
			return
		}
	} else if breq.Method == "put" {
		bres.Data, err = dbPut(
			breq.BaoId,
			breq.Path.ResourceName,
			breq.Path.Id,
			breq.Data,
		)
		return
	} else if breq.Method == "delete" {
		err = dbDelete(
			breq.BaoId,
			breq.Path.ResourceName,
			breq.Path.Id,
		)
		return
	}
	err = errors.New("no valid method specified")
	return
}

func dbGet(baoId, resource string, id int) (data []Data, err error) {
	var bdmeta BDMeta
	err = bd.C(baoId).Find(
		bson.M{"resource": resource, "id": id, "deleted": false},
	).One(&bdmeta)
	data = []Data{bdmeta.Data}
	return
}

func dbGetAll(baoId, resource string) (data []Data, err error) {
	var bdmeta []BDMeta
	err = bd.C(baoId).Find(
		bson.M{"resource": resource},
	).All(&bdmeta)
	for _, row := range bdmeta {
		data = append(data, row.Data)
	}
	return
}

func dbPut(baoId, resource string, id int, newData Data) (dataOut []Data, err error) {

	bdmeta := BDMeta{
		Resource: resource,
		Id:       id,
		Data:     newData,
		Deleted:  false, // unecessary, for clarity
	}
	// new item
	if id == 0 {
		var lastIdItem BDMeta
		var newId int
		err = bd.C(baoId).Find(
			bson.M{"resource": resource},
		).Sort("id").One(&lastIdItem)
		if err != nil && err.Error() == "not found" {
			newId = 1
		} else if err != nil {
			return
		} else {
			newId = lastIdItem.Id + 1
		}

		bdmeta.Id = newId
		bdmeta.Data["id"] = strconv.Itoa(bdmeta.Id)

		err = bd.C(baoId).Insert(bdmeta)
		dataOut = []Data{bdmeta.Data}
		return
	} else {
		var existingItem BDMeta
		err = bd.C(baoId).Find(
			bson.M{"resource": resource, "id": id, "deleted": false},
		).One(&existingItem)
		if err != nil {
			return
		}

		for key, value := range newData {
			existingItem.Data[key] = value
		}
		err = bd.C(baoId).Update(
			bson.M{"_id": existingItem.Oid},
			existingItem,
		)
		dataOut = []Data{existingItem.Data}
		if err != nil {
			return
		}
	}

	return
}

func dbDelete(baoId, resource string, id int) (err error) {
	var itemToDelete BDMeta
	err = bd.C(baoId).Find(
		bson.M{"resource": resource, "id": id},
	).One(&itemToDelete)
	if err != nil {
		return
	}

	itemToDelete.Deleted = true

	err = bd.C(baoId).Update(
		bson.M{"_id": itemToDelete.Oid},
		itemToDelete,
	)
	if err != nil {
		return
	}

	return
}

type Data map[string]string

type BDMeta struct {
	Oid      bson.ObjectId `bson:"_id,omitempty"`
	Resource string        `bson:"resource,omitempty"`
	Id       int           `bson:"id,omitempty"`
	Deleted  bool          `bson:"deleted"`
	Data     Data          `bson:"data,omitempty"`
}

type PathData struct {
	ResourceName string
	Id           int
}

type BaoRequest struct {
	Method string   `json:"method"`
	Secret string   `json:"secret"`
	BaoId  string   `json:"baoid"`
	Path   PathData `json:"path"`
	Data   Data     `json:"data"`
}

type BaoResponse struct {
	Message string `json:"message"`
	Data    []Data `json:"datas"`
	Error   string `json:"error"`
}
