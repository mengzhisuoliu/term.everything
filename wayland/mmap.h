
#include <unistd.h>
#include <stdbool.h>

void *mmap_fd(int fd, size_t size);
bool unmap(void* addr, size_t size);
void* remap(int fd, void* addr, size_t size, size_t new_size);
void* map_failed(void);