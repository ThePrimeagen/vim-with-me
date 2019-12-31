add_subdirectory(system-commands)
list(APPEND EXTRA_LIBS system-commands)
list(APPEND EXTRA_LIBS json-c)

# add the executable
add_executable(vim-with-me vim-with-me.c)

target_link_libraries(vim-with-me PUBLIC
    ${EXTRA_LIBS}
)

configure_file(vim-with-me.h.in vim-with-me.h)
configure_file(.config .config COPYONLY)

set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED True)
target_include_directories(vim-with-me PUBLIC
    "${PROJECT_BINARY_DIR}"
)

#set(JSON_C_ROOT ${CMAKE_BINARY_DIR}/thirdparty/hidapi)
#set(JSON_C_LIBDIR ${JSON_C_ROOT}/.lib)
#set(HIDAPI_INCLUDE_DIR ${JSON_C_ROOT}/bin/include)
#
#include(ExternalProject)
#
#ExternalProject_Add(hidapi_external
#                    PREFIX ${JSON_C_ROOT}
#                    GIT_REPOSITORY "https://github.com/json-c/json-c.git"
#                    GIT_TAG "ee34939e7443f042fe52154ab88a414b2b310211"
#                    UPDATE_COMMAND ""
#                    PATCH_COMMAND ""
#                    BINARY_DIR ${JSON_C_ROOT}/
#                    SOURCE_DIR ${JSON_C_ROOT}/
#                    INSTALL_DIR ${JSON_C_ROOT}/bin
#                    CONFIGURE_COMMAND ./configure --prefix=<INSTALL_DIR>
#                    BUILD_COMMAND make && make install
#                    BUILD_BYPRODUCTS ${JSON_C_LIBDIR}/libjson-c.a)
#
## hidapi requires an initial execution of "./bootstrap" to run autoreconf
#ExternalProject_Add_Step(hidapi_external
#                         bootstrap
#                         COMMAND ./bootstrap
#                         DEPENDEES download
#                         DEPENDERS configure
#                         WORKING_DIRECTORY ${JSON_C_ROOT}/src/hidapi)
#
#add_library(libhidapi STATIC IMPORTED)
#set_target_properties(libhidapi PROPERTIES IMPORTED_LOCATION ${HIDAPI_LIB_DIR}/libhidapi.a)
#add_dependencies(libhidapi hidapi_external)
