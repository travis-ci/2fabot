# 2fabot

**NOTE**: This project is under development, and should be considered alpha state.

2fabot is a Slack bot that will regularly check the 2fa status of a list of
users on various services and will send them a Slack message if they are
missing 2fa on some of the services, letting them know how they can enable it.

## Setup

You need a Slack bot user that the reminder messages will be sent from. This
bot user can be created by going to your team's settings page and [creating a
new bot user](https://my.slack.com/services/new/bot).

You will also need configuration for each individual service that you want to
check two-factor status for. Documentation for that is listed below.

## Configuration

2fabot is configured through environment variables. There are two global
environment variables and then a set of environment variables for each service
you wish to add.

* `SERVICES`: 
* `SLACK_TOKEN`: A token for a bot user on the Slack team you wish to send the
  reminders on. See the "Setup" section above for how to make a bot user.
* `MESSAGE_TEMPLATE`: A template for the Slack message to be sent to users with
  at least one service missing two-factor authentication. For example, it could look something like this:

    Hi! It looks like you're missing two-factor authentication on some services, could you please look into enabling it once you have a free moment?
    {{ range . }}- {{ .Description }} (how to enable: {{ .URL }})
    {{ end }}

## License

2fabot is released under the MIT license, see the `LICENSE` file in this
repository for details.
