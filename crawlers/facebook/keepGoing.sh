#!/bin/bash
while :
do
	echo "Clear"
	#killall youtube_maxList 
	#killall chrome
	#sleep 7200
	NUM_FB=$(ps -ef|grep "facebookrss.py"|grep -v "grep" |awk 'BEGIN{}{}END{print $2}')
	#if [ $NUM_YOUTUBE -eq 0 ]
	if [[ -z "$NUM_FB" ]]
	then
		echo "EXECUTE"
		python3 ./facebookrss.py

	fi
															  sleep 2
done
