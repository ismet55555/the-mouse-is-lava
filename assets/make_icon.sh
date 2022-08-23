#/bin/sh

#######################################################################

# This script generates icon code

# Run:
#   ./make_icon <INPUT FILE PNG> <OBJECT OUTPUT NAME>
#   ./make_icon on.png iconON
#   ./make_icon off.png iconOFF

# Inputs:
#   1 - input file name
#   2 - output object

#######################################################################

if [ -z "$GOPATH" ]; then
    echo GOPATH environment variable not set
    exit
fi

if [ ! -e "$GOPATH/bin/2goarray" ]; then
    echo "Installing 2goarray..."
    go get github.com/cratonica/2goarray
    if [ $? -ne 0 ]; then
        echo Failure executing go get github.com/cratonica/2goarray
        exit
    fi
fi

if [ -z "$1" ]; then
    echo Please specify a PNG file
    exit
fi

if [ ! -f "$1" ]; then
    echo $1 is not a valid file
    exit
fi

OUTPUT="$2unix.go"
echo Generating file: $OUTPUT  -  Object: $2
echo "//+build linux darwin" > $OUTPUT
echo >> $OUTPUT
cat "$1" | $GOPATH/bin/2goarray Data $2 >> $OUTPUT
if [ $? -ne 0 ]; then
    echo Failure generating $OUTPUT
    exit
fi
echo Finished
