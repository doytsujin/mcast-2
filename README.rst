====================================
Sockets, syscalls, goroutines and Go
====================================

With educational purposes:
   - implement multicast TCP with goroutines and map
   - use multicast UDP using low-level api (syscalls and sockets)
   - 


Run server:

   .. code-block::

      go run server.go

For TCP clients:

   .. code-block::

      telnet 127.0.0.1 3000

For UDP clients:

   .. code-block::

      echo "message" | nc -U 127.0.0.1 3000
