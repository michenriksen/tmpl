---
icon: material/download
---

# Getting started with tmpl

This guide walks you through the process of installing tmpl and creating your first session configuration file. Tmpl
is written in Go and distributed as a single, stand-alone binary with no external dependencies other than tmux, so
there is no language runtime or package managers to install.

!!! info "Beta software"
    Tmpl is currently in beta. It's usable, but there are still some rough edges. It's possible that breaking changes
    are introduced while tmpl is below version 1.0.0. If you find any bugs, please [open an issue][new-issue], and if
    you have any suggestions or feature requests, please start a new [idea discussion][new-idea].

## Installation

### Binaries

Go to the [releases page] and download the latest version for your operating system. Unpack the archive and move the
`tmpl` binary to a directory in your PATH, for example, `/usr/local/bin`:

```console title="Installing the binary"
user@host:~$ tar xzf tmpl_*.tar.gz
user@host:~$ sudo mv tmpl /usr/local/bin/
```

<!-- vale Google.Contractions = NO -->
??? failure "macOS *'cannot be opened'* dialog"
    <!-- vale on -->
    macOS may prevent you from running the pre-compiled binary due to the built-in security feature called Gatekeeper.
    This is because the binary isn't signed with an Apple Developer ID certificate.

    **If you get an error message saying that the binary is from an unidentified developer or something similar, you can
    allow it to run by doing one of the following:**

    1. :material-apple-finder: **Finder:** right-click the binary and select "Open" from the context menu and confirm
       that you want to run the binary. Gatekeeper remembers your choice and allows you to run the binary in the
       future.
    2. :material-console: **Terminal:** add the binary to the list of allowed applications by running the following
       command:

    ```console
    user@host:~$ spctl --add /path/to/tmpl
    ```
<!-- vale Google.Contractions = YES-->

### From source

If you have Go installed, you can also install tmpl from source:

```console title="Installing from source"
user@host:~$ go install {{ package }}@latest
```

<div class="next-cta" markdown>
[Next: configuring your session :material-arrow-right-circle:](configuration.md){ .md-button .md-button--primary }
</div>

[new-issue]: <{{ repo_url }}/issues/new/choose>
[new-idea]: <{{ repo_url }}/discussions/categories/ideas>
[releases page]: <{{ repo_url }}/tmpl/releases>
