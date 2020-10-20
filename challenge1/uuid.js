var crypto = require('crypto');

function uuidV5(namespace, name) {
  var hexNm = namespace.replace(/[{}\-]/g, '');
  var bytesNm = new Buffer(hexNm, 'hex');
  var bytesName = new Buffer(name, 'utf8');
  var hash = crypto.createHash('sha1')
      .update(bytesNm).update(bytesName)
      .digest('hex');
  return _hash(5,hash);
}

function _hash(version, hash) {
  return hash.substr(0,8) + '-' +
      hash.substr(8,4) + '-' +
      (version + hash.substr(13,3)) + '-' +
      ((parseInt(hash.substr(16, 2), 16) | 0x80) & 0xBF).toString(16) + hash.substr(18,2) + '-' +
      hash.substr(20,12);
}

var lineReader = require('readline').createInterface({
input: require('fs').createReadStream('names.txt')
});

var namespaceString = "d9b2d63d-a233-4123-847a-76838bf2413a";
if(process.argv[3] != null){
	namespace = process.argv[3];
}

lineReader.on('line', function (line) {
        
uuidhash = uuidV5(namespaceString,line);
	if(process.argv[2] == uuidhash) {
		process.stdout.write(line);
	}
});;
