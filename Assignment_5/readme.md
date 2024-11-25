How to run application/demonstration:

1.      Setting up the servers
1.1     Run 3 different instances of the server.go file in the server folder
1.2     Provide a port number for each of the different instances. The port numbers are to be: 5050 , 5051 , 5052
(NOTE: This step can also be done automatically on windows machines by running the run_servers.bat file)

2.      Running a client
2.1     Run the client.go file in the client folder.
2.2     Provide the client with a name to identify by.
(NOTE: A single client can also be quickly run on windows machines by running the run_clients.bat file)

3.      Interacting with the bidding system as a client:
3.a     To see what the status of the auction system / see the higest bidder type the command "Result".
3.b     To bid type the command "Bid <Amount>". Whenever the servers recieve a bid for the first time they begin the auction process and clients have 100s to bid before it ends.