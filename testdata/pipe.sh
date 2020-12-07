#!/bin/sh
echo "$2" | nc -w0 -U "$1"
