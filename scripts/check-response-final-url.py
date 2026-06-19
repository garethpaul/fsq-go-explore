#!/usr/bin/env python3
import re
import sys
from pathlib import Path


if len(sys.argv) != 4:
    raise SystemExit("usage: check-response-final-url.py API TEST PLAN")

source = Path(sys.argv[1]).read_text()
tests = Path(sys.argv[2]).read_text()
plan = Path(sys.argv[3]).read_text()


def function_body(name: str, next_name: str) -> str:
    marker = f"func {name}"
    if source.count(marker) != 1:
        raise SystemExit(f"expected one {name} function")
    body = source.split(marker, 1)[1]
    if next_name:
        body = body.split(f"func {next_name}", 1)[0]
    return body


validator = function_body("isExpectedFoursquareResponseURL", "isFoursquareJSONResponse")
required_validator_fragments = (
    'responseURL.Scheme == "https"',
    "strings.EqualFold(responseURL.Hostname(), foursquareAPIHost)",
    "responseURL.User == nil",
    'responseURL.Port() == ""',
    "responseURL.EscapedPath() == expectedEscapedPath",
    'responseURL.Fragment == ""',
)
if any(validator.count(fragment) != 1 for fragment in required_validator_fragments):
    raise SystemExit("Foursquare final response URL validation must preserve the exact endpoint boundary.")

operations = (
    (
        "(fsqs *FoursquareService) Search",
        "(fsqs *FoursquareService) VenueDetails",
        "successfulFoursquareStatus(r.StatusCode)",
        "isExpectedFoursquareResponseURL(r, foursquareSearchPath)",
        "isFoursquareJSONResponse(r)",
        "decodeFoursquareResponse(r.Body, venues)",
    ),
    (
        "(fsqs *FoursquareService) VenueDetails",
        "(fsqs *FoursquareService) VenueEdit",
        "successfulFoursquareStatus(r.StatusCode)",
        "isExpectedFoursquareResponseURL(r, venuePath)",
        "isFoursquareJSONResponse(r)",
        "decodeFoursquareResponse(r.Body, venue)",
    ),
    (
        "(fsqs *FoursquareService) VenueEdit",
        "(fsqs *FoursquareConfig) userParams",
        "successfulFoursquareStatus(resp.StatusCode)",
        "isExpectedFoursquareResponseURL(resp, venueEditPath)",
        "discardFoursquareResponse(resp.Body)",
    ),
)

for operation in operations:
    body = function_body(operation[0], operation[1])
    contract = operation[2:]
    positions = [body.find(fragment) for fragment in contract]
    if any(body.count(fragment) != 1 for fragment in contract):
        raise SystemExit(f"{operation[0]} must keep one ordered final URL boundary.")
    if -1 in positions or positions != sorted(positions):
        raise SystemExit(f"{operation[0]} must validate the final URL before reading its response body.")

required_tests = (
    "func TestExpectedFoursquareResponseURL",
    "func TestFoursquareOperationsRejectUnexpectedFinalURLBeforeRead",
    'response body reads = %d, want 0',
)
if any(tests.count(fragment) != 1 for fragment in required_tests):
    raise SystemExit("Foursquare final URL regression coverage must remain executable and unique.")

frontmatter = plan.split("---", 2)[1]
statuses = re.findall(r"^status: .+$", frontmatter, flags=re.MULTILINE)
verification = plan.split("## Verification Completed\n", 1)[-1]
required_evidence = (
    "focused final URL tests passed",
    "six hostile mutations were rejected",
    "external-directory Make gate passed",
    "No live Foursquare request",
)
if (
    statuses != ["status: completed"]
    or "## Verification Completed\n" not in plan
    or any(item not in verification for item in required_evidence)
    or re.search(r"\b(?:pending|todo|tbd|not run)\b", verification, re.IGNORECASE)
):
    raise SystemExit("Foursquare final URL plan must record completed status and actual verification.")

print("Foursquare final response URL contract passed.")
