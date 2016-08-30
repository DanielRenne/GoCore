package model

import (
	// "encoding/base64"
	"encoding/json"
	"errors"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Transactions struct{}

var mongoTransactionsCollection *mgo.Collection

func init() {
	go func() {

		for {
			if dbServices.MongoDB != nil {
				log.Println("Building Indexes for MongoDB collection Transactions:")
				mongoTransactionsCollection = dbServices.MongoDB.C("Transactions")
				ci := mgo.CollectionInfo{ForceIdIndex: true}
				mongoTransactionsCollection.Create(&ci)
				var obj Transactions
				obj.Index()
				return
			}
			<-dbServices.WaitForDatabase()
		}
	}()
}

// type Transaction struct {
// 	Id  bson.ObjectId `json:"id" bson:"_id,omitempty"`
// 	Trx string        `json:"trx" bson:"trx"`
// 	Nm  string        `json:"nm" bson:"nm"`
// 	OId string        `json:"oId" bson:"oId"`
// 	Dta string        `json:"dta" bson:"dta"`
// 	UId string        `json:"uId" bson:"uId"`
// 	Typ int           `json:"typ" bson:"typ"`
// }

type Transaction struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserId         string        `json:"userId" dbIndex:"index" bson:"userId"`
	CreateDate     time.Time     `json:"createDate" bson:"createDate"`
	LastUpdate     time.Time     `json:"lastUpdate" bson:"lastUpdate"`
	CompleteDate   time.Time     `json:"completeDate" bson:"completeDate"`
	RollbackDate   time.Time     `json:"rollbackDate" bson:"rollbackDate"`
	Committed      bool          `json:"committed" bson:"committed"`
	Error          string        `json:"error" bson:"error"`
	Collections    []string      `json:"collections" bson:"collections"`
	Details        string        `json:"details" bson:"details"`
	RolledBack     bool          `json:"rolledBack" bson:"rolledBack"`
	RolledBackBy   string        `json:"rolledBackBy" bson:"rolledBackBy"`
	RollbackReason string        `json:"rollbackReason" bson:"rollbackReason"`
}

func (self *Transactions) Single(field string, value interface{}) (retObj Transaction, e error) {
	if field == "id" {
		query := mongoTransactionsCollection.FindId(bson.ObjectIdHex(value.(string)))
		e = query.One(&retObj)
		return
	}
	m := make(bson.M)
	m[field] = value
	query := mongoTransactionsCollection.Find(m)
	e = query.One(&retObj)
	return
}

func (obj *Transactions) Search(field string, value interface{}) (retObj []Transaction, e error) {
	var query *mgo.Query
	if field == "id" {
		query = mongoTransactionsCollection.FindId(bson.ObjectIdHex(value.(string)))
	} else {
		m := make(bson.M)
		m[field] = value
		query = mongoTransactionsCollection.Find(m)
	}

	e = query.All(&retObj)
	return
}

func (obj *Transactions) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []Transaction, e error) {
	var query *mgo.Query
	if field == "id" {
		query = mongoTransactionsCollection.FindId(bson.ObjectIdHex(value.(string)))
	} else {
		m := make(bson.M)
		m[field] = value
		query = mongoTransactionsCollection.Find(m)
	}

	if limit == 0 && skip == 0 {
		e = query.All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = query.Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 {
		e = query.Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if skip > 0 {
		e = query.Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	return
}

func (obj *Transactions) All() (retObj []Transaction, e error) {
	e = mongoTransactionsCollection.Find(bson.M{}).All(&retObj)
	if len(retObj) == 0 {
		retObj = []Transaction{}
	}
	return
}

func (obj *Transactions) AllAdvanced(limit int, skip int) (retObj []Transaction, e error) {
	if limit == 0 && skip == 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if skip > 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	return
}

func (obj *Transactions) AllByIndex(index string) (retObj []Transaction, e error) {
	e = mongoTransactionsCollection.Find(bson.M{}).Sort(index).All(&retObj)
	if len(retObj) == 0 {
		retObj = []Transaction{}
	}
	return
}

func (obj *Transactions) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Transaction, e error) {
	if limit == 0 && skip == 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).Sort(index).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).Sort(index).Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).Sort(index).Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if skip > 0 {
		e = mongoTransactionsCollection.Find(bson.M{}).Sort(index).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	return
}

func (obj *Transactions) Range(min, max, field string) (retObj []Transaction, e error) {
	var query *mgo.Query
	m := make(bson.M)
	m[field] = bson.M{"$gte": min, "$lte": max}
	query = mongoTransactionsCollection.Find(m)
	e = query.All(&retObj)
	return
}

func (obj *Transactions) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Transaction, e error) {
	var query *mgo.Query
	m := make(bson.M)
	m[field] = bson.M{"$gte": min, "$lte": max}
	query = mongoTransactionsCollection.Find(m)
	if limit == 0 && skip == 0 {
		e = query.All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = query.Limit(limit).Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if limit > 0 {
		e = query.Limit(limit).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	if skip > 0 {
		e = query.Skip(skip).All(&retObj)
		if len(retObj) == 0 {
			retObj = []Transaction{}
		}
		return
	}
	return
}

func (obj *Transactions) Index() error {
	for key, value := range dbServices.GetDBIndexes(Transaction{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}

		err := mongoTransactionsCollection.EnsureIndex(index)
		if err != nil {
			log.Println("Failed to create index for Transaction." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for Transaction." + key)
		}
	}
	return nil
}

func (obj *Transactions) New(userId string) (*Transaction, error) {
	t := Transaction{}
	t.UserId = userId
	err := t.Begin()

	return &t, err
}

func (self *Transaction) Save() error {
	if mongoTransactionsCollection == nil {
		return errors.New("Collection Transactions not initialized")
	}
	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	changeInfo, err := mongoTransactionsCollection.UpsertId(objectId, &self)
	if err != nil {
		log.Println("Failed to upsertId for Transaction:  " + err.Error())
		return err
	}
	if changeInfo.UpsertedId != nil {
		self.Id = changeInfo.UpsertedId.(bson.ObjectId)
	}
	return nil
}

func (self *Transaction) Delete() error {
	return mongoTransactionsCollection.Remove(self)
}

func (self *Transaction) Begin() error {

	self.CreateDate = time.Now()
	self.LastUpdate = time.Now()
	self.CompleteDate = time.Unix(0, 0)

	err := self.Save()
	if err != nil {
		log.Println("Failed to Insert new Transaction to DB:  " + err.Error())
		return err
	}

	transactionQueue.Lock()
	var persistObj transactionsToPersist
	persistObj.t = self
	persistObj.startTime = time.Now()
	transactionQueue.queue[self.Id.Hex()] = &persistObj
	transactionQueue.Unlock()

	return nil
}

func (self *Transaction) Resume() error {

	self.LastUpdate = time.Now()

	err := self.Save()
	if err != nil {
		log.Println("Failed to Resume Transaction to DB:  " + err.Error())
		return err
	}

	transactionQueue.Lock()

	persistObj := transactionQueue.queue[self.Id.Hex()]
	persistObj.t = self
	persistObj.startTime = time.Now()

	persistObj.t = self
	persistObj.startTime = time.Now()

	//Load the queue up with the original Data and the new data.

	transactionQueue.Unlock()

	return nil
}

func (self *Transaction) Commit() error {

	transactionQueue.RLock()
	tPersist, ok := transactionQueue.queue[self.Id.Hex()]
	transactionQueue.RUnlock()

	if ok == false {
		return errors.New("Transaction not present in queue.  Make sure to perform a Begin on your transaction.")
	}

	//Attempt to Persist the items in the transaction.
	rollBack := false
	var rollBackErrorMessage string

	for _, entityTran := range tPersist.newItems {

		if entityTran.changeType == TRANSACTION_CHANGETYPE_DELETE {
			err := entityTran.entity.Delete()
			if err != nil {
				rollBack = true
				rollBackErrorMessage = "Failed to delete object in transaction collection.  Rolling back transaction id " + self.Id.Hex() + ".\n" + err.Error()
				break
			}
		} else {
			err := entityTran.entity.Save()
			if err != nil {
				rollBack = true
				rollBackErrorMessage = "Failed to persist object in transaction collection.  Rolling back transaction id " + self.Id.Hex() + ".\n" + err.Error()
				break
			}
		}
		entityTran.committed = true
	}

	//Attempt to persist the Historical Records
	if !rollBack {

		for _, entityTran := range tPersist.originalItems {

			err := entityTran.entity.Save()
			if err != nil {
				rollBack = true
				rollBackErrorMessage = "Failed to persist object in history collection.  Rolling back transaction id " + self.Id.Hex() + ".\n" + err.Error()
				break
			}
			entityTran.committed = true
		}
	}

	//If everything persists successfully just return nil

	if !rollBack {

		//Cleanup the transactionObjects Table
		var transactionObjects TransactionObjects
		tObjects, err := transactionObjects.Search("tId", self.Id.Hex())

		if err == nil {
			for _, tObj := range tObjects {
				tObj.Delete() //May create Orphan transaction Objects if delete Fails.  A query could be performed to find orphans.
			}
		}

		transactionQueue.Lock()
		delete(transactionQueue.queue, self.Id.Hex())
		transactionQueue.Unlock()

		self.Committed = true
		self.CompleteDate = time.Now()
		self.Collections = removeDuplicates(self.Collections)
		err = self.Save()
		if err != nil {
			return errors.New("Failed to Finalized Transaction Record.")
		}

		return nil
	}

	//Rollback all the original data to the tables and return the rollback error.
	for i, entityTran := range tPersist.newItems {
		if entityTran.committed {
			if entityTran.changeType == TRANSACTION_CHANGETYPE_INSERT {

				err := entityTran.entity.Delete()
				if err != nil {
					rollBackErrorMessage = rollBackErrorMessage + "\nFailed to Rollback data in rollback process.  " + err.Error()
				}

			} else {

				originalEntity := tPersist.originalItems[i]
				err := originalEntity.entity.Save()
				if err != nil {
					rollBackErrorMessage = rollBackErrorMessage + "\nFailed to Rollback data in rollback process.  " + err.Error()
				}

			}
		}
	}

	//Rollback all the original History data to the tables and return the rollback error.
	for _, entityTran := range tPersist.originalItems {
		if entityTran.committed {

			err := entityTran.entity.Delete()
			if err != nil {
				rollBackErrorMessage = rollBackErrorMessage + "\nFailed to Rollback data in rollback process.  " + err.Error()
			}
		}
	}

	self.Error = rollBackErrorMessage
	self.Save()

	return errors.New(rollBackErrorMessage)

}

func (self *Transaction) Rollback(userId string, reason string) error {

	for _, collection := range self.Collections {
		col := resolveHistoryCollection(collection)
		if col == nil {
			continue
		}

		err := col.Rollback(self.Id.Hex())
		if err != nil {
			return err
		}
	}

	self.RolledBack = true
	self.RollbackDate = time.Now()
	self.RollbackReason = reason
	self.RolledBackBy = userId
	self.Save()

	return nil
}

func (obj *Transaction) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Transaction) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}
