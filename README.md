# ferryman

Ferryman is intended to be a HTTP/HTTPS Loadbalancer with the following features:

## Version 1 will be achived if :

* [] TCP/IP Connection pooling.
* [] URI Rewriting.
* [] Sticky Session Support.
* [] Response Body content processing (ie chaging a resource host from an intenral uri to external uri.).
* [] Route failure fallback support.
* [] DSL to add routing/rules using configuration.
* [] Hot config/rule reloads.
* [] Caching of http responses.
* [] Support for websockets.
* [] Generating NginX configuration.

It is currently in development as a pet project.

## Getting Started

Nothing to see here keep on moving.

## Content consumed while learning
* [![Datastructure Theory - ART]](https://db.in.tum.de/~leis/papers/ART.pdf)
* [![Datastructure Theory - Tree Comparison]](http://daslab.seas.harvard.edu/classes/cs265/files/presentations/CS265_presentation_Sinyagin.pdf)

## Other HTTP Routing/Proxy Projects scoured and their limitations
* https://github.com/julienschmidt/httprouter, no: regex matching, reverse proxy, content rewrite
* https://github.com/buaazp/fasthttprouter, no: regex matching, reverse proxy, content rewrites
* https://github.com/containous/traefik, no: content rewrites
