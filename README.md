# kaimono

kaimono is a shopping cart library, that can be integrated into an existing server or as a standalone microservice.

It is a proof-of-concept and not currently used in production (you might want to reach for the battle-tested `medusa` project (link)[https://github.com/medusajs/medusa]).

The word itself means "shopping" in japanese.

## Usage 

### Assumptions 

The library makes only a few assumptions about the library consumer's backend:

- each user, logged in or anonymous, will have an associated session token.
- carts can be mapped to sessions.
- when sessions for logged-in users expire, Carts will be migrated to new sessions.
- when a user with an existing, valid session initiates a new session (e.g: different device), their existing Cart will be copied to the session (i.e: 1 cart to N sessions). 
- each product has associated with it an ID, a price, a title, a description.



### Core types

#### Service 

The Service type is the main type used to interact with the library. 

It is backed by three interfaces: a DB interface, an Authorizer interface and a UserContextFetcher interface. These will be explained later in more detail, but the gist is that the DB interface provides CRUD methods for storage backend of the Cart while the UserContextFetcher allows the library to specify how the session token should be extracted from the request object. The Authorizer comes into play with the admin routes, authorizing (or not) a user for a given operation. 

Service exposes two methods for every CRUD operation: one only acts within the scope of the request's associated user/session while the other skips checking the session and acts direclty on the cart specified by the ID.

The idea is that one set of methods is used to expose shopping cart functionality to a website, while the other is used for admin purposes.


