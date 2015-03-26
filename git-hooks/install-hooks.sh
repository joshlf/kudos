#!/bin/bash

# source: http://stackoverflow.com/a/957978/836390
ROOT=$(git rev-parse --show-toplevel) || exit 1

if [ "$ROOT" == "" ]; then
	echo "`git rev-parse --show-toplevel` returned empty root path" >&2
	exit 1
fi

cd $ROOT/.git/hooks || exit 1

for fpath in ../../git-hooks/*; do
	fname="$(echo $fpath | sed -e 's/\/$//' | rev | cut -d / -f 1 | rev)"
	if [ "$fname" != "install-hooks.sh" -a "$fname" != "README" ]; then
		ln -s "$fpath" "$fname" || exit 1
	fi
done
