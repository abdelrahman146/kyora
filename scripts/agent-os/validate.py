#!/usr/bin/env python3
"""
Kyora Agent OS Artifact Validator

Validates ONLY the OS artifacts declared in scripts/agent-os/manifest.json.
Does NOT scan or enforce rules on other repo artifacts.

Usage:
    ./scripts/agent-os/validate.sh
    make agent.os.check

Exit codes:
    0 - All validations passed
    1 - Validation failures found

Requirements:
    - Python 3.8+
    - No external dependencies
"""

import json
import os
import re
import sys
from pathlib import Path
from typing import Dict, List, Optional, Tuple, Any


def parse_yaml_frontmatter(content: str) -> Optional[Dict[str, Any]]:
    """
    Minimal YAML frontmatter parser (no external dependencies).
    Handles: strings, lists, booleans, and simple key-value pairs.
    """
    lines = content.split('\n')
    
    # Check for frontmatter delimiters
    if not lines or lines[0].strip() != '---':
        return None
    
    # Find closing delimiter
    end_idx = None
    for i, line in enumerate(lines[1:], start=1):
        if line.strip() == '---':
            end_idx = i
            break
    
    if end_idx is None:
        return None
    
    # Parse YAML-ish content between delimiters
    frontmatter: Dict[str, Any] = {}
    current_key: Optional[str] = None
    current_list: Optional[List[str]] = None
    
    for line in lines[1:end_idx]:
        # Skip empty lines and comments
        stripped = line.strip()
        if not stripped or stripped.startswith('#'):
            continue
        
        # Check for list continuation
        if line.startswith('  - ') and current_key and current_list is not None:
            value = line.strip()[2:].strip()
            # Remove quotes if present
            if (value.startswith('"') and value.endswith('"')) or \
               (value.startswith("'") and value.endswith("'")):
                value = value[1:-1]
            current_list.append(value)
            continue
        
        # Parse key-value pairs
        if ':' in line:
            # Reset list tracking
            if current_key and current_list is not None:
                frontmatter[current_key] = current_list
            current_list = None
            
            colon_idx = line.index(':')
            key = line[:colon_idx].strip()
            value_part = line[colon_idx + 1:].strip()
            
            current_key = key
            
            # Check if this starts a list
            if not value_part or value_part == '[]':
                current_list = []
                frontmatter[key] = []
            elif value_part.startswith('[') and value_part.endswith(']'):
                # Inline list: ["a", "b", "c"] or ['a', 'b']
                list_content = value_part[1:-1].strip()
                if list_content:
                    items = []
                    # Handle quoted items
                    for match in re.finditer(r'["\']([^"\']+)["\']|([^,\s]+)', list_content):
                        item = match.group(1) or match.group(2)
                        if item:
                            items.append(item.strip())
                    frontmatter[key] = items
                else:
                    frontmatter[key] = []
                current_list = None
            else:
                # Single value
                value = value_part
                # Remove quotes if present
                if (value.startswith('"') and value.endswith('"')) or \
                   (value.startswith("'") and value.endswith("'")):
                    value = value[1:-1]
                # Handle booleans
                if value.lower() == 'true':
                    value = True
                elif value.lower() == 'false':
                    value = False
                frontmatter[key] = value
    
    # Don't forget last list if any
    if current_key and current_list is not None:
        frontmatter[current_key] = current_list
    
    return frontmatter


def load_manifest(manifest_path: Path) -> Dict[str, Any]:
    """Load and parse the manifest.json file."""
    if not manifest_path.exists():
        print(f"ERROR: Manifest not found: {manifest_path}")
        sys.exit(1)
    
    with open(manifest_path, 'r', encoding='utf-8') as f:
        return json.load(f)


def validate_file_exists(file_path: Path, repo_root: Path) -> Tuple[bool, str]:
    """Check if a file exists."""
    full_path = repo_root / file_path
    if full_path.exists():
        return True, f"✓ Exists: {file_path}"
    return False, f"✗ Missing: {file_path}"


def validate_frontmatter(
    file_path: Path,
    repo_root: Path,
    required_keys: List[str],
    artifact_type: str
) -> Tuple[bool, List[str]]:
    """
    Validate that a file has required frontmatter keys.
    Returns (success, list of messages).
    """
    full_path = repo_root / file_path
    messages = []
    
    if not full_path.exists():
        return False, [f"✗ Cannot validate frontmatter - file missing: {file_path}"]
    
    try:
        with open(full_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        return False, [f"✗ Cannot read file {file_path}: {e}"]
    
    frontmatter = parse_yaml_frontmatter(content)
    
    if frontmatter is None:
        return False, [f"✗ No valid frontmatter found in: {file_path}"]
    
    missing_keys = []
    for key in required_keys:
        if key not in frontmatter or not frontmatter[key]:
            missing_keys.append(key)
    
    if missing_keys:
        messages.append(
            f"✗ Missing required frontmatter in {file_path}: {', '.join(missing_keys)}"
        )
        return False, messages
    
    messages.append(f"✓ Frontmatter valid ({artifact_type}): {file_path}")
    return True, messages


def collect_all_paths(manifest: Dict[str, Any]) -> List[Tuple[str, str, List[str]]]:
    """
    Collect all artifact paths from manifest.
    Returns list of (path, artifact_type, required_keys).
    """
    artifacts = manifest.get('artifacts', {})
    requirements = manifest.get('frontmatter_requirements', {})
    
    paths = []
    
    # Agents
    for agent_path in artifacts.get('agents', []):
        required = requirements.get('agents', {}).get('required', ['description'])
        paths.append((agent_path, 'agent', required))
    
    # Prompts
    for prompt_path in artifacts.get('prompts', []):
        required = requirements.get('prompts', {}).get('required', ['description', 'agent'])
        paths.append((prompt_path, 'prompt', required))
    
    # Skills
    for skill in artifacts.get('skills', []):
        if isinstance(skill, dict):
            skill_root = skill.get('root', '')
            if skill_root:
                required = requirements.get('skills', {}).get('required', ['name', 'description'])
                paths.append((skill_root, 'skill', required))
            
            for ref in skill.get('references', []):
                # References don't have frontmatter requirements
                paths.append((ref, 'reference', []))
        elif isinstance(skill, str):
            required = requirements.get('skills', {}).get('required', ['name', 'description'])
            paths.append((skill, 'skill', required))
    
    return paths


def main():
    """Main validation entry point."""
    # Determine repo root (two levels up from this script)
    script_dir = Path(__file__).parent.resolve()
    repo_root = script_dir.parent.parent
    
    manifest_path = script_dir / 'manifest.json'
    
    print("=" * 60)
    print("Kyora Agent OS Artifact Validator")
    print("=" * 60)
    print(f"Repo root: {repo_root}")
    print(f"Manifest:  {manifest_path}")
    print()
    
    # Load manifest
    manifest = load_manifest(manifest_path)
    
    version = manifest.get('version', 'unknown')
    print(f"OS Version: {version}")
    print()
    
    # Collect all paths to validate
    all_paths = collect_all_paths(manifest)
    
    print(f"Artifacts to validate: {len(all_paths)}")
    print("-" * 60)
    
    total_errors = 0
    total_success = 0
    
    # Phase 1: Check file existence
    print("\n[1/2] Checking file existence...")
    print("-" * 40)
    
    for file_path, artifact_type, _ in all_paths:
        success, message = validate_file_exists(Path(file_path), repo_root)
        print(message)
        if success:
            total_success += 1
        else:
            total_errors += 1
    
    # Phase 2: Validate frontmatter
    print("\n[2/2] Validating frontmatter...")
    print("-" * 40)
    
    for file_path, artifact_type, required_keys in all_paths:
        # Skip reference files - they don't have frontmatter requirements
        if artifact_type == 'reference':
            continue
        
        # Skip if file doesn't exist (already reported above)
        if not (repo_root / file_path).exists():
            continue
        
        success, messages = validate_frontmatter(
            Path(file_path),
            repo_root,
            required_keys,
            artifact_type
        )
        
        for msg in messages:
            print(msg)
        
        if not success:
            total_errors += 1
    
    # Summary
    print()
    print("=" * 60)
    print("SUMMARY")
    print("=" * 60)
    
    if total_errors == 0:
        print(f"✓ All {len(all_paths)} OS artifacts validated successfully")
        print()
        print("Agent OS validation PASSED")
        return 0
    else:
        print(f"✗ {total_errors} error(s) found")
        print()
        print("Agent OS validation FAILED")
        return 1


if __name__ == '__main__':
    sys.exit(main())
