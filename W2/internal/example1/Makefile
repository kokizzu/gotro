######################
### local development

compose:
	docker-compose up

apiserver:
	killall example.exe || true
	air -c air.toml # localhost:9090

webclient:
	#sudo chown -R `whoami` /usr/lib/node_modules/
	cd svelte && npm install && npm run dev # http://localhost:5500
	# https://github.com/kokizz/svelte-mpa

reverseproxy:
	caddy fmt -overwrite Caddyfile
	caddy run -config Caddyfile

run:
	go run *.go cli

##########
### setup

setup-deps:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s && mv bin/air ~/go/bin
	go install github.com/alvaroloes/enumer@latest
	go install -v github.com/kokizzu/gomodifytags@v1.12.1
	go install github.com/tinylib/msgp@latest
	go install github.com/kokizzu/replacer@latest
	go install github.com/robertkrimen/godocdown/godocdown@latest
	go install github.com/akbarfa49/farify@latest

setup-webserver-ubuntu:
	sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
	curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo tee /etc/apt/trusted.gpg.d/caddy-stable.asc
	curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
	sudo apt update
	sudo apt install caddy
	sudo systemctl stop caddy
	sudo systemctl disable caddy
	sudo setcap 'cap_net_bind_service=+ep' `which caddy`
	curl -fsSL https://deb.nodesource.com/setup_17.x | sudo -E bash -
	sudo apt-get install -y nodejs
	sudo setcap CAP_NET_BIND_SERVICE=+eip `which caddy`

setup-webserver-mac:
	brew install caddy
	brew install nodejs

#setup-local-https:
#	sudo apt install -y libnss3-tools
#	go install filippo.io/mkcert@latest
#	mkcert -install
#	mkdir -p ./deploy && cd ./deploy && mkcert "localhost" # && mkcert "*.candlestick.com"
	#grep -qxF '127.0.0.1 localapi.candlestick.com' /etc/hosts || echo '127.0.0.1 localapi.candlestick.com' | sudo tee -a /etc/hosts
	# https://medium.com/@devahmedshendy/traditional-setup-run-local-development-over-https-using-caddy-964884e75232

###############
### generators

gen-route:
	cd domain ; rm -f *MSG.GEN.go 
	cd domain ; go test -bench=Benchmark_Generate_WebApiRoutes_CliArgs 0_generator_test.go
	cd domain ; cat *.go | grep '//go:generate ' | cut -d ' ' -f 2- | sh -x
	cd domain ; go test -bench=Benchmark_Generate_SvelteApiDocs 0_generator_test.go

gen-apidoc:
	cd domain ; go test -bench=Benchmark_Generate_SvelteApiDocs 0_generator_test.go
 
gen-orm:
	./gen-orm.sh

#############
### security

gokart-scan:
	go install github.com/praetorian-inc/gokart@latest
	gokart scan -v -d

gosec-scan:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec ./...

###########
### docker

docker-clean:
	docker stop $$(docker ps -q)
	docker rm $$(docker ps -a -f status=exited -q)
	docker system prune -f
	docker system prune -f --volumes

docker-clean-containers:
	docker rm -f $$(docker ps -a -q)
	docker rm $$(docker ps -a -f status=exited -q)

docker-clean-images:
	docker rmi -f $$(docker images -a -q)

docker-clean-volumes:
	docker volume rm $$(docker volume ls -q)

docker-clean-networks:
	docker network prune

###########
### builder

release:
	./release.sh

###########
### cleaner

clean_db:
	sudo rm -rf ./backup/tarantool-data
