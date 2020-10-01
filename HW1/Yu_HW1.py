import csv
import uuid
from hashlib import sha256
from datetime import datetime
from sys import stdin

NAMESPACE = uuid.UUID('d9b2d63d-a233-4123-847a-76838bf2413a')

names = []
namesDict = {}
passwords = []
passwordDict = {}
times = []

f = open("dictionary.txt", "r")
password_list = f.read().rstrip().split("\n")
f.close()

f = open("names.txt", "r")
names_list = f.read().rstrip().split("\n")
f.close()

# creating a dictionary/hash table in the form of uuid -> name:
for name in names_list:
    namesDict[str(uuid.uuid5(NAMESPACE, name))] = name

# creating a dictionary/hash table in the form of sha256 -> password:
for pw in password_list:
    passwordDict[sha256(pw.encode('utf-8')).hexdigest().upper()] = pw

# ---------------- read database_dump.csv file ----------------
forSkip = 0
with open('database_dump.csv', mode='r') as csvfile:
    entries = csv.reader(csvfile, delimiter=',')
    for row in entries:
        if forSkip == 0:    # skipping first row: "username,password,last_access"
            forSkip += 1
        else:
            # names
            names.append(namesDict[row[0]].replace('\n', '').replace(' ', ''))    
            # passwords
            passwords.append(passwordDict[row[1]].replace('\n', '').replace(' ', ''))
            # times
            time = datetime.fromtimestamp(int(row[2])-3600).strftime('%Y-%m-%dT%H:%M:%S') 
            times.append(time.replace('\n', '').replace(' ', ''))

# ---------------- write decrypted.csv file ----------------
with open('database_dump_answer.csv', 'w',newline='') as decrypted_file:
    entries = csv.writer(decrypted_file, delimiter=',')
    entries.writerow(["username", "password", "last_access"])
    for i in range(len(names)):
        entries.writerow([names[i], passwords[i], times[i]])




        
