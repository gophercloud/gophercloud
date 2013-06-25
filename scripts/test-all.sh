#!/bin/bash
#
# This script is responsible for executing all the acceptance tests found in
# the acceptance/ directory.

# Find where _this_ script is running from.
SCRIPTS=$(dirname $0)
SCRIPTS=$(cd $SCRIPTS; pwd)

# Locate the acceptance test / examples directory.
ACCEPTANCE=$(cd $SCRIPTS/../acceptance; pwd)

# Run all acceptance tests sequentially.
# If any test fails, we fail fast.
for T in $(ls -1 $ACCEPTANCE); do
  if [ -x $ACCEPTANCE/$T ]; then
    echo "$ACCEPTANCE/$T -quiet ..."
    if ! $ACCEPTANCE/$T -quiet ; then
      echo "- FAILED.  Try re-running w/out the -quiet option to see output."
      exit 1
    fi
  fi
done

