#!/bin/bash
set -x

LOG_DIR=${LOG_DIR:-/tmp/devstack-logs}
mkdir -p $LOG_DIR
sudo journalctl -o short-precise --no-pager &> $LOG_DIR/journal.log
free -m > $LOG_DIR/free.txt
dpkg -l > $LOG_DIR/dpkg-l.txt
pip freeze > $LOG_DIR/pip-freeze.txt
sudo find $LOG_DIR -type d -exec chmod 0755 {} \;
sudo find $LOG_DIR -type f -exec chmod 0644 {} \;
