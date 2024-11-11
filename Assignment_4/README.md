How to run the application:

- Firstly run the server and enter a port number. (You should establish 3 servers on ports:5050,5051,5052)
    - This can also be done by executing the run_servers.bat file

- Secondly run the client and enter a port number to make a pair with. (Each server established should have a client connected to it.)
    - This can also be done by executing the run_clients.bat file

When done:

- The applications should be running and all nodes will make an attempt to access the critical section

- Check the logs for when each node enters the critical section.


The application uses a peer2peer node system and a node is simulated by a server and a client. The server has the responsibility of first of foremost establishing a stream to communicate with other nodes. And the client has the responsibility of connecting to other nodes/server ports.

