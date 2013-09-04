# Shep (Simple Git Based HTTP Fileserver)

## Requirements
1. [libgit2](https://github.com/libgit2/libgit2)
2. [git2go](https://github.com/libgit2/git2go)

## To Run
	go run *.go -g ShepBareRepo [-p port]
	
## Features

*  Git push to shep via over http (shepHost:shepPort/_git) (See Current Issues)
*  Request anything at the HEAD of the repo master branch

## Further Development
* Request file at any revision via a commit-ish line in the url e.g. localhost:8080/myfile.txt@HEAD^

## Current Issues
* Try pushing something other than a small text file… it's not working very well.
 
