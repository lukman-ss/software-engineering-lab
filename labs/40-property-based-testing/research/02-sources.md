# Sources

## Source 1

Title: QuickCheck: An Automatic Testing Tool for Haskell (project page)
Publisher: Chalmers University of Technology (Koen Claessen & John Hughes)
URL: https://www.cse.chalmers.se/~rjmh/QuickCheck/
Published: 2000 (original paper ICFP 2000); page last updated unknown
Accessed: 2026-09-28
Source Tier: Tier 1 (original authors, primary source)
Relevance: Foundational QuickCheck tool and paper archive; describes the original property-based testing system for Haskell

## Source 2

Title: QuickCheck Manual (Hackage documentation for QuickCheck-2.14.3)
Publisher: Haskell Community / QuickCheck package maintainers
URL: https://hackage.haskell.org/package/QuickCheck-2.14.3/docs/Test-QuickCheck.html
Published: Package version 2.14.3 (current as of 2026)
Accessed: 2026-09-28
Source Tier: Tier 1 (official package documentation)
Relevance: Authoritative reference for QuickCheck API: Arbitrary typeclass, Gen monad, shrinking combinators, property combinators, Args/Result types

## Source 3

Title: Beginner's Luck: A Language for Property-Based Generators
Publisher: arXiv (academic preprint, later published at POPL 2017)
URL: https://arxiv.org/abs/1607.05443
Published: 2016 (v1), last revised 2019 (v3)
Accessed: 2026-09-28
Source Tier: Tier 1 (academic paper)
Relevance: Formal language (Luck) for writing property-based generators from predicates; addresses generator difficulty; authors include John Hughes and Benjamin Pierce

## Source 4

Title: Hypothesis Documentation (main index)
Publisher: HypothesisWorks (David R. MacIver et al.)
URL: https://hypothesis.readthedocs.io/en/latest/
Published: Continuously updated; current as of 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official library documentation)
Relevance: Primary reference for Hypothesis, the most widely used PBT library for Python; covers strategies, stateful testing, settings, and internals

## Source 5

Title: What is Property-Based Testing?
Publisher: hypothesis.works (blog by David R. MacIver)
URL: https://hypothesis.works/articles/what-is-property-based-testing/
Published: May 14, 2016
Accessed: 2026-09-28
Source Tier: Tier 2 (authoritative practitioner explanation by library creator)
Relevance: Definitional essay distinguishing PBT from fuzzing; argues PBT is about constructing tests whose failures reveal problems not findable by direct fuzzing

## Source 6

Title: How Hypothesis Works
Publisher: hypothesis.works (blog by David R. MacIver)
URL: https://hypothesis.works/articles/how-hypothesis-works/
Published: December 10, 2016
Accessed: 2026-09-28
Source Tier: Tier 2 (internal architecture explanation by library creator)
Relevance: Describes Conjecture byte-stream fuzzer, strategy layer, and compositional shrinking design; explains why byte-stream shrinking outperforms type-based shrinking

## Source 7

Title: Domain and distribution (Hypothesis explanation)
Publisher: Hypothesis documentation
URL: https://hypothesis.readthedocs.io/en/latest/explanation/domain.html
Published: Current as of 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official documentation)
Relevance: Explains Hypothesis philosophy: library controls distribution, user controls domain; discusses swarm testing, alternative backends (hypofuzz), and why distribution control is withheld from users

## Source 8

Title: Stateful tests (Hypothesis documentation)
Publisher: Hypothesis documentation
URL: https://hypothesis.readthedocs.io/en/latest/stateful.html
Published: Current as of 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official documentation)
Relevance: Describes RuleBasedStateMachine, bundles, rules, invariants, preconditions; key reference for stateful/state-machine PBT

## Source 9

Title: Proptest Book - Introduction
Publisher: AltSysrq (proptest maintainers)
URL: https://altsysrq.github.io/proptest-book/intro.html
Published: Current as of 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official library documentation for Rust)
Relevance: Rust PBT framework inspired by Hypothesis; explains property testing basics, shrinking, and compares to QuickCheck

## Source 10

Title: Shrinking Basics (Proptest tutorial)
Publisher: proptest-book
URL: https://altsysrq.github.io/proptest-book/proptest/tutorial/shrinking-basics.html
Published: Current as of 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official tutorial)
Relevance: Demonstrates ValueTree::simplify() / complicate() binary search over input space; shows concrete shrinking example finding boundary condition (501)

## Source 11

Title: quick package - testing/quick
Publisher: Go Project (golang.org/x)
URL: https://pkg.go.dev/testing/quick
Published: Go 1.27.1 (Sep 2026); package marked as frozen
Accessed: 2026-09-28
Source Tier: Tier 1 (official Go standard library)
Relevance: Go's built-in PBT utility; implements Check, CheckEqual, Value generation via reflect; noted as frozen/not accepting new features

## Source 12

Title: GOPTER - GOlang Property TestER (GitHub README)
Publisher: leanovate (GitHub)
URL: https://github.com/leanovate/gopter
Published: Active development; current as of 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (primary source for Go PBT library)
Relevance: Go PBT library inspired by ScalaCheck/QuickCheck; features: tighter generator control, shrinkers, regex generators, stateful tests support

## Source 13

Title: What is Property-Based Testing? (fast-check documentation)
Publisher: fast-check.dev (Nicolas Dubien)
URL: https://fast-check.dev/docs/introduction/what-is-property-based-testing/
Published: Last updated Feb 7, 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official library documentation for JavaScript/TypeScript)
Relevance: Defines three core features: random input generation, multiple runs per test, counterexample shrinking; contrasts with example-based testing

## Source 14

Title: Why Property-Based Testing? (fast-check documentation)
Publisher: fast-check.dev
URL: https://fast-check.dev/docs/introduction/why-property-based/
Published: Last updated Feb 7, 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official documentation)
Relevance: Argues PBT detects edge cases better, is "designed for bugs" (upweights boundaries, duplicates, security-sensitive values like __proto__), documents code behavior, improves maintainability; recommends hybrid approach with example-based tests

## Source 15

Title: Track Record (fast-check documentation)
Publisher: fast-check.dev
URL: https://fast-check.dev/docs/introduction/track-record/
Published: Last updated Feb 7, 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official documentation, empirical evidence)
Relevance: Catalog of real bugs found by fast-check in production projects: underscore.js, jest, react, jasmine, js-yaml, query-string, left-pad, yaml, jsonwebtoken, numpy; includes specific failing inputs and CVE recoveries

## Source 16

Title: The Encode/Decode Invariant (Hypothesis blog)
Publisher: hypothesis.works (David R. MacIver)
URL: https://hypothesis.works/articles/encode-decode-invariant/
Published: April 16, 2016
Accessed: 2026-09-28
Source Tier: Tier 2 (practitioner article with concrete examples)
Relevance: Demonstrates roundtrip invariant (Decode(Encode(data)) == data); shows bug discovery in run-length encoding (empty string, missing count reset); cites Mercurial UTF8b bugs and Qutebrowser JS escaping bugs

## Source 17

Title: The Hypothesis Corpus (Hugging Face dataset)
Publisher: HypothesisWorks / Liam DeVoe
URL: https://huggingface.co/datasets/HypothesisWorks/Hypothesis-Corpus-2026
Published: October 2025
Accessed: 2026-09-28
Source Tier: Tier 1 (original dataset from PBT community)
Relevance: 28,928 Hypothesis tests across 1,529 repositories; empirical data on test execution, coverage, failure rates; largest known dataset of real-world PBT usage

## Source 18

Title: How many times will Hypothesis run my test?
Publisher: Hypothesis documentation
URL: https://hypothesis.readthedocs.io/en/latest/explanation/test-case-count.html
Published: Current as of 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (official documentation)
Relevance: Explains max_examples semantics, search space exhaustion, assume()/filter() retry behavior, entropy cap, flakiness verification replay