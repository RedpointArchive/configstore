package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
)

func (s *transactionProcessor) processTransaction(
	ctx context.Context,
	schema *Schema,
	req *MetaTransaction,
) (*MetaTransactionResult, error) {
	resp := &MetaTransactionResult{}
	resp.OperationResults = make([]*MetaOperationResult, len(req.Operations), len(req.Operations))

	err := s.client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		opProcessor := createOperationProcessor(s.client, tx)

		readStates := make([]interface{}, len(req.Operations), len(req.Operations))
		readErrors := make([]error, len(req.Operations), len(req.Operations))

		var mutatedKeys []*firestore.DocumentRef
		var deletedKeys []*firestore.DocumentRef

		for i, operation := range req.Operations {
			if opReq := operation.GetListRequest(); opReq != nil {
				readState, err := opProcessor.operationListRead(ctx, schema, opReq)
				readStates[i] = readState
				readErrors[i] = err
			}
			if opReq := operation.GetGetRequest(); opReq != nil {
				readState, err := opProcessor.operationGetRead(ctx, schema, opReq)
				readStates[i] = readState
				readErrors[i] = err
			}
			if opReq := operation.GetUpdateRequest(); opReq != nil {
				readState, err := opProcessor.operationUpdateRead(ctx, schema, opReq)
				readStates[i] = readState
				readErrors[i] = err
			}
			if opReq := operation.GetCreateRequest(); opReq != nil {
				readState, err := opProcessor.operationCreateRead(ctx, schema, opReq)
				readStates[i] = readState
				readErrors[i] = err
			}
			if opReq := operation.GetDeleteRequest(); opReq != nil {
				readState, err := opProcessor.operationDeleteRead(ctx, schema, opReq)
				readStates[i] = readState
				readErrors[i] = err
			}
		}

		for i, operation := range req.Operations {
			var operationResult *MetaOperationResult

			if opReq := operation.GetListRequest(); opReq != nil {
				if readErrors[i] != nil {
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_ListResponse{
						ListResponse: nil,
					}, readErrors[i])
				} else {
					result, err := opProcessor.operationListWrite(ctx, schema, opReq, readStates[i])
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_ListResponse{
						ListResponse: result,
					}, err)
				}
			}
			if opReq := operation.GetGetRequest(); opReq != nil {
				if readErrors[i] != nil {
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_GetResponse{
						GetResponse: nil,
					}, readErrors[i])
				} else {
					result, err := opProcessor.operationGetWrite(ctx, schema, opReq, readStates[i])
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_GetResponse{
						GetResponse: result,
					}, err)
				}
			}
			if opReq := operation.GetUpdateRequest(); opReq != nil {
				if readErrors[i] != nil {
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_UpdateResponse{
						UpdateResponse: nil,
					}, readErrors[i])
				} else {
					opResp, err := opProcessor.operationUpdateWrite(ctx, schema, opReq, readStates[i])
					if err == nil {
						ref, err := convertMetaKeyToDocumentRef(
							s.client,
							opResp.Entity.Key,
						)
						if err == nil {
							mutatedKeys = append(
								mutatedKeys,
								ref,
							)
						}
					}
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_UpdateResponse{
						UpdateResponse: opResp,
					}, err)
				}
			}
			if opReq := operation.GetCreateRequest(); opReq != nil {
				if readErrors[i] != nil {
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_CreateResponse{
						CreateResponse: nil,
					}, readErrors[i])
				} else {
					opResp, err := opProcessor.operationCreateWrite(ctx, schema, opReq, readStates[i])
					if err == nil {
						ref, err := convertMetaKeyToDocumentRef(
							s.client,
							opResp.Entity.Key,
						)
						if err == nil {
							mutatedKeys = append(
								mutatedKeys,
								ref,
							)
						}
					}
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_CreateResponse{
						CreateResponse: opResp,
					}, err)
				}
			}
			if opReq := operation.GetDeleteRequest(); opReq != nil {
				if readErrors[i] != nil {
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_CreateResponse{
						CreateResponse: nil,
					}, readErrors[i])
				} else {
					opResp, err := opProcessor.operationDeleteWrite(ctx, schema, opReq, readStates[i])
					if err == nil {
						ref, err := convertMetaKeyToDocumentRef(
							s.client,
							opResp.Entity.Key,
						)
						if err == nil {
							deletedKeys = append(
								deletedKeys,
								ref,
							)
						}
					}
					operationResult = toOperationResult(ctx, schema, &MetaOperationResult_DeleteResponse{
						DeleteResponse: opResp,
					}, err)
				}
			}

			if operationResult != nil {
				resp.OperationResults[i] = operationResult
			}
		}

		if len(mutatedKeys) > 0 || len(deletedKeys) > 0 {
			ref := s.client.Collection("Transaction").NewDoc()
			transaction := make(map[string]interface{})
			transaction["mutatedKeys"] = mutatedKeys
			transaction["deletedKeys"] = deletedKeys
			transaction["dateSubmitted"] = time.Now()
			transaction["description"] = req.Description
			tx.Create(ref, transaction)
		}

		return nil
	})
	if err != nil {
		return nil, err
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
