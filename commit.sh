#!/bin/bash

# Add changes and commit with a provided message
if [ -z "$1" ]; then
    echo "Error: Commit message is required."
    exit 1
fi

git add .
if ! git commit -m "$1"; then
    echo "Error: Commit failed. Please resolve any issues and try again."
    exit 1
fi

# Determine the branch
branch=${2:-$(git rev-parse --abbrev-ref HEAD)} # Default to the current branch if not provided

# Check if the branch exists locally
if ! git rev-parse --verify "$branch" >/dev/null 2>&1; then
    echo "Error: Branch '$branch' does not exist locally."
    echo "Please create the branch locally before pushing changes to it."
    exit 1
fi

# Warn about pushing to main and prompt for confirmation
if [ "$branch" == "main" ]; then
    read -p "You are about to push to 'main'. Are you sure? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        echo "Push to 'main' aborted."
        exit 0
    fi
fi

# Push changes
if ! git push origin "$branch"; then
    echo "Error: Push failed. Please check your connection or branch permissions."
    exit 1
fi

echo "Changes successfully pushed to '$branch'."
