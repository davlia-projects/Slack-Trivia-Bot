let fs = require('fs');
let request = require('request');
let cheerio = require('cheerio');

const base_url = "http://dota2.gamepedia.com/"

let content = JSON.read(fs.readFileSync('./heroes.json'));
let heroURLs = [];
content.heroes.forEach((e,i) => {
  heroURLs.push(e.localized_name.replace(" ", "_"));
});

console.log(`Hero URLs ${heroURLs}`);
for
