#!/bin/bash
set -euo pipefail

LATEST_TAG=$(git --no-pager describe --abbrev=0)
CURRENT_BRANCH=$(git --no-pager rev-parse --abbrev-ref HEAD)

if [ "$CURRENT_BRANCH" != "master" ]; then
  echo "⚠️  WARNING: You are not in the "master" branch! Current branch: ${CURRENT_BRANCH}"
  read -p "Press <return> if you're sure you want to release from this branch" -n 1 -r
fi

echo "Latest tag: ${LATEST_TAG}"
echo "Commits since latest tag:"
git --no-pager log --oneline "${LATEST_TAG}..HEAD"

# Get new tag and changelog from user
# TODO: Changelog could be generated automatically from commits
echo
read -p "New tag (format: vX.X.X): " -r NEW_TAG

echo 'Changelog for release (ctrl-d when done):'
CHANGELOG=$(cat)
echo

echo -e "About to create new tag: ${NEW_TAG} with changelog:\n${CHANGELOG}"
read -p "Press <return> if you're sure" -n 1 -r

echo "Creating new tag ${NEW_TAG}..."
git tag -a "$NEW_TAG" -m "$CHANGELOG"
echo 'And pushing to remote...'
git push --follow-tags
echo '🎉 Done!'
echo 'CI will create a GitHub release automatically'
