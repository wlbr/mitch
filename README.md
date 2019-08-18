# mitch

Simple Slack Bot, not for serious usage. Was mainly built as an exercise.

Skills:
- `help` creates this help.
- `stock <arg>` shows the currents stock price for the stock id `arg`. 
   Try `stock AAPL` or `stock UTDI.de`. _Update 2019: Skill needs rework since Yahoos stock API vanished._
- `hi` is a kind of HelloWorld.
- `echo <arg>*` echos the arguments.
- `uptime` shows the time elapsed since the startup of the bot.
- `version` version and build timestamp of the bot.
- `timein <arg>` shows the current time in city `arg`. Try `timein singapore` or `timein würzburg`. _Update 2019: Skill needs 
   reworkd to use Googles "new" account enforcements._
- `wheatherin <args>` shows the current time in city `arg`. Try `timein singapore' or `timein würzburg`. _This 
   skill is not complete._
