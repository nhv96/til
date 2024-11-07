# Full text search in PostgreSQL
The purpose of full text search is to look for words or groups of words that can be found in a text, full text search is more of a *contains* operation.

In PostgreSQL, full text search can be done using ***GIN*** index. The index will dissect a text, extract valuable lexemes (tokens of words) string, and index those elements rather than the underlying text.

Example:
```
# SELECT to_tsvector('english','A car, I want a car. I would not even mind having many cars');
                          to_tsvector                          
---------------------------------------------------------------
 'car':2,6,14 'even':10 'mani':13 'mind':11 'want':4 'would':8
(1 row)
```

The function `to_tsvector` take the string, apply English language rules, and perform a stemming process, the process remove the stop words, and stem individual words. So `cars` will become `car`, `many` become `mani`.

So that when search with a word such as `wanted` against that string, we will be able to find it.
```
SELECT to_tsvector('english','A car, I want a car. I would not even mind having many cars') @@ to_tsquery('english','wanted');
 ?column? 
----------
 t
(1 row)
```

# How to enable full text search
## Using functional index with GIN
```
CREATE INDEX idx_name ON table_name USING gin(to_tsvector('english', column_name));
```
This approach is slower but doesn't cost more space.
## Using materialised column
This approach takes more space, but with better runtime.
```
ALTER TABLE table_name ADD COLUMN ts tsvector;
ALTER TABLE

CREATE TRIGGER tsvectorupdate
    BEFORE INSERT OR UPDATE ON table_name
    FOR EACH ROW
    EXECUTE PROCEDURE
        tsvector_update_trigger(some_name, 'pg_catalog.english', 'column name');
```
