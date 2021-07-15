package bot

import (
	"fmt"
	"log"

	"github.com/ibrokemypie/fedibooks-go/internal/db"
	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
)

func InitBot() {
	db, err := db.InitialiseDB()
	if err != nil {
		log.Fatal(err)
	}

	statuses := GetFollowingStatuses()

	// Create a write transaction
	txn := db.Txn(true)
	for _, status := range statuses {
		err = txn.Insert("status", status)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Commit the transaction
	txn.Commit()

	// Create read-only transaction
	txn = db.Txn(false)
	defer txn.Abort()
	it, err := txn.Get("status", "id")
	if err != nil {
		log.Fatal(err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		status := obj.(fedi.Status)
		fmt.Println("id: " + status.ID + ", author_id: " + status.Account.ID + ", text: " + status.Text)
	}

}
