#!/bin/bash
# zig cc wrapper for darwin/arm64 cross-compilation from Linux.
#
# Explicitly adds zig's macOS SDK paths because zig's MachO linker does not
# auto-detect them when invoked as a linker driver by Go's build toolchain.
#
# Also filters flags incompatible with zig's MachO linker:
#   -Wl,--compress-debug-sections  GNU ld only, not supported on macOS
#   -lresolv  Part of libSystem.B on macOS 12+; safe to drop, symbols are
#             available at runtime via libSystem which is always loaded

ZIGLIB=$(zig env 2>/dev/null | python3 -c 'import json,sys; print(json.load(sys.stdin).get("lib_dir",""))' 2>/dev/null || echo "")

args=()
for arg; do
    case "$arg" in
        -Wl,--compress-debug-sections*|-lresolv) ;;
        *) args+=("$arg") ;;
    esac
done

sdk_flags=()
if [ -n "$ZIGLIB" ] && [ -d "${ZIGLIB}/libc/darwin" ]; then
    sdk_flags+=("-F${ZIGLIB}/libc/darwin/System/Library/Frameworks")
    sdk_flags+=("-L${ZIGLIB}/libc/darwin/usr/lib")
fi

exec zig cc -target aarch64-macos -I/tmp/pcsc-macos "${sdk_flags[@]}" "${args[@]}"
