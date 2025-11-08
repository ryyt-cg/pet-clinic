# GORM

## Join vs Preload

the computational cost is different for joins / vs lookups.




## When to Use Preload
* Avoiding N+1 queries: When you need to load associated records but don’t need filters or sorting.
* Large datasets: Suitable for large datasets to avoid JOIN-related overhead.
* Multiple associations: Works efficiently for loading multiple associations.

**Limitations**
* No filtering: You can’t filter or sort the associated records.
* No JOIN optimization: If filtering or sorting is required, consider eager_load or includes.


## When to Use Eager Load
Filtering or sorting: When you need to filter or sort based on the associated model’s fields.
* Single query: Useful for reducing query overhead by combining everything into one query.

**Limitations**
* Complex joins: For very complex data, eager_load queries might become slow. Ensure indexes are in place.