#!/usr/bin/env bash

# Run from any directory to link the Bruno collection to the project's .env.
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
project_root="$(cd -- "$script_dir/.." && pwd -P)"
env_file="$project_root/.env"
collection_dir="$project_root/bruno"
env_link="$collection_dir/.env"

if [[ ! -f "$env_file" ]]; then
  printf 'Missing %s. Create it from .env.example and fill in your keys.\n' "$env_file" >&2
  exit 1
fi

if [[ ! -f "$collection_dir/opencollection.yml" ]]; then
  printf 'Bruno collection not found: %s\n' "$collection_dir" >&2
  exit 1
fi

if [[ -L "$env_link" ]]; then
  if [[ "$env_link" -ef "$env_file" ]]; then
    printf 'Bruno .env link is already correct: %s\n' "$env_link"
    exit 0
  fi
elif [[ -e "$env_link" ]]; then
  if [[ ! -f "$env_link" ]]; then
    printf 'Cannot replace a non-file entry: %s\n' "$env_link" >&2
    exit 1
  fi

  # Keep existing credentials; the backup also matches the *.env Git ignore rule.
  backup_file="$(mktemp "$env_link.backup.XXXXXX.env")"
  mv -- "$env_link" "$backup_file"
  printf 'Existing Bruno .env saved to: %s\n' "$backup_file"
fi

# -n replaces a symlink to a directory instead of creating a link inside it.
ln -sfn -- ../.env "$env_link"
printf 'Bruno .env link configured: %s -> ../.env\n' "$env_link"
