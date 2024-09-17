a) What are packages in your implementation? What data structure do you use to transmit data and meta-data?
We use structs to simulate tcp packages, and send sequence number, acknowledge number and message

b) Does your implementation use threads or processes? Why is it not realistic to use threads?
For this implementation we are using Goroutines, which are a lightweight thread, to answer to question - we are using threads.
Goroutines are concurrent, very reliable and can share memory. 
In the real world we have to assume a hostile environment, where faults can happen at any point.
Code on different machines might not take the same time to run, communication is not near-instant and can also take time/be lost,
Therefore it is not realistic to use threads.

We use go threads, even though it might not be realistic because they dont run in parallel, but line for line while waiting for other parts of the code.

c) In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
You could send the message where each has an index then sort it on retrieval.

d) In case messages can be delayed or lost, how does your implementation handle message loss?
Use timeouts if message not retrieved resend message.

e) Why is the 3-way handshake important?
It makes sure that both sides are in sync and tries to establish a connection only when both are ready.
