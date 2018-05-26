SERVICE := qpclient
VERSION := 0.0.2


build:
	@echo "building..." && \
	go install && \
	echo "build done"

run: build 
	$(SERVICE) -v -k -proxy http://localhost:3443 -l :18080



	


