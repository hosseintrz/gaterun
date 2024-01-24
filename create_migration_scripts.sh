#!/bin/bash

# Your database connection string
DATABASE_URL="postgresql://postgres:pass1234@localhost:5432/gaterun?sslmode=disable"

# Directory to store migration scripts
MIGRATIONS_DIR="./internal/migrations"

# Number of digits in migration numbers
DIGITS=6

# Get the current migration number
CURRENT_NUMBER=$(ls -1 "$MIGRATIONS_DIR" | wc -l)
NEXT_NUMBER=$(printf "%0${DIGITS}d" "$((CURRENT_NUMBER + 1))")

# Prompt user for migration name
read -p "Enter migration name: " MIGRATION_NAME

# Generate migration filenames
UP_SCRIPT="${MIGRATIONS_DIR}/${NEXT_NUMBER}_${MIGRATION_NAME}.up.sql"
DOWN_SCRIPT="${MIGRATIONS_DIR}/${NEXT_NUMBER}_${MIGRATION_NAME}.down.sql"

# Create migration scripts
touch "$UP_SCRIPT" "$DOWN_SCRIPT"

# Print success message
echo "Migration scripts created successfully:"
echo "Up Script: $UP_SCRIPT"
echo "Down Script: $DOWN_SCRIPT"

# Optionally, open the scripts in the default editor
code "$UP_SCRIPT" "$DOWN_SCRIPT"
