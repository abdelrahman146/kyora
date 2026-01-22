#!/usr/bin/env bash
#
# Kyora Agent OS Artifact Validator
#
# Validates ONLY the OS artifacts declared in scripts/agent-os/manifest.json.
# Does NOT scan or enforce rules on other repo artifacts.
#
# Usage:
#   ./scripts/agent-os/validate.sh
#   make agent.os.check
#
# Exit codes:
#   0 - All validations passed
#   1 - Validation failures found
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VALIDATOR="${SCRIPT_DIR}/validate.py"

# Check Python 3 is available
if ! command -v python3 &> /dev/null; then
    echo "ERROR: python3 is required but not found in PATH"
    exit 1
fi

# Check validator script exists
if [[ ! -f "${VALIDATOR}" ]]; then
    echo "ERROR: Validator script not found: ${VALIDATOR}"
    exit 1
fi

# Run the Python validator
exec python3 "${VALIDATOR}" "$@"
