#include "common.h"

int main(int argc, char **argv)
{
    int shmid = shmget(ftok(FTOK_FILE, FTOK_ID), 0, SHM_RW_PERMISSION);
    shmctl(shmid, IPC_RMID, NULL);

    return 0;
}
