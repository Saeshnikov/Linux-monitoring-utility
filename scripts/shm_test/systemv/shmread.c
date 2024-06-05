#include "common.h"

int main(int argc, char **argv)
{
    int shmid;
    unsigned char *shmadd;
    unsigned char v;
    struct shmid_ds buf;
    int error = 0;
    int i;

    shmid   = shmget(ftok(FTOK_FILE, FTOK_ID), 0, SHM_RW_PERMISSION);
    shmadd  = shmat(shmid, NULL, 0);
    shmctl(shmid, IPC_STAT, &buf);

    for (i = 0; i < buf.shm_segsz; i++)
    {
        v = *shmadd++;

        if (v != (i % 256))
        {
            printf("error: shmadd[%d] = %d\n", i, v);
            error++;
        }
    }

    if (error == 0)
    {
        printf("all of read is ok\n");
    }

    return 0;
}
