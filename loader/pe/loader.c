
#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include "loader.h"
#include "pe_loader.h"


void pe(unsigned char *image,unsigned int imageSize) {
    peLoader((BYTE *)image, (DWORD)imageSize);
}
