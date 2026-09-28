# Source Audit

## Source 1

Claimed Title: QuickCheck: An Automatic Testing Tool for Haskell (project page)
Claimed Publisher: Chalmers University of Technology (Koen Claessen & John Hughes)
URL: https://www.cse.chalmers.se/~rjmh/QuickCheck/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. URL is reachable and verified via web fetch. Page explicitly presents ICFP 2000 paper and original Haskell tool details.

Assessment:
PASS

---

## Source 2

Claimed Title: QuickCheck Manual (Hackage documentation for QuickCheck-2.14.3)
Claimed Publisher: Haskell Community / QuickCheck package maintainers
URL: https://hackage.haskell.org/package/QuickCheck-2.14.3/docs/Test-QuickCheck.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. URL is reachable and verified via web fetch. Confirms `Arbitrary`, `Gen`, `maxSuccess` default of 100, `Args`, `Result`, and `shrink`.

Assessment:
PASS

---

## Source 3

Claimed Title: Beginner's Luck: A Language for Property-Based Generators
Claimed Publisher: arXiv (academic preprint, later published at POPL 2017)
URL: https://arxiv.org/abs/1607.05443

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. arXiv ID 1607.05443 is a standard POPL 2017 paper by Hughes, Pierce, et al.

Assessment:
PASS

---

## Source 4

Claimed Title: Hypothesis Documentation (main index)
Claimed Publisher: HypothesisWorks (David R. MacIver et al.)
URL: https://hypothesis.readthedocs.io/en/latest/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Standard official Sphinx documentation for Hypothesis.

Assessment:
PASS

---

## Source 5

Claimed Title: What is Property-Based Testing?
Claimed Publisher: hypothesis.works (blog by David R. MacIver)
URL: https://hypothesis.works/articles/what-is-property-based-testing/

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative practitioner essay by creator of Hypothesis.

Assessment:
PASS

---

## Source 6

Claimed Title: How Hypothesis Works
Claimed Publisher: hypothesis.works (blog by David R. MacIver)
URL: https://hypothesis.works/articles/how-hypothesis-works/

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Explains Conjecture byte-stream fuzzer and shrinking architecture.

Assessment:
PASS

---

## Source 7

Claimed Title: Domain and distribution (Hypothesis explanation)
Claimed Publisher: Hypothesis documentation
URL: https://hypothesis.readthedocs.io/en/latest/explanation/domain.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official documentation covering domain vs distribution philosophy.

Assessment:
PASS

---

## Source 8

Claimed Title: Stateful tests (Hypothesis documentation)
Claimed Publisher: Hypothesis documentation
URL: https://hypothesis.readthedocs.io/en/latest/stateful.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official documentation for RuleBasedStateMachine, Bundles, rules, invariants.

Assessment:
PASS

---

## Source 9

Claimed Title: Proptest Book - Introduction
Claimed Publisher: AltSysrq (proptest maintainers)
URL: https://altsysrq.github.io/proptest-book/intro.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official book for Rust proptest framework.

Assessment:
PASS

---

## Source 10

Claimed Title: Shrinking Basics (Proptest tutorial)
Claimed Publisher: proptest-book
URL: https://altsysrq.github.io/proptest-book/proptest/tutorial/shrinking-basics.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official tutorial showing `ValueTree::simplify()` / `complicate()`.

Assessment:
PASS

---

## Source 11

Claimed Title: quick package - testing/quick
Claimed Publisher: Go Project (golang.org/x)
URL: https://pkg.go.dev/testing/quick

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official Go pkg docs confirming frozen status, Check/CheckEqual API, and reflection usage.

Assessment:
PASS

---

## Source 12

Claimed Title: GOPTER - GOlang Property TestER (GitHub README)
Claimed Publisher: leanovate (GitHub)
URL: https://github.com/leanovate/gopter

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official repository README for gopter Go PBT framework.

Assessment:
PASS

---

## Source 13

Claimed Title: What is Property-Based Testing? (fast-check documentation)
Claimed Publisher: fast-check.dev (Nicolas Dubien)
URL: https://fast-check.dev/docs/introduction/what-is-property-based-testing/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official JS/TS fast-check documentation.

Assessment:
PASS

---

## Source 14

Claimed Title: Why Property-Based Testing? (fast-check documentation)
Claimed Publisher: fast-check.dev
URL: https://fast-check.dev/docs/introduction/why-property-based/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Explains biased edge-case generation (0, 1, -1, `__proto__`) and hybrid testing recommendation.

Assessment:
PASS

---

## Source 15

Claimed Title: Track Record (fast-check documentation)
Claimed Publisher: fast-check.dev
URL: https://fast-check.dev/docs/introduction/track-record/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Empirical catalog of open source bugs caught by fast-check with issue links.

Assessment:
PASS

---

## Source 16

Claimed Title: The Encode/Decode Invariant (Hypothesis blog)
Claimed Publisher: hypothesis.works (David R. MacIver)
URL: https://hypothesis.works/articles/encode-decode-invariant/

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Practical essay on roundtrip invariant testing and bugs in RLE, Mercurial, Qutebrowser.

Assessment:
PASS

---

## Source 17

Claimed Title: The Hypothesis Corpus (Hugging Face dataset)
Claimed Publisher: HypothesisWorks / Liam DeVoe
URL: https://huggingface.co/datasets/HypothesisWorks/Hypothesis-Corpus-2026

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Dataset date listed as 2026 / Oct 2025; reflects recent open dataset publication. Directly provides empirical data on 28,928 tests.

Assessment:
PASS

---

## Source 18

Claimed Title: How many times will Hypothesis run my test?
Claimed Publisher: Hypothesis documentation
URL: https://hypothesis.readthedocs.io/en/latest/explanation/test-case-count.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official explanation of max_examples, search space exhaustion, flakiness check replay.

Assessment:
PASS
