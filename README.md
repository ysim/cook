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

1. Recommended: source `~/.bash_completion` in one of your startup files, like
`~/.bashrc`:

        if [[ -f ~/.bash_completion ]] ; then
            source ~/.bash_completion
        fi

    To check that the completion is working, source this file, then run:

        complete -p

    The following line should appear in the output:

        complete -F _cook cook

### Usage

##### Create/port recipe files

In a directory at `$HOME/.recipes`, create some recipes files with the
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

##### Validate the formatting of your recipe files

Once you're finished porting/creating your recipe files, use the `validate`
command to check their formatting:

    cook validate

This will walk through `$COOK_RECIPES_DIR` and list the files with formatting
that cannot be parsed by `cook`, along with a reason.

##### View a recipe

To view a single recipe, run `cook [recipe name]`. Take advantage of the bash
completion available. For now, this will print the recipe HTML to the screen
(more features will be added soon to either render this markup in the command
line or to open it in a webpage). Feel free to make use of
[lynx](http://lynx.browser.org/) to view recipes for now, e.g.

    cook chicken-pot-pie | lynx -stdin

##### Searching

The following search syntaxes are supported:

* `cook search tags:soup`: show all recipes with the tag `soup`
* `cook search tags:soup,stew`: show all recipes with the tag `soup` OR
`stew`
* `cook search tags:soup+vegetarian`: show all recipes with the tag `soup`
AND `vegetarian`

More powerful search syntax coming soon!

##### Edit a recipe

    cook edit [recipe name]

This will open the recipe in text editor (defaults to `vim`). Set the `$EDITOR`
environment variable to open it with a different text editor.

### Customizations

* Recipe directory location (default `$HOME/.recipes`): to change this, export
`COOK_RECIPES_DIR`

* Recipe extension (default `.md`): to change this, export `COOK_RECIPES_EXT`

* Text editor for editing recipes (default `vim`): to change this, export
`EDITOR`.
