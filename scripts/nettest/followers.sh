echo 'tail -f out1.txt | grep ^[^2]'
factomd -prefix="_" -count=1 -networkPort="34341" -port="8092" -logPort="6062" -db=Map -peers="127.0.0.1:34340" -network=LOCAL -blktime=60 -net=alot+ > out1.txt


