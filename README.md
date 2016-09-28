vpc-peering
===========

A script that automatically accepts vpc peering connection
based from ACCOUNT IDs without the need to specify peering
connection ID which is a bit tasky.

Usage:
-------
Please export AWS token and default region for this to work.
Either specify OWNER ID and/or ACCOUNT_IDS as environment
variables or through .env file

```
$ go run main.go --help
```

Compile:

```
$ make ARCH=darwin
```

Execute binary:

```
$ <packagename> -o ACCOUNT_ID -l ACCOUNT_01 -l ACCOUNT_02 -l ACCOUNT_03 -r ap-southeast-2
```

Cross compile for Linux:

```
$ make ARCH=linux
```
