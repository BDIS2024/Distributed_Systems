How to run application/demonstration:

## Setting up the servers
<ol>
    <li>Run 3 different instances of the server.go file in the server folder</li>
    <li>Provide a port number for each of the different instances. The port numbers are to be: 5050 , 5051 , 5052</li>
</ol>
(NOTE: This step can also be done automatically on windows machines by running the run_servers.bat file)

## Running a client
<ol>
    <li>Run the client.go file in the client folder</li>
    <li>Provide the client with a name to identify by.</li>
</ol>   
(NOTE: A single client can also be quickly run on windows machines by running the run_clients.bat file)

## Interacting with the bidding system as a client:
- To see what the status of the auction system / see the higest bidder type the command "Result".
- To bid type the command "Bid <Amount>". Whenever the servers recieve a bid for the first time they begin the auction process and clients have 100s to bid before it ends.