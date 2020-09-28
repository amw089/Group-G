// Global Variables
const NAMESPACE = "d9b2d63d-a233-4123-847a-76838bf2413a";
const DICTIONARY = 'dictionary.txt'
const NAMES = 'names.txt'

// UUIDV5 hashing funtion using crypto module to hash sha1
// Parsing namespace string, and converting hex and ascii to bytes for hashing
function uuidV5(namespace, name) {
  var hexNm = namespace.replace(/[{}\-]/g, '');
  var bytesNm = new Buffer(hexNm, 'hex');
  var bytesName = new Buffer(name, 'utf8');
  var hash =  require("crypto").createHash('sha1')
      .update(bytesNm).update(bytesName)
      .digest('hex');
  return hash.substr(0,8) + '-' +
      hash.substr(8,4) + '-' +
      (5 + hash.substr(13,3)) + '-' +
      ((parseInt(hash.substr(16, 2), 16) | 0x80) & 0xBF).toString(16) + hash.substr(18,2) + '-' +
      hash.substr(20,12);
}

// Timestamp hash. Changing the time zone to GMT-6 timezone with moment module.
function timeStampHash(timestamp) {
  var moment = require('moment-timezone');
  var timezone = "America/Regina";  
  return moment.tz(timestamp*1000, timezone).format();
}

// Hash tables for handling dictionaries
const pass_hashTable = new Map();
const uuidV5_hashTable = new Map();
var fs = require('fs');

// Populating the dictionary hashmap. Handling syncronization with promises.
// Create a dictionary with sha256 hash words for decoding   
function populate_Dictionary(fileName) {
  return new Promise((resolve, reject) => {
	fs.readFile(fileName, 'utf8', function (error, data) {
      	if (error) return reject(error);
      	
	var n = data.split("\n");
        for(var x in n){ 
		passhash = require("crypto").createHash("sha256").update(n[x]).digest("hex").toUpperCase();
		pass_hashTable.set(passhash,n[x])
	}

      resolve();
    })
  });
}
// Populate names hashmap for UUIDV5 hashing decoding. Handling syncronization with promises.
function populate_Names(fileName) {
  return new Promise((resolve, reject) => {
	fs.readFile(fileName, 'utf8', function (error, data) {
      	if (error) return reject(error);

      	var n = data.split("\n");
        for(var x in n){ 
		uuidhash = uuidV5(NAMESPACE,n[x]);
		uuidV5_hashTable.set(uuidhash,n[x]);
	}

      resolve();
    })
  });
}
// Reading the dumpfile and traversing through the dicionaries for decoding it.
function printout(fileName) {
  return new Promise((resolve, reject) => {
	fs.readFile(fileName, 'utf8', function (error, data) {
      	if (error) return reject(error);
      	
	var n = data.split("\n")
	n.length--

        for(var x in n){ 
		var stringArray = n[x].split(",");
		var linePrint = "";
	
		if(stringArray[0] == "username") {
			linePrint = "username,password,last_access";
		}
		else {
			linePrint = ""+ uuidV5_hashTable.get(stringArray[0]) +","+ pass_hashTable.get(stringArray[1]) +","+ timeStampHash(stringArray[2]);
		}
	
		console.log(linePrint);	
	}

      resolve();
    })
  });

}
// async function for controling the different threads in order for execution
async function run() {
  await populate_Dictionary(DICTIONARY); 
  await populate_Names(NAMES);
  await printout(process.argv[2])
}

// Start of program, the dumpifle is read as an argument in the command line
if (process.argv[2] == null) {
		console.log("Usage "+ process.argv[1] +"<dump_file.csv>");
		process.exit(1);
}
// Run async function to start the decoding process
run()
