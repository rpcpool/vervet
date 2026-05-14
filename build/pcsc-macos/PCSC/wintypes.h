/*
 * Minimal macOS PCSC wintypes stub for cross-compilation via zig cc.
 * Only used at compile time on Linux build runners; the real PCSC.framework
 * is linked at runtime on macOS.
 */
#ifndef _WINTYPES_H_
#define _WINTYPES_H_

#include <stdint.h>
#include <stddef.h>

typedef uint8_t  BYTE;
typedef uint16_t WORD;
typedef uint32_t DWORD;
typedef int32_t  LONG;
typedef int      BOOL;

typedef BYTE    *LPBYTE;
typedef DWORD   *LPDWORD;
typedef LONG    *LPLONG;
typedef char    *LPSTR;
typedef const char *LPCSTR;
typedef void    *LPVOID;
typedef const void *LPCVOID;

#ifndef TRUE
#define TRUE  1
#endif
#ifndef FALSE
#define FALSE 0
#endif

#endif /* _WINTYPES_H_ */
