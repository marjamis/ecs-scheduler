# How to contribute

NOTE: This is a modified version of  https://github.com/opengovernment/opengovernment/edit/master/CONTRIBUTING.md

I'm really glad you're reading this, because we need volunteer developers to help this project come to fruition.

Here are some important resources:

## Testing

## Submitting changes

Please send a [GitHub Pull Request to marjamis/ecs-scheduler](https://github.com/marjamis/ecs-scheduler/pull/new/master) with a clear list of what you've done (read more about [pull requests](http://help.github.com/pull-requests/)). We can always use more test coverage. Please follow our coding conventions (below) and make sure all of your commits are atomic (one feature per commit).

Always write a clear log message for your commits. One-line messages are fine for small changes, but bigger changes should look like this:

    $ git commit -m "A brief summary of the commit
    >
    > A paragraph describing what changed and its impact."

## Coding conventions

Start reading our code and you'll get the hang of it. We optimize for readability:

  * We indent using two spaces (soft tabs)
  * We avoid logic in views, putting HTML generators into helpers
