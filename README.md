# testtoken

[![Build Status](https://travis-ci.com/axamon/testtoken.svg?branch=master)](https://travis-ci.com/axamon/testtoken)
[![Maintainability](https://api.codeclimate.com/v1/badges/034fec85a0c6070b9ef2/maintainability)](https://codeclimate.com/github/axamon/testtoken/maintainability)
[![Go Report Card](https://goreportcard.com/badge/github.com/axamon/testtoken)](https://goreportcard.com/report/github.com/axamon/testtoken)

testtoken checks the username and password passed as arguments against db and a json file in parallel.

If the credentials are found in any of the two storages than a pseudo token is returned.

# syntax
testtoken -u username -p password -t timeout