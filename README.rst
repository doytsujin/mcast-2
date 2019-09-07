================================
Sockets, syscalls and goroutines
================================

With educational purposes:
   - implement multicast TCP chat using goroutines and map (IPv4, IPv6)
   - implement multicast UDP chat using low-level calls: syscalls (IPv4, IPv6)
   - play with params: backlog, time-wait, blocking vs non-blocking sockets

Usage
*****

Run server:

.. code-block::

   go run server.go

For TCP clients:

.. code-block::

   telnet 127.0.0.1 3000

For UDP clients:

.. code-block::

   echo "message" | nc -U 127.0.0.1 3000
