# Changes Made

## Revision 1

**Audit Issue**: 
CRITICAL — Missing Research Findings document (Gap 1)

**Files Changed**:
- `research/02-findings.md` — **CREATE**

**Action**:
- Created comprehensive findings document answering all 10 research questions
- Provided explicit mathematical formulas with variable definitions
- Included authoritative source citations (RFCs, Little 1961, AWS, Google SRE)
- Detailed evidence extraction with quotes from verified sources

**Verification**:
- Source checked: RFC 6585, RFC 9110, RFC 2697, RFC 6598 — all verified reachable
- Source checked: Little (1961) DOI — verified
- Source checked: AWS Architecture Blog — verified
- Source checked: Google SRE Book — verified

**Status**: RESOLVED

---

## Revision 2

**Audit Issue**: 
CRITICAL — Little's Law Misattribution Risk (Claim 3, Gap 2)

**Files Changed**:
- `research/02-findings.md` §Claim 3

**Action**:
- Explicitly separated Little's Law ($L = \lambda W$) from deterministic queue buildup ($\Delta Q = (r_{in} - r_{out}) \cdot \Delta t$)
- Provided proper mathematical formulation with variable definitions
- Cited Little (1961) original paper with DOI link
- Added note that deterministic formula is a transient approximation, not Little's Law
- Added Google SRE Book citation supporting queue memory exhaustion risks

**Verification**:
- Formula verification: Little's Law requires stable system, steady-state; fluid model is transient
- Source checked: Operations Research journal DOI — responds correctly
- Original claim in plan was incorrect; corrected with precise distinction

**Status**: RESOLVED

---

## Revision 3

**Audit Issue**: 
MEDIUM — Unsubstantiated Queue Age Claim (Grant 3)

**Files Changed**:
- `research/02-findings.md` §Claim 7
- `research/03-sources.md` (Source #8, #9)

**Action**:
- Added AWS SQS CloudWatch metrics: `ApproximateAgeOfOldestMessage`
- Added Apache Kafka documentation: Consumer Lag monitoring
- Referenced vendor-specific implementations for queue age metric
- Clarified that queue age is superior early warning but both metrics should be monitored

**Verification**:
- AWS documentation confirmed for both metrics
- Kafka monitoring documentation confirms consumer lag is standard practice

**Status**: RESOLVED

---

## Revision 4

**Audit Issue**: 
MEDIUM — Generalization of Autoscaling Failure Modes (Grant 4)

**Files Changed**:
- `research/02-findings.md` §Claim 10
- `research/05-report.md` §7

**Action**:
- Added specific architectural framing from Google SRE Book §22.2
- Documented cascade failure pattern: upstream scaling → downstream saturation → retries → system crash
- Provided remediation recommendations (circuit breakers, bulkheads)

**Verification**:
- Google SRE Book cited for specific failure pattern
- Original claim was too general; now explicitly conditioned on asymmetric scaling

**Status**: RESOLVED

---

## Revision 5

**Audit Issue**: INTERNAL CONTRADICTION — Little's Law Formula

**Files Changed**:
- `research/04-contradictions.md`

**Action**:
- Documented resolution of contradiction between incorrect formula and risk flag
- Provided explicit correction in findings document

**Verification**:
- Original plan claimed formula DID represent Little's Law
- Risk flag correctly identified misattribution
- Findings document now provides correct distinction

**Status**: RESOLVED

---

## Revision 6

**Audit Issue**: INTERNAL — Irrelevant RFC Citations

**Files Changed**:
- `research/03-sources.md`

**Action**:
- Removed RFC 8305 (Happy Eyeballs — not relevant to rate limiting)
- Removed RFC 5321 (SMTP — not relevant to HTTP rate limiting)
- Maintained RFC 6585, RFC 9110, RFC 2697, RFC 6598 as correct primary sources

**Verification**:
- RFC 8305 title matches "Happy Eyeballs Algorithm" (IPv6 connection establishment)
- RFC 5321 is SMTP protocol specification (email transport)
- Neither relates to HTTP rate limiting

**Status**: RESOLVED

---

## Revision Checklist

- [x] All CRITICAL issues addressed
- [x] All HIGH issues addressed or explicitly unresolved
- [x] Unsupported major claims fixed
- [x] Weak major sources improved
- [x] Contradictions resolved or properly documented
- [x] Implementation/documentation mismatch fixed (research artifacts)
- [x] Failing tests fixed or blocker documented (N/A - Pipeline Override)
- [x] README does not apply - research only revision
- [x] Source list updated
- [x] Revision log written
- [x] Validation executed via source verification
- [x] Result marked READY_FOR_RESEARCH_REAUDIT