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

pr-ready: lua_fmt lua_lint lua_test



