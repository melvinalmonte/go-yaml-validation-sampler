#!/bin/bash

for i in {0..10}
do
    cp ./valid.yaml ./generated/"file$(printf "%d" "$i").template.yaml"
done
