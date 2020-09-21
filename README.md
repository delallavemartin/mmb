# MMB - Martin Message Broker

The goal behind this project is to explore **go concurrency basics, and golang in general**.

An HTTP custom message broker as part of the Pub/Sub pattern was implemented in order to achive this.

# Project structure
* **consumer**: Consumer/Subscriber program.
* **mb**: Message broken program & performance metrics.
* **incremental-mb**: Step by step code evolution of **mb**. it has the purpose to show the incremental proccess involve to get the final solution. Each step brings a new basic go concurreny concept along with performance metrics.
	

# Architecture Design!
[Pub./Sub. Architecture Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/publisher-subscriber)
![Architecture Design](publish-subscribe.png)



I'm using HTTP as the communication protocol, so any HTTP notification request sent to the Message Broker will act as the *Publisher*.
The *Message Broker* will dispatch each request to each *Subscriber* in an __indepedent__ and __non-blocking way__. **Explore the Project for more fun!...**


