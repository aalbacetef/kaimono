# kaimono

kaimono is a shopping cart library, that can be integrated into an existing server or as a standalone microservice.

It is a proof-of-concept implementation and not currently used in production (you might want to reach for the battle-tested `medusa` project [link](https://github.com/medusajs/medusa)).

The word itself means "shopping" in japanese.

**NOTE: currently in development/work-in-progress and not officially released**.

## Roadmap

- [ ] Finish Readme / Documenting
- [x] Implement standard routes
- [ ] Implement Admin routes
- [ ] Extensive tests
- [ ] Write some example code

## Usage 

### Assumptions 

The library makes only a few assumptions about the library consumer's backend:

- each user, logged in or anonymous, will have an associated session token.
- carts can be mapped to sessions.
- when sessions for logged-in users expire, Carts will be migrated to new sessions.
- when a user with an existing, valid session initiates a new session (e.g: different device), their existing Cart will be copied to the session (i.e: 1 cart to N sessions). 
- each product has associated with it an ID, a price, a title, a description.



### Core types

#### Cart, CartItem, Discount

These are the main data types pass to and from the API.

Their definitions are:

```go
type Cart struct {
	ID        string     `json:"id"`
	Items     []CartItem `json:"items"`
	Discounts []Discount `json:"discounts"`
}

type CartItem struct {
	ID        string     `json:"id"`
	Quantity  int        `json:"quantity"`
	Discounts []Discount `json:"discounts"`
	Price     Price      `json:"price"`
}

type Price struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type DiscountType string

const (
	PercentageDiscount  DiscountType = "percentage"
	FixedAmountDiscount DiscountType = "fixed-amount"
)

type Discount struct {
	ID    string       `json:"id"`
	Type  DiscountType `json:"type"`
	Value float64      `json:"value"`
}
```

#### Service 

The Service type is the main type used to interact with the library. 

It is backed by three interfaces: a DB interface, an Authorizer interface and a UserContextFetcher interface. These will be explained later in more detail, but the gist is that the DB interface provides CRUD methods for storage backend of the Cart while the UserContextFetcher allows the library to specify how the session token should be extracted from the request object. The Authorizer comes into play with the admin routes, authorizing (or not) a user for a given operation. 

Service exposes two methods for every CRUD operation: one only acts within the scope of the request's associated user/session while the other skips checking the session and acts direclty on the cart specified by the ID.

The idea is that one set of methods is used to expose standard shopping cart functionality to a website, while the other is used for admin purposes.

#### Standard Routes

Services exposes a router function for getting the standard route router:

```go
standardRouter := svc.Router("/cart")
```


The responses have the format:

```json
{
    "data": { /* depends on endpoint */ },
    "error": "<check depending on status code>"
}
```


Check the documentation at: [pkg.go.dev/github.com/aalbacetef/kaimono](https://pkg.go.dev/github.com/aalbacetef/kaimono) for full details of usage.
