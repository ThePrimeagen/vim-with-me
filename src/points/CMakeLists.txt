add_library(primepoints primepoints.c)

list(APPEND EXTRA_LIBS json-c)

target_link_libraries(primepoints
    ${EXTRA_LIBS}
)

target_include_directories(primepoints INTERFACE
    ${CMAKE_CURRENT_SOURCE_DIR}
)

target_include_directories(primepoints PRIVATE
    ${PROJECT_LIB_FOLDER}
)

