# testtoken

testtoken checks the username and password passed as arguments against db and a json file in parallel.

If the credentials are found in any of the two storages than a fake token is returned.

testtoken -u <username> -p <password>