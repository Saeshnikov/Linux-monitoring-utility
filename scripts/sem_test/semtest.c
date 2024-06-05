#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/sem.h>

#define SEM_KEY 0x12345

int main(int argc, char *argv[])
{
    pid_t pid;
    key_t key;
    int semid;
    int counter;
    struct sembuf buf[2];
    key = ftok("./one_two", 1);
    semid = semget(SEM_KEY, 2, 0666 | IPC_CREAT);
    buf[0].sem_num = 0;
    buf[0].sem_flg = SEM_UNDO;
    buf[1].sem_num = 1;
    buf[1].sem_flg = SEM_UNDO;

    // set two semaphores values, don't use semop now
    semctl(semid, 0, SETVAL, 1);
    semctl(semid, 1, SETVAL, 0);

    counter = 1;
    pid = fork();
    if (pid > 0) {
        while (counter < 10) {
            // wait 0 to be 0
            buf[0].sem_op = 0;
            // only read first semaphore
            semop(semid, buf, 1);

            write(2, "0", 2);

            // decrease 1, increase 0");
            buf[0].sem_op = 1;
            buf[1].sem_op = -1;

            semop(semid, buf, 2);

            counter++;

        }

    } else {

        while (counter < 10) {

            // wait 1 to be 0
            buf[1].sem_op = 0;
            // only read second semaphore
            semop(semid, buf+1, 1);

            write(2, "1", 2);

            // decrease 0, increase 1
            buf[0].sem_op = -1;
            buf[1].sem_op = 1;
            semop(semid, buf, 2);

            counter++;
            if (counter > 10) {
                puts("1: finished");
                return 1;
            }
        }

    }
    return 0;
}
