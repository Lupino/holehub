HoleHUB
=======

A cluster server of holed.

An easy way to run [hole](https://github.com/Lupino/hole)

Dirctory
--------

* [/config](https://github.com/Lupino/holehub/tree/master/config) the config of holehub.com.
* [/front](https://github.com/Lupino/holehub/tree/master/front) the site front of holehub.com
* [/holehubd](https://github.com/Lupino/holehub/tree/master/holehubd) the backend server of holehub.com.
* [/holehub](https://github.com/Lupino/holehub/tree/master/holehub) the client of holehub.com.


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
