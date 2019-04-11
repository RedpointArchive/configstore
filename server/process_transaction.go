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

		if opReq := operation.GetListRequest(); opReq != nil {
			opResp, err := operationList(ctx, schema, opReq)
			operationResult = toOperationResult(opResp, err)
		}
		if opReq := operation.GetGetRequest(); opReq != nil {
			opResp, err := operationGet(ctx, schema, opReq)
			operationResult = toOperationResult(opResp, err)
		}
		if opReq := operation.GetUpdateRequest(); opReq != nil {
			opResp, err := operationUpdate(ctx, schema, opReq)
			operationResult = toOperationResult(opResp, err)
		}
		if opReq := operation.GetCreateRequest(); opReq != nil {
			opResp, err := operationCreate(ctx, schema, opReq)
			operationResult = toOperationResult(opResp, err)
		}
		if opReq := operation.GetDeleteRequest(); opReq != nil {
			opResp, err := operationDelete(ctx, schema, opReq)
			operationResult = toOperationResult(opResp, err)
		}

		resp.OperationResults = append(
			resp.OperationResults,
			operationResult,
		)
	}
	return resp, nil
}

func toOperationResult(
	ctx context.Context, 
	schema *Schema, 
	opResp isMetaOperationResult_Operation,
	err error,
) *MetaOperationResult {
	if err != nil {
		return &MetaOperationResult{
			Error: &MetaOperationResultError{
				ErrorMessage: fmt.Sprintf("%v", err)
			},
		}
	} else {
		return &MetaOperationResult{
			Operation: opResp,
		}
	}
}