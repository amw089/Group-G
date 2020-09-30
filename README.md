# Group-G
Team Pink/Dark Red

Each of the folders in this repository holds an assignment or challenge given in our CSC443 - Forensics class.

* Assignment/Challenge 1 - Folder "HW1" - Decoding a csv dump file 
    - The program to grade is "ricardo_for_grading.js". It was coded using node js.
      To run it enter "node ricardo_for_grading.js database_dump.csv". The dump file is read in the arguments of the commandline. 
      If you wish to redirect it to a csv file, type: "node ricardo_for_grading.js databas_dump.csv > filename.csv"
      To check if the output file is identical to the "database_dump_answer.csv", you can type in the linux commandline "diff -s filename.csv database_dump_answer.csv"
