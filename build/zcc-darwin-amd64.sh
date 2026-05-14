#!/bin/bash
# zig cc wrapper for darwin/amd64 cross-compilation from Linux.
# Filters flags incompatible with zig's MachO linker:
#   -Wl,--compress-debug-sections  GNU ld only, not supported on macOS
#   -lresolv  Part of libSystem.B on macOS 12+; dropped to avoid missing stub
args=()
for arg; do
    case "$arg" in
        -Wl,--compress-debug-sections*|-lresolv) ;;
        *) args+=("$arg") ;;
    esac
done
exec zig cc -target x86_64-macos -I/tmp/pcsc-macos "${args[@]}"
