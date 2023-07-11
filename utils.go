package nearlake

import (
	"github.com/DmitriyHellyeah/near-lake-framework-go/types"
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func isTxReceiverWatched(tx types.IndexerTransactionWithOutcome, watchingList []string) bool {
	return contains(watchingList, tx.Transaction.ReceiverId)
}

func remove(slice *[]string, element string) []string {
	index := -1
	for i, val := range *slice {
		if val == element {
			index = i
			break
		}
	}

	if index != -1 {
		*slice = append((*slice)[:index], (*slice)[index+1:]...)
	}

	return *slice
}