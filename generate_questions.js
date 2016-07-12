let fs = require('fs');
let utils = require('./lib/common');

const questionsPath = './questions/questions.json';
const itemsPath = './scraper/items.json';
const skillsPath = './scraper/skills.json';
const itemNamesPath = './scraper/item_names.json';
const heroNamesPath = './scraper/heroes.json';
const quotesPath = './scraper/quotes.json';


let questionId = 0;

let generateQuoteQuestions = () => {
  let quotes = JSON.parse(fs.readFileSync(quotesPath));
  let quoteQuestions = [];
  let questionType = 'Hero Quote';
  quotes.forEach((elem, index) => {
    quoteQuestions.push({
      id: questionId++,
      type: questionType,
      prompt: elem.prompt,
      answer: elem.answer
    });
  });
  return quoteQuestions;
};

let generateItemQuestions = () => {
  let items = JSON.parse(fs.readFileSync(itemsPath));
  let questionType = 'Item';
  let itemQuestions = [];
  items.forEach((elem, index) => {
    // unscramble
    elem.name = elem.name.replace(/\s\(.*\)/g, '');
    let shuffled = utils.shuffle(elem.name).toLowerCase().replace(/\s/g, '').replace(/-/g, '');
    itemQuestions.push({
      id: questionId++,
      type: questionType,
      prompt: `Unscramble the item name: ${shuffled}`,
      answer: `${elem.name}`
    });

    // price
    if (elem.price === 'No Cost') {
      elem.price = 0;
    }
    itemQuestions.push({
      id: questionId++,
      type: questionType,
      prompt: `What is the price of ${elem.name}?`,
      answer: `${elem.price}`
    });
    // cooldown
    if (elem.cooldown) {
      itemQuestions.push({
        id: questionId++,
        type: questionType,
        prompt: `What is the cooldown of ${elem.name}?`,
        answer: `${elem.cooldown}`
      });
    }
    // manacost
    if (elem.manacost) {
      itemQuestions.push({
        id: questionId++,
        type: questionType,
        prompt: `What is the mana cost of ${elem.name}?`,
        answer: `${elem.manacost}`
      });
    }
    // description
    if (elem.description) {
      itemQuestions.push({
        id: questionId++,
        type: questionType,
        prompt: `Item name by description: ${elem.description}`,
        answer: `${elem.name}`
      });
    }
    // item stats
    elem.stats.forEach((stat) => {
      itemQuestions.push({
        id: questionId++,
        type: questionType,
        prompt: `How much ${stat.label.toLowerCase()} does ${elem.name} give?`,
        answer: `${stat.value}`
      });
    });
  });
  return itemQuestions;
};

const quantifier = {
  'duration': 'long',
  'radius': 'large',
  'range': 'many units',
  'width': 'many units',
  'distance': 'far',
  'cast': 'long',
  'area': 'large',
};

const blacklist = [
  'cast point',
];

let generateHeroQuestions = () => {
  let heroes = JSON.parse(fs.readFileSync(skillsPath));
  let questionType = 'Hero';
  let heroQuestions = [];
  heroes.forEach((hero) => {
    hero.skills.forEach((skill) => {
      // unscramble hero name
      let shuffled = utils.shuffle(hero.hero).toLowerCase().replace(/\s/g, '');
      heroQuestions.push({
        id: questionId++,
        type: questionType,
        prompt: `Unscrable the hero name: ${shuffled}`,
        answer: `${hero.hero}`
      });
      // hero stats
      heroQuestions.push({
        id: questionId++,
        type: questionType,
        prompt: `What is the movement speed of ${hero.hero}`,
        answer: `${hero.stats["Movement speed"]}`
      });
      // skill description
      heroQuestions.push({
        id: questionId++,
        type: questionType,
        prompt: `${hero.hero} Skill Name: ${skill.description}`,
        answer: `${skill.name}`
      });

      // skill stats
      skill.stats.forEach((stat) => {
        let level = 0;
        while (true) {
          if (!stat.values || level === stat.values.length || utils.hasOneOf(stat.label.toLowerCase(), blacklist)) {
            break;
          }
          let levelText = stat.values.length > 1 ? `level ${level + 1} ` : '';
          // let quantified = "";
          // for (let q in quantifier) {
          //   if (stat.label.toLowerCase().includes(q)) {
          //     quantified = quantifier[q];
          //     break;
          //   }
          // }
          let isNegative = stat.values[0] === '-' ? ' (The answer is negative)': '';
          heroQuestions.push({
            id: questionId++,
            type: questionType,
            prompt: `What is the ${stat.label.toLowerCase()} for ${hero.hero}'s ${levelText}${skill.name}${isNegative}?`,
            answer: `${stat.values[level]}`
          });
          level++;
        }
      });
    });
  });
  return heroQuestions;
};

let postProcess = (arr, func) => {
  for (let q in arr) {
    arr[q].answer = arr[q].answer.replace(/\.0%?$/g, '');
  }
  return arr;
};

let compileQuestions = () => {
  return postProcess(
    [].concat(
      generateQuoteQuestions(),
      generateHeroQuestions(),
      generateItemQuestions()
    )
  );
};

fs.writeFileSync(questionsPath, JSON.stringify(compileQuestions(), null, 4));
