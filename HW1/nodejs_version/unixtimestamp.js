var moment = require('moment-timezone');

var timezone = "America/Regina";

if (process.argv[2] != null) {
	var timestamp = process.argv[2];
	process.stdout.write(moment.tz(timestamp*1000, timezone).format());
}

