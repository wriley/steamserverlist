steamserverlist: steamserverlist.go
	@go build

install: steamserverlist
	@install -m 0755 steamserverlist /usr/local/bin/

clean:
	@/bin/rm -f steamserverlist
