#!/bin/bash

# -------- ci tool for hedzr/store -------------------------------------
#
# After main module of hedzr/store published, do these
# works step by step:
#
#  0. release the main branch to new version
#     for exmaple:
#     git tag v1.0.3 && git push --tags origin master
#     and waiting for its done at remote.
#  1. publish providers
#     $0 publish-providers
#     This command will iterate all providers, upgrade go.mod,
#     and commit each of them, tag them (with git tag).
#  2. push providers to remote and release them
#     $0 push
#  3. publish codecs
#     $0 publish-codecs
#     This command will iterate all codecs, upgrade go.mod,
#     and commit each of them, tag them (with git tag).
#  4. push codecs to remote and release them
#     $0 push
#  5. sync all (optional)
#     $0 update
#     After it done, examples and tests submodules should be upgraded.
#  6. commit and push examples and tests
#     $0 commit-codecs-tests
#     $0 commit-tests
#     $0 commit-examples
#  7. push all of them
#
# All folks.
#

BUILD_DIR=./build
INSTALL_TMP_DIR=./bin/install
OSN=store
VER=0.1.0

[ -f .version.cmake ] && {
	VER=$(echo $(grep -oE ' \d+\.\d+\.\d+' .version.cmake))
	# grep -oE ' [0-9]+.[0-9]+.[0-9]+' .version.cmake
	# echo "VERSION = $VER"
}

[ -f doc.go ] && {
	VER=$(echo v$(grep -iE 'Version[ ]*=.*' doc.go | grep -oE '\d+\.\d+\.\d+'))
	echo "VERSION = $VER"
}

#

alias cmake="cmake -Werror=dev --warn-uninitialized"

#

#

#

build-push() {
	git push origin master && git push --tags origin master
}

build-drop-tags-providers() {
	local sm d f
	for sm in providers codecs; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [[ ! "$d" == */test ]]; then
				if [ -d "$d" ]; then
					local pre="${d/.\//}"
					git tag --delete "$pre/$VER"
				fi
			fi
		done
	done
}

build-tag-codecs() { build-publish-children codecs; }
build-tag-providers() { build-publish-children providers; }
build-pub() { build-publish "$@"; }
build-publish() { build-publish-children "$@"; }
build-publish-codecs() { build-publish-children codecs; }
build-publish-providers() { build-publish-children providers; }
build-publish-tests() { build-publish-children tests; }
build-publish-examples() { build-publish-children examples; }
build-publish-children() {
	local sm d f
	local which="$1"
	for sm in "$which"; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [[ ! "$d" == */test ]]; then
				if [ -d "$d" ]; then
					do-update-dep "$d"
					commit-submodule "$d"
				fi
			fi
		done
	done
}

build-commit-codecs-tests() { build-commit-children-test-only codecs; }
build-commit-tests() { build-commit-children tests; }
build-commit-examples() { build-commit-children examples; }

build-commit-children-test-only() {
	local sm d f
	local which="$1"
	for sm in "$which"; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [[ "$d" == */test ]]; then
				if [ -d "$d" ]; then
					do-update-dep "$d"
					commit-dir "$d"
				fi
			fi
		done
	done
}

build-commit-children() {
	local sm d f
	local which="$1"
	for sm in "$which"; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [ -d "$d" ]; then
				do-update-dep "$d"
				commit-dir "$d"
			fi
		done
	done
}

commit-dir() {
	local d="$1"
	pushd "$d" >/dev/null
	if is_git_dirty; then
		git add .
		git commit -m "updated $d"
	fi
	popd >/dev/null
}

commit-submodule() {
	local d="$1"
	local pre="${d/.\//}"
	bump-and-tag "$pre" "$VER" "$d"
	# 	# drop-tag "$pre/$VER"
	# 	git tag --delete "$pre/$VER"
}

bump-and-tag() {
	local tag="$1/$2"
	local d="$3"
	pushd "$d" >/dev/null

	if [[ "$d" == tests || "$d" == examples ]]; then
		if is_git_dirty; then
			git add .
			git commit -m "upgraded $d"
		fi
	elif is_git_dirty; then
		git add .
		git commit -m "bump to $tag" && git tag "$tag"
	else
		git tag "$tag"
	fi

	popd >/dev/null
}

drop-tag() {
	local tagn=${1:-}
	if [ "$tagn" == "" ]; then
		echo "Usage: $0 tag-name [--push]"
		exit 1
	fi

	git push --delete origin $tagn
	git tag --delete $tagn

	case "$2" in
	--push | -p | --reset | -r)
		git tag $tagn
		git push origin
		;;
	*)
		:
		;;
	esac
}

#
# upgrade go.mod dependencies
#

build-update() {
	for sm in "${1:-.}"; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [ -d "$d" ]; then
				do-update-dep "$d"
			fi
		done
	done
	echo
}

build-update-cmdr() {
	for sm in ../cmdr ../tool.rd ../tool.zag; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [ -d "$d" ]; then
				do-update-dep "$d"
			fi
		done
	done
	echo
}

build-update-deps() { build-update-dep "$@"; }
build-upgrade-dep() { build-update-dep "$@"; }
build-upgrade-deps() { build-update-dep "$@"; }
build-update-dep() {
	local sm d f
	do-update-dep "."
	echo
	for sm in codecs providers; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [ -d "$d" ]; then
				do-update-dep "$d"
			fi
		done
	done
	echo
	for sm in tests examples; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [ -d "$d" ]; then
				do-update-dep "$d"
			fi
		done
	done
}

# https://gosamples.dev/update-all-packages/
do-update-dep() {
	local d="$1"
	pushd "$d" >/dev/null
	echo
	echo "==== go mod tidy, dir='$d' =========="
	go get -v -t -u && go mod tidy
	popd >/dev/null
}

#
#
#
# sync project files with iCloud backup point --------------------------
#
#

# iCloud="$HOME/Library/Mobile Documents/com~apple~CloudDocs"

#

#

#

##
##

cmd_exists() { command -v $1 >/dev/null; } # it detects any builtin or external commands, aliases, and any functions
fn_exists() { LC_ALL=C type $1 2>/dev/null | grep -qE '(shell function)|(a function)'; }
fn_builtin_exists() { LC_ALL=C type $1 2>/dev/null | grep -q 'shell builtin'; }
fn_aliased_exists() { LC_ALL=C type $1 2>/dev/null | grep -qE '(alias for)|(aliased to)'; }

is_git_clean() { git diff-index --quiet $* HEAD -- 2>/dev/null; }
is_git_dirty() { is_git_clean && return -1 || return 0; }

headline() { printf "\e[0;1m$@\e[0m:\n"; }
headline_begin() { printf "\e[0;1m"; } # for more color, see: shttps://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
headline_end() { printf "\e[0m:\n"; }  # https://misc.flogisoft.com/bash/tip_colors_and_formatting
debug() { in_debug && printf "\e[0;38;2;133;133;133m$@\e[0m\n" || :; }
debug_begin() { printf "\e[0;38;2;133;133;133m"; }
debug_end() { printf "\e[0m\n"; }
dbg() { ((DEBUG)) && printf ">>> \e[0;38;2;133;133;133m$@\e[0m\n" || :; }
tip() { printf "\e[0;38;2;133;133;133m>>> $@\e[0m\n"; }
err() { printf "\e[0;33;1;133;133;133m>>> $@\e[0m\n" 1>&2; }
mvif() {
	local src="$1" dstdir="$2"
	if [ -d "$dstdir" ]; then
		mv "$src" "$dstdir"
	fi
}

cmd="$1" && (($#)) && shift
fn_exists "$cmd" && {
	eval $cmd "$@"
	unset cmd
} || {
	xcmd="golang-$cmd" && fn_exists "$xcmd" && eval $xcmd "$@" || {
		xcmd="build-$cmd" && fn_exists "$xcmd" && eval $xcmd "$@" || {
			xcmd="build-c$cmd" && fn_exists "$xcmd" && eval $xcmd "$@"
		}
	}
	unset cmd xcmd
}
