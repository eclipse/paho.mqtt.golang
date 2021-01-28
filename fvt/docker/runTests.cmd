@ECHO OFF
REM Windows CMD file to run golang Paho tests with docker mosquitto instance
cls
:start
docker-compose up -d
REM Docker for windows does not support publishing to 127.0.0.1 so set the address for the tests to use.
set TEST_FVT_ADDR=0.0.0.0
REM `--count 1` prevents the system from using cached results. Note that running the tests multiple times may fail
REM because the broker state will not be as expected
go test --count 1 -v ../../
rem go test -race -v ../../
IF ERRORLEVEL   1 GOTO failed
GOTO successful

:failed
docker-compose down
echo "Error"
exit

:successful
docker-compose down
REM goto start