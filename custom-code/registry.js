'use strict';

const async = require('async');
const request = require('request');
const vm = require('vm');

const context = require('./context');
const beatHost = process.env['BEAT_HOST'] || 'localhost';

const customCodePath = `http://${beatHost}/api/custom-codes/`;

var COMPILE_CUSTOM_CODE_FOOT = (
    'if (Backstage.executeHookName && Backstage.defines[Backstage.executeHookName]) {\n'+
        'Backstage.defines[Backstage.executeHookName](Backstage.request, Backstage.response);\n'+
    '} else {\n'+
        'Backstage.response.success();\n'+
    '}');


const config = {
    syncTimeout: process.env.CUSTOM_CODE_SYNC_TIMEOUT || 100,
    asyncTimeout: process.env.CUSTOM_CODE_ASYNC_TIMEOUT || 5000
};

class Registry {
    constructor() {
        this.codes = {};
    }

    populateNewCustomCodes(customCodes, fn) {
        var newCustomCodes = customCodes.filter((customCodeId) => {
            return this.codes[customCodeId] === undefined;
        });

        async.filter(newCustomCodes, fetchCustomCode.bind(this), function(err, customCodeIds) {
            fn(err);
        });
    }

    put(customCodeId, customCode) {
        this.codes[customCodeId] = customCode;
    }
}

function fetchCustomCode(customCodeId, callback) {
    request({url: `${customCodePath}${customCodeId}`, json: true}, (error, res, body) => {
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
