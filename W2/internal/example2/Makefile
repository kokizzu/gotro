
setup:
	go get -u -v github.com/kokizzu/gotro@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/fatih/gomodifytags@latest
	go install github.com/kokizzu/replacer@latest
	go install github.com/akbarfa49/farify@latest
	go install golang.org/x/tools/cmd/goimports@latest
	curl -fsSL https://bun.sh/install | bash

local-tarantool:
	docker exec -it example2-tarantool1-1 tarantoolctl connect example2T:example2PT@127.0.0.1:3301
	# box.space -- list all tables
	# box.execute [[ SELECT * FROM "users" LIMIT 1 ]]
	# \set language sql
	# \set delimiter ;

local-clickhouse:
	docker exec -it example2-clickhouse1-1 clickhouse-client --host 127.0.0.1 --port 9000 --user example2C --password example2PC
	# SHOW TABLES -- list all tables
	# SELECT * FROM "actionLogs" LIMIT 1;

modtidy:
	sudo chmod -R a+rwx tmpdb && go mod tidy

orm:
	# generate ORM
	./gen-orm.sh

views:
	# generate views and routes
	./gen-views.sh