const fs = require('fs');

let js = fs.readFileSync('src/api/meta_pb.js', 'utf8')
if (js.indexOf("/* eslint-disable */") === -1) {
    js = "/* eslint-disable */\n" + js
}
fs.writeFileSync('src/api/meta_pb.js', js)

js = fs.readFileSync('src/api/meta_grpc_web_pb.js', 'utf8')
if (js.indexOf("/* eslint-disable */") === -1) {
    js = "/* eslint-disable */\n" + js
}
fs.writeFileSync('src/api/meta_grpc_web_pb.js', js)