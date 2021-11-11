install:
	go install

run: install
	checkip 91.228.166.47
	checkip 209.141.33.65
	checkip 218.92.0.158
	checkip -j 91.228.166.47 | jq -r '.[] | select(.IPaddrMalicious==true) | .Name' | sort
	checkip -j 209.141.33.65 | jq -r '.[] | select(.IPaddrMalicious==true) | .Name' | sort
	checkip -j 218.92.0.158 | jq -r '.[] | select(.IPaddrMalicious==true) | .Name' | sort
