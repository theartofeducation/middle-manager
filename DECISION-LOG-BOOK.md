# Decision Log Book

Taken from [Decision Management in Software Engineering](https://medium.com/swlh/decision-management-in-software-engineering-ca60f9d40e02)

---

## 2020-12-18: How do we bridge the gap between ClickUp and Clubhouse

### Decision Makers

* Christopher Lamm
* Chris Poulton
* Bob Yexley
* Thomas Jean
* Judy Mosley

### Context

We need to be able to manager and work stories through our normal process in Clubhouse while allowing Chris to manage the product side in ClickUp.

### Solution

Developer a service that will automatically updated ClickUp and Clubhouse based on events in the opposite tool.

#### Why This Solution

This will give us the most flexibility to perform the automation we need.

Chris can manager product in ClickUp.

Engineering can manager work in Clubhouse.

#### Limitation

Some actions might not be able to be automated and will have to be performed manually.

### Rejected Solutions

* Just use [ClickUp](http://clickup.com/)
    * While most workflow issues seemed to have a hackable solution, the integration with GitHub was going to be too much work with every team member having to setup 4 integrations for every repo.
* [Zapier](https://zapier.com/)
    * Their integrations only flowed from ClickUp to Clubhouse.
    * Very simple events supported and would not fit our use case.
