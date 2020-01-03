




add_library(comms comms.c)

target_include_directories(comms INTERFACE
    ${CMAKE_CURRENT_SOURCE_DIR}
)

target_include_directories(comms PRIVATE
    ${PROJECT_LIB_FOLDER}
)
