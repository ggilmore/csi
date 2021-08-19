#include <ctype.h>
#include <dirent.h>
#include <errno.h>
#include <grp.h>
#include <math.h>
#include <pwd.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/_types/_errno_t.h>
#include <sys/_types/_va_list.h>
#include <sys/dirent.h>
#include <sys/stat.h>
#include <sys/syslimits.h>
#include <time.h>
#include <unistd.h>

#define NUM_UNITS 9
#define USER_NAME_MAX 25
#define SIX_MONTHS_SECONDS 15780000

// doesn't seem to be longer than 12 characters for actual 'ls'
#define TIME_SIZE 12

// Can't have a filesize that's larger than INT_MAX
#define FILE_SIZE_MAX 11

typedef struct
{
  char name[NAME_MAX + 1];

  bool isHidden;

  int size;
  int hardLinks;

  uid_t uid;
  gid_t gid;

  struct timespec modTime;

} file_t;

typedef struct
{
  bool isHuman;
  bool isLong;
  bool printHidden;
} options_t;

void
listFiles(char* path, options_t* options);
void
listFilesDirectory(DIR* d, options_t* options);
file_t*
processFile(char* name);
void
printFile(file_t* file, options_t* options);

void
humanUID(uid_t uid, char* buf, int len);
void
humanGUID(gid_t gid, char* buf, int len);
void
humanSize(int bytes, char* buf, int len);
void
humanTimespec(struct timespec t, char* buf, int len);
int
compareNames(const void* a, const void* b);

bool
hasPrefix(char* base, char* prefix);
void
centerString(char* str, int width);
void
error(char* fmt, ...);

int
main(int argc, char** argv)
{
  options_t options;
  char* path;
  int i, numFiles;
  int c;
  bool printDir;

  opterr = 0;

  options.isHuman = false;
  options.isLong = false;
  options.printHidden = false;

  while ((c = getopt(argc, argv, "lah")) != -1) {
    switch (c) {
      case 'l':
        options.isLong = true;
        break;
      case 'a':
        options.printHidden = true;
        break;
      case 'h':
        options.isHuman = true;
        break;
      default:
        error("unknown option %c", c);
        exit(2);
    }
  }

  numFiles = argc - optind;
  if (numFiles == 0) {
    listFiles(".", &options);
    exit(0);
  }

  printDir = false;
  if (numFiles > 1) {
    // There is more than folder to run against
    printDir = true;
  }

  for (i = optind; i < argc; i++) {
    path = argv[i];
    if (printDir) {
      printf("%s:\n", path);
    }

    listFiles(path, &options);
    if (printDir) {
      printf("\n");
    }
  }
}

void
listFiles(char* path, options_t* options)
{
  char cwd[PATH_MAX];
  DIR* d;

  struct stat stbuf;
  file_t* file;

  if (stat(path, &stbuf) == -1) {
    error("unable to open %s: %s", path, strerror(errno));
    exit(1);
  }

  if ((stbuf.st_mode & S_IFMT) != S_IFDIR) {
    file = processFile(path);
    printFile(file, options);

    free(file);
    return;
  }

  if (getcwd(cwd, sizeof(cwd)) == NULL) {
    error("unable to get current working directory: %s", strerror(errno));
    exit(1);
  }

  if (chdir(path) == -1) {
    error("unable to chdir to %s: %s", path, strerror(errno));
    exit(1);
    return;
  }

  if ((d = opendir(".")) == NULL) {
    error("unable to open %s: %s", path, strerror(errno));
    return;
  }

  listFilesDirectory(d, options);
  closedir(d);

  if (chdir(cwd) == -1) {
    error("unable to chdir back to %s: %s", cwd, strerror(errno));
    exit(1);
    return;
  }
}

void
listFilesDirectory(DIR* d, options_t* options)
{
  struct dirent* dir;
  file_t** files;
  int size;
  file_t file;
  int i;

  size = 0;
  while ((dir = readdir(d))) {
    size++;
  }
  rewinddir(d);

  files = (file_t**)malloc(sizeof(file_t*) * size);
  if (files == NULL) {
    error("unable to allocate space for %d files", size);
    exit(1);
  }

  while ((dir = readdir(d))) {
    files[i++] = processFile(dir->d_name);
  }

  qsort(files, size, sizeof(file_t*), compareNames);

  for (i = 0; i < size; i++) {
    printFile(files[i], options);

    free(files[i]);
  }

  free(files);
}

file_t*
processFile(char* name)
{
  file_t* file;
  struct stat stbuf;

  if (stat(name, &stbuf) == -1) {
    error("unable to open %s: %s", name, strerror(errno));
    exit(1);
  }

  file = (file_t*)malloc(sizeof(file_t));
  if (file == NULL) {
    error("unable to allocate memory for file '%s'", name);
    exit(1);
  }

  strncpy(file->name, name, NAME_MAX + 1);

  file->isHidden = hasPrefix(file->name, ".");

  file->size = stbuf.st_size;
  file->hardLinks = stbuf.st_nlink;

  file->gid = stbuf.st_gid;
  file->uid = stbuf.st_uid;

  file->modTime = stbuf.st_mtimespec;

  return file;
}

void
printFile(file_t* file, options_t* options)
{

  if (file->isHidden && !options->printHidden) {
    return;
  }

  if (!options->isLong) {
    printf("%s\n", file->name);
    return;
  }

  char username[USER_NAME_MAX + 1];
  char groupName[USER_NAME_MAX + 1];
  char size[FILE_SIZE_MAX + 1];
  char modificationTime[TIME_SIZE + 1];

  humanUID(file->uid, username, USER_NAME_MAX + 1);
  humanGUID(file->gid, groupName, USER_NAME_MAX + 1);
  humanTimespec(file->modTime, modificationTime, TIME_SIZE + 1);
  if (options->isHuman) {
    humanSize(file->size, size, FILE_SIZE_MAX + 1);
  } else {
    snprintf(size, FILE_SIZE_MAX + 1, "%d", file->size);
  }

  printf("%5d", file->hardLinks);
  printf(" ");
  centerString(username, USER_NAME_MAX + 1);
  printf(" ");
  centerString(groupName, USER_NAME_MAX + 1);
  printf(" ");
  centerString(size, FILE_SIZE_MAX + 1);
  printf(" ");
  centerString(modificationTime, TIME_SIZE + 1);
  printf(" ");
  printf("%s\n", file->name);
}

void
humanSize(int bytes, char* buf, int len)
{
  int numUnits;
  int i;
  float floatBytes;

  char* units[] = { "B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB" };

  if (bytes < 1024) {
    snprintf(buf, len, "%dB", bytes);
    return;
  }

  numUnits = sizeof(units) / sizeof(char*);
  i = 0;
  floatBytes = bytes;

  while (floatBytes > 1024 && i < numUnits) {
    floatBytes /= 1024;
    i++;
  }

  snprintf(buf, len, "%5.1f%s", floatBytes, units[i]);
}

void
humanUID(uid_t uid, char* buf, int len)
{
  struct passwd* pass;

  if ((pass = getpwuid(uid)) == NULL) {
    error("failed to read uid: %s", strerror(errno));
    exit(1);
  }

  if (pass->pw_name[0] == '\0') {
    snprintf(buf, len, "%d", uid);
    return;
  }

  strncpy(buf, pass->pw_name, USER_NAME_MAX);
}

void
humanGUID(gid_t gid, char* buf, int len)
{
  struct group* g = getgrgid(gid);
  if (g == NULL) {
    error("failed to read gid: %s", strerror(errno));
    exit(1);
  }

  if (g->gr_name[0] == '\0') {
    snprintf(buf, len, "%d", gid);
    return;
  }

  strncpy(buf, g->gr_name, len);
}

void
humanTimespec(struct timespec t, char* buf, int len)
{
  char* format;
  time_t sixMonthsAgo, now, sixMonthsFromNow;

  time(&now);

  sixMonthsAgo = now - SIX_MONTHS_SECONDS;
  sixMonthsFromNow = now + SIX_MONTHS_SECONDS;

  if (t.tv_sec < sixMonthsAgo || t.tv_sec > sixMonthsFromNow) {
    format = "%Y %R";
  } else {
    format = "%b %d %R";
  }

  len = strftime(buf, len, format, localtime(&t.tv_sec));
  if (len == 0 && buf[0] != '\0') {
    error("strftime failed");
    exit(1);
  }
}

bool
hasPrefix(char* base, char* prefix)
{
  return strncmp(base, prefix, strlen(prefix)) == 0;
}

void
centerString(char* str, int width)
{
  int padding = (width - strlen(str)) / 2;
  int leftPadding = padding;
  int rightPadding = ((width % 2) == 0 ? padding : padding + 1);

  printf("%*s%s%*s", leftPadding, "", str, rightPadding, "");
}

int
compareNames(const void* a, const void* b)
{
  file_t* f1 = (file_t*)a;
  file_t* f2 = (file_t*)b;

  return strcmp(f1->name, f2->name);
}

void
error(char* fmt, ...)
{
  va_list args;

  va_start(args, fmt);
  fprintf(stderr, "error: ");
  vfprintf(stderr, fmt, args);
  fprintf(stderr, "\n");
  va_end(args);
}
