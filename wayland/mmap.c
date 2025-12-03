#include <fcntl.h>
#include <sys/mman.h>
#include <unistd.h>
#include <stdio.h>
#include <stdbool.h>


void *mmap_fd(int fd, size_t size)
{
    int prot = PROT_READ | PROT_WRITE;
    int flags = MAP_SHARED;
    void *addr = mmap(NULL, size, prot, flags, fd, 0);
    if (addr == MAP_FAILED)
    {
        perror("mmap");
        return MAP_FAILED;
    }

    return addr;
}

bool unmap(void* addr, size_t size) {
	if (addr != MAP_FAILED) {
		if (munmap(addr, size) == -1) {
			perror("munmap in unmap");
            return false;
		}
	}
    return true;
}

void* remap(int fd, void* addr, size_t size, size_t new_size) {
 	if (new_size == size)
    {
        return addr;
    }
    if (addr == MAP_FAILED)
    {
        return MAP_FAILED;
    }
    bool unmap_success = unmap(addr, size);
    if (!unmap_success) {
        return MAP_FAILED;
    }
    addr = mmap_fd(fd, new_size);
    if (addr == MAP_FAILED)
    {
        perror("mmap in remap");
        return MAP_FAILED;
    }
    size = new_size;
    return addr;
}

void* map_failed(void) {
	return MAP_FAILED;
}