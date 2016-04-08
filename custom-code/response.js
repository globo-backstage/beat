'use strict';

class CustomCodeResponse {
    constructor(callback) {
        this.finished = false;
        this.callback = (err, body) => {
            if (this.finished) {
                return;
            }

            this.finished = true;
            callback(err, body);
        };
    }

    success(body) {
        if (body) {
            this.callback(null, body);
        } else {
            this.callback(null);
        }
    }

    validationError(msg) {
        this._error(422, msg);
    }

    internalServerError(msg) {
        this._error(500, msg);
    }

    notFoundError(msg) {
        this._error(404, msg);
    }

    _error(statusCode, msg) {
        var errMsg = msg;
        
        if (msg && msg.constructor.name === 'Error') {
            errMsg = msg.message;
        }
        
        var err = new Error(errMsg);
        err.statusCode = statusCode;

        this.callback(err);
    }
}

module.exports = CustomCodeResponse;
