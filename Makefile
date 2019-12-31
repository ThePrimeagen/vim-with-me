all: twitchirc

twitchirc: src/twitchirc.c
	$(CC) -o bin/twitchirc src/twitchirc.c src/system-commands/sys-commands.c
	cp src/.config bin/.config

run: twitchirc
	./bin/twitchirc

clean:
	rm ./bin/*

test: twitchirc
	valgrind bin/twitchirc
