# Sports News

How it works:

- poller polls the data from url
- it fills data in the storage from list while marking them to get details

What does not work:

- If hash of 2 different news is the same the poller will be confused and will treat them as same articles
- Extending articles struct will cause articles id to not much
- There is not pruning data built into the system
- Getting more than 100 latest news form hullcity