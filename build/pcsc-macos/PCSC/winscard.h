/*
 * Minimal macOS PCSC winscard stub for cross-compilation via zig cc.
 * Only used at compile time on Linux build runners; the real PCSC.framework
 * is linked at runtime on macOS.
 *
 * Matches the macOS PCSC.framework API naming (_A suffix, uint32_t handles).
 */
#ifndef _WINSCARD_H_
#define _WINSCARD_H_

#include <PCSC/wintypes.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

#define PCSCLITE_VERSION_NUMBER "1.9.9"
#define MAX_ATR_SIZE 33

typedef int32_t  SCARDCONTEXT;
typedef int32_t  SCARDHANDLE;
typedef SCARDCONTEXT *LPSCARDCONTEXT;
typedef SCARDHANDLE  *LPSCARDHANDLE;

typedef struct {
    LPCSTR szReader;
    LPVOID pvUserData;
    DWORD  dwCurrentState;
    DWORD  dwEventState;
    DWORD  cbAtr;
    BYTE   rgbAtr[MAX_ATR_SIZE];
} SCARD_READERSTATE_A, *LPSCARD_READERSTATE_A;

typedef SCARD_READERSTATE_A  SCARD_READERSTATE;
typedef LPSCARD_READERSTATE_A LPSCARD_READERSTATE;

typedef struct {
    uint32_t dwProtocol;
    uint32_t cbPciLength;
} SCARD_IO_REQUEST, *LPSCARD_IO_REQUEST;
typedef const SCARD_IO_REQUEST *LPCSCARD_IO_REQUEST;

/* Scope */
#define SCARD_SCOPE_USER     0x0000
#define SCARD_SCOPE_TERMINAL 0x0001
#define SCARD_SCOPE_SYSTEM   0x0002

/* Share mode */
#define SCARD_SHARE_EXCLUSIVE 0x0001
#define SCARD_SHARE_SHARED    0x0002
#define SCARD_SHARE_DIRECT    0x0003

/* Disposition */
#define SCARD_LEAVE_CARD   0x0000
#define SCARD_RESET_CARD   0x0001
#define SCARD_UNPOWER_CARD 0x0002
#define SCARD_EJECT_CARD   0x0003

/* Protocol */
#define SCARD_PROTOCOL_UNDEFINED 0x0000
#define SCARD_PROTOCOL_T0        0x0001
#define SCARD_PROTOCOL_T1        0x0002
#define SCARD_PROTOCOL_RAW       0x0004
#define SCARD_PROTOCOL_ANY       (SCARD_PROTOCOL_T0 | SCARD_PROTOCOL_T1)

/* Reader state flags */
#define SCARD_STATE_UNAWARE     0x0000
#define SCARD_STATE_IGNORE      0x0001
#define SCARD_STATE_CHANGED     0x0002
#define SCARD_STATE_UNKNOWN     0x0004
#define SCARD_STATE_UNAVAILABLE 0x0008
#define SCARD_STATE_EMPTY       0x0010
#define SCARD_STATE_PRESENT     0x0020
#define SCARD_STATE_ATRMATCH    0x0040
#define SCARD_STATE_EXCLUSIVE   0x0080
#define SCARD_STATE_INUSE       0x0100
#define SCARD_STATE_MUTE        0x0200
#define SCARD_STATE_UNPOWERED   0x0400

#define SCARD_AUTOALLOCATE ((DWORD)(-1))
#define INFINITE           0xFFFFFFFF

/* Return codes */
#define SCARD_S_SUCCESS              ((LONG)0x00000000L)
#define SCARD_F_INTERNAL_ERROR       ((LONG)0x80100001L)
#define SCARD_E_CANCELLED            ((LONG)0x80100002L)
#define SCARD_E_INVALID_HANDLE       ((LONG)0x80100003L)
#define SCARD_E_INVALID_PARAMETER    ((LONG)0x80100004L)
#define SCARD_E_INVALID_TARGET       ((LONG)0x80100005L)
#define SCARD_E_NO_MEMORY            ((LONG)0x80100006L)
#define SCARD_E_TIMEOUT              ((LONG)0x8010000AL)
#define SCARD_E_SHARING_VIOLATION    ((LONG)0x8010000BL)
#define SCARD_E_NO_SMARTCARD         ((LONG)0x8010000CL)
#define SCARD_E_PROTO_MISMATCH       ((LONG)0x8010000FL)
#define SCARD_E_NOT_READY            ((LONG)0x80100010L)
#define SCARD_E_INVALID_VALUE        ((LONG)0x80100011L)
#define SCARD_E_NO_READERS_AVAILABLE ((LONG)0x8010002EL)
#define SCARD_E_NO_SERVICE           ((LONG)0x8010001DL)
#define SCARD_E_READER_UNAVAILABLE   ((LONG)0x80100017L)
#define SCARD_W_REMOVED_CARD         ((LONG)0x80100069L)
#define SCARD_W_RESET_CARD           ((LONG)0x80100068L)
#define SCARD_W_UNPOWERED_CARD       ((LONG)0x80100067L)
#define SCARD_W_UNRESPONSIVE_CARD    ((LONG)0x80100066L)
#define SCARD_W_UNSUPPORTED_CARD     ((LONG)0x80100065L)

const char *pcsc_stringify_error(int32_t error_code);

LONG SCardEstablishContext(uint32_t dwScope, const void *pvReserved1, const void *pvReserved2, LPSCARDCONTEXT phContext);
LONG SCardReleaseContext(SCARDCONTEXT hContext);
LONG SCardIsValidContext(SCARDCONTEXT hContext);
LONG SCardCancel(SCARDCONTEXT hContext);

LONG SCardConnect(SCARDCONTEXT hContext, LPCSTR szReader, uint32_t dwShareMode, uint32_t dwPreferredProtocols, LPSCARDHANDLE phCard, uint32_t *pdwActiveProtocol);
LONG SCardReconnect(SCARDHANDLE hCard, uint32_t dwShareMode, uint32_t dwPreferredProtocols, uint32_t dwInitialization, uint32_t *pdwActiveProtocol);
LONG SCardDisconnect(SCARDHANDLE hCard, uint32_t dwDisposition);

LONG SCardBeginTransaction(SCARDHANDLE hCard);
LONG SCardEndTransaction(SCARDHANDLE hCard, uint32_t dwDisposition);

LONG SCardStatus(SCARDHANDLE hCard, LPSTR mszReaderName, uint32_t *pcchReaderLen, uint32_t *pdwState, uint32_t *pdwProtocol, BYTE *pbAtr, uint32_t *pcbAtrLen);
LONG SCardGetStatusChange(SCARDCONTEXT hContext, uint32_t dwTimeout, void *rgReaderStates, uint32_t cReaders);

LONG SCardTransmit(SCARDHANDLE hCard, void *pioSendPci, const BYTE *pbSendBuffer, uint32_t cbSendLength, void *pioRecvPci, BYTE *pbRecvBuffer, uint32_t *pcbRecvLength);
LONG SCardControl(SCARDHANDLE hCard, uint32_t dwControlCode, const void *pbSendBuffer, uint32_t cbSendLength, void *pbRecvBuffer, uint32_t cbRecvLength, uint32_t *lpBytesReturned);

LONG SCardListReaders(SCARDCONTEXT hContext, LPCSTR mszGroups, LPSTR mszReaders, uint32_t *pcchReaders);
LONG SCardListReaderGroups(SCARDCONTEXT hContext, LPSTR mszGroups, uint32_t *pcchGroups);
LONG SCardFreeMemory(SCARDCONTEXT hContext, const void *pvMem);

LONG SCardGetAttrib(SCARDHANDLE hCard, uint32_t dwAttrId, uint8_t *pbAttr, uint32_t *pcbAttrLen);
LONG SCardSetAttrib(SCARDHANDLE hCard, uint32_t dwAttrId, const uint8_t *pbAttr, uint32_t cbAttrLen);

#ifdef __cplusplus
}
#endif

#endif /* _WINSCARD_H_ */
