Team Pink/Dark Red

Each of the folders in this repository holds an assignment or challenge given in our CSC443 - Forensics class.

Assignment/Challenge 1 - Folder "HW1" - Decoding a csv dump file

The program to grade is "ricardo_for_grading.js". It was coded using node js. To run it enter "node ricardo_for_grading.js database_dump.csv". The dump file is read in the arguments of the commandline. If you wish to redirect it to a csv file, type: "node ricardo_for_grading.js databas_dump.csv > filename.csv" To check if the output file is identical to the "database_dump_answer.csv", you can type in the linux commandline "diff -s filename.csv database_dump_answer.csv"
Assignment/Challenge 2 - Folder "HW2" - Analyzing mbr, gpt, and FAT.

Homework 2

The program to grade is "hw2-1.go". It was coded using golang. To compile it, enter "go build hw2-1.go". To run it, enter "go run hw2-1.go [mode] [filename]". The Image file is read in the arguments of the commandline. There are three modes: -mbr mbr analysis -gpt gpt analysis -fat FAT analysis When running a FAT analysis, you will be prompt to enter an offset to specify the partition to analyze. To learn of partitions, run an mbr or gpt analysis. Global variables can be specified to filter results by extension, change sector size, and debug mode.

Homework 3:

fileCarver.py runs with either -f or -d 
-d : searches data section and finds 44 files (1 fragmented only first sector) files found denoted with file_at_-d.png
-f: searches the FAT and finds 36 files (1 corrupted) files found denoted file_at_-f.png

Challenge 3:

Finds the files iterativly, I believe conflicting ending signaturs results in slightly corrupted text documents
