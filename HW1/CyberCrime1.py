###################################
#Jillian Stalder
#Professor John Spurgeon
#9/27/2020 
#######################################

import hashlib
import datetime
import uuid
import sys
import csv
#import uuid

namesFile = open('C:\\Users\\Jill Stalder\\Desktop\\Cyber\\names.txt')
names=[]
for namesLine in namesFile:    
    names.append(namesLine.strip())
namesFile.close()

dictionaryFile = open('C:\\Users\\Jill Stalder\\Desktop\\Cyber\\dictionary.txt')
dictionary=[]
for dictionaryLine in dictionaryFile:    
    dictionary.append(dictionaryLine.strip())
dictionaryFile.close()

dumpFile = 'C:\\Users\\Jill Stalder\\Desktop\\Cyber\\dump.csv'
dump=[]
with open(dumpFile) as dumpLine:
    dumpIn = dumpLine.read()
dumpReader = csv.reader(dumpIn.splitlines(), delimiter=',', quotechar='|')
for row in dumpReader:
    if (row[0] != "username"):        
        dump.append(row);
    else:
        fields = row

myuuid = uuid.UUID('d9b2d63d-a233-4123-847a-76838bf2413a')
for namesLine in names:
    nameUUID = uuid.uuid5(myuuid, namesLine)    
    for dumpLine in dump:
        if (dumpLine[0] == str(nameUUID)):
            dumpLine[0] = namesLine

for dictionaryLine in dictionary:
    modSHA = hashlib.sha256(dictionaryLine.encode('utf-8'))
    for dumpLine in dump:
        if (dumpLine[1] == modSHA.hexdigest().upper()):
            dumpLine[1] = dictionaryLine

for dumpLine in dump:
    dumpLine[2] = datetime.datetime.fromtimestamp(int(dumpLine[2])-3600).strftime('%Y-%m-%dT%H:%M:%S-0600')


#for dictionaryLine in dictionary:
#    print(dictionaryLine)

for dumpLine in dump:
    print(dumpLine)

dumpFileOut = 'C:\\Users\\Jill Stalder\\Desktop\\Cyber\\dump_answer.csv'
with open(dumpFileOut, 'w', newline ='') as csvfile:
    writer = csv.DictWriter(csvfile, fieldnames = fields)
    writer.writeheader()
    for dumpLine in dump:        
        writer.writerow({'username' : dumpLine[0].strip(), 'password' : dumpLine[1], 'last_access' : dumpLine[2].strip()})



