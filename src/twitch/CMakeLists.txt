



add_library(twitch irc.c pubsub.c)
target_include_directories(twitch INTERFACE
    ${CMAKE_CURRENT_SOURCE_DIR}
)

target_include_directories(twitch PRIVATE
    ${PROJECT_LIB_FOLDER}
)
