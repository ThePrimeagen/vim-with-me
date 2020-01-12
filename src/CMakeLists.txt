add_subdirectory(system-commands)
add_subdirectory(hashmap)
add_subdirectory(twitch)
add_subdirectory(irc)
add_subdirectory(points)
#add_subdirectory(comms)

list(APPEND EXTRA_LIBS system-commands)
list(APPEND EXTRA_LIBS json-c)
list(APPEND EXTRA_LIBS twitch)
list(APPEND EXTRA_LIBS hashmap)
list(APPEND EXTRA_LIBS websockets)
list(APPEND EXTRA_LIBS irc)
list(APPEND EXTRA_LIBS primepoints)

# add the executable
add_executable(vim-with-me vim-with-me.c)

find_package(Threads REQUIRED)

target_link_libraries(vim-with-me PUBLIC
    ${EXTRA_LIBS}
    Threads::Threads
)

configure_file(vim-with-me.h.in vim-with-me.h)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED True)

target_include_directories(vim-with-me PUBLIC
    "${PROJECT_BINARY_DIR}"
)

target_include_directories(vim-with-me PRIVATE
    "${PROJECT_LIB_FOLDER}"
)

