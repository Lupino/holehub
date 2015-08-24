HoleHUB front
=============

The site front of holehub.com

Install
-------

    npm install
    make
    make dist

Nginx config
-----------
    server {
        listen       [::]:80;
        listen       80;
        server_name  holehub.com www.holehub.com;

        access_log  /path/to/access.log;
        error_log   /path/to/errors.log;

        location / {
            root /path/to/front/dist;
            index  index.html index.htm;
        }
        location /api/ {
            proxy_pass         http://holehubd_host:holehubd_port;
        }
    }
