# File: ibrahim.py
# Name: This program was developed for CSC-543-001: Digital Forensics and Cyber Crime for Fall Quarter 2020 at Louisiana Tech University.
#       By: Ibrahim AL-Agha
# Date: 09/29/2020
# Desc: This program uses a database dump of a "users" table in a database in the "database_dump.csv" file. 
#       This program deconstructs the hashes, uuids, and timestamps to arrive at this plaintext version. 
#       This program utilizes a given "dictionary.txt" file, and just in case the users don't exist in the dictionary, a list of the names has been provided in "names.txt".
#       A hardcoded UUID namespace has been provided: d9b2d63d-a233-4123-847a-76838bf2413a

#   IMPORTING LIBRARIES
import csv
import pandas as pd
import hashlib as hl
import uuid as ud
from time import strftime, gmtime

#   INSTANTIATING GLOBAL VARIABLES, DATA FRAMES, AND ARRAYS.
GLOBAL_UUID = ud.UUID("d9b2d63d-a233-4123-847a-76838bf2413a") # Initiating the global uuid namespace.

INPUT_DUMP = pd.read_csv('database_dump.csv') # Reading in the dump file as a data frame called, "INPUT_DUMP."
INPUT_DICTIONARY = open('dictionary.txt', 'r').read().split('\n') # Reading in the dictionary file.
INPUT_NAMES = open('names.txt', 'r').read().split('\n') # Reading in the names file.

USERNAMES = INPUT_DUMP['username'].tolist() # list that contains all encrypted USERNAMES.
PASSWORDS = INPUT_DUMP['password'].tolist() # list that contains all encrypted PASSWORDS.
ACCESS_DATE = INPUT_DUMP['last_access'].tolist() # list that contains all encrypted time-stamps.

# Function that uses UUID5 function that translates a plain-text object and compares it to an encrypted object.
# 
# Inputs:
#   @param plain_text_array - Array of plain-text string objects.
#   @param encrypted_text_array - Array of encrypted string objects.
# 
# Output:     List of string username objects that have been translated.
def translate_uuid (plain_text_array, encrypted_text_array):
    translated_output = []
    for uuid_name in encrypted_text_array:
        for name in plain_text_array:
            if (str(uuid_name) == str((ud.uuid5(GLOBAL_UUID, name)))):
                translated_output.append(name)
                break
    return (translated_output)

# Function that uses SHA256 hash function that translates a plain-text object and compares it to an encrypted object.
# 
# Inputs:
#   @param plain_text_array - Array of plain-text string objects.
#   @param encrypted_text_array - Array of encrypted string objects.
# 
# Output:     List of string password objects that have been translated into plain-text.
def translate_sha2 (plain_text_array, encrypted_text_array):
    translated_output = []
    for encrypted_password in encrypted_text_array:
        for plain_password in plain_text_array:
            if (str(encrypted_password) == str((hl.sha256(plain_password.encode()).hexdigest().upper()))):
                translated_output.append(plain_password)
                break
    return translated_output

# Function that converts milliseconds since POSIX time format to GMT format.
# 
# Inputs:
#   @param time_array - Array of milliseconds since POSIX time stamps.
# 
# Output:     List of converted time stamps.
def convert_time (time_array):
    converted_output = []
    for time_stamp in time_array:
        converted_output.append(strftime('%Y-%m-%dT%H:%M:%S-0600', gmtime(int(time_stamp)-21600)))
    return converted_output

#   FUNCTION CALLS
translated_user = translate_uuid(INPUT_NAMES, USERNAMES) # Calling the "translate_user()" function.
translated_passwords = translate_sha2(INPUT_DICTIONARY, PASSWORDS) # Calling the "translate_sha2()" function.
converted_time_stamps = convert_time(ACCESS_DATE) # Calling the "convert_time()" function.

#   WRITING TO *.csv FILE.
final_dataframe = zip(translated_user, translated_passwords, converted_time_stamps)
with open("ibrahim-plain-text_data_dump.csv", "w", newline = '') as f:
    writer = csv.writer(f)
    writer.writerow(['username', 'password', 'last_access'])
    for row in final_dataframe:
        writer.writerow(row)