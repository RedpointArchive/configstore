package main

import (
	"context"
	"fmt"
)

func (s *transactionProcessor) processTransaction(
	ctx context.Context,
	schema *Schema,
	req *MetaTransaction,
) (*MetaTransactionResult, error) {
	resp := &MetaTransactionResult{}
	opProcessor := &operationProcessor{
		client: s.client,
	}
	for _, operation := range req.Operations {
		var operationResult *MetaOperationResult

		if opReq := operation.GetListRequest(); opReq != nil {
			opResp, err := opProcessor.operationList(ctx, schema, opReq)
			operationResult = toOperationResult(ctx, schema, &MetaOperationResult_ListResponse{
				ListResponse: opResp,
			}, err)
		}
		if opReq := operation.GetGetRequest(); opReq != nil {
			opResp, err := opProcessor.operationGet(ctx, schema, opReq)
			operationResult = toOperationResult(ctx, schema, &MetaOperationResult_GetResponse{
				GetResponse: opResp,
			}, err)
		}
		if opReq := operation.GetUpdateRequest(); opReq != nil {
			opResp, err := opProcessor.operationUpdate(ctx, schema, opReq)
			operationResult = toOperationResult(ctx, schema, &MetaOperationResult_UpdateResponse{
				UpdateResponse: opResp,
			}, err)
		}
		if opReq := operation.GetCreateRequest(); opReq != nil {
			opResp, err := opProcessor.operationCreate(ctx, schema, opReq)
			operationResult = toOperationResult(ctx, schema, &MetaOperationResult_CreateResponse{
				CreateResponse: opResp,
			}, err)
		}
		if opReq := operation.GetDeleteRequest(); opReq != nil {
			opResp, err := opProcessor.operationDelete(ctx, schema, opReq)
			operationResult = toOperationResult(ctx, schema, &MetaOperationResult_DeleteResponse{
				DeleteResponse: opResp,
			}, err)
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
				ErrorMessage: fmt.Sprintf("%v", err),
			},
		}
	} else {
		return &MetaOperationResult{
			Operation: opResp,
		}
	}
}
