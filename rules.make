# https://stackoverflow.com/questions/49102987/gnu-make-doesnt-match-dotted-filenames/49105330#49105330
.SUFFIXES:

# http://timmurphy.org/2015/09/27/how-to-get-a-makefile-directory-path/
PATH := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))/bin:$(PATH)

%: %.makotemplate
	expand_makotemplate -i $< -o $@

%: %.shtemplate
	expand_shtemplate < $< > $@

# if a variables_override.makotemplate file is missing, copy the corresponding default variables.makotemplate file
# TODO copy variables.makotemplate from $@'s path rather than .
variables_override.makotemplate:
	if [ ! -f "$@" ]; then cp variables.makotemplate "$@"; fi
