
Sanat
=====

Translates a “master” string resource translation file that looks something like this:

    === Login flow ===
    
      LoginView.Title
        en = Log in
        fi = Kirjaudu sisään
        jp = ログイン

…into string resource files appropriate for use on several different software platforms:

                                 all-translations.sanat
                                 ~~~~~~~~~~~~~~~~~~~~~~
                                           |
                                        (sanat)
               ____________________________|_______________________________ ...
              |                            |                           |
           Android                     iOS / OS X                   Windows
    __________|__________    ______________|_____________    __________|_________
    values-en/strings.xml    en.lproj/Localizable.strings    AppResources-en.resx
    values-fi/strings.xml    fi.lproj/Localizable.strings    AppResources-fi.resx
    values-jp/strings.xml    jp.lproj/Localizable.strings    AppResources-jp.resx

Run the main program with the `--help` argument to see “usage” information.


Translation File Syntax
------------------------

### Section titles

_(Sections are optional — translation files don't have to contain sections.)_

Section titles begin with at least three `=` symbols at the beginning of a line:

    === Login view

Both leading and trailing `=` symbols are okay:

    ======== Login view ========


### Translations

__Translation keys__ (e.g. `LoginView.Title` below) are indented by two spaces.
__Translation values__ are indented by four spaces:

      LoginView.Title
        en = Log in
        fi = Kirjaudu sisään

Each __translation value__ line must begin with an _ISO 639-1 language name_, followed by a `=` sign, followed by the actual text content of the translation (for the specified language.)

The translation text content may be double quoted:

        fi = "Kirjaudu sisään "

#### Platform limits

Translations can be limited to certain platforms like so:

      LoginView.Title
        en = Log in
        fi = Kirjaudu sisään
        platforms = apple, android

Translations that specify platforms will only be rendered in the translation output files for those platforms (and not for others.)

The currently supported values are:

- `apple` (Apple platforms; iOS and OS X)
- `android`
- `windows` (.NET)

#### Comments

Translations can be given comments like so:

      LoginView.Title
        en = Log in
        fi = Kirjaudu sisään
        comment = The title to the login view

These kinds of comments will be included in the generated platform-specific string resource files whenever possible.

Lines that start with `#` are comments that are _not_ included in the generated resource files:

      # This comment is _not_ included in the output
      LoginView.Title
        en = Log in
        fi = Kirjaudu sisään
        comment = This comment is included in the output


### Format specifiers

Translation text can contain format specifiers like this:

    Hello {s}, today it’s {f.2} degrees celsius.

Format specifiers are delineated by the `{}` signs, and they can contain the following:

    { 3: f .2 }
      ^^ ^ ^^____ (optional) decimal count
       |  \___ data type
        \___ (optional) order index

The __data type__ can be one of the following:

- `@`: Object
- `s`: String
- `d`: Integer
- `f`: Floating-point number

(These are mapped to the closest corresponding platform-specific format specifiers.)

The __decimal count__ specifies the number of decimals to show for floating-point numbers, and is only relevant if the data type is `f`.

The __order index__ specifies the 1-based index of the “printf argument” to apply for this format specifier. (This is necessary for cases where the word order for the same sentence differs between languages.)


Preprocessors
-------------

Translations can be preprocessed by specifying one of the following preprocessors:

### `markdown`

Translates the text in each translation from Markdown to HTML.

### `smartypants`

Processes the text in each translation using [Smartypants].


[Smartypants]: https://daringfireball.net/projects/smartypants/



Design Principles
-----------------

This program is written in Go. This helps satisfy these requirements:

- **Speed:** This program needs to be fast if it’s going to be included in build processes — we don’t want it to make builds noticeably slower.
- **Portability:** This program needs to run on multiple operating systems (at least OS X, Linux, and Windows) so that it can be used on any developer workstation or build server.
- **Single-file executable binary:** This program needs to be a single-file executable binary with no external dependencies so that it’d be straightforward and simple to install, or commit into source control (so that it wouldn’t need to be installed as a separate step.)



License
-------

GPLv2.

Please see the `LICENSE` file.





