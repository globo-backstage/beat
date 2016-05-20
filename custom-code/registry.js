'use strict';

const async = require('async');
const request = require('request');
const vm = require('vm');

const context = require('./context');
const config = require('./config');

const customCodePath = `http://${config.beatHost}/api/custom-codes/`;

const COMPILE_CUSTOM_CODE_FOOT = (
    'if (Backstage.executeHookName && Backstage.defines[Backstage.executeHookName]) {\n'+
        'Backstage.defines[Backstage.executeHookName](Backstage.request, Backstage.response);\n'+
    '} else {\n'+
        'Backstage.response.success();\n'+
    '}');

class Registry {
    constructor() {
        this._codes = {};
    }

    populateNewCustomCodes(customCodes, fn) {
        var newCustomCodes = customCodes.filter((customCodeId) => {
            return this._codes[customCodeId] === undefined;
        });

        async.filter(newCustomCodes, fetchCustomCode.bind(this), function(err, customCodeIds) {
            fn(err);
        });
    }

    put(customCodeId, customCode) {
        this._codes[customCodeId] = customCode;
    }

    get(customCodeId) {
        return this._codes[customCodeId];
    }

    getScript(customCodeId) {
        var item = this.get(customCodeId);
        return item ? item.script : null;
    }
}

function fetchCustomCode(customCodeId, callback) {
    const url = `${customCodePath}${customCodeId}`;
    request({url, json: true}, (error, res, body) => {
        if (error) {
            callback(error);
            return;
        }

        if (res.statusCode != 200) {
            callback(new Error(`status code ${res.statusCode}`));
            return;
        }

        compileAndPutCustomCode.call(this, body);
        callback(null);
    });
}

function compileAndPutCustomCode(customCode) {
    this.put(customCode.id, {
        hooks: customCode.hooks,
        script: compileCode(customCode)
    });
}

function compileCode(code) {
    console.log(`${process.pid}: Compiling Custom Code ${code.id}`);

    var scriptArguments = Object.keys(context.sandbox).join(', ');
    var closureStart = '(function(' + scriptArguments + ') {\n';
    var closureEnd = '\n})(' + scriptArguments + ');\n';
    var text = closureStart + code.code + closureEnd + COMPILE_CUSTOM_CODE_FOOT;

    return new vm.Script(text, {filename: 'custom-code-'+code.id+'.js', displayErrors: true, timeout: config.syncTimeout});
}


module.exports = new Registry();
