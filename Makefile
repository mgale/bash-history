.PHONY: test

cmds = $(shell find cmd/* -maxdepth 0 -type d)
pkgs = $(shell find pkg/* -type d)
protos = $(shell find proto/* -maxdepth 0 -type d)
ROOT_DIR = $(shell pwd)

all: test linter release-test

clean:
	rm -rf dist
	docker-compose down
	@echo [INFO] Clean cmds
	@for cmd in $(cmds); do \
			cd $${cmd}; \
			echo [INFO] Clean $${cmd}; \
			make clean; \
			cd ${ROOT_DIR}; \
	done

linter:
	@echo [INFO] Run linter
	docker run --rm -v ${ROOT_DIR}:/code -w /code golangci/golangci-lint:v1.43.0 golangci-lint run -v

test:
	@echo [INFO] Test cmds
	@for cmd in $(cmds); do \
			cd $${cmd}; \
			if [ `find . -maxdepth 1 -type f -name '*_test.go' | wc -l` -gt 0 ]; then \
					echo [INFO] Test $${cmd}; \
					make test || exit 1; \
			fi; \
			cd ${ROOT_DIR}; \
	done
	@echo [INFO] Test pkgs
	@for pkg in $(pkgs); do \
			cd $${pkg}; \
			if [ `find . -maxdepth 1 -type f -name '*_test.go' | wc -l` -gt 0 ]; then \
					echo [INFO] Test $${pkg}; \
					go test -v -cover ./... || exit 1; \
			fi; \
			cd ${ROOT_DIR}; \
	done

build:
	@echo [INFO] Build cmds
	@for cmd in $(cmds); do \
			cd $${cmd}; \
			if [ `find . -maxdepth 1 -type f -name 'main.go' | wc -l` -gt 0 -a `find skip-make 2>/dev/null | wc -l` -eq 0 ]; then \
					echo [INFO] Build $${cmd}; \
					make build || exit 1; \
			fi; \
			cd ${ROOT_DIR}; \
	done

release:
	@echo [INFO] Build cmds
	@for cmd in $(cmds); do \
			cd $${cmd}; \
			if [ `find . -maxdepth 1 -type f -name 'main.go' | wc -l` -gt 0 -a `find skip-make 2>/dev/null | wc -l` -eq 0 ]; then \
					echo [INFO] Build $${cmd}; \
					make release; \
			fi; \
			cd ${ROOT_DIR}; \
	done

release-test:
	@echo [INFO] Build cmds
	@for cmd in $(cmds); do \
			cd $${cmd}; \
			if [ `find . -maxdepth 1 -type f -name 'main.go' | wc -l` -gt 0 -a `find skip-make 2>/dev/null | wc -l` -eq 0 ]; then \
					echo [INFO] Build $${cmd}; \
					make release-test || exit 1; \
			fi; \
			cd ${ROOT_DIR}; \
	done
