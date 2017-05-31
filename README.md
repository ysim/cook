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

        cp completion/bash_completion ~/.bash_completion

1. Symlink the bash completion script for cook to `~/.bash_completion.d/`:

        ln -s "$(pwd)/completion/cook.bash-completion" ~/.bash_completion.d/

### Usage

1. In a directory at `$HOME/.recipes`, create some recipes files with the
`.md` (Markdown) extension, and with YAML front matter block. The file contents
should look something like this:

    ---
    name: roasted cauliflower soup
    tags: [soup, vegetarian]
    ingredients: [cauliflower]
    ---

    ### INGREDIENTS

    * 1 cauliflower, cut into florets
    * ...rest of ingredients...

    ### INSTRUCTIONS

    1. Step one

    1. Step two

    1. ...rest of the steps...

### Customizations

* Recipe directory location (default `$HOME/.recipes`): to change this, export
`COOK_RECIPES_DIR`

* Recipe extension (default `.md`): to change this, export `COOK_RECIPES_EXT`

More to come!
