#!/usr/bin/env bash
echo "will do all package test"

echo "mode: set" > acc.out
for Dir in $(find ./* -maxdepth 10 -type d ! -path './vendor*'); 
do
	if ls $Dir/*.go &> /dev/null;
	then
		returnval=`go test -coverprofile=profile.out $Dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.out ]
    		then
        		cat profile.out | grep -v "mode: set" >> acc.out 
    		fi
    	else
    		exit 1
    	fi	
    fi
done

if [ "$TRAVIS_PULL_REQUEST" = "false" ]; then 
	goveralls -coverprofile=acc.out -service travis-ci -repotoken $COVERALLS_TOKEN
fi
