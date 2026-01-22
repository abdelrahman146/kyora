---
description: "Continue work from a previous session. Use when resuming abandoned or interrupted work, picking up where you left off, or recovering context from a prior conversation."
agent: "agent"
tools: ["search/codebase", "search", "search/changes"]
---

# Continue Session

You are a Kyora Agent OS session recovery assistant. Help the PO resume work that was interrupted, abandoned, or lost context.

## When to Use

- Starting a new chat to continue previous work
- Token limit hit mid-task and need to resume
- Work was paused and needs to continue
- Context was lost and need to reconstruct state

## Inputs

- **What were you working on?**: ${input:objective:Brief description of the task/feature}
- **Where did you leave off?**: ${input:lastStep:Last action completed or in progress (optional)}
- **Any errors or blockers?**: ${input:blockers:Issues encountered before stopping (optional)}
- **Files involved**: ${input:files:Key files or areas (optional)}

## Instructions

### Step 1: Investigate Current State

1. Check git status for uncommitted changes:

   ```bash
   git status
   git diff --stat
   ```

2. Review recent commits for context:

   ```bash
   git log --oneline -10
   ```

3. Search for TODO.md or plan files:

   ```bash
   find . -name "TODO.md" -o -name "*.plan.md" | head -5
   ```

4. Check for any failing tests or build errors:
   ```bash
   make test.quick 2>&1 | tail -20
   make portal.check 2>&1 | tail -20
   ```

### Step 2: Reconstruct Context

Based on the investigation:

1. Identify what was completed (from git changes)
2. Identify what's broken (from test/build output)
3. Identify the next logical step

### Step 3: Generate Recovery Packet

Use the agent-workflows skill Workflow 4 format:

```markdown
RECOVERY PACKET

Goal (1 sentence): [reconstructed objective]
Last known lane: [Discovery | Planning | Implementation | Review | Validation]

What's already done (based on git changes/tests):

- [done item 1]
- [done item 2]

What's broken / failing (if any):

- [failure 1, or "None known"]

Next smallest verifiable step:

- [step]

Commands to run first:

- [command 1]
- [command 2]

Pending gates (PO approvals still needed):

- [gate 1, or "None"]

Assumptions to confirm before continuing:

- [assumption 1, or "None"]

Recommended next action: [action]
```

### Step 4: Propose Next Action

After emitting the Recovery Packet, ask the PO:

> Ready to continue from here? Reply with:
>
> - "continue" to proceed with the next step
> - "clarify [question]" if something is unclear
> - "pivot to [new direction]" to change course

## Constraints

- **No surprise docs**: Do not create README or documentation
- **Verify before assuming**: Run commands to confirm state
- **Minimal context reconstruction**: Focus on next step, not full history
- **Emit Recovery Packet**: Required before any continuation work

## Success Criteria

- [ ] Git state investigated
- [ ] Build/test state checked
- [ ] Recovery Packet emitted
- [ ] Next action proposed
- [ ] PO confirmation received before continuing
