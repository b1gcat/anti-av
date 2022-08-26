[SECTION .data]

global _Q360AllocateVirtualMemory

global _WhisperMain
extern _SW2_GetSyscallNumber

[SECTION .text]

BITS 32

_WhisperMain:
    pop eax                        ; Remove return address from CALL instruction
    call _SW2_GetSyscallNumber     ; Resolve function hash into syscall number
    add esp, 4                     ; Restore ESP
    mov ecx, [fs:0c0h]
    test ecx, ecx
    jne _wow64
    lea edx, [esp+4h]
    INT 02eh
    ret
_wow64:
    xor ecx, ecx
    lea edx, [esp+4h]
    call dword [fs:0c0h]
    ret

_Q360AllocateVirtualMemory:
    push 0C494DA01h
    call _WhisperMain

