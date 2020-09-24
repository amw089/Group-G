//hashtables
const hash = (key, size) => {
  let hashedKey = 0
  for (let i = 0; i < key.length; i++) {
    hashedKey += key.charCodeAt(i)
  }
  return hashedKey % size
}

class HashTable {
  constructor() {
    this.size = 101
    this.tables = Array(this.size) 

    let i = 0;
    while (i < this.tables.length) {
      this.tables[i] = new Map()
      i++;
    }
  }

  insert(key, value) {
    let idx = hash(key, this.size) 
    this.tables[idx].set(key, value)
  }

  search(key) {
    let idx = hash(key, this.size)
    return this.tables[idx].get(key)
  }

  size() {
    return this.tables.size
  }

}

//uuidv5
var crypto = require('crypto');

function uuidV5(namespace, name) {
  var hexNm = namespace.replace(/[{}\-]/g, '');
  var bytesNm = new Buffer(hexNm, 'hex');
  var bytesName = new Buffer(name, 'utf8');
  var hashi = crypto.createHash('sha1')
      .update(bytesNm).update(bytesName)
      .digest('hex');
  return hash5(hashi);
}

function hash5(hasho) {
  return hasho.substr(0,8) + '-' +
      hasho.substr(8,4) + '-' +
      (5 + hasho.substr(13,3)) + '-' +
      ((parseInt(hasho.substr(16, 2), 16) | 0x80) & 0xBF).toString(16) + hasho.substr(18,2) + '-' +
      hasho.substr(20,12);
}

var namespaceString = "d9b2d63d-a233-4123-847a-76838bf2413a";
const uuidV5_hashTable = new HashTable();

var lineReader = require('readline').createInterface({
input: require('fs').createReadStream('names.txt')
});


lineReader.on('line', function (line) {

uuidhash = uuidV5(namespaceString,line);
uuidV5_hashTable.insert(uuidhash,line);

});;

const pass_hashTable = new HashTable();
var lineReader = require('readline').createInterface({
input: require('fs').createReadStream('dictionary.txt')
});

lineReader.on('line', function (line) {

passhash = require("crypto").createHash("sha256").update(line).digest("hex").toUpperCase();
pass_hashTable.insert(passhash,line);

});

//timestamp hash
function timeStampHash(timestamp) {
  var moment = require('moment-timezone');
  var timezone = "America/Regina";
  
  return moment.tz(timestamp*1000, timezone).format();
}

// Start of program after hastables creation

if (process.argv[2] == null) {
	console.log("Usage "+ process.argv[1] +"<dump_file.csv>");
	process.exit(1);
}

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
		linePrint = "" + uuidV5_hashTable.search(stringArray[0]) +","+ pass_hashTable.search(stringArray[1]) +","+ timeStampHash(stringArray[2]);
	}
	console.log(linePrint);	
});
