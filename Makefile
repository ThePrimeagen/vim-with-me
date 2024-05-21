lua_fmt:
	echo "===> Formatting"
	stylua lua/ --config-path=.stylua.toml

lua_lint:
	echo "===> Linting"
	luacheck lua/ --globals vim

lua_test:
	echo "===> Testing"
	nvim --headless --noplugin -u scripts/tests/minimal.vim \
        -c "PlenaryBustedDirectory lua/vim-with-me {minimal_init = 'scripts/tests/minimal.vim'}"

lua_clean:
	echo "===> Cleaning"
	rm /tmp/lua_*

go-test:
	echo "===> Testing"
	go test ./pkg/v2/...
	go test ./examples/v2/...

go-relay-test:
	echo "===> Testing"
	go run ./pkg/v2/relay/cmd/test/main.go

go-fmt:
	echo "===> Format"
	go fmt github.com/theprimeagen/...

pr-ready: go-relay-test go-test go-fmt



