SERVICE := qpserver
VERSION := 0.0.2

build:
	@echo "building..." && \
	go install && \
	echo "build done"

run: build 
	$(SERVICE) -v -l :3443 -cert '../cert/cert.pem' -key '../cert/domain.key'



branch := `git rev-parse --abbrev-ref HEAD`

deploy:
	@echo "deploy remote server..." && \
	git push origin HEAD && \
	ssh vultr-la "source /etc/profile && cd /go/src/github.com/liudanking/quic-proxy/qpserver &&  git pull && git checkout $(branch) && \
	git branch && make && supervisorctl restart $(SERVICE)" && \
	echo "deploy done"
	



	


