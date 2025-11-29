FROM debian:11.6

RUN apt update
RUN apt install git -y

# ooz
RUN apt install build-essential -y
RUN apt install cmake -y
RUN apt install pkg-config -y
RUN apt install libsodium-dev -y
RUN apt install libunistring-dev -y


# Install ooz
WORKDIR /ooz
# RUN git clone --depth=1 --recurse-submodules --shallow-submodules https://github.com/zao/ooz.git .
RUN git clone --depth=1 --recurse-submodules --shallow-submodules https://github.com/Sternrassler/ooz.git .
WORKDIR /ooz/build

RUN cmake .. -D CMAKE_BUILD_RPATH='$ORIGIN'
RUN cmake --build .
ENV PATH="${PATH}:/ooz/build"

WORKDIR /workspace