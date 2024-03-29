package model

import (
	// "encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/store"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"sync"
)

var Transactions modelTransactions

type modelTransactions struct{}

var mongoTransactionsCollection *mgo.Collection
var collectionTransactionMutex *sync.RWMutex

func init() {
	store.RegisterHistoryStore(&Transactions)
	collectionTransactionMutex = &sync.RWMutex{}
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

func (obj modelTransactions) SetCollection(mdb *mgo.Database) {

	collectionTransactionMutex.Lock()
	mongoTransactionsCollection = mdb.C("Transactions")
	ci := mgo.CollectionInfo{ForceIdIndex: true}
	mongoTransactionsCollection.Create(&ci)
	collectionTransactionMutex.Unlock()
}

func (obj modelTransactions) Query() *Query {
	var query Query

	for {
		collectionTransactionMutex.RLock()
		collection := mongoTransactionsCollection
		collectionTransactionMutex.RUnlock()

		if collection != nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}

	collectionTransactionMutex.RLock()
	query.collection = mongoTransactionsCollection
	collectionTransactionMutex.RUnlock()

	return &query
}

func (obj *modelTransactions) Index() error {
	log.Println("Building Indexes for MongoDB collection Transactions:")
	for key, value := range dbServices.GetDBIndexes(Transaction{}) {
		index := mgo.Index{
			Key:        []string{key},
			Unique:     false,
			Background: true,
		}

		if value == "unique" {
			index.Unique = true
		}

		collectionTransactionMutex.RLock()
		err := mongoTransactionsCollection.EnsureIndex(index)
		collectionTransactionMutex.RUnlock()
		if err != nil {
			log.Println("Failed to create index for Transaction." + key + ":  " + err.Error())
		} else {
			log.Println("Successfully created index for Transaction." + key)
		}
	}
	return nil
}

func (obj *modelTransactions) Start() (*Transaction, error) {
	t := Transaction{}
	err := t.Begin()
	return &t, err
}

func (obj *modelTransactions) New(userId string) (*Transaction, error) {
	t := Transaction{}
	t.UserId = userId
	err := t.Begin()

	return &t, err
}

func (self *Transaction) Save() error {
	collectionTransactionMutex.RLock()
	collection := mongoTransactionsCollection
	collectionTransactionMutex.RUnlock()

	objectId := bson.NewObjectId()
	if self.Id != "" {
		objectId = self.Id
	}
	changeInfo, err := collection.UpsertId(objectId, &self)
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
	collectionTransactionMutex.RLock()
	collection := mongoTransactionsCollection
	collectionTransactionMutex.RUnlock()
	return collection.Remove(self)
}

func (self *Transaction) Begin() error {

	self.Id = bson.NewObjectId()
	self.CreateDate = time.Now()
	self.LastUpdate = time.Now()
	self.CompleteDate = time.Unix(0, 0)

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

	affectedRecordCount := make(map[string]int, 0)
	for bsonId, entityTran := range tPersist.newItems {
		pieces := strings.Split(bsonId, "_")
		// Concatenate a human readable map of counts of what we just did
		typeAction := "(s) were updated/inserted -> "
		if entityTran.changeType == TRANSACTION_CHANGETYPE_DELETE {
			typeAction = "(s) were deleted -> "
		}
		if len(pieces) == 2 {
			_, ok := affectedRecordCount[pieces[0]+typeAction]
			if ok {
				affectedRecordCount[pieces[0]+typeAction]++
			} else {
				affectedRecordCount[pieces[0]+typeAction] = 1
			}
		}
		if entityTran.changeType == TRANSACTION_CHANGETYPE_DELETE {
			err := entityTran.entity.Delete()
			if err != nil && err.Error() != "not found" {
				rollBack = true
				rollBackErrorMessage = "Failed to delete object in transaction collection.  Rolling back transaction id " + self.Id.Hex() + ".\n" + err.Error() + "\n\nRecord:\n\n" + core.Debug.GetDump(entityTran.entity)
				break
			}
		} else {
			err := entityTran.entity.Save()
			if err != nil {
				rollBack = true
				rollBackErrorMessage = "Failed to persist object in transaction collection.  Rolling back transaction id " + self.Id.Hex() + ".\n" + err.Error() + "\n\nRecord:\n\n" + core.Debug.GetDump(entityTran.entity)
				break
			}
		}
		entityTran.committed = true
	}
	log.Println("Transaction Id " + self.Id.Hex() + " committed these record counts")
	log.Printf("%+v", affectedRecordCount)

	//Attempt to persist the Historical Records
	if !rollBack {

		for _, entityTran := range tPersist.originalItems {

			//err := entityTran.entity.Save()
			//if err != nil {
			//	rollBack = true
			//	rollBackErrorMessage = "Failed to persist object in history collection.  Rolling back transaction id " + self.Id.Hex() + ".\n" + err.Error()
			//	break
			//}
			entityTran.committed = true
		}
	}

	//If everything persists successfully just return nil

	if !rollBack {

		transactionQueue.Lock()
		delete(transactionQueue.queue, self.Id.Hex())
		transactionQueue.Unlock()

		self.Committed = true
		self.CompleteDate = time.Now()
		self.Collections = removeDuplicates(self.Collections)
		err := self.Save()
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

func (self *Transaction) RollbackTransaction() error {
	return self.Rollback("", "")
}

func (self *Transaction) Rollback(userId string, reason string) error {

	for _, collection := range self.Collections {
		col := ResolveHistoryCollection(collection)
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

func (obj *Transaction) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *Transaction) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}
