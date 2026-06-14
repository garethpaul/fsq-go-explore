#!/usr/bin/env python3
import sys
from pathlib import Path


source = Path(sys.argv[1]).read_text(encoding="utf-8")
tests = Path(sys.argv[2]).read_text(encoding="utf-8")
plan = Path(sys.argv[3]).read_text(encoding="utf-8")

required_source = [
    '"mime"',
    "func isFoursquareJSONResponse(response *http.Response) bool",
    'mime.ParseMediaType(response.Header.Get("Content-Type"))',
    'mediaType == "application/json"',
    'strings.HasSuffix(mediaType, "+json")',
]
for fragment in required_source:
    if fragment not in source:
        raise SystemExit("Foursquare content-type boundary missing: " + fragment)

for method, result in (("Search", "venues"), ("VenueDetails", "venue")):
    start = source.find("func (fsqs *FoursquareService) " + method)
    end = source.find("\nfunc ", start + 1)
    body = source[start : None if end == -1 else end]
    status = body.find("successfulFoursquareStatus")
    media = body.find("isFoursquareJSONResponse")
    decode = body.find("decodeFoursquareResponse")
    if -1 in (start, status, media, decode) or not status < media < decode:
        raise SystemExit(method + " must validate status and media type before decoding.")
    if ("return " + result) not in body[media:decode]:
        raise SystemExit(method + " must fail closed before decoding non-JSON content.")

required_tests = [
    "func TestFoursquareJSONResponseMediaTypes",
    '"application/vnd.foursquare+json"',
    '"text/html"',
    "func TestSearchRejectsNonJSONResponseBeforeDecode",
    "func TestVenueDetailsRejectsNonJSONResponseBeforeDecode",
]
for fragment in required_tests:
    if fragment not in tests:
        raise SystemExit("Foursquare content-type regression missing: " + fragment)

for evidence in ("status: completed", "hostile mutations were rejected", "make check"):
    if evidence not in plan:
        raise SystemExit("Foursquare content-type plan missing: " + evidence)

print("Foursquare response content-type checks passed.")
