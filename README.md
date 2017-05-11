# cook

**UNDER ACTIVE DEVELOPMENT. USE AT YOUR OWN RISK!**

With this program, you can keep recipes in the form of Markdown files, render
recipes to the command line, and search through arbitrary user-defined fields
in the front matter, such as `tags: [vegetarian, spicy]` or
`ingredients: [salmon]`.

This program was inspired by Jason A. Donenfeld's
[pass](https://www.passwordstore.org/).

### Installation

1. Clone the repo.

1. Build the binary:

        make build

1. Move the binary somewhere on your `$PATH`.

*Bash completion:*

1. Create a directory at `~/.bash_completion.d`.

1. Copy this script to `~/.bash_completion` to source all the scripts
within `~/.bash_completion.d/`:

        cp `completion/bash_completion` ~/.bash_completion

1. Symlink the bash completion script for cook to `~/.bash_completion.d/`:

        ln -s "$(pwd)/completion/cook.bash-completion" ~/.bash_completion.d/
