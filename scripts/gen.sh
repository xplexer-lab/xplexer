#!/bin/bash

for file in $(git grep -l -a "ENUM(" | grep -E "\.go$"); do
  go tool go-enum -f $file --output-suffix=_enum_gen
done
