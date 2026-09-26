# Research Gap Analysis: Zero-Downtime Deployment

**Target Lab:** `labs/20-zero-downtime-deployment`  
**Research Run:** `2026-09-26-zero-downtime-deployment`

---

## Gap 1

Type: SCOPE_ERROR  
Severity: MEDIUM  
Location: `05-report.md: Finding 10`  
Problem: `ALTER TABLE ... ADD COLUMN` constant default optimization in PostgreSQL is presented without prominent warning that volatile defaults (`gen_random_uuid()`, `clock_timestamp()`) trigger full table rewrites and exclusive table locks.  
Required Revision: Emphasize in deployment guidelines that column defaults must be strictly non-volatile constants or nullable columns without defaults.  
Can Be Approved Without Fix: YES  

---

## Gap 2

Type: IMPLEMENTATION_GAP  
Severity: MEDIUM  
Location: `05-report.md: Finding 9` & `06-open-questions.md: #1, #2`  
Problem: The report highlights Laravel Horizon's graceful shutdown via `horizon:terminate` and Supervisor `stopwaitsecs`, but plain PHP-FPM web workers do not natively handle SIGTERM without custom configuration (`process_control_timeout` in `php-fpm.conf`) and a preStop hook in Kubernetes.  
Required Revision: Note the distinction between Queue worker shutdown (Horizon) and HTTP worker shutdown (PHP-FPM) in future engineering designs.  
Can Be Approved Without Fix: YES  

---

## Gap 3

Type: UNVERIFIED_CLAIM  
Severity: LOW  
Location: `06-open-questions.md: Weak Evidence Areas #1`  
Problem: Lack of empirical load-testing metrics demonstrating zero dropped connections during simultaneous DDL and rolling update in this specific stack.  
Required Revision: Empirical testing will be validated during the engineering execution stage.  
Can Be Approved Without Fix: YES  
