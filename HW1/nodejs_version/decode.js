// Splits a given file into smaller subfiles by line number
var infileName = 'dictionary.txt';
var fileCount = 1;
var count = 0;
var fs = require('fs');
var outStream;
var outfileName = infileName + '.' + fileCount;
newWriteStream();
var inStream = fs.createReadStream(infileName);

var lineReader = require('readline').createInterface({
    input: inStream
});

function newWriteStream(){
    outfileName = infileName + '.' + fileCount;
    outStream = fs.createWriteStream(outfileName);
    count = 0;
}

lineReader.on('line', function(line) {
    count++;
    outStream.write(line + '\n');
    if (count >= 7600) {
        fileCount++;
        outStream.end();
        newWriteStream();
    }
});

lineReader.on('close', function() {
    if (count > 0) {
    }
    inStream.close();
    outStream.end();
    console.log('Done');
});


//uuidv5
var crypto = require('crypto');

function uuidV5(namespace, name) {
  var hexNm = namespace.replace(/[{}\-]/g, '');
  var bytesNm = new Buffer(hexNm, 'hex');
  var bytesName = new Buffer(name, 'utf8');
  var hash = crypto.createHash('sha1')
      .update(bytesNm).update(bytesName)
      .digest('hex');
  return hash.substr(0,8) + '-' +
      hash.substr(8,4) + '-' +
      (5 + hash.substr(13,3)) + '-' +
      ((parseInt(hash.substr(16, 2), 16) | 0x80) & 0xBF).toString(16) + hash.substr(18,2) + '-' +
      hash.substr(20,12);
}

const pass_hashTable = new Map();
const uuidV5_hashTable = new Map();

function populate() {
	var cntr = 1
	while(cntr<fileCount) {
		var lineReader = require('readline').createInterface({
		input: require('fs').createReadStream('dictionary.txt.'+cntr)
		});

		lineReader.on('line', function (line) {
		passhash = require("crypto").createHash("sha256").update(line).digest("hex").toUpperCase();
		pass_hashTable.set(passhash,line)
		});

		cntr++	
	}
	
	var namespaceString = "d9b2d63d-a233-4123-847a-76838bf2413a";

	var lineReader = require('readline').createInterface({
	input: require('fs').createReadStream('names.txt')
	});

	lineReader.on('line', function (line) {

	uuidhash = uuidV5(namespaceString,line);
	uuidV5_hashTable.set(uuidhash,line);
	});;
}

//timestamp hash
function timeStampHash(timestamp) {
  var moment = require('moment-timezone');
  var timezone = "America/Regina";
  
  return moment.tz(timestamp*1000, timezone).format();
}

function io_dump() {
	var file = process.argv[2];

	var lineReader = require('readline').createInterface({
	input: require('fs').createReadStream(file)
	});

	lineReader.on('line', function (line) {
	var stringArray = line.split(",");
	var linePrint = "";

	if(stringArray[0] == "username") {
		linePrint = "username,password,last_access";
	}
	else {
		linePrint = "" + uuidV5_hashTable.get(stringArray[0]) +","+ pass_hashTable.get(stringArray[1]) +","+ timeStampHash(stringArray[2]);
	}
	
		console.log(linePrint);	
	});
}


// Start of program after hastables creation

setTimeout(populate,400)

if (process.argv[2] == null) {
	console.log("Usage "+ process.argv[1] +"<dump_file.csv>");
	process.exit(1);
}
setTimeout(io_dump,500)

