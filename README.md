# Prerequisite

1. go 1.22

# How to run?

1. Go to root folder
2. run `go get .` to install the dependencies
3. run `go run main.go date_helper.go`. It will print the batched insert statments.
4. run `go test` to run the unit tests.

# Assumptions:

1. For now, considering only a .csv file.
2. In case there is a wrong/mis-matched line item, I will be ignoring and storing it for review separately.
3. The solution will only return statements for insert and not actually insert anything into a DB
4. The dates are given in UTC time
5. Interval can be 5 min, 15 mins and 30 mins only (based on manual)

# Psuedo Code:

- read from a csv
  - read each line
  - push into an array and return
- for each line
  - filter for 200 and 300 lines
  - validate the structure?
    - should contain nmi
    - should contain interval length
    - should contain a date
    - should contain sufficient consumption values based on time interval
    - if not valid, ignore and push into an error array for reconsideration
  - generate insert statements for each item in array
- output: insert into a table statements

# Scale Issues:

- ~1 million records
  - can we store so much data in memory?
  - looping through them can increase latency => keeping it O(n) to reduce latency
  - batch insert or individual insert? => batching it with a variable INSERT_BATCH_SIZE

# Other Potential Issues:

- handle db insert failures (not part of the scope)
- what to do with malformed requests? => ignore and store them in-memory separately for review

# Test Cases (WIP):

- happy path with correct values
- empty file
- invalid file, assuming only csv files for now
- missing NEM value
- missing interval value
- wrong interval value?
- missing date value
- wrong date format
- insufficient number of consumption values according to interval?
- large number of records: x ?
- no 200 record found
- no 300 record found
- wrong order, 300 records are above 200
