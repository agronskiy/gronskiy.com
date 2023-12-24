---
title: Two ways to declare Rust binaries in Bazel
slug: 2023-12-declaring-rust-binaries-in-bazel-two-ways
date: 2023-12-24
lang: en
tags:
    - rust
    - bazel
summary: >-
  A short write-up with explaining two ways of declaring `rust_binary` targets when 
  working with Rust from Bazel, while maintaining Cargo-compliant structure. Bazel 
  allows way more freedom, which comes at a cost.
mathjax: false
hljs: true

---

> **TL;DR** \
> In Rust, according to the documentation, one can have several binaries per crate. 
    Classical `Cargo`-compliant structure is quite strict and relatively flat.
>   
>    Since I use `Bazel` as my build system for the reasons explained in the course of 
    the post, I explored a bit how to declare build targets in such a way that they 
    enjoy best of two worlds: place `BUILD` files close to the place of their utilization,
    while maintaining a `Cargo`-compliant structure of the project.

{{% toc %}}

## Premises

### Why Bazel

Expanding on why I use `Bazel`, in addition to great resources one can find[^rust-bazel], here's my own motivation, as a software developer.

`Cargo`, being a perfect tool for the dependency management, quickly becomes limiting in cases where I want to perform a sequence of dependent steps:
- building a docker on top of a binary is a simple yet good example;
- another, more complex, example, consists of an integration test between binaries built in different languages, say a Rust gRPC and a simple Python client calling it;
- to make things even more complex, and integration test consisting of building multiple dockers representing middleware and backend, running a docker-compose and testing the whole 
  interaction[^plan].

All the above share several traits:
- usually, *multiple* frameworks are involved (Rust + Docker, Rust + Python + Testing, Rust + Python + Docker + Testing)
- usually, *dependent* binaries/libraries are used
Both are addressed by a DAG-backed build system, whose nodes can be representing *any* task in *any* framework. 
There are more of those (in particular, I like Pantsbuild[^pants]), but `Bazel`, among them, is what I think is truly battle-tested. 
As anything battle-tested, it is *ugly* or at least not that *shiny* when it comes to diving deep into it, but... that's what we have to deal with.


### Constraints I self-impose
I still want to be able to use `Cargo` and it's typical structure, this is beneficial for several reasons:
	- compatibility with `Cargo` ecosystem
	- ability to use Rust Language Server a.k.a. `rust_analyzer`

## Boilerplate: main binary, sub-binaries and setup

The example of the repo is located at the [Github][github]. From the `Cargo` prospective, it follows the standard tree structure:
{{% codecaption "Source tree we use, see [GitHub][github]" %}}
```console
.
├── BUILD
├── Cargo.lock
├── Cargo.toml
├── main_bin_crate
│  ├── BUILD
│  ├── Cargo.toml
│  └── src
│     ├── bin
│     │  ├── binary_one  
│     │  │  └── main.rs  # <-- first sub-binary
│     │  └── binary_two  
│     │     ├── BUILD
│     │     └── main.rs  # <-- second sub-binary
│     └── main.rs        # <-- main  binary
└── WORKSPACE
```
{{% /codecaption %}}

The content of the binaries is trivial `println!`. Running `Cargo` with all the three shows the expected:
{{% codecaption "Example of running via `Cargo`" %}}
```console
❱ cargo run --bin main_bin_crate
    Finished dev [unoptimized + debuginfo] target(s) in 0.00s
     Running `target/debug/main_bin_crate`
Hello from main_bin_crate

❱ cargo run --bin binary_one
    Finished dev [unoptimized + debuginfo] target(s) in 0.00s
     Running `target/debug/binary_one`
Hello from binary_one

❱ cargo run --bin binary_two
    Finished dev [unoptimized + debuginfo] target(s) in 0.00s
     Running `target/debug/binary_two`
Hello from binary_two
```
{{% /codecaption %}}

Additionally, we see that `Cargo` creates the following output structure (this will be important below):
{{% codecaption "Output structure when running via `Cargo`" %}}
```console
.
└── target/
    └── debug/
        ├── binary_one*
        ├── binary_two*
        └── main_bin_crate*
```
{{% /codecaption %}}

In other words, all the three binaries are compiled as *top-level* ones despite residing on different levels
of the source tree.

## Bazel implementation

We operate more or less standard boilerplate of the `WORKSPACE` file, which in a nutshell loads the 
necessary metadata from both `Cargo.toml` files:

```python
load("@rules_rust//crate_universe:defs.bzl", "crates_repository")

crates_repository(
    name = "crate_index",
    cargo_lockfile = "//:Cargo.lock",
    lockfile = "//:Cargo.Bazel.lock",
    manifests = [
        "//:Cargo.toml",
        "//main_bin_crate:Cargo.toml",
    ],
)

load("@crate_index//:defs.bzl", "crate_repositories")

crate_repositories()
```



### Option 1: The closest to `Cargo` implementation

The first option consists in placing all the `Bazel` targets into the crate-level `./main_bin_crate/BUILD` file:

{{% codecaption "BUILD file from the top-level crate `./main_bin_crate/BUILD`, code from [GitHub][github]" %}}

```python
load("@rules_rust//rust:defs.bzl", "rust_binary")
load("@crate_index//:defs.bzl", "all_crate_deps")

rust_binary(
    name = "main_bin_crate",
    srcs = ["src/main.rs"],
    deps = all_crate_deps(),
)

# This is the first option to define another binary
rust_binary(
    name = "binary_one",
    srcs = ["src/bin/binary_one/main.rs"],
    deps = all_crate_deps(),
)

```
{{% /codecaption %}}

After building the project with 
```console
❱ bazel build //main_bin_crate/...
```
we can observe the output tree structure that `Bazel` puts it in (compare with the above one created by `Cargo`)
other output such as manifests and runfiles omitted:
```plaintext
bazel-bin/main_bin_crate/
├── binary_one*
└── main_bin_crate*
```
as we see, this is a perfect option to fully imitate the way `Cargo` does it, with all the binaries 
at the top level of the build.

However, this clearly goes somewhat against the habit from the `Bazel` world, where one places `BUILD` files 
as close to it's sources level as possible. To that end, we proceed to the second option.

### Option 2: `Bazel` spirit, put `BUILD` where the sources are

Initially placing the following code under `./main_bin_crate/src/bin/binary_two/BUILD`
```python
load("@rules_rust//rust:defs.bzl", "rust_binary")
load("@crate_index//:defs.bzl", "all_crate_deps")

# Wrong way of placing BUILD near the code
rust_binary(
    name = "binary_two",
    srcs = ["main.rs"],
    deps = all_crate_deps(),
)
```
results in an error
```console
❱ bazel build //main_bin_crate/src/bin/binary_two
...
Error in fail: Tried to get all_crate_deps for package main_bin_crate/src/bin/binary_two 
  but that package had no Cargo.toml file
WARNING: Target pattern parsing failed.
```
Likely, the root cause is the fact that `rules_rust` considers each package[^package] to map into a crate, and 
as we saw in the `WORKSPACE`, there's indeed not `Cargo.toml` for the `./main_bin_crate/src/bin/binary_two`. The 
necessary correction would be
```python
load("@rules_rust//rust:defs.bzl", "rust_binary")
load("@crate_index//:defs.bzl", "all_crate_deps")

# This is the second option
rust_binary(
    name = "binary_two",
    srcs = ["main.rs"],
    deps = all_crate_deps(
        package_name = "main_bin_crate",  # <--- this is the added line
    ),
)
```
This compiles, and the compiled output structure is as follows (again, auxiliary files 
like manifests and runfiles are omitted):
```plaintext
bazel-bin/main_bin_crate/
├── main_bin_crate*
└── src/
    └── bin/
        └── binary_two/
            └── binary_two*
```
As we see, the only inconvenience consists in the fact that the output binary tree layout is now dissimilar 
to that of what `Cargo` created. This is a minor inconvenience provided, especially if you're developing
a non-publishable crate for some internal usage (e.g. dockerizing it later).

## Discussion

The following takeaways of this simple experiment, as I see them, would be

- it is possible to maintain `Cargo`-compliant project structure and yet apply `Bazel` 
  in its natural way -- by placing `BUILD` files close to their sources;
- `Bazel` philosophy is the superset of `Cargo`, allowing more flexibility about where to
  place various output units of the code (binaries, libraries). It definitely helps  for 
  a better structuring of your project;
- the above becomes especially important in larger monorepos, and in repos with multiple languages/frameworks 
  involved (think examples in the beginning, such as Rust + Python + Docker);
- the standardized approach that `Bazel` exhibits gives, in the end of the day, a great 
  freedom to focus on the tasks one solves regardless of the language/stack of a particular 
  project

[^pants]: Actually, a nice and beautiful reincarnation of a Pants version 1: https://www.pantsbuild.org

[^package]: Too ambiguous terminology here, *package* in `Bazel` world is actually any place in the source tree where a `BUILD` file
  would reside.

[^rust-bazel]: To list some: a) a great blog post by Roman Kashitsyn: [Scaling Rust builds with Bazel](https://mmapped.blog/posts/17-scaling-rust-builds-with-bazel.html)
    b) another one by Ilya Polyakovskiy: [Building a Rust workspace with Bazel](https://www.tweag.io/blog/2023-07-27-building-rust-workspace-with-bazel)

[^plan]: One of the near-future write-ups will be around building Docker with a Rust binary using cross-compilation, stay tuned.

[github]: https://github.com/agronskiy/rust-etudes/tree/main/rs_multiple_binaries
