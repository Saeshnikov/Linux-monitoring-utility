#ifndef _COMMON_H_
#define _COMMON_H_

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/shm.h>

#define FTOK_FILE          "/home/anna/ftok.file"
#define FTOK_ID            1

#define SHM_RD_PERMISSION  0444
#define SHM_WR_PERMISSION  0222
#define SHM_RW_PERMISSION  (SHM_RD_PERMISSION | SHM_WR_PERMISSION)

#endif
