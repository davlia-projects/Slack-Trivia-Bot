'use strict';
let DotaBot = require('../lib/dotabot');
let Bot = require('slackbots');

let token = process.env.BOT_API_KEY;
let name = process.env.BOT_NAME;

let settings = {
    token,
    name,
};
let dotabot = new DotaBot(settings);

dotabot.run();
