MAKO_EXPANDED_FILES = $(patsubst %.mako,%,$(wildcard *.mako))

install: $(MAKO_EXPANDED_FILES)
	chmod +x *.sh

clean:
	rm -f $(MAKO_EXPANDED_FILES)

include ../../rules.make
