# Deployment Sequence and Rollback Strategies

## 1. Rolling Deployment Realities
During rolling deployments, there is always a window where:
- Instance A runs **Version N** (Old code)
- Instance B runs **Version N+1** (New code)

Both instances share the same database concurrently.

### Fundamental Rule of Migration Safety
**The database schema must ALWAYS be compatible with both Version N and Version N+1 simultaneously.**

If Version N+1 alters a column in a way that breaks Version N, Instance A will throw runtime database errors, crashing production during the rolling deployment window.

## 2. Safe Deployment Sequence
To ensure compatibility across versions, migrations must be broken into separate deployments:

1. **Deployment 1 (Expand Schema)**: Apply database migrations to create new tables/columns (all non-breaking, nullable, or with defaults).
2. **Deployment 2 (Dual Write & Support New Schema)**: Deploy Application version supporting dual writes and backfill workers. Old instances continue reading old fields without crashing.
3. **Run Backfill**: Run data backfill script in the background to sync legacy data to the new schema.
4. **Deployment 3 (Switch Read Path)**: Deploy Application version that reads exclusively from the new schema (or uses a feature flag to ramp up reads).
5. **Deployment 4 (Stop Dual Write)**: Deploy Application version that ceases writes to the legacy schema.
6. **Deployment 5 (Contract Schema)**: Drop the legacy column or table in the database.

## 3. Rollback Scenarios
What happens if Application N+1 is deployed, encounters an issue, and must be rolled back to Application N?

- **If Schema was Expanded (Non-destructive)**:
  Application N still functions properly because the legacy columns and tables were untouched. Safe to roll back.
- **If Rollback occurs during Dual Write**:
  Data remains safe if both schemas were maintained synchronously.
- **If Rollback occurs after stopping Dual Write**:
  **Data Loss Risk!** Application N expects to find data in the old schema, which was not updated during the N+1 runtime.
- **Rule**: Never execute the Contract phase until the new application version has demonstrated stable production operation over a sustained observation period.
