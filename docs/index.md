# tmpl - simple tmux session management

Tmpl streamlines your tmux workflow by letting you describe your sessions in simple YAML files and have them
launched with all the tools your workflow requires set up and ready to go. If you often set up the same windows and
panes for tasks like coding, running unit tests, tailing logs, and using other tools, tmpl can automate that for you.

## Highlights

- **Simple and versatile configuration:** easily set up your tmux sessions using straightforward YAML files, allowing
  you to create as many windows and panes as needed. Customize session and window names, working directories, and
  start-up commands.

- **Inheritable environment variables:** define environment variables for your entire session, a specific window, or a
  particular pane. These variables cascade from session to window to pane, enabling you to set a variable once and
  modify it at any level.

- **Custom hook commands:** customize your setup with on-window and on-pane hook commands that run when new windows,
  panes, or both are created. This feature is useful for initializing a virtual environment or switching between
  language runtime versions.

- **Non-intrusive workflow:** while there are many excellent session managers out there, some of them tend to be quite
  opinionated about how you should work with them. Tmpl allows configurations to live anywhere in your filesystem and
  focuses only on launching your session. It's intended as a secondary companion, and not a full workflow replacement.

- **Stand-alone binary:** Tmpl is a single, stand-alone binary with no external dependencies, except for tmux. It's easy
  to install and doesn't require you to have a specific language runtime or package manager on your system.

<div class="next-cta" markdown>
[Get started with tmpl :material-arrow-right-circle:](getting-started.md){ .md-button .md-button--primary }
</div>
