#!/bin/bash
SOURCE_DIR="/home/tmp/minecraft-server/world"
DEST_DIR="/home/tmp/backups"

# tar -czf "$DEST_DIR/$(date +%Y-%m-%d_%H-%M-%S).tar.gz" "$SOURCE_DIR"
#
# Using tar
# cd "$SOURCE_DIR" && tar -czf "$DEST_DIR/backup-$(date +%Y-%m-%d_%H-%M-%S).tar.gz" .
# Using zip
cd "$SOURCE_DIR" && zip -r "$DEST_DIR/backup-$(date +%Y-%m-%d_%H-%M-%S).zip" .

if [ $(find "$DEST_DIR" -maxdepth 1 -type f | wc -l) -gt 48 ]; then
   find "$DEST_DIR" -maxdepth 1 -type f -printf '%T@ %p\0' 2>/dev/null | sort -z -r -n | head -1 | tr -d '\0' | xargs rm -f
fi
