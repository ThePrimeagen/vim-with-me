cmake_minimum_required(VERSION 3.10)

# set the project name
project(vim-with-me VERSION 1.0)

# find_library(LIB_WEB_SOCKETS libwebsockets.so)

set(PROJECT_LIB_FOLDER ${CMAKE_CURRENT_SOURCE_DIR}/lib)

add_subdirectory(3rdparty/json-c)
add_subdirectory(3rdparty/libwebsockets)

add_subdirectory(src)

configure_file(.config src/.config COPYONLY)


#add_library(libwebsockets SHARED IMPORTED) # or STATIC instead of SHARED
#set_target_properties(libwebsockets PROPERTIES
#  IMPORTED_LOCATION "/usr/local/lib/libwebsockets.so"
#)

# add the executable
# add_executable(vim-with-me src/vim-with-me.c)

# add the executable
# add_executable(vim-with-me ./src/vim-with-me.c)
#target_link_libraries(${PROJECT_NAME} PUBLIC 3rdparty/json-c)

