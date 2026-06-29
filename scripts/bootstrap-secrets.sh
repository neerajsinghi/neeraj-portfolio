#!/usr/bin/env bash
# bootstrap-secrets.sh
#
# Push API keys into AWS SSM Parameter Store as SecureStrings.
# Run once after provisioning the AWS infrastructure, then again on key rotation.
#
# Usage:
#   export ANTHROPIC_API_KEY=sk-ant-...
#   export OPENAI_API_KEY=sk-...
#   export XAI_API_KEY=xai-...
#   export GEMINI_API_KEY=AI...
#   export ALLOWED_ORIGIN=https://your-frontend.amplifyapp.com
#
#   ./scripts/bootstrap-secrets.sh [--region us-east-1] [--no-overwrite]
#
# Prerequisites: aws CLI configured with a profile that has ssm:PutParameter

set -euo pipefail

# ── Defaults ─────────────────────────────────────────────────────────────────
REGION="${AWS_REGION:-us-east-1}"
OVERWRITE="--overwrite"
SSM_PREFIX="/neeraj-portfolio"

# ── Argument parsing ──────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case "$1" in
    --region)       REGION="$2";  shift 2 ;;
    --no-overwrite) OVERWRITE=""; shift   ;;
    *) echo "Unknown flag: $1"; exit 1    ;;
  esac
done

# ── Validate prerequisites ────────────────────────────────────────────────────
command -v aws &>/dev/null || { echo "ERROR: aws CLI not found"; exit 1; }

missing=()
for var in ANTHROPIC_API_KEY OPENAI_API_KEY XAI_API_KEY GEMINI_API_KEY ALLOWED_ORIGIN; do
  [[ -z "${!var:-}" ]] && missing+=("$var")
done

if [[ ${#missing[@]} -gt 0 ]]; then
  echo "ERROR: The following environment variables are not set:"
  printf '  %s\n' "${missing[@]}"
  exit 1
fi

# ── Helper ────────────────────────────────────────────────────────────────────
put_param() {
  local name="$1"
  local value="$2"

  aws ssm put-parameter \
    --region "$REGION" \
    --name "$SSM_PREFIX/$name" \
    --value "$value" \
    --type SecureString \
    ${OVERWRITE} \
    --no-cli-pager \
    --output text \
    --query 'Version' 2>/dev/null \
  | { read -r version; echo "  ✓ $SSM_PREFIX/$name  (version $version)"; } \
  || echo "  ⚠ $SSM_PREFIX/$name  (already exists, use --overwrite to update)"
}

# ── Write parameters ──────────────────────────────────────────────────────────
echo ""
echo "Storing secrets in SSM Parameter Store (region: $REGION)"
echo "Prefix: $SSM_PREFIX"
echo "---"

put_param "anthropic-api-key" "$ANTHROPIC_API_KEY"
put_param "openai-api-key"    "$OPENAI_API_KEY"
put_param "xai-api-key"       "$XAI_API_KEY"
put_param "gemini-api-key"    "$GEMINI_API_KEY"
put_param "allowed-origin"    "$ALLOWED_ORIGIN"

echo "---"
echo "All parameters stored. Verifying (names only):"
aws ssm get-parameters-by-path \
  --region "$REGION" \
  --path "$SSM_PREFIX" \
  --query 'Parameters[].Name' \
  --output table \
  --no-cli-pager
