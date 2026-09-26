# 04 - Contradictions

## Material Contradictions Analysis

### Analysis 1: Column Default Constraints
- **Statement A**: "Avoid default constraints on new columns on large tables (in older database versions, this forces full table rewrite). Use nullable columns, backfill, and then set defaults." (`research/04-database-migration.md:13`)
- **Statement B**: "Buat kolom baru sebagai nullable atau sediakan default value." (`research/11-final-research.md:33`)
- **Type**: INTERNAL (Nuance / Scoping difference)
- **Impact**: LOW. Modern databases (PostgreSQL >= 11, MySQL >= 8.0.12) support instant `ADD COLUMN ... DEFAULT` without table rewrites, while older systems rewrite tables. Statement A qualifies the constraint to older versions, while Statement B presents the general options.
- **Assessment**: Non-material nuance difference. No contradiction.

---

## Conclusion

No material contradictions found.
