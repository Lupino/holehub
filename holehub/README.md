HoleHUB
=======

The client of holehub.com.

Install
-------

    go get -v github.com/Lupino/holehub/holehub

Quick start
-----------

Go to holehub.com then signup or signup by curl:

    # signup on holehub.com
    curl -d username=yourusername -d password=yourpassword -d email=youremail http://holehub.com/api/signup/
    # then go to you email inbox active you account

Process the client:

    # Login holehub.com
    holehub login

    # run a app
    holehub run --rm -n sshd -lp 22
