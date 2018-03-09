# https://stackoverflow.com/questions/49102987/gnu-make-doesnt-match-dotted-filenames/49105330#49105330
.SUFFIXES:

# http://timmurphy.org/2015/09/27/how-to-get-a-makefile-directory-path/
PATH := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))/bin:$(PATH)

%: %.mako
	expand_mako -i $< -o $@

%: %.shtemplate
	expand_shtemplate < $< > $@

# if a variables_override.mako file is missing, copy the corresponding default variables.mako file
# TODO copy variables.mako from $@'s path rather than .
variables_override.mako:
	if [ ! -f "$@" ]; then cp variables.mako "$@"; fi
