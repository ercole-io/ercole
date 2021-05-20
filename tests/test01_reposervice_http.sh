#!/bin/sh
# Copyright (c) 2020 Sorint.lab S.p.A.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

# SCOPE: This test try to list files inside repo service and download foo.bar via HTTP

# TEST config
TEST_ROOT_DIR=`mktemp -d /tmp/ercole-XXXXXXXXXXXXXXXX`
CONFIG_FILENAME="$TEST_ROOT_DIR/config.toml" 
PORT=`tests/freeport.sh`
DISTRIBUTED_FILES="$TEST_ROOT_DIR/distributed_files"
DOWNLOADED_FILES="$TEST_ROOT_DIR/downloaded_files"
ESCAPED_ROOT_DIR=$( echo $TEST_ROOT_DIR | sed -e 's/[\/&]/\\&/g' )
RESULT=100

# Materialize the config
cp -r tests/test_distributed_directory $DISTRIBUTED_FILES
cat tests/only_reposervice_http_config.toml | sed -E "s/%TEST_DIR%/$ESCAPED_ROOT_DIR/" | sed -e "s/\"%PORT%\"/$PORT/" > $CONFIG_FILENAME
mkdir -p $DOWNLOADED_FILES

# Start ercole
./ercole -c $CONFIG_FILENAME serve --enable-repo-service > /dev/null &
ERCOLE_PID=$!
sleep 0.1

# Perform some operations
wget --quiet -O $DOWNLOADED_FILES/foo.bar http://127.0.0.1:$PORT/foo.bar 
if cmp --silent $DOWNLOADED_FILES/foo.bar $DISTRIBUTED_FILES/foo.bar; then
    echo "OK"
    RESULT=0
else
    echo "FAILED"
    RESULT=2
fi


# Cleanup
rm -r $TEST_ROOT_DIR
kill $ERCOLE_PID

# Exit
exit $RESULT
