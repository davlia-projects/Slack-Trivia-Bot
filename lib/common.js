module.exports = {
  rawString: str => {
    return str.toLowerCase().trim().replace(/[.,\/#!$%\^&\*;:{}=\-_`~()']/g,'');
  },

  hasOneOf: (str, list) => {
    for (let e in list) {
      if (str.includes(list[e])) {
        return true;
      }
    }
    return false;
  },

  shuffle: str => {
    let a = str.split(''),
    n = a.length;
    for(let i = n - 1; i > 0; i--) {
        let j = Math.floor(Math.random() * (i + 1));
        let tmp = a[i];
        a[i] = a[j];
        a[j] = tmp;
    }
    return a.join('');
  },
  replace: (str, replaceBy) => {
    let chars = str.split('');
    for (let key in replaceBy) {
      chars[key] = replaceBy[key];
    }
    return chars.join('');
  }
};
