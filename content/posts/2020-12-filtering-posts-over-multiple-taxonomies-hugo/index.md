---
title: "Filtering posts over multiple taxonomies in Hugo"
slug: filtering-posts-over-multiple-taxonomies-hugo
date: 2020-12-23
lang: en
tags:
    - web
    - hugo
    - go
thumbnail: tags.png
mathjax: true
hljs: true
---

> **TL;DR** \
> We will find out how to implement lists of posts based on a filter, which consists of several tags from several taxonomies – something Hugo does not provide out of the box. This is an interesting exercise in itself, and it represents a widespread use-case.

{{% toc %}}

## Motivation

Here I'll describe the implementation of filtering by a pair of tags, such as in
the screenshot below:

{{% includeimage img="tags.png" style="width: 50%" %}}

Note that here the "Language" is a *tag* and is not related to Hugo's [multilanguage
feature](https://gohugo.io/content-management/multilingual/): as opposed to that,
on my web page I wanted posts in all languages to be mixed by default, with an option
to filter those – just as with normal tag.

## What Hugo can do

Operating on *taxonomies* – which are hierarchies of tags – is where Hugo shines. Page kinds *taxonomy* and *taxonomyTerm* (explained below in a second) are first-class citizens in Hugo’s system of types.￼

To bring an example: with minimal modifications of your `config.yaml￼` (to be precise,in this particular example no modifications are needed – since “tags” is one of two default taxonomies￼), it is enough to add e.g.
```yaml
tag:
  - aviation
  - hugo
```
to some page’s front matter, and Hugo will create
- *taxonomy* page with URL `example.com/tags` and
- *taxonomyTerm* pages with URLs `example.com/tags/aviation` and `example.com/tags/hugo`,
and those pages will have associated list templates. Pretty automatic.

As already stated, I wanted filtering by *pair* of tags simultaneously from *two* different taxonomies, and Hugo’s built-in capabilities are not enough, unfortunately.

## Existing solutions

To the best of my knowledge (and even after extensive search I might have not hit some
obvious pointers), the only solution that implements filtering by multiple tags is
[Pointy](https://pointy.netlify.app/filter/).

{{% includeimage img="pointly.png"
    caption=`One can see that multiple tags from the same taxonomy can be selected.` %}}

This is a JavaScript-based solution, where filtering itself happens dynamically after loading all the posts from all the taxonomies.

Pros:

- solution is very elegant and clean
- solution is *inter*-taxonomy-scalable: it works on any amount of taxonomies
- solution is *intra*-taxonomy-scalable: it works when one wants to filter by two tags from the same taxonomy

Cons:

- it is unclear what performance implication dynamic filtering has on large amounts of posts (should not be an issue for small amounts, however)
- there's no immediate way to have permalinks for filtered lists
- per-tag (similarly, per-{set-of-tags}) RSS is not immediately available.

## Static solution

I wanted (well, partially out of some curiosity, not necessarily driven by pragmatism)
to develop some fully static solution in the spirit of Hugo. Inevitably creating some of
its own cons (discussed later), this solution would mitigate cons of the above solution's ones.

Below, we will create a solution which allows to list posts on a cross-product of `tag` and `lang` dimensions under such URLs as:

{{% table %}}
| URL | meaning |
|-------|---|
| `/posts` | all tags in all languages |
| `/posts/ru` | all tags in `ru` language |
| `/tags/aviation` | tagged as `aviation` in all languages |
| `/tags/aviation/ru` | tagged as `aviation` in `ru` language |
{{% /table %}}

### Directory structure

Let's introduce two subdirectories under Hugo's `/content` directory,

```plaintext
content/
├── posts/
│   └── some-post/
│       ├── index.md        => URL: posts/some-post
│       └── picture.png
└── tags/
    ├── _index.md           => URL: posts/      [*]
    ├── ru/
    │   └── _index.md       => URL: posts/ru    [*]
    ├── en/
    │   └── _index.md       => URL: posts/en    [*]
    ├── aviation/
    │   ├── _index.md       => URL: tags/aviation
    │   ├── ru/
    │   │   └── _index.md   => URL: tags/aviation/ru
    │   └── en/
    │       └── _index.md   => URL: tags/aviation/en
    └── hugo/
        ├── _index.md       => URL: tags/hugo
        ├── ru/
        │   └── _index.md   => URL: tags/hugo/ru
        └── en/
            └── _index.md   => URL: tags/hugo/en
```
Note `[*]`-marked lines: they signify non-trivial permalinks. This is an optional "sugar" feature I wanted to have: "all tags" should be available on `/posts` URL for simplicity, not on `/tags`.

### Frontmatter of list pages

Each of these list pages should contain only front matter which is later used by the `list` template to
perform actual filtering. E.g.

{{% codecaption "Content of `/tags/aviation/en/_index.md`" %}}
```yaml
---
tag: aviation
lang: en
---
```
{{% /codecaption %}}

or, for versions where either `tag` or `lang` is not filtered (i.e. we want to display all),

{{% codecaption "Content of `/tags/en/_index.md`" %}}
```yaml
---
tag: all        # Note `all` here.
lang: en
---
```
{{% /codecaption %}}

Now, remember the highlighted lines in the directory structure above: we want `/tags/_index.md`, `/tags/en/_index.md` and `/tags/ru/_index.md` to translate into special URL so for *these* files only we add, respectively, `url: /posts`, `url: /posts/ru` or `url: /posts/en`, such as

{{% codecaption "Content of `/tags/_index.md`" %}}
```yaml
---
tag: all
lang: all
url: /posts     # Note `url:` here.
---
```
{{% /codecaption %}}

> **CAVEAT**: \
    there was a nontrivial consequence of redirecting permalinks from one section (`/tags`) to another
    (`/posts`), which is related to templates. I will talk about it later.

### Generating tag pages

Now we need to automate the process of creating the structure under `/tags`, because it is obviously very cumbersome to do it manually. For that, I implemented a small Go module[^go] that concurrently scans all the content pages, checks front matter for tags and creates the necessary `_index.md` under `/tags`.

Overall structure follows:

```plaintext
               +--------------+
               |  inputQueue  |      file paths
               +--------------+
                      |
        +--------------------------+
        |      |      |     |      |
        |      |      |     |      |
    +---v------v------v-----v------v---+
    |                                  |
    |         parallel workers         |
    |                                  |
    +----------------------------------+
        |      |      |     |      |
        |      |      |     |      |
        ^------v------------v------^
                      |
                      |
            +---------v-----------+
            |  intermediateQueue  |
            +---------------------+
                      |
                      |
              +-------v-------+
              |  outputQueue  |       tags
              +---------------+
```

This is a separate interesting topic of correctly implementing concurrency in Go, and this topic
itself is worth looking into, so will dive into it in another post. Here, just for completeness, I'll provide short skeleton snippets of the main logic.

{{% codecaption `Defining trivial types for input of the
    pipeline and output thereof` %}}
```golang
type (
    inputPath  string
    outputTag  string
)
```
{{% /codecaption %}}

{{% codecaption `Concurrency logic behind each worker that
    parses files` %}}
```golang
// worker takes the page file paths from the `inputQueue` and sends them to
// `intermediateQueue`. Additionally, it keeps notifying the runner when it
// stops via `counterCh`.
func worker(
    inputQueue <-chan inputPath,
    intermediateQueue chan<- outputTag,
    counterCh chan<- int,
) {
    // The way to report to the runner when the worker starts and finishes. This is
    // done so as to stop the pipeline gracefully only after all the workers drain
    // the inputPath and finish.
    counterCh <- 1
    defer func() {
        counterCh <- -1
    }()

    //
    for path := range inputQueue {

        // ... here happens the extraction of the YAML frontmatter from the `path`
        // ... and outputting to the streamlined `intermediateQueue`.
    }
}
```
{{% /codecaption %}}

{{% codecaption `Main runner orchestrating concurrency and gracefully
    closing output channels` %}}
```golang
// makeRunner creates `inputQueue` and `outputQueue` channels, inbetween which several concurrent
// workers are launched. Control of the pipeline is done via closing the `inputPath`
// channel, and the workers are closed gracefully after draining all the inputPath tasks.
func makeRunner() (chan<- inputPath, <-chan outputTag) {
    var (
        numCPU = runtime.NumCPU()

        inputQueue        = make(chan inputPath, numCPU)
        intermediateQueue = make(chan outputTag, numCPU)
        outputQueue       = make(chan outputTag, numCPU)

        counterCh = make(chan int)

        numOpenWorkers = 0
    )

    // This function allows to first drain the intermediate queue after the
    // inputPath queue was closed.
    stopGracefully := func() {
        for {
            // First, drain remaining results, and only then stop.
            select {
            case out := <-intermediateQueue:
                outputQueue <- out
            default:
                close(outputQueue)
                return
            }
        }
    }

    // Create actual runners. More precisely:
    //  1.  spawns workers
    //  2a. listens to their outputTag
    //  2b. does bookkeeping, counts how many are still working and closes
    //      outputTag channel
    go func() {
        for i := 0; i < numCPU; i++ {
            go worker(inputQueue, intermediateQueue, counterCh)
        }

        for {
            select {
            case out := <-intermediateQueue:
                outputQueue <- out
            case n := <-counterCh:
                numOpenWorkers += n
                if numOpenWorkers > 0 {
                    continue
                }
                stopGracefully()
            }
        }
    }()

    return inputQueue, outputQueue
}
```
{{% /codecaption%}}

{{% codecaption `Output processing and main` %}}
```golang
func processOutput(outputQueue <-chan outputTag) {

    // This will create files under 'tags/...'
    for res := range outputQueue {
        // ... process each tag and write files for each language.
    }
}

func main() {
    inputQueue, outputQueue := makeRunner()

    go func() {
        paths := make([]string, 0)
        pathsGlob, err := filepath.Glob("../../content/posts/*/*.md")
        if err != nil {
            return
        }
        paths = append(paths, pathsGlob...)
        close(inputQueue)  // This controls the closing.
    }()

    processOutput(outputQueue)
}
```
{{% /codecaption %}}

### Permalink conflict for the tag page mapped to `/posts`

In the above, remember the highlighted part from the directory structure above which signifies what I call *permalink
conflict*.
{{% codecaption `Relevant extract from the content directory above` %}}
```plaintext {hl_lines=[2]}
content/
├── posts/           =>  ??? conflict here
│   ├── some-post/
│   └── some-other-post/
└── tags/
    ├── _index.md    => URL: posts/
```
{{% /codecaption %}}
Hugo uses quite complex [layout lookup rules](https://gohugo.io/templates/lookup-order/) to determine the mapping between

- each page from `/content/` folder and
- the way to render it, called *layout* amd residing under `/themes/<theme_name>/layouts`.

In the example above, conflict arises from the following facts:
- the fact that page `/tags/_index.md` has a permalink of `/posts`, and
- it has the `list` layout template (because it represents a so called *branch bundle*[^bundle]), among others from `/layouts/tags/list.html` and
- the fact that there is an *implicitly generated* section page `/posts`, and
- this implicitly generated section page uses, among others, the `list` template, among others from `/layouts/posts/list.html` too.

As we see, two pages with the same type of layout - list - are mapped to `/posts` URL. Or, to put it differently, there is a page from `/tags` section which "steals" the URL of the section page of `/posts`, making it unclear which layout is used. I had to create a bit of
experimentation playground[^experiment], which explores various combinations of existing and non-existing layouts. In short,
`/layouts/posts/list.html` is preferred, and if one deletes it, Hugo falls back to `/layouts/tags/list.html` *but* throws a warning
when building the website:

```terminfo
➙  hugo serve
Start building sites …
WARN 2020/12/15 10:41:26 found no layout file for "HTML" for kind "section":
You should create a template file which matches Hugo Layouts Lookup Rules for this combination.
```

which is technically not true, because as I said, Hugo still finds the layout for it and renders it.
The solution I have found[^dis] consists in creating an *empty* `/layouts/posts/list.html` (not even a whitespace there!)
which effectively disables using this layout. To cite the Hugo's discussion thread:

>  If `layouts/foo/single.html` and `layouts/section/foo.html` are zero-length, no files will be published for section
`foo`, but the content can still be retrieved through taxonomies, etc. This was briefly broken in 0.20-0.20.2, but was fixed in 0.20.3.

This solution is worth looking into, because it reveals an interesting property of Hugo: one can see undocumented but stable features to reach flexibility.

### Disabling default taxonomies

Last step is needed to avoid potential confusions for Hugo – which, as noted in the beginning – has some predefined taxonomies[^defaulttax]: `tags` and `categories`. By a coincidence, our method intersects with the former, so to avoid spurious undefined behavior such as implicit generation of taxonomy pages[^tax-gen], we need to disable that. This can be done in the config.

{{% codecaption "Section in the `config.yaml`" %}}
```yaml
taxonomies: []
```
{{% /codecaption %}}

As a result, Hugo won't try creating any potentially conflicting taxonomy list pages.

## Pros and cons, summary

### Pros

- this approach is fully static, no JavaScript involved;
- each combination of `tag` and `lang` has associated static permalink;
- each combination of `tag` and `lang` has associated resources such as RSS.

### Cons

- main "con" comes from the respective "pro": static pages means that number of those increases as
    the product of numbers of tags in each of two taxonomies. For example, having $N_{\mathrm{T}}$ tags of
    type `tag` and $N_{\mathrm{L}}$ of type `lang`, the total amount of generated pages will be
    $O(N_{\mathrm{T}} \times N_{\mathrm{L}})$. \
    \
    For out case, this does not pose significant problems, since $N_{\mathrm{L}} = \mathrm{const}$, but the
    approach is generally non-scalable;
- still no way to filter by more than one tag from a single taxonomy, e.g. if we wanted to filter all the
    posts with `tag` in `[aviation, hugo]`.

### Summary

All in all:

- we managed to imitate Hugo's implimentation of taxonomies, allowing e.g. `/tags/aviation` to list
    all posts with tag `aviation`;
- in addition, we managed to extend it by allowing such permalinks as `/tags/aviation/en` to list
    all posts with tag `aviation` and language `en`;
- the above approach allows to mix more than two taxonomies;
- why all this? Well – because we can :)

## References

[^go]: Available in [Github repo](https://github.com/agronskiy/website-public/tree/master/bin/tags-gen) of this website.

[^bundle]: More about bundles in the [Hugo documentation](https://gohugo.io/content-management/page-bundles/).

[^experiment]: Experiment is available in my [Github](https://github.com/agronskiy/hugo-experimentation).

[^dis]: See [Hugo discourse thread](https://discourse.gohugo.io/t/would-like-to-ignorefiles-but-use-the-underlying-content/6368/10).

[^defaulttax]: See ["Default taxonomies"](https://gohugo.io/content-management/taxonomies/#default-taxonomies) in the Hugo documentation.

[^tax-gen]: See ["Default destinations"](https://gohugo.io/content-management/taxonomies/#default-destinations) in the Hugo documentation.
