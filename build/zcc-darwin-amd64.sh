#!/bin/bash
# zig cc wrapper for darwin/amd64 cross-compilation from Linux.
#
# Zig's MachO linker does not auto-detect its bundled macOS SDK paths when
# invoked as a linker driver by Go's build toolchain. We find the SDK by
# navigating relative to the zig binary (sibling lib/ directory).
#
# Also filters flags incompatible with zig's MachO linker:
#   -Wl,--compress-debug-sections  GNU ld only, not supported on macOS
#   -lresolv  Part of libSystem.B on macOS 12+; safe to drop, symbols are
#             available at runtime via libSystem which is always loaded

# Find zig's lib dir: try zig env first, then fall back to binary-relative path
ZIGLIB=$(zig env 2>/dev/null | sed -n 's/.*"lib_dir"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1)
if [ -z "$ZIGLIB" ] || [ ! -d "$ZIGLIB" ]; then
    ZIG_BIN=$(readlink -f "$(which zig)" 2>/dev/null || which zig 2>/dev/null)
    ZIGLIB="$(dirname "$ZIG_BIN")/lib"
fi

args=()
for arg; do
    case "$arg" in
        -Wl,--compress-debug-sections*|-lresolv) ;;
        *) args+=("$arg") ;;
    esac
done

sdk_flags=()
if [ -d "${ZIGLIB}/libc/darwin" ]; then
    sdk_flags+=("-F${ZIGLIB}/libc/darwin/System/Library/Frameworks")
    sdk_flags+=("-L${ZIGLIB}/libc/darwin/usr/lib")
fi

exec zig cc -target x86_64-macos -I/tmp/pcsc-macos \
    -F/tmp/macos-sdk/System/Library/Frameworks \
    -Wl,-undefined,dynamic_lookup \
    "${sdk_flags[@]}" "${args[@]}"
