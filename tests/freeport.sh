#!/bin/bash
comm -23 <(seq "40000" "50000" | sort) <(ss -Htan | awk '{print $4}' | cut -d':' -f2 | sort -u) | shuf 2>/dev/null | head -n "1" 2>/dev/null
