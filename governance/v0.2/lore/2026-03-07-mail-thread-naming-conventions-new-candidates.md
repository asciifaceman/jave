# Internal Mail Thread: Naming Conventions For Newly Approved v0.2 Candidates

Date: 2026-03-07
List: jave-governance@orcal.internal
Thread: "Naming freeze pressure: env/logging/timing/networking"
Status: Archived excerpt

## Message 1

From: Priya Sol (Developer Learning and Documentation)
Subject: Naming freeze pressure: env/logging/timing/networking

If naming is vague, docs absorb the fallout for months.

Proposal:
- environment reads under a stable namespace,
- logging helpers with explicit output-channel semantics,
- timing helpers that do not pretend to be full profilers,
- networking names that reflect baseline scope only.

## Message 2

From: Elia Mercer (Runtime)
Subject: Re: Naming freeze pressure: env/logging/timing/networking

Agreed on baseline scope, not on over-abstract naming.

Runtime will not ship ten layers of wrappers to satisfy naming aesthetics. We need direct operators, clear signatures, and deterministic failure behavior.

## Message 3

From: Miko Dane (Tooling)
Subject: Re: Naming freeze pressure: env/logging/timing/networking

Tooling position: ambiguous naming is effectively an ABI hazard for automation.

If helper semantics are unstable, docs and CI scripts become untrustworthy. I will block release-readiness claims until naming and exit behavior are pinned.

## Message 4

From: Nessa Thorne (Product)
Subject: Re: Naming freeze pressure: env/logging/timing/networking

Directive:
- publish one naming-convention proposal record,
- include at least one rejected alternative per candidate area,
- route disputes through steering packet updates, not side-channel commitments.

No team has unilateral naming authority once candidate status is approved.
