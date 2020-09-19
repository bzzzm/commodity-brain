package db

import (
	"encoding/json"
	"fmt"
	"github.com/bzzzm/commodity-brain/pkg/utils"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const (
	BoardBucket = "board"
	CVBucket    = "cv"
)

// setupDB create a new bbolt database and creates the buckets
// required by the app to run
func SetupDB(path string) (*bbolt.DB, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {

		_, err = tx.CreateBucketIfNotExists([]byte(BoardBucket))
		if err != nil {
			return fmt.Errorf("could not create board bucket: %v", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(CVBucket))
		if err != nil {
			return fmt.Errorf("could not create cv bucket: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not setup db: %v", err)
	}
	return db, nil
}

func DumpBucket(db *bbolt.DB, bucket string) error {
	output := make(map[string]string)
	err := db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			switch len(v) {
			case 1:
				output[string(k)] = utils.BoolByteToString(v[0])
			}

		}
		b, err := json.Marshal(output)
		if err != nil {
			return err
		}
		zap.S().Infof("current db: %v", string(b))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
