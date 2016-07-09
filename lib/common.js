module.exports = {

  rawString: str => {
    return str.toLowerCase().trim().replace(/[.,\/#!$%\^&\*;:{}=\-_`~()]/g,"")
  },
};
