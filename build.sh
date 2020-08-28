#!/bin/sh

# local test
#make clean && make && docker run -it --rm -v $(pwd)/bin:/in ripx80/upx -9 -o in/brewman_upx in/brewman
# compress with upx
make clean && make && docker run -it --rm -v $(pwd)/bin:/in ripx80/upx -9 -o in/brewman_arm_upx in/brewman_arm scp bin/brewman_arm_upx pi:~/brewman/brewman
# without upx
#make clean && make && docker run -it --rm -v in:/in ripx80/upx -9 -o in/brewman_arm_upx in/brewman_arm scp bin/brewman_arm_upx pi:~/brewman/brewman
exit 0
