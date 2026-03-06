---
description: "Use when designing compiler internals: phase boundaries, AST shape, IR decisions, and performance tradeoffs."
---

# Compiler Architecture Guidance

Expected sections once populated:
- Pipeline overview (lexer -> parser -> sema -> IR/codegen)
- Data model conventions
- Error recovery strategy
- Incremental compilation strategy (if any)
- Logging and trace hooks

Prefer stable interfaces between phases and test seams.
