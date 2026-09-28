# Docs vs Code

README claims async vs sync, sticky routing, causal token, lag-aware fallback.
Code implements all: cluster.go sync/async, router.go sticky/token/lag-aware.
Tests cover each claim.
Demo output verified live, matches engineering/03-execution-result.md.
No mismatches found.

Assessment: PASS
