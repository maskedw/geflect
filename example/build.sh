#!/bin/bash

if [ "$(uname)" == 'Darwin' ]; then
  OS=darwin_amd64
elif [ "$(expr substr $(uname -s) 1 5)" == 'Linux' ]; then
  OS=linux_amd64
elif [ "$(expr substr $(uname -s) 1 5)" == 'MINGW' ]; then
  OS=windows_amd64
else
  echo "Your platform ($(uname -a)) is not supported."
  exit 1
fi

cd "$(dirname "${BASH_SOURCE:-$0}")"

if [ ! -d .git ]; then
    git init
    git add main.c
    git commit -m "first commit"
    git add build.sh
    git commit -m "first release"
    git tag v1.0.0
    git add CMakeLists.txt
    git commit -m "update"
fi

if [ ! -d build ]; then
    mkdir build
fi

if [ ! -d build/CMakeFiles ]; then
    (cd build; cmake ../ -G "Unix Makefiles")
fi

bin.${OS}/geflect --ignore-git-errors -o gitmeta.c gitmeta.template

(cd build; make)
