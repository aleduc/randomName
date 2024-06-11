## Behaviour 
* Any error - stop processing, exact result or nothing.
* Date param can be empty
* Additional new line symbol in output, just to simplify.
## How to run 
* ./csv_count -f test.csv -d 2018-12-09
## Principles
* KISS 
* DRY 
* SOLID 
* YAGNI
  * Flat package structure, flat interfaces structure.
    **Example**: it is possible to have more structures, like separate "calculator" for row.
      would have been better for **S**OLID.
     But it breaks KISS, YAGNI.
* Avoid premature optimization(With one exception, since it is test task)

## Additional assumptions
* Any assumption during the next stage can be changed.
  **Examples**:
    * "You don't have enough memory to store the contents of the whole file"
    * File structure might be changed 
    * Columns format might be changed
    * Output format might be changed 
* Due to KISS/YAGNI, I didn't prepare code for everything. 