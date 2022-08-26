.686
.XMM 
.MODEL flat, c 
ASSUME fs:_DATA 

.data

.code

EXTERN SW2_GetSyscallNumber: PROC

WhisperMain PROC
    pop eax                        ; Remove return address from CALL instruction
    call SW2_GetSyscallNumber      ; Resolve function hash into syscall number
    add esp, 4                     ; Restore ESP
    mov ecx, fs:[0c0h]
    test ecx, ecx
    jne _wow64
    lea edx, [esp+4h]
    INT 02eh
    ret
_wow64:
    xor ecx, ecx
    lea edx, [esp+4h]
    call dword ptr fs:[0c0h]
    ret
WhisperMain ENDP

Q360AllocateVirtualMemory PROC
    push 0C494DA01h
    call WhisperMain
Q360AllocateVirtualMemory ENDP

end