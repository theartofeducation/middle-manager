# Decision Log Book

Taken from [Decision Management in Software Engineering](https://medium.com/swlh/decision-management-in-software-engineering-ca60f9d40e02)

---

## 2020-11-17: What ORM should we use

### Decision Makers

* Christopher Lamm
* Thomas Jean

### Context

We will be using Go for several services and need to pick an ORM to use across them all, so we have uniformity.

### Solution

[ent](https://github.com/facebook/ent)

#### Why This Solution

This tool is backed by Facebook and its community which will give us a lot of examples and support.

Schema will be code and not hidden in tags.

#### Limitation

No known limitations at this time.

### Rejected Solutions

* [gorm](https://gorm.io/)
    * More traditional ORM with expected overhead

---

## 2020-11-17: What web framework should we use

### Decision Makers

* Christopher Lamm
* Thomas Jean

### Context

We will be using Go for several services and need to pick a framework to use across them all, so we have uniformity.

### Solution

[gorilla/mux](https://github.com/gorilla/mux)

#### Why This Solution

Mux and Chi are pretty similar, and we can easily switch one out for the other.

As of this writing Chi is not a compatible Go module while Mux is, so we're basing our decision on that.

#### Limitation

No limitations at the time of writing.

### Rejected Solutions

* [Gin](https://github.com/gin-gonic/gin)
    * Requires special middleware pattern
    * Uses special Context struct over standard library
* [Echo](https://github.com/labstack/echo)
    * Requires special middleware pattern
    * Uses special Context struct over standard library
* [Chi](https://github.com/go-chi/chi)
    * Not yet a compatible Go modules
