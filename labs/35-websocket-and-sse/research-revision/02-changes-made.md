# Changes Made

## Revision 1

Audit Issue:
LOW — URL typo in `05-report.md` Finding 8 (`rfc7541` instead of `rfc7540` / `rfc9113`).

Files Changed:
- research/05-report.md

Action:
- Corrected URL from `https://datatracker.ietf.org/doc/html/rfc7541` to `https://datatracker.ietf.org/doc/html/rfc9113`.
- Updated Finding 8 source citation to RFC 9113 (obsoletes RFC 7540).

Verification:
- RFC 9113 confirmed as the current HTTP/2 standard (June 2022).
- URL is valid.

Status:
RESOLVED

---

## Revision 2

Audit Issue:
LOW — RFC 7540 cited without noting RFC 9113 obsoletion (June 2022).

Files Changed:
- research/02-sources.md
- research/05-report.md
- research/04-contradictions.md

Action:
- Added obsolescence note to Source 6 in `02-sources.md`.
- Finding 7 and Finding 8 now cite RFC 9113 as the current standard, with RFC 7540 noted as the predecessor.
- Added currency note to `04-contradictions.md`.

Verification:
- RFC 9113 Abstract: "This document obsoletes RFCs 7540 and 8740."
- Stream multiplexing and connection-header rules remain unchanged.

Status:
RESOLVED

---

## Revision 3

Audit Issue:
MEDIUM — Evidence 12 and 13 cited internal lab specification as the primary source for 100k scaling claims.

Files Changed:
- research/03-evidence.md
- research/05-report.md
- research/02-sources.md

Action:
- Replaced lab-spec citations in Evidence 12 with Linux kernel documentation (`epoll(7)`, `fs.file-max`, `ulimit -n`).
- Replaced lab-spec citations in Evidence 13 with Redis Pub/Sub documentation.
- Updated Finding 11 evidence and sources accordingly.
- Added Source 7 (epoll(7) / sysctl) to `02-sources.md`.
- Renumbered NGINX source to Source 8.

Verification:
- `https://man7.org/linux/man-pages/man7/epoll.7.html` is official Linux man page.
- `https://redis.io/docs/manual/pubsub/` is official Redis documentation.
- Confidence remains MEDIUM (no empirical 100k benchmark).

Status:
RESOLVED

---

## Revision 4

Audit Issue:
LOW — HTTP/3 (QUIC / WebTransport) flagged as an open question without a brief note in the report.

Files Changed:
- research/06-open-questions.md

Action:
- Added a clarifying sentence to Question 4 noting HTTP/3 and WebTransport as future alternatives.

Verification:
- Open question retained; not elevated to a finding because it was not deeply researched.

Status:
RESOLVED
