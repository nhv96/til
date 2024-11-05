# Slow performance with LIKE query
Like queries or %, usually have the worst performance problem.

In Postgres we can speed up this kind of query using an extension called `pg_trgm`.

This extension supports fuzzy search, it calculates the distance between two strings.
```
SELECT 'abcde' <-> 'abceacb';
?column?
--------
0.83333
--------
(1 row)
```
The lower the number, the more similar the two strings are.

This is perfect for use cases such as a searching query, when input keyword can sometime not exactly what the user meant, but they also can't remember the exact word.

# How it works
`pg_trgm` comes with an operator class called `GiST_trgm_ops`, which is design to do similarity searches.

First, we enable this extension.
```
CREATE EXTENSION pg_trgm;
```

Then, for the column that we are running LIKE query against, we add an index.
```
CREATE INDEX `idx_trgm` ON table_name USING GiST(column name GiST_trgm_ops);
```

Now, when running query against the column `name`, the query planner will chose the `idx_trgm` index to work on, this index now use the operator class `GiST_trgm_ops` to do the search, hence speed up the query.

Trigram indexes also can help speed up regular expression queries.