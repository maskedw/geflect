#include <stdio.h>
#include <stdbool.h>

extern const char GIT_HASH[];
extern const char GIT_BRANCH[];
extern const char GIT_TAG[];
extern const char GIT_DESCRIBE[];
extern const char GIT_SHOART_HASH[];
extern const bool GIT_IS_CLEAN;
extern const bool GIT_IS_CLEAN_NO_UNTRACED_FILES;


int main(int argc, char *argv[])
{
    printf("GIT_HASH                        = %s\n", GIT_HASH);
    printf("GIT_BRANCH                      = %s\n", GIT_BRANCH);
    printf("GIT_TAG                         = %s\n", GIT_TAG);
    printf("GIT_DESCRIBE                    = %s\n", GIT_DESCRIBE);
    printf("GIT_SHOART_HASH                 = %s\n", GIT_SHOART_HASH);
    printf("GIT_IS_CLEAN                    = %s\n", GIT_IS_CLEAN ? "true" : "false");
    printf("GIT_IS_CLEAN_NO_UNTRACED_FILES  = %s\n", GIT_IS_CLEAN_NO_UNTRACED_FILES ? "true" : "false");

    return 0;
}
