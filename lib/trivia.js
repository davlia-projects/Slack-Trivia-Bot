let Database = require('./db');

const CORRECT = 1;
const INCORRECT = 0;
const USERS_DB_PATH = './data/users.json';
const QUESTIONS_DB_PATH = './data/questions.json';

class Trivia {
	constructor() {
		this.maxEarnable = 5;
		this.maxHint = 3;
		this.db = new Database(USERS_DB_PATH, QUESTIONS_DB_PATH);
		this.date = new Date();
		this.reset();
		this.questions = this.db.getAllQuestions();
		this.participants = this.db.getAllUsers();
		// console.log(this.questions);
	}

	reset() {
		this.roundStarted = false;
		this.question = null;
		this.answer = null;
		this.hint = null;
		this.hintCount = 0;
		this.startTime = 0;
		for (let user in this.participants) {
			this.participants[user].guesses = 0;
		}
		this.questions = this.db.getAllQuestions();
	}

	startNewRound() {
		this.reset();
		this.roundStarted = true;
		this.startTime = this.date.getTime();
		let res = this._sampleNewQuestion();
		this.question = res.prompt;
		this.answer = res.answer;
	}

	makeGuess(name, guess) {
		let p = this.participants;
		if(!p[name]) {
			p[name] = {
				points: 0,
				guesses: 0,
				streak: 0,
			};
		}
		if (guess.toLowerCase().trim() === this.answer.toLowerCase().trim()) {
			let earned = Math.max(1, this.maxEarnable - p[name].guesses);
			this.clearOtherStreak(name);
			p[name].points += earned;
			p[name].streak++;
			this.roundStarted = false;
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

	getQuestion() {
		return this.question;
	}

	nextHint() {
		if (this.hintCount == this.maxHint) {
			return 0;
		}
		switch (this.hintCount) {
			case 0: {
				this.hint = this.answer.replace(/\w/ig,'*');
				break;
			}
			default: {
				let newHint = [];
				let answer = this.answer.split(' ');
				let tokens = this.hint.split(' ');
				tokens.forEach((i,x) => {
					let chars = i.split('');
					for (let i = 0; i < this.hintCount; i++) {
						let rand = Math.floor(Math.random() * chars.length);
						chars[rand] = answer[x][rand];
					}
					newHint.push(chars.join(''))
				})
				this.hint = newHint.join(' ');
			}
		}
		this.hintCount = Math.min(this.hintCount + 1, this.maxHint);
		return this.hint;
	}

	getStats(name) {
		return this.participants[name];
	}

	getPoints(name) {
		return this.participants[name].points;
	}

	getAnswer() {
		return this.answer;
	}

	endRound() {
		this.reset();
	}

	clearOtherStreak(name) {
		for (let n in this.participants) {
			if (n !== name) {
				this.participants[n].streak = 0;
			}
		}
	}

	elapsed() {
		return this.date.getTime() - this.startTime;
	}

	_sampleNewQuestion() {
		let length = Object.keys(this.questions).length;
		return this.questions[Math.floor(Math.random() * length)];
	}
	_fetchUserStats() {
		return this.db.getAllUsers();
	}
}

module.exports = Trivia;
