#!/usr/bin/env bash

for element in $(ls -1 | grep "split_"); do
  echo "Running: $element"
  go run . < $element
done
