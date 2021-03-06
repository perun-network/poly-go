#!/bin/bash

# SPDX-License-Identifier: Apache-2.0

set -e
exit_code=0

# This script checks that all new commiters email adresses are contained in the NOTICE
# file.
# To check all commits on the current branch:       ./check-notice-authors.sh
# To check all commits newer than the last merge:   ./check-notice-authors.sh $(git log --pretty=format:"%H" --merges -n 1)

# Call with an ancestor whereas all commits newer than the ancestor are checked.
base="$1"
if [ -z "$base" ]; then
    commits="$(git rev-list --no-merges --reverse HEAD)"
else
    commits="$(git rev-list --no-merges --reverse $base..HEAD)"
fi

echo "Current commit: $(git rev-parse HEAD)"
# Authors found in commits and NOTICE.
declare -A known_authors
# Authors found only in commits but not NOTICE file.
declare -A assumed_authors

for c in $commits; do
    echo "Checking commit: $c"
    author=$(git show -s --format='%an <%ae>' $c)

    # Get the notice file of the commit and check that the author is in it.
    notice=$(git show $c:NOTICE 2> /dev/null || true)
    for k in "${known_authors[@]}"; do
        a="${known_authors[$k]}"
        if ! echo "$notice" | grep -wq "$a"; then
            echo "Author '$a' was deleted from NOTICE in commit $c"
            exit_code=1
        fi
    done
    if [ -n "${assumed_authors[$author]}" ]; then
        continue
    fi
    # This must be the first commit of this author, since he should add himself
    # to the NOTICE file here.
    if ! echo "$notice" | grep -wq "$author"; then
        echo "Author '$author' is missing from NOTICE file and should have been added in commit $c."
        assumed_authors[$author]="$author"
        unset "known_authors[$author]"
        exit_code=1
    else
        known_authors[$author]="$author"
        unset "assumed_authors[$author]"
    fi
done

exit $exit_code
