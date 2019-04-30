package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
)

const traceDefaultServiceName = "configstore"

// ==== tracing code begin ====
// ==== tracing.go and generator_gosdk_template.gotxt should be identical here ====

var hasInited bool
var traceServiceName string
var traceEnabled bool

func doInit() {
	traceServiceName = os.Getenv("CONFIGSTORE_SERVICE_NAME")
	if traceServiceName == "" {
		traceServiceName = traceDefaultServiceName
	}
	traceEnabled = os.Getenv("CONFIGSTORE_ENABLE_TRACE") == "1"
	hasInited = true
}

func getTraceServiceName() string {
	if !hasInited {
		doInit()
	}
	return traceServiceName
}

func isTracingEnabled() bool {
	if !hasInited {
		doInit()
	}
	return traceEnabled
}

func serializeMetaEntityForTrace(entity *MetaEntity) string {
	var lines []string
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  key = %s", SerializeKey(entity.Key)))
	for _, value := range entity.Values {
		serValue := "(unknown)"
		switch value.Type {
		case ValueType_double:
			serValue = fmt.Sprintf("%f", value.DoubleValue)
		case ValueType_int64:
			serValue = fmt.Sprintf("%d", value.Int64Value)
		case ValueType_string:
			serValue = value.StringValue
		case ValueType_timestamp:
			serValue = fmt.Sprintf("%+v", value.TimestampValue)
		case ValueType_boolean:
			if value.BooleanValue {
				serValue = "true"
			} else {
				serValue = "false"
			}
		case ValueType_bytes:
			serValue = "(bytes)"
		case ValueType_key:
			if value.KeyValue == nil {
				serValue = "(nil)"
			} else {
				serValue = SerializeKey(value.KeyValue)
			}
		case ValueType_uint64:
			serValue = fmt.Sprintf("%d", value.Uint64Value)
		}
		lines = append(lines, fmt.Sprintf("  %d = %s", value.Id, serValue))
	}
	return strings.Join(lines, "\n")
}

func RecordTrace(trace *ConfigstoreTraceEntry) {
	if !isTracingEnabled() {
		return
	}

	RecordTraceWithCustomEntity(trace, nil)
}

func RecordTraceWithCustomEntity(trace *ConfigstoreTraceEntry, entity interface{}) {
	if !isTracingEnabled() {
		return
	}

	var entityString string
	if entity != nil {
		metaEntity, ok := entity.(*MetaEntity)
		if ok {
			entityString = serializeMetaEntityForTrace(metaEntity)
		} else {
			entityString = fmt.Sprintf("%+v", entity)
		}
	} else {
		if trace.Entity != nil {
			entityString = serializeMetaEntityForTrace(trace.Entity)
		}
	}

	switch trace.Type {
	case ConfigstoreTraceEntry_INITIAL_STATE_SEND_BEGIN:
		log.Printf("%s: initial state: send begin", trace.OperatorId)
	case ConfigstoreTraceEntry_INITIAL_STATE_SEND_ENTITY:
		log.Printf("%s: initial state: send entity: %s", trace.OperatorId, entityString)
	case ConfigstoreTraceEntry_INITIAL_STATE_SEND_END:
		log.Printf("%s: initial state: send end", trace.OperatorId)
	case ConfigstoreTraceEntry_INITIAL_STATE_RECEIVE_BEGIN:
		log.Printf("%s: initial state: receive begin", trace.OperatorId)
	case ConfigstoreTraceEntry_INITIAL_STATE_RECEIVE_ENTITY:
		log.Printf("%s: initial state: receive entity: %s", trace.OperatorId, entityString)
	case ConfigstoreTraceEntry_INITIAL_STATE_RECEIVE_END:
		log.Printf("%s: initial state: receive end", trace.OperatorId)
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_SEND_BEGIN:
		log.Printf("%s: transaction batch: %s: send begin", trace.OperatorId, trace.TransactionId)
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_SEND_MUTATED_ENTITY:
		log.Printf("%s: transaction batch: %s: send mutated entity: %s", trace.OperatorId, trace.TransactionId, entityString)
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_SEND_DELETED_ENTITY_KEY:
		log.Printf("%s: transaction batch: %s: send deleted entity key: %s", trace.OperatorId, trace.TransactionId, SerializeKey(trace.Key))
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_SEND_END:
		log.Printf("%s: transaction batch: %s: send end", trace.OperatorId, trace.TransactionId)
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_RECEIVE_BEGIN:
		log.Printf("%s: transaction batch: %s: receive begin", trace.OperatorId, trace.TransactionId)
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_RECEIVE_MUTATED_ENTITY:
		log.Printf("%s: transaction batch: %s: receive mutated entity: %s", trace.OperatorId, trace.TransactionId, entityString)
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_RECEIVE_DELETED_ENTITY_KEY:
		log.Printf("%s: transaction batch: %s: receive deleted entity key: %s", trace.OperatorId, trace.TransactionId, SerializeKey(trace.Key))
	case ConfigstoreTraceEntry_TRANSACTION_BATCH_RECEIVE_END:
		log.Printf("%s: transaction batch: %s: receive end", trace.OperatorId, trace.TransactionId)
	case ConfigstoreTraceEntry_IN_MEMORY_STORE_ENTITY:
		log.Printf("%s: in memory: store %s", trace.OperatorId, entityString)
	case ConfigstoreTraceEntry_IN_MEMORY_DELETE_ENTITY:
		log.Printf("%s: in memory: delete %s", trace.OperatorId, SerializeKey(trace.Key))
	case ConfigstoreTraceEntry_TRANSACTION_ARRIVED:
		log.Printf("%s: transaction queue: %s: transaction arrived", trace.OperatorId, trace.TransactionId)
	case ConfigstoreTraceEntry_TRANSACTION_FINISHED_PROCESSING:
		log.Printf("%s: transaction queue: %s: finished processing transaction (%d transactions left to process)", trace.OperatorId, trace.TransactionId, trace.RemainingTransactionQueueCount)
	case ConfigstoreTraceEntry_TRANSACTION_STALLED:
		log.Printf("%s: transaction queue: %s: can't process transaction, waiting on entity snapshot with key: %s", trace.OperatorId, trace.TransactionId, SerializeKey(trace.Key))
	case ConfigstoreTraceEntry_CONFIGSTORE_CONSISTENT:
		log.Printf("%s: consistent and ready to serve transactions", trace.OperatorId)
	case ConfigstoreTraceEntry_TRANSACTION_MUTATED_ENTITY_KEY:
		log.Printf("%s: transaction queue: %s: transaction contained mutated entity key: %s", trace.OperatorId, trace.TransactionId, SerializeKey(trace.Key))
	case ConfigstoreTraceEntry_TRANSACTION_DELETED_ENTITY_KEY:
		log.Printf("%s: transaction queue: %s: transaction contained deleted entity key: %s", trace.OperatorId, trace.TransactionId, SerializeKey(trace.Key))
	case ConfigstoreTraceEntry_TRANSACTION_RECONSTRUCT_APPEND_MUTATED_ENTITY:
		log.Printf("%s: transaction reconstructor: %s: append mutated entity with key %s: %s", trace.OperatorId, trace.TransactionId, SerializeKey(trace.Key), entityString)
	case ConfigstoreTraceEntry_CLIENT_CURRENTLY_DISCONNECTED_ATTEMPTING_RECONNECT:
		log.Printf("%s: currently disconnected, attempting to reconnect...", trace.OperatorId)
	case ConfigstoreTraceEntry_CLIENT_FAILED_RECONNECT:
		log.Printf("%s: failed to reconnect: %s", trace.OperatorId, trace.ErrorString)
	case ConfigstoreTraceEntry_CLIENT_CONNECTION_REESTABLISHED:
		log.Printf("%s: connection re-established", trace.OperatorId)
	case ConfigstoreTraceEntry_CLIENT_GOT_EOF_ATTEMPTING_RECONNECTING:
		log.Printf("%s: got EOF from configstore endpoint, attempting to reconnect...", trace.OperatorId)
	case ConfigstoreTraceEntry_CLIENT_GOT_NIL_BUG_IGNORING:
		log.Printf("%s: got nil response from endpoint (this is a bug), ignoring...", trace.OperatorId)
	case ConfigstoreTraceEntry_CLIENT_GOT_UNEXPECTED_CODE_ATTEMPTING_RECONNECT:
		log.Printf("%s: got %s code from configstore endpoint, attempting to reconnect...", trace.OperatorId, trace.ReconnectionCodeString)
	case ConfigstoreTraceEntry_SERVER_STARTUP_GRPC_PORT:
		log.Printf("%s: running gRPC server on port %d...", trace.OperatorId, trace.Port)
	case ConfigstoreTraceEntry_SERVER_STARTUP_HTTP_PORT:
		log.Printf("%s: running HTTP server on port %d...", trace.OperatorId, trace.Port)
	}
}

// ==== tracing code end ====

// SerializeKey is compatibility with the generated Go SDK, so we can keep the tracing code
// in sync.
func SerializeKey(key *Key) string {
	return serializeKey(key)
}

func RecordPendingChangesByTimestamp(pendingChanges map[string]map[string]*firestore.DocumentChange) {
	if !isTracingEnabled() {
		return
	}

	for ts, entities := range pendingChanges {
		var lines []string
		lines = append(lines, fmt.Sprintf("  at %s:", ts))
		for entityKey, change := range entities {
			if change.Doc != nil {
				lines = append(lines, fmt.Sprintf("   - %s = %s", entityKey, serializeRef(change.Doc.Ref)))
			} else {
				lines = append(lines, fmt.Sprintf("   - %s = (nil doc)", entityKey))
			}
		}
		log.Printf("pending changes summary:\n%s", strings.Join(lines, "\n"))
	}
}
