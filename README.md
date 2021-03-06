IMPORTANT
=========

go-cjdns was recently restructured and has probably broken any programs that
were using it. Sorry about that, but it's for the best, I promise. The code is
much better organized now and will continue to be improved.

[![GoDoc](http://godoc.org/github.com/inhies/go-cjdns?status.png)](http://godoc.org/github.com/inhies/go-cjdns)
[![Build Status](https://travis-ci.org/inhies/go-cjdns.png?branch=development)](https://travis-ci.org/inhies/go-cjdns)
[![Coverage Status](https://coveralls.io/repos/inhies/go-cjdns/badge.png?branch=development)](https://coveralls.io/r/inhies/go-cjdns?branch=development)


go-cjdns
========

This is a group of packages for working with cjdns using Go. 


go-cjdns/key
------------
Package key provides access to cjdns encryption keys and associated methods.


go-cjdns/config
---------------

Package config allows easy loading, manipulation, and saving of cjdns
configuration files. It can load, parse, and save the standard cjdns file as
well as a file with custom data added. This enables you to store additional
information with your peer connections and server details.


go-cjdns/admin
--------------

Package admin provides methods to access a running cjdns instance via the admin
udp socket. It not only allows you to send any command and receive the response
but it also provides convenience functions. It relies on go-cjdns/config for
loading of the configuration data and getting the IP address, port, and password
for the admin interface.

