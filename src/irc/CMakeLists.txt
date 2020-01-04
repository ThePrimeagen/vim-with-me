
add_library(irc irc.c)

list(APPEND EXTRA_LIBS system-commands)
list(APPEND EXTRA_LIBS twitch)
list(APPEND EXTRA_LIBS hashmap)

target_link_libraries(irc PUBLIC
    ${EXTRA_LIBS}
)

target_include_directories(irc INTERFACE
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_include_directories(irc PRIVATE
    ${PROJECT_LIB_FOLDER}
)
