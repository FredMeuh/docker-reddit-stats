# docker-reddit-stats
Docker + pup + InfluxDB + Go code + Grafana to graph reddit user stats

Dummy project to play with a few technologies and finaly get a nice graph of the number of users logged in and viewing a subreddit.
Goals:
- learn a bit of:
  - docker
  - go (and pup)
  - InfluxDB
  - Grafana
  - git
- and share this experience if this can be useful for someone else.

Disclaimer: I'm neither a developper nor a devops and I don't intend to be one! Please bear with me. This is an exercice and not production code.

*** LEARNING BY TRYING THINGS - v1 ***

(0) Pre-requisites

You should know a bit of docker, a bit of HTML/CSS, a bit of bash scripting and a few other things. Also, on your Linux host, you should have docker installed and ready for prime time. If you follow the script below, everything should be fine and you will end up with an env. up and running. Some lines will have to be drawn between the dots ... experiment!

(0.1) What is docker?

https://www.docker.com/

(0.2) Where to test all this?

On your environment (workstation)? On a linux server in your lab? On a CoreOS image on AWS EC2 (some "free" resources on AWS account for 1 year)? 

http://aws.amazon.com/ec2/

(1) Intro

(1.1) What is reddit?

https://www.reddit.com/about/

(1.2) What's the end goal of this project?

to collect the number of users looking at a subreddit, and plot this metric on a graph over time. Example:

https://www.reddit.com/r/Guitar/
241 Guitarists (looking at the page)

(2) Capturing the data

(2.1) curl is your (old) friend

```
# curl -s http://www.reddit.com/r/Guitar
```

and see the HTML code here:

<p class="users-online " title="logged-in users viewing this subreddit in the past 15 minutes">
	<span class="number">192</span>&#32;<span class="word">users here now</span>
</p>

Now, Q: how to I fetch a small block of info from this HTML page?
A: using pup, a nice tool written in Go.

(2.2) pup is your (new) friend

https://github.com/ericchiang/pup

Using CSS selectors you can grab a specific part of the HTML page. Example:

```
# /go/bin/pup ".users-online .number text{}"
```

(2.2.1) But ... wait ... pup is not installed on my machine.

This is where docker is now also a friend. Let's install and run Go in a container that you'll be able to dispose / customize when you want.

```
# docker pull golang
```

Now you have Go on your host, within a docker image ready to run. You can play with it. Let's go straight to the point and reuse an image with pup already installed.

```
# docker pull fredmeuh/pup:latest

# docker run --rm fredmeuh/pup /bin/sh -c '/usr/bin/curl -s http://www.reddit.com/r/Guitar | /go/bin/pup ".users-online .number text{}"'
```

Cool ... I want more! What's next?

(3) Retaining the data 

(3.1) InfluxDB may be my new friend

Many different DB solutions out there. I don't think there is one solution which fits all needs. Here we want to retain time series data (<timestamp>, <value>) and InfluxDB is a new solution in this space. Let's test it!

(0.9 is the latest stable release, the rest of the exercice won't work on pre-0.9)

https://influxdb.com/docs/v0.9/introduction/overview.html

Those guys were probably the first to publish a docker image of the 0.9 stable version of InfluxDB:
https://registry.hub.docker.com/u/savoirfairelinux/influxdb/

```
# docker pull savoirfairelinux/influxdb

# docker run -d --name influxdb -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 savoirfairelinux/influxdb
```

connect to the UI (port 8083, root:root) and create a DB called "reddit" with a user (name:rw / pwd:rwrw)

run a request: 
	select value from nb_users;

no data ...

(3.2) Let's push some data into this DB.

(3.2.1) InfluxDB client library in Go

Why Go? Because InfluxDB is written in Go. Because I don't know Go. Because it's time to learn.

Nice intro page to the Go client library for InfluxDB:

https://github.com/influxdb/influxdb/tree/master/client

In my case I want to read the data from STDIN, so I can pipe the data into the client.

(3.2.2) Installing the client lib (reusing the pup image)

Maybe it's better to run pup in a container and the client InfluxDB in another. I'm not sure it's the case and I prefer to put all the "data collection" layer in one single place.

run a request:
	select value from nb_users;

nice ...

(3.2.3) Some (hugly) glue

Putting all those pieces together with some glue (shell script). Everything (looping, fetching the web page, parsing the web page, waiting for the next interval) could have been done in the Go program but for now I let this to you as an exercice ;-)

docker run --rm -v /home/core:/go/tmp -e INFLUXDB_HOST=`docker inspect --format '{{ .NetworkSettings.IPAddress }}' influxdb` -it fredmeuh/influxdb-go-lib /go/tmp/docker-reddit-stats/influxdb/client/collect.sh Guitar 300

(4) Presenting the data

(4.1) Why Grafana?

Because it's the default dashboard UI for InfluxDB. Because it's close to Kibana that I want to investigate later.

http://docs.grafana.org/datasources/influxdb/

(4.2) Installing ... in a container I guess?

Yep, another container.

https://registry.hub.docker.com/u/grafana/grafana/

```
docker pull grafana/grafana:latest
```

(you need version 2.0+ to interface with influxdb 0.9)

(5) What do we have accomplish?

curl to get a HTML web page from a web site
pup to parse the HTML page and grab some data
a Go program to load the data into InfluxDB
Grafana to display the data
all those services put into docker images and running as containers to simplify deployment

(6) Next steps

- deploy a fleet of containers using CoreOS fleet to collect multiple stats 
- put all this data into a InfluxDB cluster 
- use a load balancer in front of Grafana

Why all this? Just for the fun of it.




