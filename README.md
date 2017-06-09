# Blom, the Content Manager (A WORK IN PROGRESS)

## In General

Blom is the content manager for [my blog](http://ratan.blog). Other blogs should probably avoid using it, until I heavily refactor the code.

Blom does the following tasks:

 * Takes file information from the blog administrator
 * Converts articles formatted in Markdown to articles formatted in HTML
 * Sends articles formatted in HTML through a template, allowing a static site to have a unified layout
 * Generates XML, Atom and JSON feeds
 * Generates a homepage with the latest content
 * Generates a chronological archive (Using the months of the [tranquility calendar](https://github.com/ratanvarghese/tqtime))
 * Generates an archive organized by tags
 * If run through a cron job, it can update every page on the site to show the current date

Notably, blom does not act as a file server. I use a seperate static server for my blog.

Unfortunately, this is written in a manner which makes it difficult for other people to use. Generic content managers are a dime a dozen after all. I wrote this for fun, for knowledge and so that I could stick Tranquility dates everywhere.

Here are some features not yet implemented:

 * The ability to use hostnames other than my blog.
 * The ability to use directory structures other than the one I currently use for my blog.
 * Automated testing. This doesn't really count as a feature, but nonetheless it is not available.
 * The ability to disable Tranquility dates.
 * Automatic GZIP compression or HTML/CSS minification. (On my blog, the document size is dwarfed by the font file size... even when using a Unicode subset)
 * Any concurrency at all. There are definitely stages where parallelism or concurrency could be used, but are not.

As of the writing of this README, I am trying to refactor/rewrite the code into a more general, testable form. However, since that work is not strictly necessary for my blog to function, my mind could wander...

The refactoring could be followed by slight changes in the directory structure of my blog, mostly regarding attachments. I will be sure to keep article permalinks functional.

## Actually Using Blom

At the root level of your site, there must be a directory for each article. There must be a `template.html` one level up from the root directory of the site. The following variables must exist in the `template.html`:

 * {{.Title}}
 * {{.Today}} (the server date)
 * {{.Date}} (the publication date of the current article)
 * {{.ContentHTML}}

The article content (not including the title) must be placed in a `content.md` or `content.html` in this directory. Make one of those directories your working directory and run the following.

    blom article -title YourTitleHere -attach comma,seperated,files,to,attach -tags comma,seperated,tags

Your attachments must be in the same folder as the associated article. This will likely change in a future update. If you ever wish to edit your article, change the `content.md`. If you never had a `content.md` and do not wish to add one, then change the `content.html`. After the editing run the following.

    blom article

You can overwrite the currently existing tags and attachments by using the associated options to `blom article`.

Once all your articles are ready to be added to the archives, feeds and possibly the home page, go to the root folder of the site and run the following.

    blom update

Run `blom update` every day to ensure that the statements about the current date are correct. (Blom is implemented to use the server time, not the time at the client ... this is a static site after all). I have not tried `blom update` on a large blog yet. It has to search every top-level directory of the site.

To hide an article from the feeds, homepage and archives, delete the `item.json` that `blom article` generates. Then run `blom update`.
