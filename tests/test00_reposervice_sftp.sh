#!/bin/sh
# Copyright (c) 2019 Sorint.lab S.p.A.
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

# SCOPE: This test try to list files inside repo service and download foo.bar via SFTP

# TEST config
TEST_ROOT_DIR=`mktemp -d /tmp/ercole-XXXXXXXXXXXXXXXX`
CONFIG_FILENAME="$TEST_ROOT_DIR/config.toml" 
PORT=`tests/freeport.sh`
DISTRIBUTED_FILES="$TEST_ROOT_DIR/distributed_files"
PRIVATE_KEY="$TEST_ROOT_DIR/test.key"
DOWNLOADED_FILES="$TEST_ROOT_DIR/downloaded_files"
TEST_OUTPUT="$TEST_ROOT_DIR/output.txt"
TEST_EXPECTED_OUTPUT="$TEST_ROOT_DIR/expected_output.txt"
ESCAPED_ROOT_DIR=$( echo $TEST_ROOT_DIR | sed -e 's/[\/&]/\\&/g' )
RESULT=100

# Materialize the config
cp -r tests/test_distributed_directory $DISTRIBUTED_FILES
openssl genrsa -out $PRIVATE_KEY 1024 > /dev/null 2>/dev/null
cat tests/only_reposervice_sftp_config.toml | sed -E "s/%TEST_DIR%/$ESCAPED_ROOT_DIR/" | sed -e "s/\"%PORT%\"/$PORT/" > $CONFIG_FILENAME
cat tests/test00_reposervice_sftp_expected_output.txt | sed -E "s/%TEST_DIR%/$ESCAPED_ROOT_DIR/" > $TEST_EXPECTED_OUTPUT
mkdir -p $DOWNLOADED_FILES

# Start ercole
./ercole -c $CONFIG_FILENAME serve --enable-repo-service > /dev/null &
ERCOLE_PID=$!
sleep 0.1

# Perform some operations
echo "ls\npwd\ncd ..\nls\npwd\nget foo.bar $DOWNLOADED_FILES/foo.bar" | sftp -P $PORT -oStrictHostKeyChecking=no -oUserKnownHostsFile=/dev/null 127.0.0.1 2>/dev/null > $TEST_OUTPUT 
if cmp --silent $TEST_OUTPUT $TEST_EXPECTED_OUTPUT; then
    if cmp --silent $DOWNLOADED_FILES/foo.bar $DISTRIBUTED_FILES/foo.bar; then
        echo "OK"
        RESULT=0
    else
        echo "FAILED"
        RESULT=2
    fi
else
    echo "FAILED"
    RESULT=1
fi

# Cleanup
rm -r $TEST_ROOT_DIR
kill $ERCOLE_PID

# Exit
exit $RESULT
