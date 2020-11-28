set -e
set -x
docker run -d -p 15221:15121 -p 15120:15120 go-pprof
