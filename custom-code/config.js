const config = {
    syncTimeout: process.env.CUSTOM_CODE_SYNC_TIMEOUT || 100,
    asyncTimeout: process.env.CUSTOM_CODE_ASYNC_TIMEOUT || 5000,
    beatHost: process.env['BEAT_HOST'] || 'localhost'
};

module.exports = config;
