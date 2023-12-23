load("@build_stack_rules_hugo//hugo:rules.bzl", "hugo_serve", "hugo_site")

HUGO_BIN = select({
    "@platforms//os:linux": "@hugo_linux//:hugo",
    "@platforms//os:macos": "@hugo_darwin//:hugo",
})

# Note, here we are using the config_dir attribute to support multi-lingual configurations.
hugo_site(
    name = "site_complex",
    config = "//:config.yaml",
    content = glob(["content/**"]),
    data = glob(["data/**"]),
    hugo = HUGO_BIN,
    quiet = False,
    static = glob(["static/**"]),
    theme = "//themes/gron",
)

# Run local development server
hugo_serve(
    name = "serve",
    dep = [":site_complex"],
    hugo = HUGO_BIN,
)
