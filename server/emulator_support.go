package main

import (
	"os"
)

func runWithoutFirestoreTransactionalQueries() bool {
	// ... of course the Firestore emulator doesn't support transactional queries :(
	//
	// When running against the emulator, we have to make any queries (non-key based get operations and non-writes) outside
	// of the transaction, because rather than hiding this from the user (by saying, just running the query anyway), the Firestore
	// emulator returns an error with no description.
	return os.Getenv("FIRESTORE_EMULATOR_HOST") != ""
}
