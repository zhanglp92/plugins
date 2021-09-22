#!/bin/bash

if [ "" == "${GOPATH}" ]; then
  echo "ERR: GOPATH is empty"
  exit 1
fi

work_dir="$GOPATH/src"
git_path="github.com/zhanglp92/plugins"

cd $work_dir
git clone -- https://${git_path}.git $git_path
if [ $? != 0 ]; then
  echo "ERR: clone fail"
  exit 1
fi

cd $git_path/imports
GO111MODULE=on go install

echo "done."
