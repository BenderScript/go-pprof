set -e
set -x

# One request
ab -n 1 http://127.0.0.1:15120/

# 10000
ab -n 10000 http://127.0.0.1:15120/
