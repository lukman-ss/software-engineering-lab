# Research Gap Analysis

Target Lab: `labs/35-websocket-and-sse`

---

## Gap 1

Type:
WEAK_SOURCE

Severity:
MEDIUM

Location:
`03-evidence.md` (Evidence 12, 13), `05-report.md` (Finding 11)

Problem:
Evidence 12 and 13 cite the "Lab specification (Senior Software Engineer Lab #35)" as their primary source rather than authoritative Linux kernel documentation, POSIX standards, or empirical socket scaling publications (e.g. C10K/C10M benchmarks).

Required Revision:
In corporate/production documentation, replace internal spec citations with Linux networking docs (`epoll(7)`, `sysctl` `fs.file-max`, socket buffer memory tuning).

Can Be Approved Without Fix:
YES (The report explicitly marks confidence as MEDIUM and identifies lack of empirical benchmarks in Limitations and Open Questions).

---

## Gap 2

Type:
OUTDATED_SOURCE

Severity:
LOW

Location:
`02-sources.md` (Source 6), `05-report.md` (Finding 7, 8)

Problem:
RFC 7540 is cited for HTTP/2 without noting that RFC 9113 obsoleted it in June 2022.

Required Revision:
Update references to cite RFC 9113 (or RFC 7540 / RFC 9113).

Can Be Approved Without Fix:
YES (Core stream multiplexing semantics remain identical).

---

## Gap 3

Type:
SCOPE_ERROR

Severity:
LOW

Location:
`05-report.md` (Finding 8, line 138)

Problem:
Finding 8 text discusses RFC 7540 but provides the URL for RFC 7541 (`https://datatracker.ietf.org/doc/html/rfc7541`).

Required Revision:
Correct the URL link to `https://datatracker.ietf.org/doc/html/rfc7540` or `https://datatracker.ietf.org/doc/html/rfc9113`.

Can Be Approved Without Fix:
YES (Minor editorial URL error).

---

## Gap 4

Type:
MISSING_CASE

Severity:
LOW

Location:
`06-open-questions.md` (Question 4, Question 7)

Problem:
HTTP/3 (QUIC / WebTransport) impact on bidirectional vs streaming protocols is flagged as an open question but not deeply explored in the report.

Required Revision:
Add a brief note in subsequent lab content clarifying that HTTP/3 and WebTransport represent future alternatives.

Can Be Approved Without Fix:
YES.
