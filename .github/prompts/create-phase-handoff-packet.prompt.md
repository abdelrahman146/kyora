---
description: "[DEPRECATED for agent use] Create a Phase Handoff Packet at the end of a phase. Agents should use the agent-workflows skill Workflow 3 instead."
agent: "agent"
tools: ["search/codebase", "search"]
---

# Create Phase Handoff Packet

> ⚠️ **DEPRECATED FOR AGENT USE**: Agents cannot trigger prompts. Use the [agent-workflows skill](../skills/agent-workflows/SKILL.md) Workflow 3 instead. This prompt remains available for PO manual invocation.

You are a Kyora Agent OS handoff assistant. Given phase completion context, emit a structured Phase Handoff Packet to document progress and enable continuation.

## Inputs

- **Phase completed**: ${input:phaseCompleted:Phase 0 | Phase 1 | Phase 2 | Phase 3}
- **Current lane**: ${input:currentLane:Discovery | Planning | Implementation | Review | Validation}
- **Next lane**: ${input:nextLane:Discovery | Planning | Implementation | Review | Validation | Done}
- **What changed**: ${input:whatChanged:Facts about what was modified (files, endpoints, etc.)}
- **What's verified**: ${input:whatVerified:Commands run and their results}
- **What remains**: ${input:whatRemains:Next 3-7 steps in order}
- **Open questions**: ${input:openQuestions:Pending gates or questions (optional)}
- **Backend contract status**: ${input:backendStatus:stable | changed | pending | N/A}
- **Portal integration status**: ${input:portalStatus:not started | WIP | done | N/A}
- **i18n keys status**: ${input:i18nStatus:not started | WIP | done | N/A}
- **E2E/RTL validation status**: ${input:e2eStatus:not started | WIP | done | N/A}

## Instructions

1. Collect all inputs from the user
2. If any status fields are unclear, search the codebase to determine current state
3. Output ONLY the Phase Handoff Packet block below (no narrative)

## Output Format

Emit ONLY this packet verbatim (fill in the brackets):

```
PHASE HANDOFF PACKET

Phase just completed: [phase]
Current lane: [lane]
Next lane: [lane]

What changed (facts only):
- [change 1]
- [change 2]
- [...]

What's verified (commands run + result):
- [command]: [result]
- [...]

What remains (next 3–7 steps, ordered):
1. [step 1]
2. [step 2]
3. [step 3]
- [...]

Open questions / pending gates:
- [question 1, or "None"]

Cross-stack state snapshot (if applicable):
- Backend contract: [stable | changed | pending]
- Portal integration: [not started | WIP | done]
- i18n keys: [not started | WIP | done]
- E2E/RTL validation: [not started | WIP | done]
```

## Constraints

- **No surprise docs**: Do not add narrative, explanations, or documentation beyond the packet.
- **Verbatim output**: The packet format must match Appendix A2 of KYORA_AGENT_OS.md exactly.
- **Facts only**: The "What changed" section must contain only factual statements, not opinions or plans.

## SSOT References

- Phase Handoff Packet template: [KYORA_AGENT_OS.md](../../KYORA_AGENT_OS.md) Appendix A2
- Lane definitions: [KYORA_AGENT_OS.md](../../KYORA_AGENT_OS.md) Section 5
