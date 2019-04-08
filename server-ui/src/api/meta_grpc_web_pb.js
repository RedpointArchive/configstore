/* eslint-disable */
/**
 * @fileoverview gRPC-Web generated client stub for meta
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.meta = require('./meta_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.meta.ConfigstoreMetaServiceClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.meta.ConfigstoreMetaServicePromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.meta.GetSchemaRequest,
 *   !proto.meta.GetSchemaResponse>}
 */
const methodInfo_ConfigstoreMetaService_GetSchema = new grpc.web.AbstractClientBase.MethodInfo(
  proto.meta.GetSchemaResponse,
  /** @param {!proto.meta.GetSchemaRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.meta.GetSchemaResponse.deserializeBinary
);


/**
 * @param {!proto.meta.GetSchemaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.meta.GetSchemaResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.meta.GetSchemaResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServiceClient.prototype.getSchema =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/GetSchema',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_GetSchema,
      callback);
};


/**
 * @param {!proto.meta.GetSchemaRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.GetSchemaResponse>}
 *     A native promise that resolves to the response
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.getSchema =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/GetSchema',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_GetSchema);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.meta.MetaListEntitiesRequest,
 *   !proto.meta.MetaListEntitiesResponse>}
 */
const methodInfo_ConfigstoreMetaService_MetaList = new grpc.web.AbstractClientBase.MethodInfo(
  proto.meta.MetaListEntitiesResponse,
  /** @param {!proto.meta.MetaListEntitiesRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.meta.MetaListEntitiesResponse.deserializeBinary
);


/**
 * @param {!proto.meta.MetaListEntitiesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.meta.MetaListEntitiesResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.meta.MetaListEntitiesResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServiceClient.prototype.metaList =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaList',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaList,
      callback);
};


/**
 * @param {!proto.meta.MetaListEntitiesRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.MetaListEntitiesResponse>}
 *     A native promise that resolves to the response
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.metaList =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaList',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaList);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.meta.MetaGetEntityRequest,
 *   !proto.meta.MetaGetEntityResponse>}
 */
const methodInfo_ConfigstoreMetaService_MetaGet = new grpc.web.AbstractClientBase.MethodInfo(
  proto.meta.MetaGetEntityResponse,
  /** @param {!proto.meta.MetaGetEntityRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.meta.MetaGetEntityResponse.deserializeBinary
);


/**
 * @param {!proto.meta.MetaGetEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.meta.MetaGetEntityResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.meta.MetaGetEntityResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServiceClient.prototype.metaGet =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaGet',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaGet,
      callback);
};


/**
 * @param {!proto.meta.MetaGetEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.MetaGetEntityResponse>}
 *     A native promise that resolves to the response
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.metaGet =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaGet',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaGet);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.meta.MetaUpdateEntityRequest,
 *   !proto.meta.MetaUpdateEntityResponse>}
 */
const methodInfo_ConfigstoreMetaService_MetaUpdate = new grpc.web.AbstractClientBase.MethodInfo(
  proto.meta.MetaUpdateEntityResponse,
  /** @param {!proto.meta.MetaUpdateEntityRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.meta.MetaUpdateEntityResponse.deserializeBinary
);


/**
 * @param {!proto.meta.MetaUpdateEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.meta.MetaUpdateEntityResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.meta.MetaUpdateEntityResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServiceClient.prototype.metaUpdate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaUpdate',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaUpdate,
      callback);
};


/**
 * @param {!proto.meta.MetaUpdateEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.MetaUpdateEntityResponse>}
 *     A native promise that resolves to the response
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.metaUpdate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaUpdate',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaUpdate);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.meta.MetaCreateEntityRequest,
 *   !proto.meta.MetaCreateEntityResponse>}
 */
const methodInfo_ConfigstoreMetaService_MetaCreate = new grpc.web.AbstractClientBase.MethodInfo(
  proto.meta.MetaCreateEntityResponse,
  /** @param {!proto.meta.MetaCreateEntityRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.meta.MetaCreateEntityResponse.deserializeBinary
);


/**
 * @param {!proto.meta.MetaCreateEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.meta.MetaCreateEntityResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.meta.MetaCreateEntityResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServiceClient.prototype.metaCreate =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaCreate',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaCreate,
      callback);
};


/**
 * @param {!proto.meta.MetaCreateEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.MetaCreateEntityResponse>}
 *     A native promise that resolves to the response
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.metaCreate =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaCreate',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaCreate);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.meta.MetaDeleteEntityRequest,
 *   !proto.meta.MetaDeleteEntityResponse>}
 */
const methodInfo_ConfigstoreMetaService_MetaDelete = new grpc.web.AbstractClientBase.MethodInfo(
  proto.meta.MetaDeleteEntityResponse,
  /** @param {!proto.meta.MetaDeleteEntityRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.meta.MetaDeleteEntityResponse.deserializeBinary
);


/**
 * @param {!proto.meta.MetaDeleteEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.meta.MetaDeleteEntityResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.meta.MetaDeleteEntityResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServiceClient.prototype.metaDelete =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaDelete',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaDelete,
      callback);
};


/**
 * @param {!proto.meta.MetaDeleteEntityRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.MetaDeleteEntityResponse>}
 *     A native promise that resolves to the response
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.metaDelete =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/MetaDelete',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_MetaDelete);
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.meta.GetDefaultPartitionIdRequest,
 *   !proto.meta.GetDefaultPartitionIdResponse>}
 */
const methodInfo_ConfigstoreMetaService_GetDefaultPartitionId = new grpc.web.AbstractClientBase.MethodInfo(
  proto.meta.GetDefaultPartitionIdResponse,
  /** @param {!proto.meta.GetDefaultPartitionIdRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.meta.GetDefaultPartitionIdResponse.deserializeBinary
);


/**
 * @param {!proto.meta.GetDefaultPartitionIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.meta.GetDefaultPartitionIdResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.meta.GetDefaultPartitionIdResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServiceClient.prototype.getDefaultPartitionId =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/GetDefaultPartitionId',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_GetDefaultPartitionId,
      callback);
};


/**
 * @param {!proto.meta.GetDefaultPartitionIdRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.GetDefaultPartitionIdResponse>}
 *     A native promise that resolves to the response
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.getDefaultPartitionId =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/meta.ConfigstoreMetaService/GetDefaultPartitionId',
      request,
      metadata || {},
      methodInfo_ConfigstoreMetaService_GetDefaultPartitionId);
};


module.exports = proto.meta;

