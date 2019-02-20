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
   * @private @const {!proto.meta.ConfigstoreMetaServiceClient} The delegate callback based client
   */
  this.delegateClient_ = new proto.meta.ConfigstoreMetaServiceClient(
      hostname, credentials, options);

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
 * @param {!Object<string, string>} metadata User defined
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
      metadata,
      methodInfo_ConfigstoreMetaService_GetSchema,
      callback);
};


/**
 * @param {!proto.meta.GetSchemaRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.GetSchemaResponse>}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.getSchema =
    function(request, metadata) {
  return new Promise((resolve, reject) => {
    this.delegateClient_.getSchema(
      request, metadata, (error, response) => {
        error ? reject(error) : resolve(response);
      });
  });
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
 * @param {!Object<string, string>} metadata User defined
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
      metadata,
      methodInfo_ConfigstoreMetaService_MetaList,
      callback);
};


/**
 * @param {!proto.meta.MetaListEntitiesRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.meta.MetaListEntitiesResponse>}
 *     The XHR Node Readable Stream
 */
proto.meta.ConfigstoreMetaServicePromiseClient.prototype.metaList =
    function(request, metadata) {
  return new Promise((resolve, reject) => {
    this.delegateClient_.metaList(
      request, metadata, (error, response) => {
        error ? reject(error) : resolve(response);
      });
  });
};


module.exports = proto.meta;

