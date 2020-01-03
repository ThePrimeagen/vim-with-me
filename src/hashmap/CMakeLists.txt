add_library(hashmap hashmap.c)
target_include_directories(hashmap INTERFACE
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_include_directories(hashmap PRIVATE
    ${PROJECT_LIB_FOLDER}
)
