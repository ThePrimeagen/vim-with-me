add_subdirectory(system-commands)
add_subdirectory(hashmap)
add_subdirectory(twitch)

list(APPEND EXTRA_LIBS system-commands)
list(APPEND EXTRA_LIBS json-c)
list(APPEND EXTRA_LIBS twitch)
list(APPEND EXTRA_LIBS hashmap)
# list(APPEND EXTRA_LIBS ${LIB_WEB_SOCKETS})

# add the executable
add_executable(vim-with-me vim-with-me.c)

target_link_libraries(vim-with-me PUBLIC
    ${EXTRA_LIBS}
)

configure_file(vim-with-me.h.in vim-with-me.h)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED True)

target_include_directories(vim-with-me PUBLIC
    "${PROJECT_BINARY_DIR}"
)

