cmake_minimum_required(VERSION 2.8)
project(geflect-example)

add_executable(geflect-example
    main.c
    gitmeta.c
)
add_custom_command(
    TARGET geflect-example POST_BUILD
    COMMAND  ${PROJECT_BINARY_DIR}/geflect-example
)
