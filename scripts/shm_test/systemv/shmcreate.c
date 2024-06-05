#include "common.h"

int main(int argc, char **argv)
{
    int length = atoi(argv[1]);
    int oflag = IPC_CREAT | SHM_RW_PERMISSION;
    int shmid = shmget(ftok(FTOK_FILE, FTOK_ID), length, oflag);

    if (shmid >= 0)
    {
        printf("shmget create success, shmid = %d\n", shmid);
    }

    return 0;
}
