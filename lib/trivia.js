let Database = require('./db');
let utils = require('./common');

const CORRECT = 1;
const INCORRECT = 0;
const USERS_DB_PATH = './data/users.json';
const QUESTIONS_DB_PATH = './questions/questions.json';

class Trivia {
	constructor() {
		this.maxEarnable = 5;
		this.db = new Database(USERS_DB_PATH, QUESTIONS_DB_PATH);
		this.reset();
		this.questions = this._fetchQuestions();
		this.participants = this._fetchUserStats();
		// console.log(this.questions);
	}

	reset() {
		this.roundStarted = false;
		this.question = null;
		this.answer = null;
		this.hint = {
			stars: undefined,
			revealed: {},
			final: false,
		};
		for (let user in this.participants) {
			this.participants[user].guesses = 0;
		}
		this.questions = this._fetchQuestions();
	}

	startNewRound() {
		this.reset();
		this.roundStarted = true;
		let res = this._sampleNewQuestion();
		this.question = res.prompt;
		this.answer = res.answer;
	}

	validateUser(name) {
		let p = this.participants;
		if(!(`${name}` in p)) {
			p[name] = {
				points: 0,
				guesses: 0,
				streak: 0,
				bestStreak: 0,
			};
		}
	}

	makeGuess(name, guess) {
		let p = this.participants;
		this.validateUser(name);
		if (utils.rawString(guess) === utils.rawString(this.answer)) {
			let earned = Math.max(1, this.maxEarnable - p[name].guesses);
			// this.clearOtherStreak(name);
			p[name].points += earned;
			p[name].streak++;
			p[name].bestStreak = Math.max(p[name].bestStreak, p[name].streak);
			this.roundStarted = false;
			this._writeUsers();
			return {
				status: CORRECT,
				earned: earned,
				points: p[name].points,
				streak: p[name].streak,
			};
		} else {
			p[name].guesses++;
			return {
				status: INCORRECT,
			};
		}
	}

	nextHint() {
		let hint = this.hint;
		if (!hint.stars) {
			hint.stars = this.answer.replace(/(\w|\d)/ig,'*');
			return hint;
		} else if (Object.keys(hint.revealed).length >= Math.ceil(this.answer.replace(/\W/g, '').length / 2)) {
			hint.final = true;
			return hint;
		} else {
			let tokens = hint.stars.split(" ");
			let len = 0, index;
			for (let i = 0; i < tokens.length; i++) {
				do {
					index = Math.floor(Math.random() * tokens[i].length) + len;
				} while (this.answer[index].match(/\W/));
				hint.revealed[index] = this.answer[index];
				len += tokens[i].length + 1;
			}
			hint.stars = utils.replace(hint.stars, hint.revealed);
			return hint;
		}
	}
	getQuestion() {
		return this.question;
	}

	getAnswer() {
		return this.answer;
	}

	getStats(name) {
		return this.participants[name];
	}

	clearOtherStreak(name) {
		let p = this.participants;
		for (let n in p) {
			if (n !== name) {
				if (name && p[n].streak > 0) {
					let ret = {name: n, streak: p[n].streak};
					p[n].streak = 0;
					return ret;
				}
				p[n].streak = 0;
			}
		}
	}

	_sampleNewQuestion() {
		let length = Object.keys(this.questions).length;
		return this.questions[Math.floor(Math.random() * length)];
	}
	_fetchUserStats() {
		return this.db.getAllUsers();
	}

	_fetchQuestions() {
		let questions;
		while (!questions) {
			try {
				questions = this.db.getAllQuestions();
			} catch (err) {
				console.log("Error retrieving questions. Trying again...");
			}
		}
		return this.db.getAllQuestions();
	}
	_writeUsers() {
		this.db.writeAllUsers(JSON.stringify(this.participants));
	}
}

module.exports = Trivia;
