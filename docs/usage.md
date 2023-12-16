---
icon: material/play-circle-outline
---

# Launching your session

After you've [configured your session](configuration.md), you can spin it up with the `tmpl` command:

```console title="Launching a session"
user@host:~/project$ tmpl
13:37:00 INF configuration file loaded path=/home/user/project/.tmpl.yaml
13:37:00 INF session created session=project
13:37:00 INF window created session=project window=project:code
13:37:00 INF window send-keys cmd=./scripts/init-env<cr> session=project window=project:code
13:37:00 INF window send-keys cmd="nvim .<cr>" session=project window=project:code
13:37:00 INF pane created session=project window=project:code pane=project:code.1 pane_width=192 pane_height=11
13:37:00 INF pane send-keys cmd=./scripts/init-env<cr> session=project window=project:code pane=project:code.1 pane_width=192 pane_height=11
13:37:00 INF pane send-keys cmd=./scripts/test-watcher<cr> session=project window=project:code pane=project:code.1 pane_width=192 pane_height=11
13:37:00 INF window created session=project window=project:shell
13:37:00 INF window send-keys cmd=./scripts/init-env<cr> session=project window=project:shell
13:37:00 INF window send-keys cmd="git status<cr>" session=project window=project:shell
13:37:00 INF window selected session=project window=project:code
13:37:00 INF switching client to session windows=2 panes=1 session=project
```

Tmpl attaches to the new session quite quickly, so you likely won't see the output. This video shows how it looks when
running the command:

<figure>
  <video controls>
    <source src="../assets/videos/demo.webm" type="video/webm" />
    <source src=",./assets/videos/demo.mp4" type="video/mp4" />
    <div class="admonition failure">
      <p class="admonition-title">Oh no, your browser doesn't support HTML video</p>
      <p>Here is a <a href="../assets/videos/demo.mp4">link to the video</a> instead.</p>
    </div>
  </video>
  <figcaption>Tmpl creating the session with Neovim and test runner ready.</figcaption>
</figure>

## Shared and global configurations

When tmpl searches for a configuration file, it scans the directory tree upward until it locates one or reaches the root
directory. This allows you to position configurations at different directory levels, serving as shared configurations if
needed.

``` title="Example directory structure" hl_lines="2 6 13"
home/user/
├── .tmpl.yaml (catch-all/default configuration)
└── projects/
    ├── work/
    │   ├── project_group/
    │   │   ├── .tmpl.yaml (used by projects in this group)
    │   │   ├── project_a/
    │   │   ├── project_b/
    │   │   └── project_c/
    │   ├── project_d/
    │   └── project_e/
    └── private/
        ├── .tmpl.yaml (used by private projects)
        ├── project_f/
        ├── project_g/
        └── project_h/
```

## Testing and verifying configurations

When creating a new configuration, it can be useful to ensure that it functions correctly without actually creating and
attaching a new session. Tmpl offers a dry-run mode for this purpose.

```console title="Launching a session in dry-run mode" hl_lines="3"
user@host:~/project$ tmpl --dry-run
13:37:00 INF configuration file loaded path=/home/user/project/.tmpl.yaml
13:37:00 INF DRY-RUN MODE ENABLED: no tmux commands will be executed and output is simulated
13:37:00 INF session created session=project dry_run=true
13:37:00 INF window created session=project window=project:code dry_run=true
13:37:00 INF window send-keys cmd=./scripts/init-env<cr> session=project window=project:code dry_run=true
13:37:00 INF window send-keys cmd="nvim .<cr>" session=project window=project:code dry_run=true
13:37:00 INF pane created session=project window=project:code pane=project:code.1 pane_width=40 pane_height=12 dry_run=true
13:37:00 INF pane send-keys cmd=./scripts/init-env<cr> session=project window=project:code pane=project:code.1 pane_width=40 pane_height=12 dry_run=true
13:37:00 INF pane send-keys cmd=./scripts/test-watcher<cr> session=project window=project:code pane=project:code.1 pane_width=40 pane_height=12 dry_run=true
13:37:00 INF window created session=project window=project:shell dry_run=true
13:37:00 INF window send-keys cmd=./scripts/init-env<cr> session=project window=project:shell dry_run=true
13:37:00 INF window send-keys cmd="git status<cr>" session=project window=project:shell dry_run=true
13:37:00 INF window selected session=project window=project:code dry_run=true
13:37:00 INF switching client to session windows=2 panes=1 session=project dry_run=true
```

!!! tip "Tip: debug mode"
    To see even more information about what tmpl is doing, including exact tmux commands being run, use the `--debug`
    flag. This also works in dry-run mode.

If you just want to verify that your configuration is valid, you can use the `check` sub-command:

```console title="Checking a configuration for errors"
user@host:~/project$ tmpl check
13:37:00 INF configuration file loaded path=/home/user/project/.tmpl.yaml
13:37:00 ERR configuration file is invalid errors=1
13:37:00 WRN session.name must only contain alphanumeric characters, underscores, dots, and dashes field=session.name
13:37:00 WRN session.windows.0.panes.0.env "my-env" is not a valid environment variable name field=session.windows.0.panes.0.env
13:37:00 WRN session.windows.1.path directory does not exist field=session.windows.1.path
```

## Command usage help

To see available commands, options, and usage examples for tmpl, you can use the `-h/--help` flag. This can also be used
on sub-commands.

```console title="Getting usage information for tmpl"
user@host:~$ tmpl --help
--8<-- "cli-usage.txt"
```

## That's it

This wraps up the getting started guide. Check out the recipe on [making a project launcher] for an example of how to
use tmpl in combination with other command-line tools to further streamline your workflow.

[making a project launcher]: <recipes/project-launcher.md>
