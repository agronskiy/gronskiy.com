---
title: Experience of moving from Jekyll to Hugo
slug: moved-my-website-jekyll-to-hugo
date: 2020-11-23
lang: en
tags:
    - web
    - hugo
    - another
thumbnail: tags.png
---

> **TL;DR** \
> Moving this website from Jekyll to Hugo was sometimes non-trivial but fun exercise.

## Implementing cross-product of tags

Here I'll describe the implementation of filtering by a pair of tags, such as in
the screenshot below:

{{% includeimage img="tags.png" style="width: 50%" %}}

Note that here, the "Language" is a *tag* and is not related to Hugo's [multilanguage
feature](https://gohugo.io/content-management/multilingual/): as opposed to that,
on my web page I wanted posts in all languages to be mixed by default, with an option
to filter those – just as with normal tag.

### What Hugo can do

Operating on *taxonomies* – Which are hierarchies of tags is where Hugo shines. Page kinds `taxonomy` and `taxonomyTerm` (explained in a second) are first-class citizens in Hugo’s system of types.￼

To bring an example: with minimal modifications of your `config.yaml￼` (to be precise,in this particular example no modifications are needed – since “tags” is one of two default taxonomies￼), it is enough to add e.g.
```yaml
tag:
  - aviation
  - science
```
to some page’s front matter, and Hugo will create
- *taxonomy* page with URL `example.com/tags` and
- *taxonomyTerm* pages with URLs `example.com/tags/aviation` and `example.com/tags/science`,
and those pages will have associated list templates. Pretty automatic.

### My usecase

As already stated, I wanted filtering by *pair* of tags from *two* different taxonomies.

```plain
tags: `aviation`, `news`, `hugo`
languages: `ru`, `en`
```


