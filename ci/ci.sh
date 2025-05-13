#!/bin/bash

# -------- ci tool for hedzr/store -------------------------------------
#
# [deprecated] After main module of hedzr/store published, do these
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
# [deprecated] END
#
# Howto release the New Version:
#
# 1. `$0 update`                                             : upgrade deps and make security patches
# 2. edit doc.go and CHANGELOG
# 3. `$0 update-main`                                        : upgrade deps to the releasing version of main lib (hedzr/store)
# 4. `make cov`                                              : ensure all tests passed
#    `go test ./... -v -race -cover -coverprofile=./logs/coverage-cl.txt -covermode=atomic -test.short -vet=off 2&>1 | tee ./logs/cover-cl.log`
# 4.99. commit the upgraded changes (go.mod & go.sum, ...)
# 5. `git push --all`                                        : commit all, and wait for remote tests passed
# 6. `git tag $VER && git push --all && git push --tags`     : bump version, push it
# 7. `$0 publish-all && git push --all && git push --tags`   : release the submodules
#
# 2024-12-12 Updates:
#
# After this updates, just one step needs before commit and push:
#
# 1. edit doc.go and CHANGELOG
# 2. `$0 update`                                             : update deps in main module, and child modules;
# 3. `make cov`                                              : ensure all tests passed
# 4. `git commit -am 'security patch' && git push --all`     : and now waiting for the remote ci passed
# 5. `git tag $VER && git push --all && git push --tags`     : bump version, push it
# 6. `$0 publish all && git push --all && git push --tags`   : release the submodules
#
# Using ci.sh to upgrade go modules under current directory, try this:
#
#     $0 update
#
# All folks.
#

BUILD_DIR=./build
INSTALL_TMP_DIR=./bin/install
OSN=store
VER=""

[ -f .version.cmake ] && {
	VER=$(echo $(grep -oE ' \d+\.\d+\.\d+' .version.cmake))
	# grep -oE ' [0-9]+.[0-9]+.[0-9]+' .version.cmake
	# echo "VERSION = $VER"
}

if [ x"$VER" == x ]; then
	notfound=1
	for f in doc.go _examples/doc.go _examples/small/doc.go slog/doc.go; do
		(($notfound)) && [ -f "$f" ] && {
			# echo "checking $f for VER..."
			VER="$(echo v$(grep -iE 'Version[ ]*=.*' "$f" | grep -oE '\d+\.\d+\.\d+'))"
			[ "$VER" != "v" ] && { echo "VERSION = $VER found!" && notfound=0; } # || { echo " ..loop next"; }
		}
	done
	unset notfound
else
	echo "VERSION = $VER found-"
fi

MODULE=$(grep 'module .*' go.mod | awk '{print $2}' | awk 'sub("git(lab|hub).com/","",$0)')

#

alias cmake="cmake -Werror=dev --warn-uninitialized"

#

#

#

build-test() {
	go test
}

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
build-publish-all() { build-publish-children all "$@"; }
build-publish-codecs() { build-publish-children codecs; }
build-publish-providers() { build-publish-children providers; }
build-publish-tests() { build-publish-children tests; }
build-publish-examples() { build-publish-children examples; }
build-publish-children() {
	local sm d f
	local which="$1"
	if [ "$which" = "all" ]; then
		pub-main
		pub-child codecs
		pub-child providers
		pub-child tests
		pub-child examples
	else
		for sm in $which; do
			for f in $(find ./$sm -type f -iname 'go.mod' -print); do
				d="$(dirname "$f")"
				if [[ ! "$d" == */test ]]; then
					if [ -d "$d" ]; then
						do-update-dep "$d"
						commit-submodule "$d"
					fi
				fi
			done
		done
	fi
}
build-setver() { build-setver-children "$@"; }
build-setver-codecs() { build-setver-children codecs; }
build-setver-providers() { build-setver-children providers; }
build-setver-tests() { build-setver-children tests; }
build-setver-examples() { build-setver-children examples; }
build-setver-children() {
	local sm d f dirty=0
	local which="$1" && shift
	for sm in "$which"; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			# if [[ ! "$d" == */test ]]; then
			if [ -d "$d" ]; then
				echo "  looking for ${d/\.\//}"
				update-submodule "${d/\.\//}" "$VER" "$d" || dirty=1
			fi
			# fi
		done
		if ((dirty)); then
			echo "  erase go.mod.bak and git commit $sm"
			find ./$sm -type f -iname 'go.mod.bak' -delete
			git commit -m "update $sm and publish them"
		fi
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			# if [[ ! "$d" == */test ]]; then
			if [ -d "$d" ]; then
				setver-submodule "${d/\.\//}" "$VER" "$d" "$@"
			fi
			# fi
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

update-submodule() {
	local tag="$1/$2"
	local d="$3" ret=0
	if [[ "$d" == */test* ]]; then
		pushd "$d" >/dev/null
		echo "    entering $d and sed $(ls -b go.mod) for tagged version $2 ..."
		sed -i.bak -E "s,(/store/)(codecs|providers)(/.*)v[0-9]+\.[0-9]+\.[0-9]+,\1\2\3$2," go.mod
		ret=$(diff go.mod go.mod.bak | wc -l)
		((ret)) && echo "      go mod tidy..." &&
			go get -v -t -u ./... && go mod tidy && git add go.mod go.sum &&
			ret=1
		popd >/dev/null
		# git tag "$tag"
	else
		echo "    git tag $tag"
		# git tag "$tag"
	fi
	return $ret
}

setver-submodule() {
	local tag="$1/$2"
	local d="$3"
	shift
	shift
	shift
	echo "    [setver] git tag $tag"
	git tag $* "$tag"
}

pub-main() {
	local ver="$VER"
	# if [ -f slog/doc.go ]; then
	# 	ver="$(grep -Eio 'Version += +\"(v?[0-9]+\.[0-9]+\.[0-9]+)\"' slog/doc.go | awk '{print $3}')"
	# elif [ -f doc.go ]; then
	# 	ver="$(grep -Eio 'Version += +\"(v?[0-9]+\.[0-9]+\.[0-9]+)\"' doc.go | awk -F$' ' '{print $3}')"
	# fi

	if [ "$ver" = "" ]; then
		echo "version tag not found, add doc.go and Version=\"1.0.0\" and retry."
	else
		ver="$(eval echo $ver)"
		echo "ver=$ver found"
		if is_git_dirty; then
			echo "repo is dirty, nothing to do before the changes are reviewed."
		else
			$ECHO git tag "$ver"
			$ECHO
		fi
	fi
}

pub-child() {
	local which="$1"
	if [ -d "$which" ]; then
		local ver="$VER"
		# if [ -f slog/doc.go ]; then
		# 	ver="$(grep -Eio 'Version += +\"(v?[0-9]+\.[0-9]+\.[0-9]+)\"' slog/doc.go | awk '{print $3}')"
		# elif [ -f doc.go ]; then
		# 	ver="$(grep -Eio 'Version += +\"(v?[0-9]+\.[0-9]+\.[0-9]+)\"' doc.go | awk -F$' ' '{print $3}')"
		# fi

		if [ "$ver" = "" ]; then
			echo "version tag not found, add doc.go and Version=\"1.0.0\" and retry."
		else
			ver="$(eval echo $ver)"
			echo "ver=$ver found"
			for sm in $which; do
				for f in $(find ./$sm -type f -iname 'go.mod' -print); do
					d="$(dirname "$f")"
					if [ -d "$d" ]; then
						# do-update-dep "$d"
						# commit-submodule "$d"
						echo "  - publishing $d "
						local pre="${d/.\//}"
						pushd "$d" >/dev/null
						if is_git_dirty; then
							echo "repo is dirty, nothing to do before the changes are reviewed."
						else
							$ECHO git tag "$pre/$ver"
							$ECHO
						fi
						popd >/dev/null
					fi
				done
			done
		fi
	fi
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

# update all go modules in .. [ cmdr.v2/ ]
build-update-all() {
	for d in ../*; do
		if [ -d "./$d" ]; then
			if [ -f "./$d/go.mod" ]; then
				pushd "./$d" >/dev/null
				# if [ -x ../libs.store/ci/ci.sh ]; then
				# 	../libs.store/ci/ci.sh update
				# elif [ -x ./ci/ci.sh ]; then
				# 	./ci/ci.sh update
				# fi
				build-update
				popd >/dev/null
			fi
		fi
	done
}

build-update() {
	headline "[update] upgrade any deps for '$MODULE'"
	for sm in "${1:-.}"; do
		for f in $(find "./$sm" -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			# if [[ "$(basename $d)" == '.'* ]]; then
			# 	echo && tip "==== ignore .github and others hidden directories: '$d'"
			# else
			if [ -d "$d" ]; then
				# tip "==== found go.mod in '$d'..."
				do-update-dep "$d"
			fi
			# fi
		done
	done
	echo
	[ -f go.work ] && build-update-main
	echo
}

# update all refs in child modules to hedzr/store's releasing version
build-update-main() {
	headline "[update-main][go.work] update dep to main module VERSION: $VER"
	local ix=0
	local mod="$MODULE"
	for sm in "${1:-.}"; do
		for f in $(find ./$sm -type f -iname 'go.mod' -print); do
			d="$(dirname $f)"
			if [ -d "$d" ]; then
				if [ "$f" != "././go.mod" ]; then
					# echo "d: $d"
					if grep -qE "$mod" $f; then
						tip "*** file: $f ***************"
						if sed -i '' -E -e 's#('$mod'.*)v[0-9]+\.[0-9]+\.[0-9]+#\1'$VER'#g' $f; then
							let ix++
							echo "   $f: $(grep -E $mod'.*v[0-9]+\.[0-9]+\.[0-9]+' $f)"
							[ -f "$f.bak" ] && rm "$f.bak"
						else
							echo "   $f: sed not ok"
						fi
						# if [[ $ix -gt 2 ]]; then
						# 	return
						# fi
					fi
				fi
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
	tip "==== go mod tidy with update all, dir='$d ($(pwd))' =========="
	go get -v -t -u ./... && go mod tidy
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

is_git_clean() { git diff-index --quiet "$@" HEAD -- 2>/dev/null; }
is_git_dirty() {
	if is_git_clean "$@"; then
		false
	else
		true
	fi
}

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
