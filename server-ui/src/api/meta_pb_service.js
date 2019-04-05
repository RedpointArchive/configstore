// package: meta
// file: meta.proto

var meta_pb = require("./meta_pb");
var grpc = require("@improbable-eng/grpc-web").grpc;

var ConfigstoreMetaService = (function () {
  function ConfigstoreMetaService() {}
  ConfigstoreMetaService.serviceName = "meta.ConfigstoreMetaService";
  return ConfigstoreMetaService;
}());

ConfigstoreMetaService.GetSchema = {
  methodName: "GetSchema",
  service: ConfigstoreMetaService,
  requestStream: false,
  responseStream: false,
  requestType: meta_pb.GetSchemaRequest,
  responseType: meta_pb.GetSchemaResponse
};

ConfigstoreMetaService.MetaList = {
  methodName: "MetaList",
  service: ConfigstoreMetaService,
  requestStream: false,
  responseStream: false,
  requestType: meta_pb.MetaListEntitiesRequest,
  responseType: meta_pb.MetaListEntitiesResponse
};

ConfigstoreMetaService.MetaGet = {
  methodName: "MetaGet",
  service: ConfigstoreMetaService,
  requestStream: false,
  responseStream: false,
  requestType: meta_pb.MetaGetEntityRequest,
  responseType: meta_pb.MetaGetEntityResponse
};

ConfigstoreMetaService.MetaUpdate = {
  methodName: "MetaUpdate",
  service: ConfigstoreMetaService,
  requestStream: false,
  responseStream: false,
  requestType: meta_pb.MetaUpdateEntityRequest,
  responseType: meta_pb.MetaUpdateEntityResponse
};

ConfigstoreMetaService.MetaCreate = {
  methodName: "MetaCreate",
  service: ConfigstoreMetaService,
  requestStream: false,
  responseStream: false,
  requestType: meta_pb.MetaCreateEntityRequest,
  responseType: meta_pb.MetaCreateEntityResponse
};

ConfigstoreMetaService.MetaDelete = {
  methodName: "MetaDelete",
  service: ConfigstoreMetaService,
  requestStream: false,
  responseStream: false,
  requestType: meta_pb.MetaDeleteEntityRequest,
  responseType: meta_pb.MetaDeleteEntityResponse
};

ConfigstoreMetaService.GetDefaultPartitionId = {
  methodName: "GetDefaultPartitionId",
  service: ConfigstoreMetaService,
  requestStream: false,
  responseStream: false,
  requestType: meta_pb.GetDefaultPartitionIdRequest,
  responseType: meta_pb.GetDefaultPartitionIdResponse
};

exports.ConfigstoreMetaService = ConfigstoreMetaService;

function ConfigstoreMetaServiceClient(serviceHost, options) {
  this.serviceHost = serviceHost;
  this.options = options || {};
}

ConfigstoreMetaServiceClient.prototype.getSchema = function getSchema(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ConfigstoreMetaService.GetSchema, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ConfigstoreMetaServiceClient.prototype.metaList = function metaList(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ConfigstoreMetaService.MetaList, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ConfigstoreMetaServiceClient.prototype.metaGet = function metaGet(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ConfigstoreMetaService.MetaGet, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ConfigstoreMetaServiceClient.prototype.metaUpdate = function metaUpdate(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ConfigstoreMetaService.MetaUpdate, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ConfigstoreMetaServiceClient.prototype.metaCreate = function metaCreate(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ConfigstoreMetaService.MetaCreate, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ConfigstoreMetaServiceClient.prototype.metaDelete = function metaDelete(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ConfigstoreMetaService.MetaDelete, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

ConfigstoreMetaServiceClient.prototype.getDefaultPartitionId = function getDefaultPartitionId(requestMessage, metadata, callback) {
  if (arguments.length === 2) {
    callback = arguments[1];
  }
  var client = grpc.unary(ConfigstoreMetaService.GetDefaultPartitionId, {
    request: requestMessage,
    host: this.serviceHost,
    metadata: metadata,
    transport: this.options.transport,
    debug: this.options.debug,
    onEnd: function (response) {
      if (callback) {
        if (response.status !== grpc.Code.OK) {
          var err = new Error(response.statusMessage);
          err.code = response.status;
          err.metadata = response.trailers;
          callback(err, null);
        } else {
          callback(null, response.message);
        }
      }
    }
  });
  return {
    cancel: function () {
      callback = null;
      client.close();
    }
  };
};

exports.ConfigstoreMetaServiceClient = ConfigstoreMetaServiceClient;

