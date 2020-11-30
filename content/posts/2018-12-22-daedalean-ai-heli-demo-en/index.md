---
title: 'Daedalean AI helicopter flight tests (demo video)'
thumbnail: 'test-area.jpg'
lang: en
date: 2018-12-22
slug: daedalean-ai-heli-demo-en
tags:
    - aviation
    - news
summary: 'Just in case you wondered what we do at Daedalean AI, working
towards the airworthy certifiable visual autopilot.'
---

> **TL;DR** \
> Video of one of numerous helicopter test flights at [daedalean.ai](https://daedalean.ai)
> to check the functionality of visual odometry and semantic terrain recognition.

{{< includeimage img="ddln.png">}}

At [Daedalean](https://daedalean.ai), we work towards creating self-flying cars.
This would sound familiar and, eghmmm... easy for literally everyone who is into AI nowadays: yes,
you constantly hear that it is just about throwing bunch of Nice Data&trade; into some
Deep Network&trade; — and you're done! Ok, I exaggerate a lot, but the optimism
is still exceptionally high.

{{< includevideo url="https://www.youtube.com/embed/KhZkHfo4Px4"
caption=`Flight in the region of heliport Schindellegi (CH) and aerodrome Schänis (CH).
Certifiable airworthy [^1] system on board, performing real-time visual (no GPS involved!)
odometry [^2] and localization **and** semantic segmentation, allowing to constantly monitor
terrain in search for emergency landing sites.` >}}

What is usually at least under-highlighted is the amount of non-trivial work
needed to get this done: gather the data, manage it, develop both ML and
traditional  algorithms, test them, tune them, and prove correctness. Add here a
hard path of certifying such systems — and you will  get an impressive list of
engineering, organizational and management challenges  on the way from an idea
to a full realization.

{{< includeimage img="test-area.jpg" caption=`Approximate region of testing
features various types of terrain with numerous types of objects and landmarks.` >}}

Since July 2018, I'm lucky to be on a team of really professional and passionate
people who work on one of the trickiest applications of AI. As part of our routines,
we have to collect unique data, test our solutions under real conditions,
identify wins and fails, and re-iterate. Daedalean has issued a short video snippet on
how our pilots and engineers work together. Enjoy!



## References

[^1]: Airworthiness — a legal term from air law defining conditions that must be
    met by a system to be allowed to be used for air operations.

[^2]: Odometry — process of estimation the changes in position over time (e.g.
    path traveled.
