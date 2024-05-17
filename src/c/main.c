#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h> 
  
char* createFiles(int num, int size, const char* path) {
    char* filename = malloc(256 * sizeof(char));
    int bufferSize = size * sizeof(char);
    char *buffer = (char *)malloc(size * sizeof(char));
    memset(buffer, '1', bufferSize);

    for (int i = 1; i <= num; ++i) {
        sprintf(filename, "%s/data%d.txt", path, i);

        FILE *file = fopen(filename, "w");
        if (file == NULL) {
            printf("Error creating file %s\n", filename);
            return NULL;
        }

        fwrite(buffer, sizeof(char), size, file);
        
        fclose(file);
    }

    free(buffer);
    return filename;
}

void syncVolume(const char *filename) {
    if (strlen(filename) > 0) {
        int fd = open(filename, O_RDONLY);
        lseek(fd, 0, SEEK_SET);
        lseek(fd, 0, SEEK_END);
        syncfs(fd);
        close(fd);
    }
}

int main() {
    const char* volume1 = getenv("VOLUME_PATH1");
    const char* volume2 = getenv("VOLUME_PATH2");
    const char* fileNumEnv = getenv("FILE_NUM");
    const char* fileSizeEnv = getenv("FILE_SIZE");

    const int fileNum = atoi(fileNumEnv);
    const int fileSize = atoi(fileSizeEnv);

    while (1) {
        printf("[INFO] Producing files ...\n");
        const char* vol1File = createFiles(fileNum, fileSize, volume1);
        const char* vol2File = createFiles(fileNum, fileSize, volume2);
        syncVolume(vol1File);
        syncVolume(vol2File);
        sleep(2);
    }
    return 0;
}
