#!/bin/bash

# SPDX-License-Identifier: Apache-2.0

ERR=0

# Call with working directory, eg: .scripts/check-vanity-imports.sh $PWD
[ "$#" -ne 1 ] && (echo "Need working directory as first argument" && exit 1)

# Loop over all subdirectories not starting with a . and alias them $pkg.
for pkg in $(find $1 -mindepth 1 -type d ! -path '*/\.*'); do
	# Count the vanity imports in the $pkg folder.
	numImports="$(find $pkg -maxdepth 1 -type f -name '*.go' -print0 | xargs -0 egrep -ho '^package [a-z0-9_]+ // import ".*$' | wc -l)"
	# Count how many non _test package definitions (eg. package my_test) the $pkg folder contains.
	numNonTestPackages="$(find $pkg -maxdepth 1 -type f -name '*.go' -print0 | xargs -0 egrep -ho '^package [a-z0-9_]+' | egrep -v '_test' | wc -l)"

	# _test packages are allowed to not have a vanity import path.
	# So if the directory contains a non-_test package, there must be an import path.
	if [ $numImports -eq 0 ] && [ $numNonTestPackages -gt 0 ]; then
		echo "Package is missing vanity import path: $pkg"
		ERR=1
	# Here we check that the whole folder does not have more than one vanity
	# import. This implies that _test packages do not have import paths when
	# they are in the same folder as a non _test package.
	elif [ $numImports -gt 1 ]; then
		echo "Package has more than one vanity import path: $pkg"
		ERR=1
	fi
done

exit $ERR
