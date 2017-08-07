# Github Bulletin

You are working in an organization, contributing everyday in repositories.What happens when any new issue is added and assigned to you? generally you get a mail(which generally gets lost in the pool of other more important mails),And that way you lose track of the issues assigned to you until you log into your github account and check for yourself. Github Bulletin is a simple application which you can integrate with your slack and that will notify you and your other subscriber friends in your slack channel with issues that got assigned to them.

## Prerequisites
In your slack organization, just set up a bot like [this](https://api.slack.com/bot-users) and generate a [bot token](https://api.slack.com/tokens) that you will require to setup github bulletin.This bot will send notifications to subscribed users
all the issue tracking status.


## Run

```shell
$ go get github.com/shreyaganguly/github-bulletin
$ github-bulletin -github-token <github-token-that-will-fetch-the-issues> -slack-token <bot-token-you-just-generated> -org <organization for which issues will be fetched> -t <time-interval>
```

**N.B.**
* The Github Token you will be providing in the flag to set up the application must have access to the organization to fetch the issues.
* Give `t` value atleast more than 60 seconds.(I hope your code does not produce issues at this high rate! ;) )

#### Send bulletins and keep you and your friends all updated! :)
