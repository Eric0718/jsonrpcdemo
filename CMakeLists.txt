cmake_minimum_required(VERSION 3.0)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(CMAKE_C_STANDARD 99)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_EXTENSIONS OFF)

project(startRpc)

#add_subdirectory(jsonrpc/cmd)

#GO
set(GO_TARGET jsonrpc_go)
set(GO_SRCS jsonrpc/cmd/main.go)
set(GO_LIB libjsonrpc.a)

set(GO_PATH ./)

add_custom_command(OUTPUT ${CMAKE_CURRENT_BINARY_DIR}/${GO_LIB}
  DEPENDS ${GO_SRCS}
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}

  #COMMAND go env -w GO111MODULE="on"
  #COMMAND cd ${GO_PATH} 
  #COMMAND go mod init
  #COMMAND go mod tidy

  COMMAND env GOPATH=${GOPATH} go build -buildmode=c-archive
  -o "${CMAKE_CURRENT_BINARY_DIR}/${GO_LIB}"
  ${CMAKE_GO_FLAGS} ${GO_SRCS}
  COMMENT "Building Go library")

add_custom_target(${GO_TARGET} DEPENDS ${GO_LIB} ${HEADER})

add_library(jsonrpc STATIC IMPORTED GLOBAL)
add_dependencies(jsonrpc ${GO_TARGET})
set_target_properties(jsonrpc
  PROPERTIES
  IMPORTED_LOCATION ${CMAKE_CURRENT_BINARY_DIR}/${GO_LIB}
  INTERFACE_INCLUDE_DIRECTORIES ${CMAKE_CURRENT_BINARY_DIR})


#CPP 
set(TARGET startRpc)
set(CPP_SRCS run.cpp)

add_executable(${TARGET} ${CPP_SRCS})
target_link_libraries(${TARGET} jsonrpc pthread)