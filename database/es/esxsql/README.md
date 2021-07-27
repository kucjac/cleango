# Esxsql

Module provides an implementation of the event store with event state tracking based on the SQL. 

`github.com/kucjac/cleango/database/es`

## Drivers

Currently, all SQL drivers are supported if the database is properly migrated.


## Migration 

Current implementation allows only PostgreSQL and MySQL automatic migration. 
PostgreSQL allows also sharding of the event table - shard by `aggregate_type` and event_state by `handler_name`.

In order to start migration use the function `esxsql.Migrate` with proper configuration of table names and 
sharding possibilities. 

In order to migrate the eventstate a configuration needs to have non-nil EventState field defined.

### Sharding

In case of sharding all aggregate types provided in the configuration would have its own partition table for the event.
The `Config` needs to have `PartitionEventTable` field set to `true`.
When a new aggregates is provided, it needs to have its own partition table created - in that case use 
`esxsql.MigrateEventPartitions` function, and provide all new aggregate types.

The table that is following event state could also be sharded. In order to migrate event state table with sharding enabled
mark `PartitionState` field as `true` in the `EventStateConfig`.   
