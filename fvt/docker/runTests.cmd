@ECHO OFF
REM Windows CMD file to run golang Paho tests with docker mosquitto instance
cls
docker-compose up -d
REM Docker for windows does not support publishing to 127.0.0.1 so set the address for the tests to use.
set TEST_FVT_ADDR=0.0.0.0
go test -v ../../
rem go test -race -v ../../
docker-compose down