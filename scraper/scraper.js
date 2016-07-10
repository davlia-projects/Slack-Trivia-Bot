  let fs = require('fs');
let request = require('request');
let cheerio = require('cheerio');
let utils = require('../lib/common');

const BASE_URL = 'http://www.dotabuff.com';
const HERO_QUERY = 'heroes';
const ITEM_QUERY = 'items';

let scrapeHeroMetadata = () => {
  let content = JSON.parse(fs.readFileSync('./heroes.json'));
  let heroURLs = [];
  content.heroes.forEach((elem, i) => {
    let hero = utils.rawString(elem.localized_name).replace(/\s/g, "-");
    heroURLs.push({
      url: `${BASE_URL}/${HERO_QUERY}/${hero}/abilities`,
      hero: elem.localized_name,
    });
  });

  let heroes = [];
  let failures = 0;
  // heroURLs = heroURLs.slice(0,1);
  heroURLs.forEach((e,i,a) => {
    let opts = {
      url: e.url,
      headers: {
        'User-Agent': 'Mozilla/5.0',
      },
    };
    request(opts, (err, resp, html) => {
      if (err) {
        console.log(err);
        return;
      }
      let $ = cheerio.load(html);
      let heroSkills = [];
      if ($('.col-8').text() === '') {
        console.log(`Grabbing data failed for ${e.hero}`);
        failures++;
      }
      // grabs skills
      $('.col-8 section').each((i, elem) => {
        let skillName = $(elem).find('header').contents().filter((i, elem) => {
          return elem.type === 'text';
        }).text();

        let skillStats = [];
        let $stats = $(elem).find('article');
        $stats.find('.stats .stat').each((i, elem) => {
          let label = $(elem).find('.label').text().replace(":", "");
          let values = [];
          $(elem).find('.values .value').each((i, elem) => {
            values.push($(elem).text());
          });
          skillStats.push({label, values});
        });
        let description = $stats.find('.description').text().replace(/\s+/g, " ").replace(/(\w)\.(\w)/g, "$1. $2");
        let manacost = [];
        $stats.find('.cooldown .number').each((i, elem) => {
          manacost.push($(elem).text());
        });
        let cooldown = [];
        $stats.find('.manacost .number').each((i, elem) => {
          cooldown.push($(elem).text());
        });
        skillStats.push({manacost});
        skillStats.push({cooldown});
        heroSkills.push({
          name: skillName,
          description: description,
          stats: skillStats,
        });
      });
      let heroStats = {};
      $('.hero_attributes .other tr').each((i, elem) => {
        let children = $(elem).children();
        let name = $(children[0]).text();
        let value = $(children[1]).text();
        heroStats[name] = value;
      });
      heroes.push({
        hero: e.hero,
        stats: heroStats,
        skills: heroSkills,
      });
      if (heroes.length === heroURLs.length - failures) {
        let out = JSON.stringify(heroes, null, 4);
        fs.writeFileSync("skills.json", out);
      }
    });
  });
};

let scrapeItems = () => {
  const totalItems = 167;
  let items = [];
  let opts = {
    url: `${BASE_URL}/${ITEM_QUERY}`,
    headers: {
      'User-Agent': 'Mozilla/5.0',
    },
  };
  request(opts, (err, resp, html) => {
    if (err) {
      console.log(err);
      return;
    }
    let $ = cheerio.load(html);
    console.log($('.cell-xlarge').html());
    $('.cell-xlarge').each((i, elem) => {
      items.push($(elem).text());
    });
    console.log(items);
    if (items.length === totalItems) {
      fs.writeFileSync("item_names.json", JSON.stringify(items, null, 4));
    }
  });
};

let scrapeItemMetadata = () => {
  let content = JSON.parse(fs.readFileSync('item_names.json'));
  let itemURLs = [];
  let items = [];
  content.forEach((elem, i) => {
    itemURLs.push({
      url: `${BASE_URL}/${ITEM_QUERY}/${utils.rawString(elem).replace(/\s/g, "-")}`,
      item: elem,
    });
  });

  itemURLs.forEach((e, i) => {
    let opts = {
      url: e.url,
      headers: {
        'User-Agent': 'Mozilla/5.0',
      },
    };
    request(opts, (err, resp, html) => {
      if (err) {
        console.log(e.item, err);
        return;
      }
      let $ = cheerio.load(html);
      let $data = $('.item-tooltip');
      let price = $data.find('.price .value').text();
      let description = $data.find('.description').text();
      let cooldown = $data.find('.cooldown').text();
      let manaCost = $data.find('.manacost').text();
      // build recipe
      let itemStats = [];
      $data.find('.attribute').each((i, elem) => {
        itemStats.push({
          label: $(elem).find('.label').text(),
          value: $(elem).find('.value').text()
        });
      });
      let itemEffect = [];
      $data.find('.effect').each((i, elem) => {
        itemEffect.push({
          label: $(elem).find('.label').text(),
          value: $(elem).find('.value').text(),
        });
      });
      items.push({
        name: e.item,
        price: price,
        description: description,
        cooldown: cooldown,
        manacost: manaCost,
        stats: itemStats,
      });
      if (items.length === itemURLs.length) {
        let out = JSON.stringify(items, null, 4);
        fs.writeFileSync("items.json", out);
      }
    });
  });
};

scrapeHeroMetadata();
// scrapeItems();
// scrapeItemMetadata();
