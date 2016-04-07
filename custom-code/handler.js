function handler(req, res) {
    bufferizeBody(req, (err, body) => {
        if (err) {
            res.writeHead(500, {'Content-Type': 'application/json'});
            res.end(JSON.stringify({error: err.toString()}));
            return;
        }

        res.writeHead(200, {
            'Content-Type': 'text/plain'
        });
        console.info('received: ', body);
        res.end('ok');
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

module.exports = handler;
