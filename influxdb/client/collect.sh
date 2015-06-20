#!/bin/bash

subreddit=$1
loopDuration=$2

while true
do
	echo `date +"%x:%X"` `/usr/bin/curl -s http://www.reddit.com/r/${subreddit} | /go/bin/pup ".users-online .number text{}"` ${subreddit} | go run /go/tmp/docker-reddit-stats/influxdb/client/load2idb.go
	sleep ${loopDuration}
done
