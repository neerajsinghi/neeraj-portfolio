#!/usr/bin/env bash
# import-lambda-env-to-ssm.sh
#
# Copy environment variables that were set directly on the backend Lambda
# (outside Terraform) into SSM Parameter Store, so Terraform reads them from
# SSM instead of resetting them on the next apply. Values never touch the
# terminal or argv: they go from the Lambda config to SSM through a pipe.
#
# Usage:
#   AWS_PROFILE=portfolio ./scripts/import-lambda-env-to-ssm.sh [--region ap-south-1] [--overwrite] [VAR ...]
#
# With no VAR arguments it imports LINKEDIN_ACCESS_TOKEN, LINKEDIN_AUTHOR_URN
# and LINKEDIN_VERSION. Each VAR maps to /neeraj-portfolio/<var-in-kebab-case>.
#
# Prerequisites: aws CLI and jq; lambda:GetFunctionConfiguration and ssm:PutParameter

set -euo pipefail

REGION="${AWS_REGION:-ap-south-1}"
OVERWRITE=false
SSM_PREFIX="/neeraj-portfolio"
FUNCTION_NAME="neeraj-portfolio-backend-api"
vars=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --region)    REGION="$2";    shift 2 ;;
    --overwrite) OVERWRITE=true; shift   ;;
    -*) echo "Unknown flag: $1"; exit 1  ;;
    *)  vars+=("$1");            shift   ;;
  esac
done

[[ ${#vars[@]} -eq 0 ]] && vars=(LINKEDIN_ACCESS_TOKEN LINKEDIN_AUTHOR_URN LINKEDIN_VERSION)

command -v aws &>/dev/null || { echo "ERROR: aws CLI not found"; exit 1; }
command -v jq  &>/dev/null || { echo "ERROR: jq not found"; exit 1; }

env_json="$(aws lambda get-function-configuration \
  --region "$REGION" \
  --function-name "$FUNCTION_NAME" \
  --query 'Environment.Variables' \
  --output json \
  --no-cli-pager)"

echo "Importing Lambda env vars into SSM (region: $REGION, prefix: $SSM_PREFIX)"
for var in "${vars[@]}"; do
  name="$(tr 'A-Z_' 'a-z-' <<<"$var")"
  if [[ "$(jq -r --arg k "$var" '.[$k] // "" | length' <<<"$env_json")" == "0" ]]; then
    echo "  • Skipping $SSM_PREFIX/$name ($var is empty on the Lambda)"
    continue
  fi
  if jq -c --arg n "$SSM_PREFIX/$name" --arg k "$var" --argjson o "$OVERWRITE" \
      '{Name: $n, Value: .[$k], Type: "SecureString", Overwrite: $o}' <<<"$env_json" \
    | aws ssm put-parameter \
        --region "$REGION" \
        --cli-input-json file:///dev/stdin \
        --no-cli-pager \
        --query 'Version' \
        --output text >/dev/null 2>&1; then
    echo "  ✓ $SSM_PREFIX/$name"
  else
    echo "  ⚠ $SSM_PREFIX/$name  (already exists, use --overwrite to update)"
  fi
done
unset env_json
