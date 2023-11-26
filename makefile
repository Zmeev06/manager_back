all:
	go build
	scp ./stupidauth scp://root@195.2.93.178//usr/local/bin/
