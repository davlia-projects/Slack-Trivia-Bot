let util = require('util');
let path = require('path');
let Bot = require('slackbots');
let Trivia = require('./trivia');

const CORRECT = 1;
const INCORRECT = 0;
const HINT_DELAY = 3000;
const QUESTION_DELAY = 2000;
const MAX_INTERVALS = 4;

class DotaBot extends Bot {
	constructor(settings) {
		super(settings);
		this.settings = settings;
		this.settings.name = this.settings.name || 'dota-2-trivia-bot';
		this.dbPath = settings.dbPath || path.resolve(process.cwd(), 'data', 'dotaquestions.db');
		this.user = null;
		this.db = null;
		this.trivia = new Trivia();
		this.timeIntervals = 0;
		this.continuous = false;
		this.missedCount = 0;
	}

	run() {
		this.on('start', this._onStart);
		this.on('message', this._onMessage);
	}

	_onStart() {
		console.log("Starting up");
		this._loadBotUser();
	}

	_loadBotUser() {
		this.user = this.users.filter(user => {
			return user.name === this.name;
		})[0];
	}

	_onMessage(message) {
		switch(message.type) {
			case 'hello':
				console.log("Ready");
				break;
			case 'channel_created':
				this._onChannelCreated(message);
				break;
			case 'message':
				this._onChannelMessage(message);
				break;
			case 'presence_change':
				break;
			case 'reconnect_url':
				break;
			case 'user_typing':
				break;
			default:
				console.log("Event not handled: %s", message.type);
				break;
		}
	}

	_onChannelMessage(message) {
		let say = t => this.postMessageToChannel(channel.name, t, {as_user: true});
		let trivia = this.trivia;
		let channel = this._getChannelById(message.channel);
		let answer = trivia.getAnswer();
		let user = this._getUserById(message.user);
		let name;
		try {
			name = user.name;
		} catch (err) {
			console.log(message.subtype);
		}
		if (message.subtype === 'message_changed') {
			return;
		}
		if (message.user === this.user.id) {
			if (message.text.includes('hint')) {
				this.hintTS = message.ts;
			}
			return;
		}
		if (message.text === 'test') {
			console.log(message);
			say(":kappa:");
		}

		switch(message.text.substr(0,2)) {
			case '!q': {
				this._startRound(channel, message);
				break;
			}
			case '!s': {
				let stats = trivia.getStats(name);
				if (!stats) {
					say(`${name} has not played yet! Press !q to play.`);
					break;
				}
				say(`${name} has ${stats.points} points with a best streak of ${stats.bestStreak}.`);
				break;
			}
			case '!h': {
				if (!trivia.roundStarted) {
					break;
				}
				let hint = trivia.hint;
				if (hint) {
					say(`hint: \`${hint}\``);
				}
				break;
			}
			case '!c': {
				this.continuous = true;
				say('Continuous mode on');
				break;
			}
			case '!o': {
				this.continuous = false;
				say('Continuous mode off');
				break;
			}
			case '!?': {
				say('`!q` to get a question or the current one. `!h` for the current hint. `!s` for your stats.');
				break;
			}
			default: {
				if (trivia.roundStarted) {
					let result = trivia.makeGuess(name, message.text);
					if (result.status === CORRECT) {
						let cleared = this.trivia.clearOtherStreak(name);
						let streak = "";
						if (result.streak > 1) {
							streak = ` win streak: ${result.streak}`;
							if (result.streak > 3) {
								streak += ' :fire: \\:cheerpogchamp:/ :fire:';
							}
						}
						let ended = "";
						if (cleared && cleared.streak > 1) {
							ended = ` ${cleared.name}'s ${cleared.streak} win streak was ended :parrot:`;
						}
						say(`Yay, ${name}! ${answer} is correct. +${result.earned} [total points: ${result.points}]${streak}${ended}`);
						this._cleanup();
						setTimeout((self => {return () => this.continuous ? this._startRound(channel, message) : 0;})(this), QUESTION_DELAY);
					}
				}
			}
		}
	}

	_startRound(channel, message) {
		let say = t => this.postMessageToChannel(channel.name, t, {as_user: true});
		let trivia = this.trivia;
		if (!trivia.roundStarted) {
			trivia.startNewRound();
			this.control = setInterval((self => {return () => self._roundControl(channel, message);})(this), HINT_DELAY);
		}
		let question = trivia.getQuestion();
		say(question);
	}

	_roundControl(channel, message) {
		let say = t => this.postMessageToChannel(channel.name, t, {as_user: true});
		let hint = this.trivia.nextHint();
		console.log(hint);
		if (this.timeIntervals >= MAX_INTERVALS) {
			say(`Time is up! :lul: The answer is \`${this.trivia.getAnswer()}\`.`);
			if (this.missedCount++ > 3) {
				this.continuous = false;
			}
			this.trivia.clearOtherStreak();
			this._cleanup();
			setTimeout((self => {return () => this.continuous ? this._startRound(channel, message) : 0;})(this), QUESTION_DELAY);
		} else if (hint.stars) {
			let hintText = `hint: \`${hint.stars}\``;
			if (!this.hintTS) {
				say(hintText);
			} else {
				this.updateChat({
					channel: channel.id,
					ts: this.hintTS,
					text: hintText,
				});
			}
			if (hint.final) {
				this.timeIntervals = MAX_INTERVALS;
			} else {
				this.timeIntervals++;
			}
		}
	}

	_cleanup() {
		this.hintTS = 0;
		this.timeIntervals = 0;
		this.trivia.reset();
		clearInterval(this.control);
	}

	_onChannelCreated(message) {
		this.channels.push(message.channel);
	}

	_getChannelById(channelId) {
	    return this.channels.filter(item => {
	        return item.id === channelId;
	    })[0];
	}

	_getUserById(userId) {
		return this.users.filter(user => {
			return user.id === userId;
		})[0];
	}
}

module.exports = DotaBot;
