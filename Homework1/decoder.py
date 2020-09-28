# Matt Johnson
# decoder.py
import csv
from uuid import uuid5, UUID
from hashlib import sha256 as hasher
from time import *

usernames = []
passwords = []
time_stamps = []

NAMESPACE = UUID('d9b2d63d-a233-4123-847a-76838bf2413a')
with open('database_dump.csv', newline='') as data_dump:
    reader = csv.reader(data_dump, dialect='excel')
    for row in reader:
        hashed_username = row[0].replace('\n', '')
        hashed_password = row[1].replace('\n', '')
        time_stamp = row[2].replace('\n', '')

        # decode and store each time stamp
        if(time_stamp != "last_access"):
                time_stamps.append(strftime('%Y-%m-%dT%H:%M:%S', gmtime(int(time_stamp)-21600)))
        
        # hashes each word in dictionary then compares 
        # with the hashed usernames and passwords
        # matches are stored in their corresponding list
        dictionary = open('dictionary.txt', 'r')
        for word in dictionary:
            word = word.replace('\n', '').replace(' ', '')

            # hash the dictionary for the username
            hashed_wordU = uuid5(NAMESPACE, word)

            # hash the dictionary for the password
            hashed_wordP = hasher(bytes(word, 'utf-8')).hexdigest().upper()

            if (hashed_username != "username" and str(hashed_wordU) == str(hashed_username)):
                usernames.append(word)
            if (hashed_password != "password" and str(hashed_wordP) == str(hashed_password)):
                passwords.append(word)
             

# display usernames, passwords, and time
for num in range(len(usernames)-1):
    print("|" +"-" * 55 + "|")
    print("|", usernames[num], " " * (12-len(usernames[num])) ,"|", passwords[num], " " * (14-len(passwords[num])), "|", time_stamps[num], "|")
print("|" +"-" * 55 + "|")