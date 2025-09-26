# Database Migrations

This project using [golang-migrate](https://github.com/golang-migrate/migrate) for safer database schema management.

## Installing golang-migrate

```bash
make migrate-install
```

Or install manually:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Using Migration Commands

### Apply all migrations
```bash
make migrate-up
```

### Rollback last migration
```bash
make migrate-down
```

### Show current migration version
```bash
make migrate-version
```

### Rollback all migrations
```bash
make migrate-reset
```

### Create new migration file
```bash
make migrate-create name=add_user_profile_table
```

### Go to specific migration version
```bash
make migrate-goto version=1
```

### Force migration version (use when encountering errors)
```bash
make migrate-force version=1
```

## Current Migration Files

### 000001_create_users_table
- **Up**: Create `users` table with soft delete support
  - UUID primary key
  - Unique indexes for username and email (only for active users)
  - Trigger to automatically update `updated_at`
  - Performance indexes

- **Down**: Drop `users` table and related objects

### 000002_create_refresh_tokens_table
- **Up**: Create `refresh_tokens` table
  - UUID primary key
  - Foreign key to users table
  - Unique index for token_hash
  - Indexes for common queries

- **Down**: Drop `refresh_tokens` table and constraints

## Migration File Structure

### Up Migration (`xxx_name.up.sql`)
```sql
-- Create tables, indexes, constraints
-- Add new columns
-- Insert initial data
```

### Down Migration (`xxx_name.down.sql`)
```sql
-- Reverse all changes from up migration
-- Drop tables, indexes, constraints
-- Remove columns
```

## Best Practices

1. **Always create both up and down migrations**
2. **Test migrations on database clone first**
3. **Backup database before production migration**
4. **Never modify migration files that have been applied**
5. **Use transactions in migrations (PostgreSQL automatic)** 
6. **Check migration version before deployment**

## Troubleshooting

### "no change" error when migrate-up
- Migration has already been applied
- Run `make migrate-version` to check

### "dirty database" error
- Migration was interrupted mid-way
- Run `make migrate-force version=X` to force version

### Rolling back migration in production
```bash
# Backup database first
make migrate-down
# Or rollback multiple steps
make migrate-goto version=1
```

## Migration vs AutoMigrate

| Feature | AutoMigrate | golang-migrate |
|---------|-------------|----------------|
| Safety | ❌ Not safe | ✅ Safe |
| Rollback | ❌ Not supported | ✅ Full support |
| Version Control | ❌ No tracking | ✅ Has tracking |
| Production Ready | ❌ Not recommended | ✅ Production ready |
| Custom SQL | ❌ Limited | ✅ Full control |