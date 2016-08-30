package model

import (
	"encoding/base64"
	"sync"
	"time"
)

const (
	TRANSACTION_DATATYPE_ORIGINAL = 1
	TRANSACTION_DATATYPE_NEW      = 2

	TRANSACTION_CHANGETYPE_INSERT = 1
	TRANSACTION_CHANGETYPE_UPDATE = 2
	TRANSACTION_CHANGETYPE_DELETE = 3

	MGO_RECORD_NOT_FOUND = "not found"
)

type modelEntity interface {
	Save() error
	Delete() error
}

type modelCollection interface {
	Rollback(transactionId string) error
}

type tQueue struct {
	sync.RWMutex
	queue map[string]*transactionsToPersist
}

type transactionsToPersist struct {
	t             *Transaction
	newItems      []entityTransaction
	originalItems []entityTransaction
	startTime     time.Time
}

type entityTransaction struct {
	changeType int
	committed  bool
	entity     modelEntity
}

var transactionQueue tQueue

func init() {
	transactionQueue.queue = make(map[string]*transactionsToPersist)
	go clearTransactionQueue()
}

//Every 12 hours check the transactionQueue and remove any outstanding stale transactions > 48 hours old
func clearTransactionQueue() {

	transactionQueue.Lock()

	for key, value := range transactionQueue.queue {

		if time.Since(value.startTime).Hours() > 48 {
			delete(transactionQueue.queue, key)
		}
	}

	transactionQueue.Unlock()

	time.Sleep(12 * time.Hour)
	clearTransactionQueue()
}

func getBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func decodeBase64(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return string(data[:]), nil
}

func getNow() time.Time {
	return time.Now()
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
