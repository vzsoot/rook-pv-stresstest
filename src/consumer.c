#include <stdlib.h>
#include <stdio.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <sys/io.h>
#include <sys/mman.h>
#include <dirent.h>
#include <string.h>
  
void readFiles(char* path) {
    DIR *dir = opendir(path);
    if (dir == NULL) {
        printf("[ERROR] Failed to open directory: %s\n", path);
        return;
    }

    struct dirent *entry;
    char *filename = malloc(256 * sizeof(char));
    unsigned char *f;
    int size;
    struct stat s;

    printf("[INFO] Read files from volume: %s\n", path);
    while ((entry = readdir(dir)) != NULL) {
        int nameLength = strlen(entry->d_name);
        if (strcmp(&(entry->d_name)[nameLength-4], ".txt")==0) {
            sprintf(filename, "%s/%s", path, entry->d_name);
            int fd = open(filename, O_RDONLY);
            int status = fstat(fd, &s);
            size = s.st_size;

            f = (char *) mmap(0, size, PROT_READ, MAP_PRIVATE, fd, 0);
            
            printf("[INFO] Last value from file: %c\n", f[size-1]);
            close(fd);
        }
    }

    free(filename);
    closedir(dir);
}

int main() {
    char* volume1 = getenv("VOLUME_PATH1");
    char* volume2 = getenv("VOLUME_PATH2");

    setbuf(stdout, NULL);
    
    while (1) {
        printf("[INFO] Reading files ...\n");
        readFiles(volume1);
        readFiles(volume2);
        sleep(1);
    }
    return 0;
}
