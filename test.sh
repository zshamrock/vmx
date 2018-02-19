#!/usr/bin/env /bin/bash

TESTNUM=1

function colorize {
    echo -e "\e[1;32m${1}\e[0m"
}

function testcase {
    colorize "#${TESTNUM} ${1}"
    TESTNUM=$((TESTNUM+1))
}

function ok {
    colorize "OK"
    echo -e "Press Enter to continue..."
    read
}

# build the latest version
go build

testcase "defined host defined command"
./vmx run messaging-prod mem
ok

testcase "defined host non defined command"
./vmx run messaging-prod pwd
ok

testcase "non defined host defined command"
./vmx run rest-prod2 mem
ok

testcase "non defined host non defined command"
./vmx run rest-prod2 "free -h"
ok

testcase "command with the confirmation and no reply"
./vmx run rest-prod1 date
ok

testcase "command with the confirmation and yes reply"
./vmx run rest-prod1 date
ok

testcase "list command"
./vmx list
ok

testcase "all hosts group"
./vmx run all mem
ok

testcase "tail-ing/following"
./vmx run rest-prod1 less-logs
ok
