#include "common.h"

int main(int argc, char **argv)
{
    int shmid;
    unsigned char *shmadd;
    struct shmid_ds buf;
    int i;

    shmid   = shmget(ftok(FTOK_FILE, FTOK_ID), 0, SHM_RW_PERMISSION);
    shmadd  = shmat(shmid, NULL, 0);
    shmctl(shmid, IPC_STAT, &buf);

    for (i = 0; i < buf.shm_segsz; i++)
    {
        *shmadd++ = i % 256;
    }

    return 0;
}
