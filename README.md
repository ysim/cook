# cook

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

##### View a recipe

    cook [recipe name]

This will print the recipe to the screen with some light styling.

##### Validate the formatting of your recipe files

    cook validate

If you see any errors parsing files while using `cook`, you can run this
command to validate the formatting of your recipe files. This will walk through
`$COOK_RECIPES_DIR` and list the files with formatting that cannot be parsed by
`cook`, along with a reason.

    cook validate [recipe name]

You can also validate a single file. Bash completion is available for this
option.

##### Searching

The following search syntaxes are supported:

* `cook search tags:soup`: show all recipes with the tag `soup`
* `cook search tags:soup,stew`: show all recipes with the tag `soup` OR
`stew`
* `cook search tags:soup+vegetarian`: show all recipes with the tag `soup`
AND `vegetarian`

Note that values with a space in them must either be quoted or escaped:

❌ `cook search tags:comfort food`  
✅ `cook search 'tags:comfort food'`  
✅ `cook search tags:'comfort food'`  
✅ `cook search tags:comfort\ food`

More powerful search syntax coming soon; next up, support for multi-field
searching (e.g. `tags:lazy AND ingredients:bacon`)

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
