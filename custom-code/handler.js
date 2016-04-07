var registry = require('./registry');

function handler(req, res) {
    bufferizeBody(req, (err, body) => {
        if (err) {
            writeError(err, res);
            return;
        }

        registry.populateNewCustomCodes(body.customCodes, function(err) {
            if (err) {
                writeError(err, res);
                return;
            }

            res.writeHead(200, {
                'Content-Type': 'text/plain'
            });

            console.info('received: ', body);
            res.end('ok');
        });
    });

}

function bufferizeBody(req, callback) {
    var body = [];

    req.on('data', function(chunk) {
        body.push(chunk);
    }).on('end', function() {
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
