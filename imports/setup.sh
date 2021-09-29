#!/bin/bash

if [ "" == "${GOPATH}" ]; then
  echo "ERR: GOPATH is empty"
  exit 1
fi

work_dir="$GOPATH/src"
git_path="github.com/zhanglp92/plugins"
git_dir="$work_dir/$git_path"


if [ -e "$git_dir" ]; then
  cd $git_dir
  git pull
  if [ $? != 0 ]; then
    echo "ERR: git pull failed, src: $git_dir"
    exit 1
  fi
else
  cd $work_dir
  git clone -- https://${git_path}.git $git_path
  if [ $? != 0 ]; then
    echo "ERR: clone fail"
    exit 1
  fi
fi

cd $git_dir/imports
GO111MODULE=on go install

echo "done."

echo "bin: $GOPATH/bin/imports"
