---
title: {{ replace .Name "-" " " }}
slug: {{ .Name }}
date: {{ dateFormat "2006-01-02" .Date }}
url: {{ md5 .Name }}
aliases:
    - %INSERT
lang: %INSERT
tags:
    - %INSERT
summary: >-
    %INSERT
thumbnail:
mathjax: false
hljs: false

# For info, these are mostly not needed:
# draft: true
# unlisted: true

---

> **TL;DR** \
> Here tl dr

