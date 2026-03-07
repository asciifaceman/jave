# Internal Mailing List Snippet: v0.2 Candidate Discussion

Date: 2026-03-07
List: jave-core@orcal.internal
Thread: "v0.2 scope after v0.1 release prep"
Status: Archived excerpt

## Message 1

From: Elia Mercer (Runtime)
Subject: Re: v0.2 scope after v0.1 release prep

We can keep proving parser correctness forever, but we still cannot open a file from Jave code.

If we want people to automate real jobs, we need runtime I/O and filesystem operations in the first wave.

I suggest we treat stdin/stdout/stderr channels and read/write file operations as one delivery train, not separate trains.

## Message 2

From: Miko Dane (Tooling)
Subject: Re: v0.2 scope after v0.1 release prep

Agreed. Also, we need deterministic exit behavior if `baggage run` is going to play well in CI.

Right now, users can read diagnostics, but orchestration systems need stable exit semantics.

Request:
- define a small failure taxonomy,
- map each class to an exit code,
- keep it documented and tested.

And before Runtime claims all schedule authority: Tooling needs the exit contract signed before file I/O reaches public examples. If that means splitting deliveries, we split them.

## Message 3

From: Priya Sol (Developer Learning and Documentation)
Subject: Re: v0.2 scope after v0.1 release prep

Please keep standard library naming review tightly scoped.

If we add helpers for string/number/collections/time, docs can move from "language tour" to "practical recipe".

But we should avoid opening ten utility surfaces at once. Better to ship a coherent minimum set with examples.

Also, if Runtime lands APIs without docs-ready shape, Documentation will block publication on those surfaces. We are not repeating the "ship now, explain later" cycle.

## Message 4

From: Gideon Wren (Compiler)
Subject: Re: v0.2 scope after v0.1 release prep

For cross-platform reliability, path behavior needs explicit guardrails.

If we add filesystem ops without a path story, issue volume will spike on Windows path edge cases.

Proposal for steering review:
- include path normalization/join helpers in v0.2 candidate slate,
- define diagnostics format for path-related runtime failures.

If Product wants to call path helpers "optional," Compiler will classify filesystem APIs as incomplete and defer approval. Windows path variance is not optional.

## Message 5

From: Nessa Thorne (Product)
Subject: Re: v0.2 scope after v0.1 release prep

Captured for steering packet:
- Required: I/O foundation, filesystem ops, richer stdlib, runtime error/exit semantics.
- Candidate for inclusion: path and working-directory utilities.
- Candidate for inclusion: deterministic exit contract details.

Please anchor follow-up docs to real examples we already run in repo (service planning, imports, conditions) and extend from there.

Clarification for all leads: candidate ownership is not veto ownership. Escalations go to steering; schedule threats do not.
