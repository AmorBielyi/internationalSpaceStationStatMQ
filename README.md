## Obtain International Space Station (ISS) real-time data and proceed with RabbitMQ
RabbitMQ producer and consumer uses Publish/Subscribe pattern.
Consumer either can tell you about ISS geolocation, orbital speed or ISS visibility in real-time.
## What is it? 
Producer continuously obtains real-time ISS data via NORAD API and publishes it into the queue,

Consumers subscribe to queue and print useful information about ISS.

**usage (producer):**
 - `go run producer/produce.go`

      It will continuously produce messages with ISS data 
      and check X-Rate-Limit usage of API.  
 
**usage (consumer - geolocation mode):**
  - `go run consumer/consume.go location`
  
       It will consume the message and send useful information about ISS 
       like its longitude, latitude and current orbital speed.  
   
       **Example output:** 
       2021/10/24 16:52:38 ISS geolocation is: longitude: 134.07502033476 and latitude: -43.301622038287;
                             orbital speed is about 27531.771607701 km/h
     
**usage (consumer - visibility mode)**
  - `go run consumer/consume.go visibility`  
  
       It will consume the message and send useful information about ISS
       like its current visibility, e.g eclipsed or daylight.
       
       **Example output:**
       2021/10/24 16:54:03 ISS visibility now is eclipsed  