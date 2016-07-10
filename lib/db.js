let fs = require('fs');
class Database {
	constructor(userdb, questiondb) {
		this.userdb = userdb;
		this.questiondb = questiondb;
		this._connectDb();
	}

	_connectDb() {
		if (!fs.existsSync(this.userdb)) {
			console.log("Users database doesn't exist. Creating...");
	        fs.closeSync(fs.openSync(this.userdb, 'w'));
	    }
    if (!fs.existsSync(this.questiondb)) {
    	console.log("Questions database doesn't exist. Creating...");
        fs.closeSync(fs.openSync(this.questiondb, 'w'));
    }
	}

	getAllQuestions() {
		let results = JSON.parse(fs.readFileSync(this.questiondb));
		return results;
	}

	getAllUsers() {
		let results = JSON.parse(fs.readFileSync(this.userdb));
		return results;
	}

	writeAllUsers(users) {
		let userJson = JSON.stringify(users);
		fs.writeFile(this.userdb, users);
	}
}
module.exports = Database;
