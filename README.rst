================================
Sockets, syscalls and goroutines
================================

With educational purposes:
   - implement multicast TCP chat using goroutines and map
   - implement multicast UDP chat using low-level api (syscalls and sockets)
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
