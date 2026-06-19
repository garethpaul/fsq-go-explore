#!/usr/bin/env python3
import re
import sys
from pathlib import Path


auth = Path(sys.argv[1]).read_text()
tests = Path(sys.argv[2]).read_text()
plan = Path(sys.argv[3]).read_text()

client = auth.split("func getHttpClient", 1)[-1].split("\n}", 1)[0]
required_auth = (
    "oauthUserRequestTimeout   = 10 * time.Second",
    "Timeout:       oauthUserRequestTimeout",
)
if any(auth.count(item) != 1 for item in required_auth):
    raise SystemExit("OAuth user requests must retain the reviewed 10-second timeout.")
if "Timeout:" not in client or "oauthUserRequestTimeout" not in client:
    raise SystemExit("The OAuth user HTTP client must apply the reviewed timeout.")

test_name = "TestGetHTTPClientBoundsOAuthUserRequests"
if tests.count(test_name) != 1 or "client.Timeout != 10*time.Second" not in tests:
    raise SystemExit("OAuth user timeout must retain its focused behavior test.")

frontmatter = plan.split("---", 2)[1]
statuses = re.findall(r"^status: .+$", frontmatter, flags=re.MULTILINE)
verification = plan.split("## Verification Completed\n", 1)[-1]
required_evidence = (
    "timeout removal mutation failed",
    "timeout drift mutation failed",
    "client assignment mutation failed",
    "focused test mutation failed",
    "guidance mutation failed",
    "plan evidence mutation failed",
    "hosted pull-request check",
)
if (
    statuses != ["status: completed"]
    or "## Verification Completed\n" not in plan
    or any(item not in verification for item in required_evidence)
    or re.search(r"\b(?:pending|todo|tbd|not run|not yet)\b", verification, re.IGNORECASE)
):
    raise SystemExit(
        "OAuth user timeout plan must remain completed with actual verification recorded."
    )
