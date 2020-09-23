var lineReader = require('readline').createInterface({
input: require('fs').createReadStream('dictionary.txt')
});

lineReader.on('line', function (line) {
        
hash = require("crypto").createHash("sha256").update(line).digest("hex").toUpperCase();

if(hash == process.argv[2]) 
	process.stdout.write(line);

});
