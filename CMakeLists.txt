cmake_minimum_required(VERSION 3.10)

# set the project name
project(vim-with-me VERSION 1.0)

add_subdirectory(3rdparty/json-c)
add_subdirectory(src)

# add the executable
# add_executable(vim-with-me src/vim-with-me.c)

# add the executable
# add_executable(vim-with-me ./src/vim-with-me.c)
#target_link_libraries(${PROJECT_NAME} PUBLIC 3rdparty/json-c)

