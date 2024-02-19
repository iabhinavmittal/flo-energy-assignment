# Questions:

1. Is the interval value in minutes?
2. Is it mandatorty to have full set of consumption values for each day? For example, is it possible to have only 10 "30-min" values for a day?
3. What is the scale that we are talking about? 10k records? 1M records? Getting a ball park number is helpful in avoiding over engineering, which I prefer.

# Assumptions:

1. For now, considering only a .csv file.
2. In case there is a wrong/mis-matched line item, I will be ignoring and storing it for review separately.
3. The solution will only return statements for insert and not actually insert anything into a DB
4. The dates are given in UTC time

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
  - looping through them can increase latency => benchmark for performance
  - batch insert or individual insert?

# Other Potential Issues:

- handle db insert failures (not part of the scope)
- what to do with malformed requests? => store them in-memory separately for review

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
