In this document I am going to document my workflow and progress with the challenge.

2022/01/13
After I came back from the interview I had a look at the challenge readme and supplied code. From the Readme, it seemed clear that it was a classic graph traversal problem. However as I explored different graph structures, it seemed that the example given in the description did not cover some cases (common ancestors, cycles, etc). I sent an e-mail asking for clarification and received a quick reply stating that the graph will actually be a tree. Which simplifies things a lot, really.

// PICTURE

2022/01/16 5PM
I started this repo, commented out some code to tighten my development feedback loop while I work (authentication, simulated network delay, random failures, just 10 pages instead of 100 etc.), run it and hit the endpoint with Postman.

// PICTURE

This code challenge is very simple, so frankly no design is required except for basic common sense. 

However I like to separate I/O from non-I/O code (e.g. domain logic, formatting, checks) for three reasons: (1) Non-i/o code is simple to unit test (2) I find shorter files that deal with one thing only easier to read (3) Coherence: usually it's easier to spot code re-use opportunities. In this case, most of the code will probably deal with network issues (I/O) except for the indexing size calculation. So the structure I will use is something like this:

```
main.go 
network.go
calculation.go
calculation_test.go
```

I might also add a service-level test at the end (`main_test.go`)


2022/01/16 18:30PM
Well, I implemented a synchronous version after some thought and reading the docs. It doesn't handle network stuff.

Well, scratch all of that. The actual calculation is trivial, the real issue is how to do it concurrently.

I think the easiest thing to do is to keep the current implementation and use a channel for the queue and the counter.

Next step is to make it work concurrently
Then handle i/o problems
Then renew token
