.PHONY: build tag commit

build:
	cd contentd/server; make rm; make build; cd - ; \
	cd dispatcherd; make rm; make build; cd - ; \
	cd inmem/server; make rm; make build; cd - ; \
	cd mo; make rm; make build; cd - ; \
	cd mt_manager; make rm; make build; cd - ;
	cd operator/pk/mobilink; make rm; make build; cd - ;
	cd pixels; make rm; make build; cd - ;
	cd qlistener; make rm; make build; cd - ;\

tag:
    if [ -z "$1" ]
      then
        echo "tag name required"
    fi; \
    if [ -z "$2" ]
          then
            echo "tag comment required"
    fi; \
	cd contentd;  git tag -a $1 -m $2; cd - ; \
	cd dispatcherd; git tag -a $1 -m $2;  cd - ; \
	cd inmem; git tag -a $1 -m $2; cd - ; \
	cd mo; git tag -a $1 -m $2; cd - ; \
	cd mt_manager; git tag -a $1 -m $2; cd - ;
	cd operator; git tag -a $1 -m $2;  cd - ;
	cd pixels; git tag -a $1 -m $2; cd - ;
	cd pixels; git tag -a $1 -m $2; cd - ;
	cd qlistener; git tag -a $1 -m $2; cd - ;\

commit:
    if [ -z "$1" ]
      then
        echo "comment required"
    fi
	cd contentd;  git commit -am $1; cd - ; \
	cd dispatcherd; git commit -am $1;  cd - ; \
	cd inmem; git commit -am $1; cd - ; \
	cd mo; git commit -am $1; cd - ; \
	cd mt_manager; git commit -am $1; cd - ;
	cd operator; git commit -am $1;  cd - ;
	cd pixels; git commit -am $1; cd - ;
	cd pixels; git commit -am $1; cd - ;
	cd qlistener; git commit -am $1; cd - ;\

d:
	ssh -i ~/.ssh/mobilink_centos.pem centos@52.66.132.93 mv /home/centos/linkit/dispatcherd-linux-amd64 /home/centos/linkit/prev/
	scp -i ~/.ssh/mobilink_centos.pem ~/linkit/dispatcherd-linux-amd64 centos@52.66.132.93:/home/centos/linkit/

deploy:
	ssh -i ~/.ssh/mobilink_centos.pem centos@52.66.132.93 mv /home/centos/linkit/*-linux-amd64 /home/centos/linkit/prev/
	scp -i ~/.ssh/mobilink_centos.pem *-linux-amd64 centos@52.66.132.93:/home/centos/linkit/