const vm = require('vm');


var sandbox = {
    Backstage: null,
    Buffer: null,
    console: null,
    exports: null,
    module: null,
    setTimeout: null,
    require: null,
    relativeRequire: null,
    define: null
};


function create(customCodeId, request, backstageOptions, customCodeCallback) {
    // only for compatibily with unit tests =)
    var exports = {};

    sandbox.Backstage = {
        defines: {},
        define: function(hookName, func) {
            sandbox.Backstage.defines[hookName] = func;
        },
        modules: {},
        request: request
    };

    for (var key in backstageOptions) {
        sandbox.Backstage[key] = backstageOptions[key];
    }

    sandbox.console = console;
    sandbox.exports = exports;
    sandbox.module = {
        exports: exports
    };
    sandbox.Buffer = Buffer;
    sandbox.setTimeout = setTimeout;

    // for compatibility
    sandbox.define = sandbox.Backstage.define;

    return vm.createContext(sandbox);
}

exports.sandbox = sandbox;
exports.create = create;
