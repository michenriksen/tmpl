---
icon: material/cog
---

# Configuring your session

After you've [installed tmpl](getting-started.md), you can create your session configuration using the `tmpl init`
command. This creates a basic `.tmpl.yaml` configuration file in the current directory.

```console title="Creating a new configuration file"
user@host:~/project$ tmpl init
13:37:00 INF configuration file created path=/home/user/project/.tmpl.yaml

```

The file sets you up with a session named after the current directory with a single window called main:

```yaml title=".tmpl.yaml"
# Note: the real file has helpful comments which are omitted for brevity.
---
session:
  name: project

  windows:
    - name: main
```

This may be all you need for a simple project, but to get the most out of tmpl you'll want to customize your session to
set up as much of your development environment as possible. The following sections describe how to use the options to
bootstrap a more interesting session.

## Windows and panes

Tmpl sessions can have as many windows and panes as you need for your workflow. This example configures the session
with two windows, one named `code` and another named `shell`. The `code` window is also configured to have a pane at the
bottom with a height of 20% of the available space.

```yaml title=".tmpl.yaml" hl_lines="5 6 7 9"
session:
  name: project

  windows:
    - name: code
      panes:
        - size: 20%

    - name: shell
```

## Commands

It's possible to configure commands to automatically run in each window and pane. This example builds on the previous by
configuring the `code` window to automatically start Neovim and its bottom pane to run a fictitious test watcher. The
`shell` window is configured to run git status.

```yaml title=".tmpl.yaml" hl_lines="6 9 12"
session:
  name: project

  windows:
    - name: code
      command: nvim .
      panes:
        - size: 20%
          command: ./scripts/test-watcher

    - name: shell
      command: git status
```

!!! tip "Tip: multiple commands"
    Tmpl also supports running a sequence of commands if needed. Each command is sent to the window or pane using the
    [tmux send-keys][send-keys] command. This means that it also works when connecting to remote systems:

    ```yaml title=".tmpl.yaml"
    session:
      windows:
        - name: server-logs
          commands:
            - ssh user@remote.host
            - cd /var/logs
            - tail -f app.log
    ```

## Environment variables

Setting up a development environment often involves configuring environment variables. Tmpl allows you to set
environment variables at different levels, such as session, window, and pane. These variables cascade from session to
window to pane, making it easy to set a variable once and make changes at any level.

This example builds on the previous by setting up a few environment variables at different levels. `APP_ENV` and `DEBUG`
are set at the session level and cascade to all windows and panes. The pane overrides `APP_ENV` to `test` to run tests
in the correct environment. Finally, the shell window overrides `HISTFILE` to maintain a project-specific command
history file.

```yaml title=".tmpl.yaml" hl_lines="3 4 5 13 14 18 19"
session:
  name: project
  env:
    APP_ENV: development
    DEBUG: true

  windows:
    - name: code
      command: nvim .
      panes:
        - size: 20%
          command: ./scripts/test-watcher
          env:
            APP_ENV: test

    - name: shell
      command: git status
      env:
        HISTFILE: ~/project/command-history
```

## Hook commands

Another frequent step in setting up a development environment involves executing project-specific initialization
commands. These commands can range from activating a virtual environment to switching between different language runtime
versions. Tmpl lets you configure commands that run in every window, pane, or both, when they're created.

This example builds on the previous by setting up a hook command to run a fictitious script in every window and pane.

```yaml title=".tmpl.yaml" hl_lines="7 8"
session:
  name: project
  env:
    APP_ENV: development
    DEBUG: true

  # Run the init-env script in every window and pane.
  on_any: ./scripts/init-env

  windows:
    - name: code
      command: nvim .
      panes:
        - size: 20%
          command: ./scripts/test-watcher
          env:
            APP_ENV: test

    - name: shell
      command: git status
      env:
        HISTFILE: ~/project/command-history
```

!!! tip "Tip: on window and on pane hooks"
    It's also possible to specify hooks that only run in window or pane contexts for more granular control:

    ```yaml title=".tmpl.yaml"
    session:
      on_window: ./scripts/init-window
      on_pane: ./scripts/init-pane
    ```

## More options

This wraps up the basic configuration options for tmpl. You can find more details on the available options in the
[configuration reference](reference.md) section if you want to learn more.

<div class="next-cta" markdown>
[Next: launching your session :material-arrow-right-circle:](usage.md){ .md-button .md-button--primary }
</div>

[send-keys]: https://man.archlinux.org/man/tmux.1#send-keys
