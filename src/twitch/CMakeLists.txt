add_library(twitch twitch.c)
target_include_directories(twitch INTERFACE
    ${CMAKE_CURRENT_SOURCE_DIR}
)

