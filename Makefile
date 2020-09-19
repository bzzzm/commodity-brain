version ?= $(shell date +%s)
app_name ?= commodity
docker_tag ?= dev
build_path ?= build
path ?= /home/mihai/work/commodity/brain
remote_ssh ?= mihai@commodity.local
remote_path ?= /home/mihai/brain
go ?= /usr/local/go/bin/go

.PHONY: rpi
rpi: upload-rpi run-rpi

.PHONY: dev
dev: run-dev


.PHONY: upload-rpi
upload-rpi:
	ssh $(remote_ssh) mkdir -p $(remote_path)
	rsync --exclude-from='rsync_exclude.txt' -a . $(remote_ssh):$(remote_path)

.PHONY: build-rpi
build-rpi:
	scp $(build_path)/$(app_name) $(remote_ssh):$(remote_path)
	env GOOS=linux GOARCH=arm GOARM=5 $(go) build -o $(build_path)/cmd/brain/main.go
	chmod +x $(build_path)/$(app_name)

.PHONY: run-rpi
run-rpi:
	ssh $(remote_ssh) $(go) run $(path)/cmd/brain/main.go

dep-rpi:
	ssh $(remote_ssh) cd $(path) && $(dep) ensure

.PHONY: run-dev
run-dev:
	$(go) run $(path)/cmd/brain/main.go -database-path bolt.db

.PHONY: clean
clean:
	rm -rf build/*