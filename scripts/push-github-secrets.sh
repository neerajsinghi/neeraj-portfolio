#!/usr/bin/env bash
# push-github-secrets.sh
#
# Push required secrets and variables into GitHub Actions using the gh CLI.
# Prerequisites:
#   - gh auth login  (already authenticated)
#   - terraform apply already run (for AWS_ROLE_ARN / AMPLIFY_APP_ID)
#   - backend/.env filled in (for API keys / ALLOWED_ORIGIN)
#
# Usage:
#   set -a; source backend/.env; set +a
#   export AWS_ROLE_ARN=$(terraform -chdir=infrastructure/terraform output -raw github_actions_role_arn)
#   export AMPLIFY_APP_ID=$(terraform -chdir=infrastructure/terraform output -raw amplify_app_id)
#   export AWS_REGION=ap-south-1
#   ./scripts/push-github-secrets.sh

set -euo pipefail

command -v gh &>/dev/null || { echo "ERROR: gh CLI not found"; exit 1; }

REPO="neerajsinghi/neeraj-portfolio"

req_secrets=(ANTHROPIC_API_KEY  ALLOWED_ORIGIN AWS_ROLE_ARN AMPLIFY_APP_ID)
missing=()
for var in "${req_secrets[@]}"; do
  [[ -z "${!var:-}" ]] && missing+=("$var")
done
if [[ ${#missing[@]} -gt 0 ]]; then
  echo "ERROR: missing required env vars before running this script:"
  printf '  %s\n' "${missing[@]}"
  exit 1
fi

echo "Setting GitHub secrets on $REPO ..."
gh secret set ANTHROPIC_API_KEY --repo "$REPO" --body "$ANTHROPIC_API_KEY"
gh secret set ALLOWED_ORIGIN    --repo "$REPO" --body "$ALLOWED_ORIGIN"
gh secret set AWS_ROLE_ARN      --repo "$REPO" --body "$AWS_ROLE_ARN"
gh secret set AMPLIFY_APP_ID    --repo "$REPO" --body "$AMPLIFY_APP_ID"

if [[ -n "${GITHUB_TOKEN:-}" ]]; then
  gh secret set GH_PAT --repo "$REPO" --body "$GITHUB_TOKEN"
fi

echo "Setting GitHub variables on $REPO ..."
gh variable set AWS_REGION           --repo "$REPO" --body "${AWS_REGION:-ap-south-1}"
gh variable set ANTHROPIC_MODEL      --repo "$REPO" --body "${ANTHROPIC_MODEL:-claude-sonnet-4-6}"
gh variable set PROFILE_GITHUB_USER  --repo "$REPO" --body "${GITHUB_USER:-neerajsinghi}"
gh variable set PORT                 --repo "$REPO" --body "${PORT:-8080}"

echo "Done. Verify at: https://github.com/$REPO/settings/secrets/actions"
