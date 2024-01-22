#!/bin/sh
echo 'making tiny'
echo '0' > data/tiny

echo 'making normal'
rm -f data/normal
for i in $(seq 1 100); do
	echo "000000000000000000" >> data/normal
done

echo 'making massive'
rm -f data/massive
for i in $(seq 1 1000000); do
	echo "000000000000000000" >> data/massive
done
