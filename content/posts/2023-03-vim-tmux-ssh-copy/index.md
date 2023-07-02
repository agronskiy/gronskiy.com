---
title: Seamless copy-paste between  tmux, vim and clipboard over ssh
slug: 2023-03-26-copy-via-vim-tmux-ssh
date: 2023-03-26
lang: en
tags:
    - tools
    - vim
summary: >-
  Allowing to copy into local clipboard over the `vim → tmux → ssh` chain (possibly multiple hops
  if working on a remote box from a remote box) increased my productivity from 95% to 100% :)
mathjax: false
hljs: true
asciinema: true

---

> **TL;DR**
>
> I spend a majority of my development in terminal: from my Mac which acts as a thin development client, I stay over `ssh` into `tmux` on
  my main Linux workstation and numerous other boxes, with my main IDE a.k.a. `neovim` running remotely. Sometimes, from
  that `tmux` I connect over another `ssh` and fire `tmux` on a box over two hops from my cosy Mac.
> 
> This worked smoothly with one common caveat: I need to be able, from time to time, to select-and-copy into my local
  clipboard (either to copy into some GUI for demos, or to copy into an `ssh`-over-`ssh` connected VM) -- keeping my 
  local clipboard an ulimate source of copy-paste storage. 
>
> This short post was born as a consequence of "shell user's block" -- reading numerous posts on the topic which make
  things look quite pessimistic in terms of the complexity, thus dicouraging me to give it a try. It all turned out to be
  quite trivial. So here, I aggregate a working set of steps -- and keep all the credits with the developers of the 
  plugins mentioned in the course of the story.

{{% includeascii file="demo.cast" id="terminal" bleed="true" opts="{cols: 122, rows: 22, preload: true, poster: 'npt:0:21'}" 
caption=`Example of copying from local <code>tmux</code> pane into <code>vim</code> on remote, and then
vice versa from remote <code>vim</code> into local <code>tmux + vim</code>.` %}}

{{% toc %}}

## Out of scope

As mentioned, we'll be setting up the ability to copy things *downstream*, i.e. over one or several `tmux`/`ssh` hops
into the *local* clipboard. It stays out of scope to operate the remote clipboards, primarily because this has never
popped up in my (quite intense) development activity -- it has always been enough to have things in my local clipboard
and paste them in `INSERT` mode into remote, if needed.

## OSC52 and OSCYanc

I was a happy user of a terminal that supports OSC52, an [ANSI Operating System Command][osc52] set of escape sequences
that, provided the support of the receiving terminal emulator, allow to copy text into the local buffer.

Another bit, [`vim-OSCYank`][yank] (and its Lua rewrite by the same author, [`nvim-osc52`][nvim]) is a wonderful plugin 
that trivializes sending the OSC52-enriched messages to the local terminal emulator.

## Recipe

### Set up the `vim-oscyank` 

In your `init.lua` or equivalent, you'd need to add the `vim-oscyank` plugin [see here][yank], and set it up along 
the following lines (here, the syntax is `Packer.nvim`s one):
```lua
use { "ojroques/vim-oscyank",
  config = function()
    -- Should be accompanied by a setting clipboard in tmux.conf, also see
    -- https://github.com/ojroques/vim-oscyank#the-plugin-does-not-work-with-tmux
    vim.g.oscyank_term = "default"
    vim.g.oscyank_max_length = 0  -- unlimited
    -- Below autocmd is for copying to OSC52 for any yank operation,
    -- see https://github.com/ojroques/vim-oscyank#copying-from-a-register
    vim.api.nvim_create_autocmd("TextYankPost", {
      pattern = "*",
      callback = function()
        if vim.v.event.operator == "y" and vim.v.event.regname == "" then
          vim.cmd('OSCYankRegister "')
        end
      end,
    })
  end,
}
```


### Set up the `tmux`

To enable interoperation of `tmux` with the clipboard, set in your `.tmux.conf`:
```bash
# Allow clipboard with OSC-52 work, see https://github.com/tmux/tmux/wiki/Clipboard
set -s set-clipboard on
```

Additionally, I prefer to set `y` in `tmux` scroll mode for copying:
```bash
# Use vim keybindings in copy mode
setw -g mode-keys vi
unbind -T copy-mode-vi MouseDragEnd1Pane

# Make `y` copy the selected text, not exiting the copy mode. For copy-and-exit
# use ordinary `Enter`
bind -T copy-mode-vi y send-keys -X copy-pipe  # Only copy, no cancel
```


## Bonuses, mouse selection in `neovim` and `tmux`

Don't forget to set the mouse selection in `init.lua`

```lua
vim.opt.mouse = "a"
```

It is annoying that when selecting things with mouse in `tmux` scroll mode, this selection sticks
until you select something else -- common UX is to drop selection on a single click anywhere:

```bash
# Clear selection on single click
bind -T copy-mode-vi MouseDown1Pane send-keys -X clear-selection \; select-pane
```

[osc52]: https://en.wikipedia.org/wiki/ANSI_escape_code#Escape_sequences
[yank]: https://github.com/ojroques/vim-oscyank
[nvim]: https://github.com/ojroques/nvim-osc52
