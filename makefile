run:
	GOOS=linux GO111MODULE=on go run -tags sqlite -ldflags "-X 'main.buildTime=`date '+%Y-%m-%dT%H:%M:%S'`' -X 'main.buildVersion=`git log --pretty=oneline|head -n1|sed 's/\"//g'|sed "s/'//g"`' -X 'main.buildGoVersion=`go version`' -X 'main.buildBy=`logname`@`hostname -f`'" .

build:
	CGO_ENABLED=0 GOOS=linux GO111MODULE=on CGO=false go build -ldflags "-X 'main.buildTime=`date '+%Y-%m-%dT%H:%M:%S'`' -X 'main.buildVersion=`git log --pretty=oneline|head -n1|sed 's/\"//g'|sed "s/'//g"`' -X 'main.buildGoVersion=`go version`' -X 'main.buildBy=`logname`@`hostname -f`'"

deployabroad: build
	export remote=pluto;\
    scp ./pluto $$remote:/opt/pluto/api/pluto_new &&\
    ssh $$remote 'chmod +x /opt/pluto/api/pluto_new &&\
    cp /opt/pluto/api/pluto /opt/pluto/api/archive/pluto.`date +'%Y%m%d.%H%M%S'` &&\
    supervisorctl stop pluto &&\
    mv /opt/pluto/api/pluto_new /opt/pluto/api/pluto &&\
    supervisorctl start pluto &&\
    supervisorctl status pluto&&\
    supervisorctl tail -f pluto'


migrateabroad:
	export remote=pluto;\
	cd ./cmd/migrate/ ;\
	go build ;\
	scp ./migrate $$remote:/opt/pluto/api/migrate &&\
	ssh $$remote 'cd /opt/pluto/api && ./migrate && rm ./migrate'
