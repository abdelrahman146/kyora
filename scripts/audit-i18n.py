#!/usr/bin/env python3
"""
Audit i18n translation files for:
1. Duplicate keys within each file
2. Missing keys across locales (en/ar parity)
3. Structural differences between locales
"""

import json
import sys
from pathlib import Path
from collections import defaultdict

def find_duplicate_keys(obj, path="", duplicates=None):
    """Recursively find duplicate keys in a JSON object."""
    if duplicates is None:
        duplicates = []
    
    if isinstance(obj, dict):
        seen_keys = {}
        for key, value in obj.items():
            full_path = f"{path}.{key}" if path else key
            if key in seen_keys:
                duplicates.append(f"Duplicate key '{key}' at path: {path}")
            seen_keys[key] = True
            find_duplicate_keys(value, full_path, duplicates)
    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            find_duplicate_keys(item, f"{path}[{i}]", duplicates)
    
    return duplicates

def get_all_keys(obj, path=""):
    """Get all keys from a nested JSON object as flat paths."""
    keys = []
    if isinstance(obj, dict):
        for key, value in obj.items():
            full_path = f"{path}.{key}" if path else key
            keys.append(full_path)
            if isinstance(value, dict):
                keys.extend(get_all_keys(value, full_path))
    return keys

def compare_structures(en_data, ar_data, namespace):
    """Compare structures between EN and AR."""
    en_keys = set(get_all_keys(en_data))
    ar_keys = set(get_all_keys(ar_data))
    
    missing_in_ar = en_keys - ar_keys
    missing_in_en = ar_keys - en_keys
    
    return {
        'namespace': namespace,
        'missing_in_ar': sorted(missing_in_ar),
        'missing_in_en': sorted(missing_in_en),
        'en_count': len(en_keys),
        'ar_count': len(ar_keys)
    }

def main():
    i18n_dir = Path(__file__).parent.parent / 'portal-web' / 'src' / 'i18n'
    en_dir = i18n_dir / 'en'
    ar_dir = i18n_dir / 'ar'
    
    # Get all namespace files
    en_files = {f.stem: f for f in en_dir.glob('*.json')}
    ar_files = {f.stem: f for f in ar_dir.glob('*.json')}
    
    all_namespaces = set(en_files.keys()) | set(ar_files.keys())
    
    print("=" * 80)
    print("i18n Translation Files Audit")
    print("=" * 80)
    
    issues_found = False
    
    # Check for missing namespace files
    missing_en = set(ar_files.keys()) - set(en_files.keys())
    missing_ar = set(en_files.keys()) - set(ar_files.keys())
    
    if missing_en:
        print(f"\n❌ Missing EN namespace files: {', '.join(missing_en)}")
        issues_found = True
    if missing_ar:
        print(f"\n❌ Missing AR namespace files: {', '.join(missing_ar)}")
        issues_found = True
    
    # Audit each namespace
    for namespace in sorted(all_namespaces):
        print(f"\n{'─' * 80}")
        print(f"Namespace: {namespace}")
        print('─' * 80)
        
        # Check EN file
        if namespace in en_files:
            with open(en_files[namespace], 'r', encoding='utf-8') as f:
                try:
                    en_data = json.load(f)
                    en_duplicates = find_duplicate_keys(en_data)
                    if en_duplicates:
                        print(f"\n❌ EN duplicates found:")
                        for dup in en_duplicates:
                            print(f"   {dup}")
                        issues_found = True
                except json.JSONDecodeError as e:
                    print(f"\n❌ EN file has invalid JSON: {e}")
                    issues_found = True
                    continue
        else:
            print(f"\n❌ EN file missing!")
            issues_found = True
            continue
        
        # Check AR file
        if namespace in ar_files:
            with open(ar_files[namespace], 'r', encoding='utf-8') as f:
                try:
                    ar_data = json.load(f)
                    ar_duplicates = find_duplicate_keys(ar_data)
                    if ar_duplicates:
                        print(f"\n❌ AR duplicates found:")
                        for dup in ar_duplicates:
                            print(f"   {dup}")
                        issues_found = True
                except json.JSONDecodeError as e:
                    print(f"\n❌ AR file has invalid JSON: {e}")
                    issues_found = True
                    continue
        else:
            print(f"\n❌ AR file missing!")
            issues_found = True
            continue
        
        # Compare structures
        comparison = compare_structures(en_data, ar_data, namespace)
        
        if comparison['missing_in_ar']:
            print(f"\n❌ Keys in EN but missing in AR ({len(comparison['missing_in_ar'])}):")
            for key in comparison['missing_in_ar'][:10]:  # Show first 10
                print(f"   {key}")
            if len(comparison['missing_in_ar']) > 10:
                print(f"   ... and {len(comparison['missing_in_ar']) - 10} more")
            issues_found = True
        
        if comparison['missing_in_en']:
            print(f"\n❌ Keys in AR but missing in EN ({len(comparison['missing_in_en'])}):")
            for key in comparison['missing_in_en'][:10]:
                print(f"   {key}")
            if len(comparison['missing_in_en']) > 10:
                print(f"   ... and {len(comparison['missing_in_en']) - 10} more")
            issues_found = True
        
        if not comparison['missing_in_ar'] and not comparison['missing_in_en']:
            print(f"\n✅ Perfect parity! ({comparison['en_count']} keys)")
    
    print(f"\n{'=' * 80}")
    if issues_found:
        print("❌ Issues found! See details above.")
        sys.exit(1)
    else:
        print("✅ All translation files are in perfect sync!")
        sys.exit(0)

if __name__ == '__main__':
    main()
