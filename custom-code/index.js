const cluster = require('cluster');
const http = require('http');
const numCPUs = require('os').cpus().length;
const handler = require('./handler');

const customCodePort = parseInt(process.env['CUSTOM_CODE_PORT'] || '8100');


if (cluster.isMaster) {
    // Fork workers.
    for (var i = 0; i < numCPUs; i++) {
        cluster.fork();
    }

    cluster.on('exit', (worker, code, signal) => {
        console.log(`worker ${worker.process.pid} died`);
    });
} else {
    http.createServer((req, res) => {
        handler(req, res);
    }).listen(customCodePort);
}
