# https://stackoverflow.com/questions/49102987/gnu-make-doesnt-match-dotted-filenames/49105330#49105330
.SUFFIXES:

# http://timmurphy.org/2015/09/27/how-to-get-a-makefile-directory-path/
PATH := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))/bin:$(PATH)

%: %.makotemplate
	expand_makotemplate -i $< -o $@

%: %.shtemplate
	expand_shtemplate < $< > $@
