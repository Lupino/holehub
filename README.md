HoleHUB
=======

A cluster server of holed server.

An easy way to run [hole](https://github.com/Lupino/hole)

Install
-------

    go get -v github.com/Lupino/holehub/holehub

Quick start
-----------

    # signup on holehub.com
    curl -d username=yourusername -d password=yourpassword -d email=youremail http://holehub.com/api/signup/
    # then go to you email inbox active you account

    # Login holehub.com
    holehub login

    # run a app
    holehub run --rm -n sshd -lp 22
