#!/usr/bin/env python3
"""
Find unused translation keys by searching the codebase.
This helps identify keys that can be safely removed.
"""

import json
import re
import subprocess
from pathlib import Path
from collections import defaultdict

def get_all_keys_flat(obj, prefix=""):
    """Get all keys from nested JSON as flat dot-notation paths."""
    keys = []
    if isinstance(obj, dict):
        for key, value in obj.items():
            full_key = f"{prefix}.{key}" if prefix else key
            keys.append(full_key)
            if isinstance(value, dict):
                keys.extend(get_all_keys_flat(value, full_key))
    return keys

def search_key_usage(key, directory):
    """Search for key usage in codebase using ripgrep or grep."""
    # Escape special characters for regex
    search_pattern = key.replace('.', r'\.')
    
    try:
        # Try ripgrep first (faster)
        result = subprocess.run(
            ['rg', '-l', f"['\"]({search_pattern})['\"]", str(directory)],
            capture_output=True,
            text=True,
            timeout=5
        )
        return result.returncode == 0
    except (subprocess.TimeoutExpired, FileNotFoundError):
        try:
            # Fallback to grep
            result = subprocess.run(
                ['grep', '-r', '-l', f"['\"]\\({search_pattern}\\)['\"]", str(directory)],
                capture_output=True,
                text=True,
                timeout=5
            )
            return result.returncode == 0
        except:
            # If both fail, assume it's used to be safe
            return True

def main():
    base_dir = Path(__file__).parent.parent
    i18n_dir = base_dir / 'portal-web' / 'src' / 'i18n' / 'en'
    src_dir = base_dir / 'portal-web' / 'src'
    
    print("=" * 80)
    print("Scanning for potentially unused translation keys...")
    print("=" * 80)
    print("Note: This is a heuristic check. Manual verification recommended.")
    print()
    
    unused_keys = defaultdict(list)
    total_keys = 0
    checked_keys = 0
    
    # Check each namespace
    for json_file in sorted(i18n_dir.glob('*.json')):
        namespace = json_file.stem
        
        with open(json_file, 'r', encoding='utf-8') as f:
            data = json.load(f)
        
        keys = get_all_keys_flat(data)
        total_keys += len(keys)
        
        print(f"\nChecking {namespace} ({len(keys)} keys)...")
        
        for key in keys:
            checked_keys += 1
            # Simple heuristic: search for the key in source code
            # We search for the leaf key (last part after final dot)
            leaf_key = key.split('.')[-1]
            
            # Skip very common words that would have too many false positives
            if leaf_key in ['name', 'title', 'description', 'status', 'actions', 
                          'add', 'edit', 'delete', 'save', 'cancel', 'close']:
                continue
            
            if not search_key_usage(leaf_key, src_dir):
                unused_keys[namespace].append(key)
        
        if unused_keys[namespace]:
            print(f"  ⚠️  Found {len(unused_keys[namespace])} potentially unused keys")
        else:
            print(f"  ✅ All keys appear to be in use")
    
    print(f"\n{'=' * 80}")
    print("Summary")
    print('=' * 80)
    print(f"Total keys scanned: {total_keys}")
    print(f"Total namespaces: {len(list(i18n_dir.glob('*.json')))}")
    
    if unused_keys:
        print(f"\n⚠️  Potentially unused keys found: {sum(len(v) for v in unused_keys.values())}")
        print("\nDetailed list by namespace:")
        for namespace, keys in sorted(unused_keys.items()):
            print(f"\n{namespace}:")
            for key in sorted(keys):
                print(f"  - {key}")
        
        print("\n⚠️  IMPORTANT: These are heuristic results!")
        print("   - Keys may be used dynamically (e.g., constructed at runtime)")
        print("   - Keys may be used in commented code or tests")
        print("   - Manual verification is required before deletion")
    else:
        print("\n✅ No obviously unused keys detected!")
        print("   (Note: This doesn't guarantee all keys are used)")

if __name__ == '__main__':
    main()
