MODULE = $(shell go list -m)
TY = $(shell date +%Y)

# https://gist.github.com/prwhite/8168133
help: ## Show this help
	@ echo 'Usage: make <target>'
	@ echo
	@ echo 'Available targets:'
	@ grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

check-aoc-cookie:  ## ensures $AOC_SESSION_COOKIE env var is set
	@ if [ -f .env ]; then \
		export $(shell cat .env | xargs); \
	fi
	@ test $${AOC_SESSION_COOKIE?env var not set}

skeleton: ## make skeleton main(_test).go files, optional: $DAY and $YEAR
	@ if [[ -n $$DAY && -n $$YEAR ]]; then \
		go run scripts/cmd/skeleton/main.go -day $(DAY) -year $(YEAR) ; \
	elif [[ -n $$DAY ]]; then \
		go run scripts/cmd/skeleton/main.go -day $(DAY); \
	else \
		go run scripts/cmd/skeleton/main.go; \
	fi

input: check-aoc-cookie ## get input, requires $AOC_SESSION_COOKIE, optional: $DAY and $YEAR
	@ if [[ -n $$DAY && -n $$YEAR ]]; then \
		go run scripts/cmd/input/main.go -day $(DAY) -year $(YEAR) -cookie $(AOC_SESSION_COOKIE); \
	elif [[ -n $$DAY ]]; then \
		go run scripts/cmd/input/main.go -day $(DAY) -cookie $(AOC_SESSION_COOKIE); \
	else \
		go run scripts/cmd/input/main.go -cookie $(AOC_SESSION_COOKIE); \
	fi

prompt: check-aoc-cookie ## get prompt, requires $AOC_SESSION_COOKIE, optional: $DAY and $YEAR
	@ if [[ -n $$DAY && -n $$YEAR ]]; then \
		go run scripts/cmd/prompt/main.go -day $(DAY) -year $(YEAR) -cookie $(AOC_SESSION_COOKIE); \
	elif [[ -n $$DAY ]]; then \
		go run scripts/cmd/prompt/main.go -day $(DAY) -cookie $(AOC_SESSION_COOKIE); \
	else \
		go run scripts/cmd/prompt/main.go -cookie $(AOC_SESSION_COOKIE); \
	fi

all: skeleton input prompt ## run skeleton, input and prompt, optional: $DAY and $YEAR

run-%: ## run day $*, optional: $YEAR
	@ if [[ -n $$YEAR ]]; then \
		go run $(YEAR)/day$*/main.go -part 1; \
		go run $(YEAR)/day$*/main.go -part 2; \
	else \
		go run $(TY)/day$*/main.go -part 1; \
		go run $(TY)/day$*/main.go -part 2; \
	fi

check-%: ## run day $*, optional: $YEAR
	@ if [[ -n $$YEAR ]]; then \
		go test $(MODULE)/$(YEAR)/day$*; \
	else \
		go test $(MODULE)/$(TY)/day$*; \
	fi

.PHONY: help skeleton input prompt run-% check-% all
