#!/bin/bash

if [ -z "$1" ]; then
	echo "Usage $0 <dump_file.csv>"
	exit
fi

if [ ! -f "$1" ]; then
	echo "$1 doesn't exist!"
	exit
fi

file=$1

while IFS= read -r line
do
	readarray -d , -t strarr <<<"$line"
	
	for (( n=0; n < ${#strarr[*]}; n++ ))  
	do  
		case $n in

  			0)
    				if [ ${strarr[n]} == "username" ]; then 
					printf "username"
					printf ","
				else 
					node uuid.js ${strarr[n]}
					printf ","
				fi
   				;;

  			1)
    				if [ ${strarr[n]} == "password" ]; then 
					printf "password"
					printf ","
				else 
					node password.js ${strarr[n]}
					printf ","
				fi
   				;;


  			2)
    				if [ ${strarr[n]} == "last_access" ]; then 
					printf '%s\n' "last_access"
				else 
					node unixtimestamp.js ${strarr[n]}
					printf '\n'
				fi
   				;;

			
			*)
    				printf '%s\n' "unknown"
    				;;
		esac
	done 

done <"$file"
