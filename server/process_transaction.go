package main

import (
	"context"
	"fmt"
)

func processTransaction(
	ctx context.Context, 
	schema *Schema, 
	req *MetaTransaction,
) (*MetaTransactionResult, error) {
	resp := &MetaTransactionResult{}
	for i, operation := range req.Operations {
		var operationResult *MetaOperationResult

		// 

		resp.OperationResults := append(
			resp.OperationResults,
			operationResult,
		)
	}
}