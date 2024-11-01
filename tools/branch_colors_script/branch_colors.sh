# Script to generate a bunch of branches to see what colors are used by other git gui tools

rm -fr .git
git init

for i in $(seq 80 0);
do
    git checkout -b "_$i"
    git commit --allow-empty -m "$i"
done

for i in $(seq 0 80);
do
    git checkout "_$i"
    git commit --allow-empty -m "$i"
done
