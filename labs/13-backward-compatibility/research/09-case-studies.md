# Case Studies in Backward Compatibility

## Case Study A — Customer Phone (1:1 to 1:N)
**Initial**: `customers (id, name, phone)`
**Requirement**: Customer can have multiple phones.
**Target**: `customer_phone_numbers (id, customer_id, phone_number)`

### Strategy:
1. **Expand**: Create `customer_phone_numbers` table. Leave `customers.phone` intact.
2. **Migrate (Code)**: Deploy app that writes to `customers.phone` (taking the primary/first phone) AND inserts into `customer_phone_numbers`.
3. **Migrate (Data)**: Backfill script copies `phone` from `customers` into `customer_phone_numbers`.
4. **Migrate (Read)**: Deploy app reading from `customer_phone_numbers`.
5. **Contract**: Stop writing to `customers.phone`. Drop column `phone`.

---

## Case Study B — CMMS Invoice Mechanics (1:1 to N:M)
**Initial**: `invoice (id, mechanic_id)`
**Requirement**: Multiple mechanics per invoice.
**Target**: `invoice_mechanics (invoice_id, mechanic_id)`

### Strategy:
1. **Expand**: Create `invoice_mechanics` table. Leave `mechanic_id` on `invoice`.
2. **Migrate (Code)**: On new invoice creation, insert into `invoice` (set `mechanic_id` to first mechanic for legacy apps) AND insert all mechanics into `invoice_mechanics`.
3. **Migrate (Data)**: Backfill script copies `mechanic_id` for existing invoices to `invoice_mechanics`.
4. **Migrate (Read)**: Update all consumers (mobile/web) to use new endpoint `/api/invoices` returning list of `mechanics` instead of singular `mechanic_id`. For API version compatibility, transform the multiple `mechanics` list down to the first item for older clients still requesting legacy formats.
5. **Contract**: Monitor logs for legacy endpoint usage. When 0, drop `invoice.mechanic_id`.

---

## Case Study C — Multi Currency (Schema Splitting)
**Initial**: `products (id, name, price)`
**Requirement**: Support multiple currencies per product.
**Target**: `product_prices (id, product_id, currency, amount)`

### Strategy:
1. **Expand**: Create `product_prices`.
2. **Migrate (Data)**: Insert existing product prices into `product_prices` assuming standard base currency (e.g., USD).
3. **Migrate (API Compatibility)**: Expose new endpoints for currency selection. For old API consumers requesting products, the API dynamically pulls `product_prices` for the default currency and flattens it to return `{"price": 100}` mimicking the old shape.
4. **Contract**: The physical database column `products.price` can be dropped once all components source data from `product_prices`, but the API contract might still output `.price` forever via API version transformers if external mobile apps cannot be forced to update.
