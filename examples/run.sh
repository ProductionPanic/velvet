#!/bin/sh
# get name of example to run and default with demo
EXAMPLE=${1:-demo}
# run example
cd $EXAMPLE
go run .