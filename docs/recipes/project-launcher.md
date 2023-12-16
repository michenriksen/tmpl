---
icon: material/rocket-launch
---

# Project launcher

This is a recipe for a shell script that combines tmpl with a few other command-line tools to create a spiffy project
launcher that presents your projects in a selection menu with fuzzy search support. Pressing <kbd>Enter</kbd> launches a
tmpl session for the selected project:

<figure>
  <video controls>
    <source src="../../assets/videos/launcher.webm" type="video/webm" />
    <source src=",./../assets/videos/launcher.mp4" type="video/mp4" />
    <div class="admonition failure">
      <p class="admonition-title">Oh no, your browser doesn't support HTML video</p>
      <p>Here is a <a href="../../assets/videos/launcher.mp4">link to the video</a> instead.</p>
    </div>
  </video>
  <figcaption>Selecting and launching a project with the launcher.</figcaption>
</figure>

## Prerequisites

This recipe assumes you have gone through the [getting started](../getting-started.md) guide and that you have a basic
understanding of shell scripting. It also assumes that your projects are located under a single directory and that
you are using a shell such as Bash or Zsh.

## Dependencies

The launcher makes use of a few other command-line tools, so you'll need to install those first. The tools are:

- [fd] - Fast and user-friendly alternative to `find`. `fd` is used to find projects to launch.
- [fzf] - General-purpose command-line fuzzy finder. `fzf` is used for the selection menu.

<!-- vale off -->
=== ":material-apple: Homebrew"

    ```console title=""
    user@host:~$ brew install fd fzf
    ```

=== ":material-ubuntu: APT"

    ```console title=""
    user@host:~$ sudo apt install fd-find fzf
    ```

=== ":material-arch: Pacman"

    ```console title=""
    user@host:~$ sudo pacman -S fd fzf
    ```

=== ":material-fedora: DNF"

    ```console title=""
    user@host:~$ sudo dnf install fd-find fzf
    ```

=== ":material-gentoo: Portage"

    ```console title=""
    user@host:~$ emerge --ask sys-apps/fd app-shells/fzf
    ```
<!-- vale on -->

<small>If your package manager isn't listed, refer to the installation instructions for [fd][fd-install] and
[fzf][fzf-install].</small>

## The script

This is the project launcher script. The script is heavily commented to explain how it works and key configuration
options are highlighted to make it easier to modify the script to your needs.

??? info "Installation instructions"

    1. Copy the script and save it as a file named `tmpo` in a directory that is included in your `$PATH`. For example,
       you could save it as `/usr/local/bin/tmpo` (may require `sudo` to write to that location).
    2. Modify `$projects_dir` to the directory where your projects are located.
    3. Make the script executable by running `chmod +x path/to/tmpo`.
    4. Verify that the script works by running `tmpo` in a terminal.

``` { .bash .copy title="Project launcher" hl_lines="14-17 19-23 25-31 45-58 60-63" }
--8<-- "recipes/project-launcher.sh"
```

## Tips and tricks

These are some optional tips and tricks to make the launcher even more useful.

### Key binding shortcut

Use the [bind-key tmux command][tmux-bind-key] to assign the launcher to a key for super quick access. Add the following
to your `tmux.conf` file to bring up the launcher with your prefix key (default <kbd>Ctrl</kbd>+<kbd>b</kbd>) followed
by <kbd>f</kbd>:

``` { .bash .copy title="tmux.conf" }
bind-key f run-shell "/path/to/project-launcher"
```

### Default tmpl configuration

The launcher script changes the working directory to the selected project root before launching tmpl. Because tmpl
traverses the directory tree upwards when looking for a configuration file, you can easily set up a default
configuration for all your projects by adding a `.tmpl.yaml` file in a shared parent directory. If some projects need
special configuration, you can override the default configuration by adding a `.tmpl.yaml` file in the project root.

### fzf preview window

A feature of fzf that's left out in the launcher script for simplicity, is the [preview window feature][fzf-preview]
which makes it possible to dynamically show the output of a command for the currently selected project. As an example,
modifying the fzf options to the following shows `git status` for the selected project:

``` { .bash .copy title="Project launcher with preview window" hl_lines="11-13" }
selected_project="$(
  get_project_data |
    fzf \
      --delimiter="\t" \
      --nth=1 \
      --with-nth=2 \
      --scheme="path" \
      --no-info \
      --no-scrollbar \
      --ansi \
      --preview="git -c color.status=always -C {1} status" \
      --preview-window="right:40%:wrap" \
      --preview-label="GIT STATUS" |
    cut -d $'\t' -f 1
)"
```

Experiment with this and find the command that works best for you.

### Project icons

If you use a [Nerd Font] in your terminal, you can add helpful icons to the list of projects. The following example
modifies the `get_project_data` function to add a GitHub icon for projects with "github.com" in their path, and a GitLab
icon for projects with "gitlab.com" in their paths:

``` { .bash .copy title="Project launcher with icons" hl_lines="11-20" }
function get_project_data() {

  ... # NOTE: lines removed for brevity

  while IFS= read -r path; do
    pretty_path="${path#"$projects_dir"/}"

    name="$style_boldblue$(basename "$pretty_path")$style_reset"
    pretty_path="$(dirname "$pretty_path")/$name"

    if [[ "$path" == *"github.com"* ]]; then
      # Prepend white GitHub logo
      pretty_path=" \e[38;5;15m\uf113$style_reset $pretty_path"
    elif [[ "$path" == *"gitlab.com"* ]]; then
      # Prepend orange GitLab logo
      pretty_path=" \e[38;5;214m\uf296$style_reset $pretty_path"
    else
      # Prepend red Git logo for anything else
      pretty_path=" \e[38;5;124m\uf1d3$style_reset $pretty_path"
    fi

    project_data+="$path\t$pretty_path\n"
  done <<< "$projects"

  ... # NOTE: lines removed for brevity

}
```

<figure>
  <img src="../../assets/images/launcher-icons.png" alt="Project launcher with icons" />
  <figcaption>Project launcher with icons.</figcaption>
</figure>

Other ideas for icon usage:

- :material-office-building: and :material-home: for work and personal projects
- :material-language-go: :material-language-javascript: :material-language-ruby: :material-language-rust: etc. for
  project languages
- :material-star: or :material-heart: for important or favorite projects
- :material-account-hard-hat: for projects with uncommitted changes

Have a look at the [Nerd Font cheat sheet][nf-cheat-sheet] for a complete list of available icons.

!!! tip "Tip: emoji icons"
    It's also possible to use emoji icons if you don't want to use Nerd Fonts, however, the selection of emoji icons is
    quite limited compared to Nerd Fonts.

[fd]: https://github.com/sharkdp/fd
[fzf]: https://github.com/junegunn/fzf
[fd-install]: https://github.com/sharkdp/fd#installation
[fzf-install]: https://github.com/junegunn/fzf#installation
[fzf-preview]: https://github.com/junegunn/fzf#preview-window
[tmux-bind-key]: https://man.archlinux.org/man/tmux.1#bind-key
[Nerd Font]: https://www.nerdfonts.com/
[nf-cheat-sheet]: https://www.nerdfonts.com/cheat-sheet
