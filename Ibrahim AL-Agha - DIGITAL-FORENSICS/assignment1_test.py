##########################################################################################################################################
# author:       Ibrahim AL-Agha                                                                                                          #
# date:         Friday, September 26, 2020                                                                                               #
# description:  This code is written for CSC-532-001: Digital Fornesics and Cyber Crime for Fall Quarter 2020.                           #
#               This program takes a database dump and uses SHA-2 hashing and UUID to convert the dump file into a plain text format.    #
##########################################################################################################################################

# Importing libraries.
import csv
import pandas as pd
import hashlib

# Instantiating global variables.
dump = pd.read_csv('database_dump.csv') # Reading in the dump file as a data frame called, "dump."

#   Breaking the columns of the "dump" dataframe into 3 deperate dataframes.
usernames = dump ['username'] # dataframe that contains all encrypted usernames.
passwords = dump ['password'] # dataframe that contains all encrypted passwords.
access_date = dump ['last_access'] # dataframe that contains all encrypted time-stamps.




