'use strict';
let DotaBot = require('../lib/dotabot');
let Bot = require('slackbots');

let token = process.env.BOT_API_KEY;
let dbPath = process.env.BOT_DB_PATH;
let name = process.env.BOT_NAME;

let settings = {
    token,
    dbPath,
    name,
};
let dotabot = new DotaBot(settings);

dotabot.run();
