# cook

With this program, you can keep recipes in the form of Markdown files, render
recipes to the command line, and search through arbitrary user-defined fields
in the front matter, such as `tags: [vegetarian, spicy]` or
`ingredients: [salmon]`.

This program was inspired by Jason A. Donenfeld's
[pass](https://www.passwordstore.org/).

## Installation

Homebrew:

1. Install:

        brew install ysim/cook/cook

1. Add the following line to your `.bashrc` to enable bash completion:

        source /usr/local/opt/pass/etc/bash_completion.d/cook

Download the binary:

1. Go to the releases page: <https://github.com/ysim/cook/releases>

1. Download the correct binary for your OS/arch and untar it somewhere on your
`$PATH`.

1. To enable bash completion, either download and uncompress or clone the source
code of a release. Go into the directory, then run:

        make install-bash-completion

    Then source `~/.bash_completion` in one of your startup files, like
`~/.bashrc`:

        if [[ -f ~/.bash_completion ]] ; then
            source ~/.bash_completion
        fi

    To check that the completion is working, source this file, then run:

        complete -p | grep cook

    You should get the output:

        complete -F _cook cook

Build from source:

1. Clone this repo.

1. Download and install [Go](https://golang.org/).

1. Download and install [dep](https://github.com/golang/dep/releases).

1. Run `dep ensure` to populate `vendor/` with dependencies.

1. Build the binary:

        make build GOOS=... GOARCH=...

    Possible values for `$GOOS`: `android darwin dragonfly freebsd linux nacl netbsd openbsd plan9 solaris windows zos`  
    Possible values for `$GOARCH`: `386 amd64 amd64p32 arm armbe arm64 arm64be ppc64 ppc64le mips mipsle mips64 mips64le mips64p32 mips64p32le ppc s390 s390x sparc sparc64`

    List of valid combinations of `$GOOS` and `$GOARCH`: <https://golang.org/doc/install/source#environment>

1. Follow the above instructions for bash completion.

## Usage

### Create/port recipe files

In a directory at `$HOME/.recipes`, create some recipes files with the `.md`
(Markdown) extension, and with YAML front matter block. The file contents
should look something like this:

        ---
        name: roasted cauliflower soup
        tags: [soup, vegetarian]
        ingredients: [cauliflower]
        ---
        # roasted cauliflower soup

        ## INGREDIENTS

        * 1 cauliflower, cut into florets
        * ...rest of ingredients...

        ## INSTRUCTIONS

        1. Step one

        1. Step two

        1. ...rest of the steps...

### View a recipe

    cook [recipe name]

This will print the recipe to the screen with some light styling.

### Create a new recipe

    cook new -filename=beer-bread -f='tags=[baking,bread]' -f='ingredients=[flour,yeast]'

This will create a recipe file at `${COOK_RECIPES_DIR}/beer-bread${COOK_RECIPES_EXTENSION}`
with the given attributes and open it for editing by default.

Only the `-filename` flag is mandatory; the rest can be omitted and
subsequently filled in while editing:

    $ cook new -filename=dandan-noodles

    ...new file opened in $EDITOR...

    ---
    name:
    ---
    #

    ## INGREDIENTS

    ## INSTRUCTIONS

    ---
    Source:

### List unique values for a key

    cook list -key=tags

This will list, in alphabetical order, all the unique values for the `tags` key
as defined in the recipe files' YAML front matter blocks.

### Validate the formatting of your recipe files

    cook validate

If you see any errors parsing files while using `cook`, you can run this
command to validate the formatting of your recipe files. This will walk through
`$COOK_RECIPES_DIR` and list the files with formatting that cannot be parsed by
`cook`, along with a reason.

    cook validate [recipe name]

You can also validate a single file. Bash completion is available for this
option.

### Searching

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

### Edit a recipe

    cook edit [recipe name]

This will open the recipe in a text editor.

## Customizations

Here are a few environment variables that can be set to override the default
settings:

| setting          | default          | environment variable |
| ---------------- | ---------------- | -------------------- |
| recipe directory | `$HOME/.recipes` | `COOK_RECIPES_DIR`   |
| recipe extension | `.md`            | `COOK_RECIPES_EXT`   |
| text editor      | `vim`            | `EDITOR`             |
