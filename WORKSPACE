load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Update these to latest
RULES_HUGO_COMMIT = "fe024c327cfad486151571fed5a109fdda504130"

RULES_HUGO_SHA256 = "4039295de2d12b041a84140446c836c3f9fa11e360f88ee5b2db934df1d4e041"

http_archive(
    name = "build_stack_rules_hugo",
    sha256 = RULES_HUGO_SHA256,
    strip_prefix = "rules_hugo-%s" % RULES_HUGO_COMMIT,
    url = "https://github.com/stackb/rules_hugo/archive/%s.zip" % RULES_HUGO_COMMIT,
)

load("@build_stack_rules_hugo//hugo:rules.bzl", "hugo_repository")

#
# Load hugo binary itself
#
# Optionally, load a specific version of Hugo, with the 'version' argument
HUGO_VERSION = "0.121.1"

hugo_repository(
    name = "hugo_darwin",
    os_arch = "darwin-universal",
    version = HUGO_VERSION,
)

hugo_repository(
    name = "hugo_linux",
    os_arch = "Linux-64bit",
    version = HUGO_VERSION,
)
