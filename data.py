# Katie Hay
# Assignment 1 - Database Dump
# Python 3.8
# To run in cmd: data.py
# This will create a CSV file named data.csv containing the decoded data
# Make sure database_dump.csv, dictionary.txt, and names.txt are in same location as data.py

# Libraries
import csv, uuid, hashlib, time
from datetime import datetime, date, time, timedelta
#from datetime import timezone

from collections import defaultdict

# Constants
namespace = uuid.UUID('{d9b2d63d-a233-4123-847a-76838bf2413a}')

# Function to find first name
def username(namespace, name):
    name = uuid.UUID(name)
    # Compare UUID created with namespace and each line in dictionary.txt with UUID in DB dump, find a match
    with open("dictionary.txt", "r") as dictionary:
        #print ("Checking dictionary.txt")
        for line in dictionary:
            temp = line.strip()
            temp_uuid = uuid.uuid5(namespace, temp)
            if (temp_uuid == name):
                return temp
                
    # If not in dictionary.txt, check names.txt
    with open("names.txt", "r") as names:
        #print ("Checking names.txt")
        for line in names:
            temp2 = line.strip()
            temp_uuid2 = uuid.uuid5(namespace, temp2)
            if (temp_uuid2 == name):
                return temp2


# Function to find password
def password(passw):
    passw = passw.lower()
    # Compare hash created on each line in dictionary.txt with password hash in DB dump, find a match
    with open("dictionary.txt", "r") as dictionary:
        for line in dictionary:
            temp = line.strip()
            temp_hash = hashlib.sha256(temp.encode()).hexdigest()
            if (temp_hash == passw):
                return temp


# Function to find last access
def last_access(timestamp):
    timestamp = int(timestamp)
    date = datetime.fromtimestamp(timestamp)
    date = date - timedelta(hours=1,minutes=0)
    date = date.strftime("%Y-%m-%dT%H:%M:%S-0600")
    return date
            

# Creating a data.csv and writing all decoded data to that new file
with open('data.csv', 'w', newline='') as file:
    writer = csv.writer(file)
    writer.writerow(["username", "password", "last_access"])
    with open('database_dump.csv', newline='') as csvfile:
        spamreader = csv.DictReader(csvfile)
        for row in spamreader:
            writer.writerow([username(namespace,(row['username'])), password(row['password']), last_access(row['last_access'])])























