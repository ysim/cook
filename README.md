# cook

With this program, you can keep recipes in the form of Markdown files, render
recipes to the command line, and search through arbitrary user-defined fields
in the front matter, such as `tags: [vegetarian, spicy]` or
`ingredients: [salmon]`.

This program was inspired by Jason A. Donenfeld's
[pass](https://www.passwordstore.org/).

### Installation

Homebrew:

1. Install:

        brew install ysim/cook/cook

1. Add the following line to your `.bashrc` to enable bash completion:

        source /usr/local/opt/pass/etc/bash_completion.d/cook

Download the binary:

1. Go to the releases page: <https://github.com/ysim/cook/releases>

1. Download the tarball of the binary (should be named `cook-vx.x.x.tar.gz`,
where `x.x.x` is the version number) and move it somewhere on your `$PATH`.

1. To enable bash completion, either download and uncompress or clone the source
code of a release. Copy the file `completion/cook.bash-completion` somewhere
sensible, for example `~/.bash_completion.d`.

        mkdir ~/.bash_completion.d
        cp completion/cook.bash-completion ~/.bash_completion.d

    Copy the file `completion/bash_completion`:

        cp completion/bash_completion ~/.bash_completion

    Then source `~/.bash_completion` in one of your startup files, like
`~/.bashrc`:

        if [[ -f ~/.bash_completion ]] ; then
            source ~/.bash_completion
        fi

    To check that the completion is working, source this file, then run:

        complete -p

    The following line should appear in the output:

        complete -F _cook cook

Build from source:

1. Clone this repo.

1. Download and install [Go](https://golang.org/).

1. Build the binary:

        make build

1. Follow the above instructions for bash completion.

### Usage

##### Create/port recipe files

In a directory at (or symlinked to) `$HOME/.recipes`, create some recipes files
with the `.md` (Markdown) extension, and with YAML front matter block. The file
contents should look something like this:

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

##### Edit a recipe

    cook edit [recipe name]

This will open the recipe in text editor (defaults to `vim`). Set the `$EDITOR`
environment variable to open it with a different text editor.

### Customizations

Here are a few environment variables that can be set to override the default
settings:

| setting          | default          | environment variable |
| ---------------- | ---------------- | -------------------- |
| recipe directory | `$HOME/.recipes` | `COOK_RECIPES_DIR`   |
| recipe extension | `.md`            | `COOK_RECIPES_EXT`   |
| text editor      | `vim`            | `EDITOR`             |
