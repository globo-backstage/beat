const async = require('async');

const registry = require('./registry');
const config = require('./config');
const context = require('./context');


function handler(req, res) {
    bufferizeBody(req, (err, body) => {
        if (err) {
            writeError(err, res);
            return;
        }

        runCustomCodes(body, (err) => {
            if (err) {
                writeError(err, res);
                return;
            }

            res.writeHead(200, {
                'Content-Type': 'text/plain'
            });

            res.end('ok');
        });
    });
}

function runCustomCodes(body, callback) {
    registry.populateNewCustomCodes(body.customCodes, (err) => {
        var finished = false;
        var finish = (err) => {
            if (finished) {
                return;
            }
            finished = true;
            callback(err);
        };

        async.eachSeries(body.customCodes, (customCodeId, customCodeCallback) => {
            if (finished) {
                customCodeCallback(null);
                return;
            }
            runScript(customCodeId, body, customCodeCallback, finish);
        }, (err, blah) => {
            console.info('blah ...', blah);
            finish(err);
        });
    });
}

function runScript(customCodeId, body, callback) {
    var script = registry.getScript(customCodeId);
    var ctx = context.create(
        customCodeId,
        body.req,
        {executeHookName: body.hook},
        callback
    );
    script.runInContext(ctx, {timeout: config.syncTimeout});
}

function bufferizeBody(req, callback) {
    var body = [];

    req
        .on('data', (chunk) => {
            body.push(chunk);
        })
        .on('end', () =>{
            body = Buffer.concat(body).toString();
            try {
                callback(null, JSON.parse(body));
            } catch(error) {
                callback(error);
            }
        });
}

function writeError(err, res) {
    res.writeHead(500, {'Content-Type': 'application/json'});
    res.end(JSON.stringify({error: err.toString()}));
}

module.exports = handler;
