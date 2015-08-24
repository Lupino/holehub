HoleHUBD
========

The backend server of holehub.com.

Install
-------

    go get -v github.com/Lupino/holehub/holehubd

Run holehubd
------------

    holehubd --config_dir=/path/to/config --hole_host=hole_host --host=holehubd_host --port=holehubd_port --min_port=10000 --sendgrid_key=your_sendgrid_key --sendgrid_user=your_sendgrid_user

Run holed process manager
-------------------------

    go get -v github.com/bradfitz/runsit
    runsit --config_dir=/path/to/config

Next
----

Now you can you the client or install the front.
