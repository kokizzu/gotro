var DS = DS || {};
DS.Currency = {
	IDR: 'Indonesian Rupiah',
	USD: 'United States Dollar',
	SGD: 'Singapore Dollar',
	JPY: 'Japanese Yen',
	THB: 'Thailand Baht',
	MYR: 'Malaysian Ringgit',
	PHP: 'Philippine Peso'
};
DS.WeekDay = {
	0: 'Sunday',
	1: 'Monday',
	2: 'Tuesday',
	3: 'Wednesday',
	4: 'Thursday',
	5: 'Friday',
	6: 'Saturday'
};
DS.PublicChatState = {
	0: 'waiting',
	1: 'ongoing',
	2: 'ended'
};
DS.CdnFrameRate = {
	'30fps': '30fps',
	'60fps': '60fps'
};
DS.CdnResolution = {
	'240p': '240p',
	'360p': '360p',
	'480p': '480p',
	'720p': '720p',
	'1080p': '1080p',
	'1440p': '1440p'
};
DS.CdnIngestionType = {
	'dash': 'Dynamic Adaptive Streaming over HTTP',
	'rtmp': 'Real Time Messaging Protocol'
};
DS.TodoStatus = {
	'': 'Unverified',
	wontfix: 'Invalid',
	someday: 'Queued',
	queued: 'Assigned',
	done: 'Completed',
	repeat: 'Reopened',
	merged: 'Duplicate',
	high: 'critical'
};
var Const_DMY = 'D MMM YYYY';
var Const_DMYHM = 'D MMM YYYY HH:mm';